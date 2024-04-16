FROM golang:1.21.3 AS builder
RUN apt-get update -y && apt-get upgrade -y \
	&& apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workdpaces/
COPY . .
RUN go mod download
RUN go build -o /go/bin/res ./cmd/piglet-auth

FROM ubuntu:22.04 AS runner
WORKDIR /app
ENV CONFIG_PATH="./config/piglet-auth/local.yaml"
COPY --from=builder /workdpaces/config/piglet-auth ./config/piglet-auth
COPY --from=builder /go/bin/res .
EXPOSE 8080
ENTRYPOINT ["./res"]