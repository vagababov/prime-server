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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

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

	if *negate {
		resp.Answer = -resp.Answer
	}
	fmt.Printf("Resp: %#v Negate: %v", resp, *negate)

	logger.Infof("Response: %#v", resp)
	stream, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error formatting answer: %#v\r\n", err)
		return
	}
	w.Write(stream)
}

func ReadResponse(r io.Reader) (*pb.Response, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		logger.Errorw("error reading the body", zap.Error(err))
		return nil, err
	}
	resp := &pb.Response{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Errorw("error unmarshaling json", zap.Error(err))
		return nil, err
	}
	return resp, nil
}

func ReadRequest(r io.Reader) (*pb.Request, error) {
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
