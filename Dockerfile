FROM golang:1.17.8-bullseye

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

EXPOSE 8000

RUN go build -o /go-gin-demo

CMD ["/go-gin-demo"]
