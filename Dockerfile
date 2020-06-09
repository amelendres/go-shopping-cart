FROM golang:1.14.4-alpine3.12

ENV APP_NAME cart
ENV PORT 5000

WORKDIR /go/src/${APP_NAME}
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["webserver"]

EXPOSE ${PORT}