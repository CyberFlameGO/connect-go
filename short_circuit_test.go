package rerpc_test

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/rerpc/rerpc"
)

type shortCircuit struct {
	err error
}

var _ rerpc.Interceptor = (*shortCircuit)(nil)

func (sc *shortCircuit) Wrap(next rerpc.Func) rerpc.Func {
	return rerpc.Func(func(_ context.Context, _ proto.Message) (proto.Message, error) {
		return nil, sc.err
	})
}

func (sc *shortCircuit) WrapHandlerStream(next rerpc.HandlerStreamFunc) rerpc.HandlerStreamFunc {
	return rerpc.HandlerStreamFunc(func(ctx context.Context, _ rerpc.Stream) {
		next(ctx, &errStream{sc.err})
	})
}

func (sc *shortCircuit) WrapCallStream(next rerpc.CallStreamFunc) rerpc.CallStreamFunc {
	return rerpc.CallStreamFunc(func(_ context.Context) rerpc.Stream {
		return &errStream{sc.err}
	})
}

// ShortCircuit builds an interceptor that doesn't call the wrapped RPC at all.
// Instead, it returns the supplied Error immediately. ShortCircuit works for
// unary and streaming RPCs.
//
// This is primarily useful when testing error handling. It's also used
// throughout reRPC's examples to avoid making network requests.
func ShortCircuit(err error) rerpc.Interceptor {
	return &shortCircuit{err}
}

type errStream struct {
	err error
}

var _ rerpc.Stream = (*errStream)(nil)

func (s *errStream) Receive(_ proto.Message) error { return s.err }
func (s *errStream) CloseReceive() error           { return s.err }
func (s *errStream) Send(_ proto.Message) error    { return s.err }
func (s *errStream) CloseSend(_ error) error       { return s.err }