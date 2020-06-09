FROM golang:1.13.12-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git bash openssh

WORKDIR /app

COPY go.mod ./

RUN go build

COPY . .

#RUN go build -o main .

#EXPOSE 5000

#CMD ["./main"]