FROM golang:latest

COPY . /go/src/github.com/ahmet/phonenumber/
WORKDIR /go/src/github.com/ahmet/phonenumber/
RUN go get ./
RUN go build -o main .
CMD ["./main"]
EXPOSE 8080
