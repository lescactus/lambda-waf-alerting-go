# lambda-waf-alerting-go

This repository contains the code for a lambda function used to send CloudWatch WAFv2 alerts in a given Slack channel.

## Build the lambda

The source code is in the `function/` directory

### From source with go

You need a working [go](https://golang.org/doc/install) toolchain (It has been developped and tested with go 1.14 only, but should work with go >= 1.11 ). Refer to the official documentation for more information (or from your Linux/Mac/Windows distribution documentation to install it from your favorite package manager).

```sh
# Clone this repository
git clone https://github.com/lescactus/lambda-waf-alerting-go.git && cd lambda-waf-alerting-go

# Build from sources. Use the '-o' flag to change the compiled binary name
cd function/
GOOS=linux CGO_ENABLED=0 go build -o main

# Zip the binary to upload it to Lambda
zip main.zip main
```

### From source with docker

If you don't have [go](https://golang.org/) installed but have docker, run the following command to build inside a docker container:

```sh
# Build from sources inside a docker container. Use the '-o' flag to change the compiled binary name
# Warning: the compiled binary belongs to root:root
docker run -e GOOS=linux -e CGO_ENABLED=0 --rm -it -v "$PWD":/app -w /app golang:1.14 go build -o main

# Zip the binary to upload it to Lambda
zip main.zip main
```