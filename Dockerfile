FROM golang:latest AS build

# Arguments that will be passed from the build command
ARG VERSION=""

WORKDIR /build

# Copy source code into build container
COPY . .

# Build the binary
RUN go build -v -ldflags="-X github.com/dtrejod/airgradient-exporter/version.version=${VERSION}" -o airgradient-exporter

# Use SCRATCH base image
FROM scratch

# Set the working directory
WORKDIR /app

# Copy your binary into the image
COPY --from=build /build/airgradient-exporter /app/airgradient-exporter

# Set the ENDPOINT environment variable
ENV ENDPOINT=""

# Command to run the exporter with the passed endpoint
CMD ["/app/airgradient-exporter", "exporter", "--endpoint", "${ENDPOINT}"]

# Add common labels
LABEL org.opencontainers.image.title="AirGradient Exporter" \
      org.opencontainers.image.description="Prometheus exporter AirGradient ONE" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.source="https://github.com/dtrejod/airgradient-exporter" \
      org.opencontainers.image.licenses="Apache-2.0"
