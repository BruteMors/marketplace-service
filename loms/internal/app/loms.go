package app

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/config"
	"github.com/BruteMors/marketplace-service/loms/internal/controller/grpcapi/middleware"
	httpmiddleware "github.com/BruteMors/marketplace-service/loms/internal/controller/httpapi/middleware"
	"github.com/BruteMors/marketplace-service/loms/internal/metric"
	"github.com/BruteMors/marketplace-service/loms/pkg/api/grpc/loms/v1"
	"github.com/BruteMors/marketplace-service/loms/pkg/closer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type LomsApp struct {
	serviceProvider   *serviceProvider
	grpcServer        *grpc.Server
	grpcGatewayServer *http.Server
	shutdownTracer    func(context.Context) error
}

func NewLoms(ctx context.Context) (*LomsApp, error) {
	a := &LomsApp{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (l *LomsApp) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := l.shutdownTracer(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := l.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := l.runGRPCGatewayServer()
		if err != nil {
			log.Fatalf("failed to run gRPC gateway server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (l *LomsApp) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		l.initConfig,
		l.initServiceProvider,
		l.initMetrics,
		l.initTracing,
		l.initGRPCServer,
		l.initGRPCGatewayServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LomsApp) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (l *LomsApp) initServiceProvider(_ context.Context) error {
	l.serviceProvider = newServiceProvider()
	return nil
}

func (l *LomsApp) initMetrics(ctx context.Context) error {
	err := metric.Init(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (l *LomsApp) initTracing(ctx context.Context) error {
	l.shutdownTracer = tracing.InitTracer()
	return nil
}

func (l *LomsApp) initGRPCServer(ctx context.Context) error {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.Panic,
			middleware.GetTraceID,
			middleware.RequestLogger,
			middleware.Validate,
			middleware.ErrorHandler,
			middleware.RequestMetric,
		),
	)

	reflection.Register(grpcServer)

	l.grpcServer = grpcServer

	loms.RegisterOrdersServer(grpcServer, l.serviceProvider.OrderGRPCApi(ctx))
	loms.RegisterStockServer(grpcServer, l.serviceProvider.StockGRPCApi(ctx))

	return nil
}

func (l *LomsApp) runGRPCServer() error {
	address := l.serviceProvider.GRPCServerConfig().Address()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	log.Printf("gRPC server is running on %s", address)

	err = l.grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (l *LomsApp) initGRPCGatewayServer(ctx context.Context) error {
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Failed to dial:", err)
		return err
	}

	gwmux := runtime.NewServeMux()

	if err = loms.RegisterOrdersHandler(ctx, gwmux, conn); err != nil {
		slog.Error("Failed to register orders gateway:", err)
		return err
	}

	if err = loms.RegisterStockHandler(ctx, gwmux, conn); err != nil {
		slog.Error("Failed to register stock gateway:", err)
		return err
	}

	fs := http.FileServer(http.Dir("./public/swagger-ui"))

	mux := http.NewServeMux()
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", fs))
	mux.HandleFunc("/api/openapiv2/loms.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./api/openapiv2/loms.swagger.json")
	})
	mux.Handle("/", gwmux)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))

	mux.Handle("/metrics", promhttp.Handler())

	grpcGatewayServer := &http.Server{
		Addr:    l.serviceProvider.HTTPServerConfig().Address(),
		Handler: httpmiddleware.RequestLogger(mux),
	}

	l.grpcGatewayServer = grpcGatewayServer

	return nil
}

func (l *LomsApp) runGRPCGatewayServer() error {
	log.Printf("HTTP grpc-gateway server is running on %s", l.serviceProvider.HTTPServerConfig().Address())

	err := l.grpcGatewayServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
