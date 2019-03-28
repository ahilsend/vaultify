FROM golang:1.12 as builder

WORKDIR /build
COPY go.* ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg

RUN CGO_ENABLED=0 go build cmd/vaultify/vaultify.go

FROM alpine:3.8

COPY --from=builder /build/vaultify /bin/vaultify
USER 65535
ENTRYPOINT ["/bin/vaultify"]
