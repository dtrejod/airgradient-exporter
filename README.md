# AirGradient Exporter

This is a simple prometheus exporter for the AirGradient air quality monitor. It uses the [AirGradient
LocalServer API](https://github.com/airgradienthq/arduino/blob/master/docs/local-server.md) to get the data.

*NOTE: Usage of LocalServer API requires the device to be running AirGradient firmware version 3.0.10 or later.*

## Usage

Whether running as a container or manually, the exporter requires you to know the `ENDPOINT` of your AirGradient device.
The `ENDPOINT` should can be easily obtained if you know the serial number of your AirGradient device. AirGradient
devices use mDNS to easily discover the device on the network. The device can be accessed at
`http://airgradient_<SERIAL>.local`. The `ENDPOINT` should be set to this URL.

Once running, the exporter, by default, will expose the metrics at `:9091/metrics`.

### Docker Image
The exporter is available as a docker image on GitHub Container Registry. You can run the docker image with the
following docker-compose configuration:

```
airgradient-exporter:
  image:  ghcr.io/dtrejod/airgradient-exporter:latest
  container_name: airgradient-exporter
  restart: always
  ports:
    - "9091:9091"
  environment:
    - ENDPOINT=http://airgradient_<SERIAL>.local
```


### Running Locally
Run the exporter with the following command:

```bash
./airgradient-exporter exporter --endpoint http://airgradient_<SERIAL>.local
```

## Development

The exporter is written in Go. The exporter can be built as a docker image or locally.

### Docker Build
To build the exporter as a docker image, run the following command:

```bash
docker build .
```

### Local Build
To build the exporter locally, run the following command:

```bash
go build
```

The exporter will then be available as `airgradient-exporter`.


