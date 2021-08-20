FROM golang:1.16.7-alpine3.14 as builder

COPY go.mod go.sum /go/src/cfbackendapp/

WORKDIR /go/src/cfbackendapp

RUN go mod download

COPY . /go/src/cfbackendapp

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/app


FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /go/src/cfbackendapp/build/app /usr/bin/cfbackendapp

EXPOSE 5000 5000

ENTRYPOINT [ "/usr/bin/cfbackendapp" ]