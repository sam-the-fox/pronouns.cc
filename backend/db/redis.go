package db

import (
	"context"
	"net"

	"emperror.dev/errors"
	"github.com/mediocregopher/radix/v4"
)

var _ radix.Client = (*dummyRedis)(nil)

type dummyRedis struct{}

func (*dummyRedis) Addr() net.Addr { return &net.IPAddr{} }
func (*dummyRedis) Close() error   { return nil }

func (*dummyRedis) Do(context.Context, radix.Action) error {
	return errors.Sentinel("this is a dummy client")
}
