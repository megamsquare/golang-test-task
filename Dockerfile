FROM golang:latest

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080
EXPOSE 8081
EXPOSE 8083

CMD ["app1", "app2"]