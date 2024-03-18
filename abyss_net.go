package abyss_net

import (
	"errors"

	"github.com/quic-go/quic-go"
)

type abyss_host struct {
	_raw_transport quic.Transport
}

func NewAbyssHost() (*abyss_host, error) {
	return nil, errors.ErrUnsupported
}
func (h *abyss_host) Init() error {
	return errors.ErrUnsupported
}
func (h *abyss_host) Connect() error { //blocking call, returns when it accepts/connects target peer
	return errors.ErrUnsupported
}
func (h *abyss_host) RequestResource() error { //blocking call, returns when response arrives
	return errors.ErrUnsupported
}
func (h *abyss_host) ReturnResource() error { //after this call, access to resource is unsafe.
	return errors.ErrUnsupported
}
func (h *abyss_host) JoinWorld() error {
	return errors.ErrUnsupported
}
func (h *abyss_host) CloseConnection() error { //this fails if not returned resources exist.
	return errors.ErrUnsupported
}
