FROM piglet-auth_base:0.1.0 AS builder

WORKDIR /workspaces/piglet
COPY ./piglet-auth .

RUN go build -o /go/bin/res ./cmd/main.go

FROM ubuntu:22.04 AS runner
WORKDIR /app
ENV CONFIG_PATH="./config/local.yaml"
COPY --from=builder /workspaces/piglet/config ./config
COPY --from=builder /go/bin/res .
EXPOSE 8080
ENTRYPOINT ["./res"]