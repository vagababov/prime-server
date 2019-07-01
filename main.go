/*
Copyright 2019 Victor Agababov
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/vagababov/prime-server/proto"
)

const (
	defaultPort      = "8080"
	portVariableName = "PORT"

	defaultGRPCPort      = "8081"
	grpcPortVariableName = "GRPC_PORT"
)

var negate = flag.Bool("negate", false, "Negates the result for display")

var logger *zap.SugaredLogger

func main() {
	flag.Parse()

	tl, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer tl.Sync()
	logger = tl.Sugar()

	port := os.Getenv(portVariableName)
	if port == "" {
		port = defaultPort
	}
	logger.Infof("Starting prime generator on port %s", port)
	http.Handle("/", http.HandlerFunc(handler))
	http.Handle("/q", http.HandlerFunc(qhandler))
	go http.ListenAndServe(":"+port, nil)

	grpcPort := os.Getenv(grpcPortVariableName)
	if grpcPort == "" {
		grpcPort = defaultGRPCPort
	}

	logger.Infof("Starting GRPC prime generator on port %s", grpcPort)
	ls, err := net.Listen("tcp", "localhost:"+grpcPort)
	if err != nil {
		logger.Fatalf("Error listening on port %s: %v", grpcPort, err)
	}
	primeSrv := &primeServer{}
	grpcServer := grpc.NewServer()
	pb.RegisterPrimeServiceServer(grpcServer, primeSrv)
	grpcServer.Serve(ls)
}
