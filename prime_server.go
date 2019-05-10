package main

import (
	"context"
	"errors"

	"github.com/vagababov/maxprimesrv/proto"
	"go.uber.org/zap"
)

type primeServer struct {
	logger *zap.SugaredLogger
}

func (ps *primeServer) Get(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	return nil, errors.New("not yet implemented")
}
