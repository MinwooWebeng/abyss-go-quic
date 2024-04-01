package abyssnet

import (
	"crypto/rsa"

	"github.com/quic-go/quic-go"
)

type AbyssPeer struct {
	InboundConnection  *quic.Connection
	OutboundConnection *quic.Connection
	PublicKey          *rsa.PublicKey
	ID                 PeerID
}
