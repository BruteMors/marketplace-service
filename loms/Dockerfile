FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /loms ./cmd/loms/main.go

FROM scratch

COPY --from=builder /loms /bin/loms
COPY --from=builder /build/.env /.env

ENTRYPOINT ["/bin/loms"]
