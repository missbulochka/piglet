FROM golang:1.21.3 AS builder

WORKDIR /workspaces/piglet
COPY ./piglet-gateway .

RUN go build -o /go/bin/res ./cmd/main.go

FROM ubuntu:22.04 AS runner
WORKDIR /app
COPY --from=builder /go/bin/res .
COPY --from=builder /workspaces/piglet/manage.env .
EXPOSE 8083
ENTRYPOINT ["./res"]