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
	"fmt"

	pb "github.com/vagababov/maxprimesrv/proto"
	"go.uber.org/zap"
)

type primeServer struct {
	logger *zap.SugaredLogger
}

func (ps *primeServer) Get(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	if req.Query <= 1 {
		return nil, fmt.Errorf("%d is not a valid input", req.Query)
	}
	a := calcPrime(req.Query)
	return &pb.Response{
		Answer: a,
	}, nil
}
