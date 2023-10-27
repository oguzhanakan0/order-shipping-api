FROM golang:alpine

WORKDIR /go/src/
COPY . .
RUN go build .

EXPOSE 8080

CMD ["go", "run", "main.go"]