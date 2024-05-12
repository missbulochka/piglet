FROM piglet-transactions_base:0.1.0 AS builder

WORKDIR /workspaces/piglet
COPY ./piglet-transactions .

RUN go build -o /go/bin/res ./cmd/main.go

FROM ubuntu:22.04 AS runner
WORKDIR /app
COPY --from=builder /go/bin/res .
COPY --from=builder /workspaces/piglet/migration ./migration
COPY --from=builder /workspaces/piglet/transactions.env .
COPY --from=builder /workspaces/piglet/pg-db.env .
EXPOSE 8081
ENTRYPOINT ["./res"]