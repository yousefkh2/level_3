# STAGE 1: build the Go binary

FROM golang:1.25-alpine AS builder

# install build dependencies (git for private repos, ca-cert. for HTTPS)
RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./

# download dependencies (cached unless go.mod/go.sum change)
RUN go mod download

# copy source code
COPY . .

# build the binary
# CGO not enabled: static binary (no C dependencies)
# -ldflags="-w -s": strip debug info (smaller binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o api ./main.go



# STAGE 2: minimal runtime image with distroless
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary from builder stage
COPY --from=builder /build/api /app/api

EXPOSE 8080

# distroless already runs as nonroot user
# K8s probes are better than Docker healthchecks be handled

# Set working directory and run
WORKDIR /app
CMD ["/app/api"]