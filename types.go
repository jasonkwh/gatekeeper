package gatekeeper

import (
	"reflect"

	"google.golang.org/grpc"
)

type incomingRequest interface {
	Validate() error
}

type recvWrapper struct {
	grpc.ServerStream
}

var (
	validators map[reflect.Type]bool
)
