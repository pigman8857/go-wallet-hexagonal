# Stage 1: build the binary. Alpine keeps this stage small; none of it ships in the final image.
FROM golang:1.27-alpine AS builder

# ca-certificates: needed later to copy real CA certs into the certless distroless stage.
# tzdata: lets the binary do timezone-aware time handling without relying on the host.
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy only dependency manifests first so this layer is cached and skipped
# on rebuilds where only source code changed, not go.mod/go.sum.
COPY go.mod go.sum ./

# similar to `npm i` — downloads deps into the build cache, reusing the layer above.
RUN go mod download

# Now copy the rest of the source; only invalidates the cache from here down.
COPY . .

# CGO_ENABLED=0: static binary, no libc dependency — required to run on distroless/scratch.
# GOOS/GOARCH pinned: reproducible build regardless of the host building the image.
# -ldflags="-s -w": strips debug symbols/DWARF info to shrink the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags="-s -w" -o /app/server ./cmd/main.go


# Stage 2: runtime image. Only the compiled binary and CA certs make it in —
# no shell, no package manager, no Go toolchain, tiny attack surface and image size.
FROM gcr.io/distroless/static-debian12:nonroot

# Distroless ships no CA bundle; without this, outbound TLS (e.g. Postgres SSL) fails.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

COPY --from=builder /app/server /server

EXPOSE 8080

# nonroot base image already runs this as an unprivileged user.
CMD ["/server"]

