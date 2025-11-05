
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o consumer cmd/consumer/main.go

FROM golang:1.24-alpine AS development

WORKDIR /app

RUN apk add --no-cache git make bash curl

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY --from=builder /app .


FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates bash curl

WORKDIR /root/

COPY --from=builder /app/api .
COPY --from=builder /app/consumer .
COPY --from=builder /app/migrations ./migrations

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

EXPOSE 8080

CMD ["./api"]