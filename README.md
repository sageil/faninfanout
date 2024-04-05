# Demo of fan in fan out pattern using Go

## Purpose

This demo aims to show how to efficiently use the fan in fan out pattern in Go when one of the pipeline stages is time-consuming. In this demo, finding the prime number stage is the slowest stage when applying the naive implementation of the pattern. The efficient implementation shows performance improvement (see screen capture below) when utilizing all machine resources effectively.

### Running the code in Docker

1. Build the docker image by running

```bash
docker image build -t faninout .
```

To verify the efficient implementation of the pattern start the container by running

```bash
docker container run --rm faninout:latest efficient
```

To verify the naive implementation of the pattern start the container by running

```bash
docker container run --rm faninout:latest naive
```

## Running the code locally

1. To verify the naive implementation of the pattern run

```bash
go run cmd/main.go naive
```

2. To verify the efficient implementation of the pattern run

```bash
go run cmd/main.go efficient
```

<p align="center">
  <img  src="assets/results.png?raw=true">
</p>
