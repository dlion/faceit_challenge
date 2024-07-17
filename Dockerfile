# Dockerfile
FROM golang:1.22 as builder

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /user-service cmd/user-service/main.go

FROM gcr.io/distroless/base as app

COPY --from=builder /user-service /user-service

EXPOSE 80

CMD ["/user-service"]
