# building

# 1.19 or latest
FROM golang:1.19 AS builder

WORKDIR /app

#copy sources
COPY cmd cmd
COPY internal internal

COPY .env .env

COPY go.mod go.mod
COPY go.sum go.sum

# go mod
#RUN go mod vendor

# build
RUN go build -o order-server cmd/main.go

# prod
#FROM alpine:latest as production
FROM gcr.io/distroless/base-debian10 as production

# copy builded app and env file
COPY --from=builder app/order-server /
COPY --from=builder app/.env /

EXPOSE 8080

ENTRYPOINT ["/order-server"]