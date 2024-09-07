# Use the distroless static-debian12 base image
FROM gcr.io/distroless/static-debian12

# Set the working directory
WORKDIR /app

# Copy your binary into the image
# Binary build path is referenced from Github workflow
COPY ./bin/airgradient-exporter-amd64-linux /app/airgradient-exporter

# Set the ENDPOINT environment variable
ENV ENDPOINT=""

# Command to run the exporter with the passed endpoint
CMD ["/app/airgradient-exporter", "exporter", "--endpoint", "${ENDPOINT}"]
