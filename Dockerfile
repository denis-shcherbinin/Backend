FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./
RUN go mod download
RUN go build -o network-backend ./cmd/app/main.go

CMD ["./network-backend"]
