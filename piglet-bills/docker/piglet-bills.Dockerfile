FROM piglet_base:0.1.0 AS builder

WORKDIR /workspaces/piglet
COPY ./piglet-bills .

RUN go build -o /go/bin/res ./cmd/main.go

FROM ubuntu:22.04 AS runner
WORKDIR /app
COPY --from=builder /go/bin/res .
EXPOSE 8080
ENTRYPOINT ["./res"]