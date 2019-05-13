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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/vagababov/maxprimesrv/proto"
)

const (
	defaultPort      = "8080"
	portVariableName = "PORT"

	defaultGRPCPort      = "8081"
	grpcPortVariableName = "GRPC_PORT"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query, err := readRequest(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing input: %v\r\n", err)
		return
	}
	logger.Infof("Request: %#v", query)
	resp := &pb.Response{
		Answer: calcPrime(query.Query),
	}
	logger.Infof("Response: %#v", resp)
	stream, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error formatting answer: %#v\r\n", err)
		return
	}
	w.Write(stream)
}

func readRequest(r io.Reader) (*pb.Request, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		logger.Errorw("error reading the body", zap.Error(err))
		return nil, err
	}
	req := &pb.Request{}
	err = json.Unmarshal(data, req)
	if err != nil {
		logger.Errorw("error unmarshaling json", zap.Error(err))
		return nil, err
	}
	return req, nil

}

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
	go http.ListenAndServe(":"+port, http.HandlerFunc(handler))

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
