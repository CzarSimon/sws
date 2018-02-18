FROM golang:1.8-alpine as build
RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR /usr/app/src/api
COPY . .

RUN go get github.com/CzarSimon/httputil \
           github.com/CzarSimon/sws/pkg/user \
           github.com/CzarSimon/sws/pkg/service \
           github.com/CzarSimon/go-endpoint \
           github.com/CzarSimon/sws/pkg/swsutil \
           github.com/lib/pq

RUN go build

FROM alpine:3.6 as run
WORKDIR /root/
COPY --from=build /usr/app/src/api/api api
CMD ["./api"]
