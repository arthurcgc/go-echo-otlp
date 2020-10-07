# Open Telemetry: Collector + OTLP Receiver + Jaeger Exporter Example

A Hello World example to showcase how to use [Open Telemetry Collector](https://github.com/open-telemetry/opentelemetry-collector) using the open telemetry standard for application tracing, in conjuction with [otelecho](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/github.com/labstack/echo/otelecho)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

To run the example, one needs docker and subsequently docker-compose installed and, of course, a version of Go
To test the integration it would be nice to have curl or PostMan installed

```
On Arch Linux:
yay -Sy docker
```

### Installing

To run the example, with all the prerequisites met, just do a simple 'make' on the root of the project:
```
make build
```

Or

```
make
```

To check if everything is in ship-shape, just run a curl in your local terminal to http://localhost:9999/hello and check the Jager UI at: [](http://localhost:16686) for the tracing info

Example:

```
curl localhost:9999/hello
hello world!
```
