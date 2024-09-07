# Use the distroless static-debian12 base image
FROM gcr.io/distroless/static-debian12

# Arguments that will be passed from the build command
ARG GOOS
ARG GOARCH

# Set the working directory
WORKDIR /app

# Copy your binary into the image
# Binary build path is referenced from Github workflow
COPY ./airgradient-exporter-${GOARCH}-${GOOS} /app/airgradient-exporter

# Set the ENDPOINT environment variable
ENV ENDPOINT=""

# Command to run the exporter with the passed endpoint
CMD ["/app/airgradient-exporter", "exporter", "--endpoint", "${ENDPOINT}"]
