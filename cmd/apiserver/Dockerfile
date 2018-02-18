FROM golang:1.8-alpine as build

WORKDIR /usr/app/src/api
COPY . .

RUN dep ensure
RUN go build

FROM scratch as run
WORKDIR /root/
COPY --from=build /usr/app/src/api/api api
CMD ["./api"]
