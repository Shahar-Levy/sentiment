FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build sentiment/main.go

EXPOSE 8080

CMD ["./main"]