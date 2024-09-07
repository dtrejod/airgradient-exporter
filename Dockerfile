# Use SCRATCH base image
FROM scratch

# Arguments that will be passed from the build command
ARG GOOS
ARG GOARCH

# Set the working directory
WORKDIR /app

# Copy your binary into the image
COPY ./airgradient-exporter-${GOARCH}-${GOOS} /app/airgradient-exporter

# Set the ENDPOINT environment variable
ENV ENDPOINT=""

# Set the USER to nobody to avoid running the exporter as root
USER nobody

# Command to run the exporter with the passed endpoint
CMD ["/app/airgradient-exporter", "exporter", "--endpoint", "${ENDPOINT}"]
