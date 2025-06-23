# AirGradient Exporter

This is a simple prometheus exporter for the AirGradient air quality monitor. It uses the [AirGradient
LocalServer API](https://github.com/airgradienthq/arduino/blob/master/docs/local-server.md) to get the data.

## Compatibility

Below is a compatibility matrix for verified versions of the exporter and the AirGradient device firmware.

| AirGradient Firmware | Exporter Version |
|----------------------|------------------|
| `x.x.x -> 3.0.9`     | Not Supported    |
| `3.0.10 -> 3.0.11`   | `>=0.2.0`        |
| `3.0.12`             | `<=0.3.0`        |

Your mileage may vary with newer versions of the firmware.

## Usage

Whether running as a container or manually, the exporter requires you to know the `ENDPOINT` of your AirGradient device.
The `ENDPOINT` should can be easily obtained if you know the serial number of your AirGradient device. AirGradient
devices use mDNS to easily discover the device on the network. The device can be accessed at
`http://airgradient_<SERIAL>.local`. The `ENDPOINT` should be set to this URL.

**NOTE: When running the exporter as a container, the `ENDPOINT` should be a static IP address or hostname that can be
resolved by the container. The container does not have access to the mDNS service.**

Once running, the exporter, by default, will expose the metrics at `:9091/metrics`.

### HomeKit Integration

The exporter supports running as a HomeKit bridge to expose the AirGradient device as a HomeKit accessory. To enable
this feature, set the `ENABLE_HOMEKIT=true` environment variable. The exporter will then expose the AirGradient
device as a HomeKit accessory.

When HomeKit is enabled, disk storage is used to persist the HomeKit state. The exporter by default will create a
directory `./data` to store the HomeKit state. This directory can be changed by setting the `DATA_DIR`.

To add this bridge, go to your Home app on iOS or macOS, and select "Add Accessory". You will add the accessory
by selecting "More options..." and selecting "AirGradient Exporter" from the discovered accessories. When pairing the
device to HomeKit, you will need to use the `PIN=1845-8232` code.

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
    - ENDPOINT=<ip of airgradient device>
    - ENABLE_HOMEKIT=true # optional
    - DATA_DIR=/data # required if ENABLE_HOMEKIT=true is set
  volumes:
    - airgradient-exporter-data:/data # required if ENABLE_HOMEKI=true is set
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


