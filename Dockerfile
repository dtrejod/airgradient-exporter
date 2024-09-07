FROM golang:latest AS build

# Arguments that will be passed from the build command
ARG VERSION=""

WORKDIR /build

# Copy source code into build container
COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -v -ldflags="-X github.com/dtrejod/airgradient-exporter/version.version=${VERSION}" -o airgradient-exporter

# Use SCRATCH base image
FROM gcr.io/distroless/static:nonroot

# Arguments that will be passed from the build command
ARG VERSION=""

# Set the working directory
WORKDIR /app

# Copy your binary into the image
COPY --from=build /build/airgradient-exporter /usr/bin/airgradient-exporter

# Set the default values for user configurable environment variables
ENV ENDPOINT=""
ENV LISTEN_ADDRESS=":9091"

# Expose the port
EXPOSE 9091

# Set entrypoint and command
ENTRYPOINT ["airgradient-exporter", "exporter"]

# Add common labels
LABEL org.opencontainers.image.title="AirGradient Exporter" \
      org.opencontainers.image.description="Prometheus exporter AirGradient ONE" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.source="https://github.com/dtrejod/airgradient-exporter" \
      org.opencontainers.image.licenses="Apache-2.0"
