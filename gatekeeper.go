package gatekeeper

import (
	"context"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	validators = make(map[reflect.Type]bool)
}

// RegisterRequest - register the request that has the Validate() method.
func RegisterRequest(request incomingRequest) {
	validators[reflect.TypeOf(request)] = true
}

// RegisterRequests - register the requests that have the Validate() method.
func RegisterRequests(requests ...incomingRequest) {
	if len(requests) == 0 {
		return
	}

	for _, request := range requests {
		RegisterRequest(request)
	}
}

// UnaryServerInterceptor - returns a new unary server interceptor that validates incoming messages.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		t := reflect.TypeOf(req)

		if validators[t] {
			r, ok := req.(incomingRequest)
			if !ok {
				return nil, status.Errorf(codes.Internal, "something wrong in the gatekeeper")
			}

			if err := r.Validate(); err != nil {
				return nil, status.Errorf(codes.Internal, "%v", err.Error())
			}
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor - returns a new streaming server interceptor that validates incoming messages.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &recvWrapper{stream}

		return handler(srv, wrapper)
	}
}

func (s *recvWrapper) RecvMsg(msg any) error {
	if err := s.ServerStream.RecvMsg(msg); err != nil {
		return err
	}

	t := reflect.TypeOf(msg)

	if validators[t] {
		r, ok := msg.(incomingRequest)
		if !ok {
			return status.Errorf(codes.Internal, "something wrong in the gatekeeper")
		}

		if err := r.Validate(); err != nil {
			return status.Errorf(codes.Internal, "%v", err.Error())
		}
	}

	return nil
}
