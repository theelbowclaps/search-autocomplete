FROM golang:latest
RUN mkdir /app 
COPY . /app/ 
WORKDIR /app 
RUN go get github.com/gomodule/redigo/redis
RUN go get github.com/rs/cors
RUN go build -o bin/main ./src
CMD ["/app/bin/main"]