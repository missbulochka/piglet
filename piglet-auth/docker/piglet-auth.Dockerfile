FROM piglet-auth_base:0.1.0 AS builder

WORKDIR /workspaces/piglet
COPY . .

RUN go build -o /go/bin/res ./cmd/piglet-auth/main.go
#CMD ["ls", "-al", ""]

FROM ubuntu:22.04 AS runner
WORKDIR /app
ENV CONFIG_PATH="./config/piglet-auth/local.yaml"
COPY --from=builder /workspaces/piglet/config/piglet-auth ./config/piglet-auth
COPY --from=builder /go/bin/res .
EXPOSE 8080
ENTRYPOINT ["./res"]