FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /cart ./cmd/cart/main.go

FROM scratch

COPY --from=builder /cart /bin/cart
COPY --from=builder /build/.env /.env

ENTRYPOINT ["/bin/cart"]
