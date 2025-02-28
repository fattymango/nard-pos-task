FROM golang:1.23.4 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/multi_tenant/main.go

FROM alpine:latest  

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]