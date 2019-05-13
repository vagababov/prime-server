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
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "github.com/vagababov/maxprimesrv/proto"
)

var (
	srvAddr  = flag.String("server_address", "127.0.0.1:8081", "The endpoint to send requests to")
	srvHost  = flag.String("srv_host", "", "The host header to set on the requests")
	insecure = flag.Bool("insecure", true, "true if we want to skip SSL certificate")
	query    = flag.Int64("query", 42, "The value to query")
)

func main() {
	flag.Parse()

	fmt.Printf("Connecting to %s with host: %q\n", *srvAddr, *srvHost)

	var opts []grpc.DialOption
	if *srvHost != "" {
		opts = append(opts, grpc.WithAuthority(*srvHost))
	}
	if *insecure {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*srvAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPrimeServiceClient(conn)

	resp, err := client.Get(
		context.Background(),
		&pb.Request{
			Query: *query,
		})
	if err != nil {
		fmt.Printf("Error calling Get: %+v\n", err)
		return
	}
	fmt.Printf("Response: %+v\n", resp)
}
