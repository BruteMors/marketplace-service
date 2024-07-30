module github.com/BruteMors/marketplace-service/notifier

go 1.22

replace github.com/BruteMors/marketplace-service/libs => ../libs

require (
	github.com/joho/godotenv v1.5.1
	github.com/BruteMors/marketplace-service/libs v0.0.0-00010101000000-000000000000
)

require (
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
)
