package app

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/BruteMors/marketplace-service/cart/internal/config"
	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/middleware"
	"github.com/BruteMors/marketplace-service/cart/internal/metric"
	"github.com/BruteMors/marketplace-service/cart/pkg/closer"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CartApp struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
	shutdownTracer  func(context.Context) error
}

func NewCart(ctx context.Context) (*CartApp, error) {
	a := &CartApp{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (c *CartApp) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := c.shutdownTracer(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := c.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (c *CartApp) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		c.initConfig,
		c.initServiceProvider,
		c.initMetrics,
		c.initTracing,
		c.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CartApp) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (c *CartApp) initServiceProvider(_ context.Context) error {
	c.serviceProvider = newServiceProvider()
	return nil
}

func (c *CartApp) initMetrics(ctx context.Context) error {
	err := metric.Init(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *CartApp) initTracing(_ context.Context) error {
	c.shutdownTracer = tracing.InitTracer()
	return nil
}

func (c *CartApp) initHTTPServer(ctx context.Context) error {

	cart := c.serviceProvider.CartHttpApi(ctx)

	mux := http.NewServeMux()
	mux.Handle("POST /user/{user_id}/cart/{sku_id}", middleware.TraceID(middleware.RequestMetric(middleware.RequestLogger(middleware.ErrorWrapper(cart.AddItem)))))
	mux.Handle("POST /user/{user_id}/cart/checkout", middleware.TraceID(middleware.RequestMetric(middleware.RequestLogger(middleware.ErrorWrapper(cart.Checkout)))))
	mux.Handle("DELETE /user/{user_id}/cart/{sku_id}", middleware.TraceID(middleware.RequestMetric(middleware.RequestLogger(middleware.ErrorWrapper(cart.DeleteItem)))))
	mux.Handle("DELETE /user/{user_id}/cart", middleware.TraceID(middleware.RequestMetric(middleware.RequestLogger(middleware.ErrorWrapper(cart.DeleteItemsByUserID)))))
	mux.Handle("GET /user/{user_id}/cart/list", middleware.TraceID(middleware.RequestMetric(middleware.RequestLogger(middleware.ErrorWrapper(cart.GetCart)))))
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))

	c.httpServer = &http.Server{
		Addr:    c.serviceProvider.HTTPServerConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (c *CartApp) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", c.serviceProvider.HTTPServerConfig().Address())

	err := c.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
