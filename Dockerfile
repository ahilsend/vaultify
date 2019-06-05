FROM golang:1.12 as builder

WORKDIR /build
COPY go.* ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg

ARG COMMIT_HASH="HEAD"

RUN CGO_ENABLED=0 go build \
      -ldflags "-X github.com/ahilsend/vaultify/pkg.CommitHash=${COMMIT_HASH}" \
      cmd/vaultify/vaultify.go

FROM alpine:3.8

COPY --from=builder /build/vaultify /bin/vaultify
USER 65535
ENTRYPOINT ["/bin/vaultify"]
