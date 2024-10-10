FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download


COPY . .

RUN CGO_ENABLED=0  GOOS=linux go build -o /app/main ./cmd/main.go


FROM alpine AS runner
WORKDIR app
COPY --from=builder app/main /app/main

EXPOSE 8022
ENTRYPOINT ["/app/main"]