# Example application for the Knative demonstration

This is a demo application used to showcase deployment of a service to a
Kubernetes cluster that uses Knative for serving using HTTP and gRPC as the L7
protocols.

## Prior Art

The app included here is a fork of the application written by Mark Chmarny and
located here: https://github.com/mchmarny/maxprime.

## Description

Here we have an application that computes the larges prime less than given
integer. 
The application obviously is for the demonstration purposes and is built to be
ineffective (i.e. to introduce latencies if desired).

### Pre-reqs

- go >=1.12
- ko
- protoc

## How to run locally?

To run locally, just execute `go run .` from the main drectory and it will start
an HTTP server on port 8080 and gRPC server on port 8081.

## How to run on knative

The incorporated YAML files presume you have `ko` installed.

To run HTTP service:

`ko apply -f ./http_service.yaml`

To run gRPC service:

`ko apply -f ./grpc_service.yaml`


## How to make changes?

### To change proto files

- edit the proto file
- run `protoc --go_out=plugins=grpc:.  ./*.proto` from the `proto` directory.
- update th go files as necessary

