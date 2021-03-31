FROM golang:1.15.4-alpine3.12
ENV APP_NAME cart
#ENV PORT 5000

WORKDIR /go/src/${APP_NAME}

RUN apk add --no-cache protoc make

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

#CMD ["webserver"]
CMD ["grpc"]

EXPOSE ${SERVER_PORT}