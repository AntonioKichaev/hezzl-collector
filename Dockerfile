FROM golang:latest as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux GOARH=amd64 go build -o server cmd/server/main.go
ENTRYPOINT ["/bin/app/server"]


FROM golang
WORKDIR /bin/app
COPY --from=builder /app/server server

EXPOSE 5050
ENTRYPOINT ["/bin/app/server"]
