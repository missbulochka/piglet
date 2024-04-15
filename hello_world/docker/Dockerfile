FROM hello-world_build-base:0.1.0 AS builder
WORKDIR /workdpaces/hello-world
COPY . .
RUN go build -o /go/bin/res ./cmd/

FROM ubuntu:22.04 AS runner
WORKDIR /app
COPY --from=builder /go/bin/res .
EXPOSE 8080
ENTRYPOINT ["./res"]