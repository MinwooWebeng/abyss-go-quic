package abyssnet

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"net"

	"github.com/quic-go/quic-go"
)

type AbyssHost struct {
	LastError    <-chan error
	MasterKey    *rsa.PrivateKey
	UptimeKey    *rsa.PrivateKey
	TlsConf      tls.Config
	QuicConf     quic.Config
	RawTransport quic.Transport
	Terminate    func()

	PeerPool
}

func GenerateRSAKeypairPKCS8() ([]byte, error) {
	new_rsa_key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	rsa_priv_pkcs8, err := x509.MarshalPKCS8PrivateKey(new_rsa_key)
	if err != nil {
		return nil, err
	}

	return rsa_priv_pkcs8, nil
}

func NewAbyssHost(host_master_key_pkcs8 []byte) (*AbyssHost, error) {
	new_host := new(AbyssHost)
	hosterr := make(chan error, 10)
	new_host.LastError = hosterr

	_pmk, err := x509.ParsePKCS8PrivateKey(host_master_key_pkcs8)
	if err != nil {
		return nil, err
	}
	masterKey, okay := _pmk.(*rsa.PrivateKey)
	if !okay {
		return nil, err
	}
	new_host.MasterKey = masterKey

	//generate random key pair
	uptime_rsa_key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	new_host.UptimeKey = uptime_rsa_key
	uptime_private_key_PKCS8, err := x509.MarshalPKCS8PrivateKey(uptime_rsa_key)
	if err != nil {
		return nil, err
	}
	uptime_private_key_pem := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: uptime_private_key_PKCS8})

	uptime_cert_draft := x509.Certificate{SerialNumber: big.NewInt(42)}
	uptime_cert_x509, err := x509.CreateCertificate(rand.Reader, &uptime_cert_draft, &uptime_cert_draft, &uptime_rsa_key.PublicKey, uptime_rsa_key)
	if err != nil {
		return nil, err
	}
	uptime_cert_pem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: uptime_cert_x509})
	uptime_tls_cert, err := tls.X509KeyPair(uptime_cert_pem, uptime_private_key_pem)
	if err != nil {
		return nil, err
	}

	new_host.TlsConf = tls.Config{
		Certificates:       []tls.Certificate{uptime_tls_cert},
		InsecureSkipVerify: true,
	}
	new_host.QuicConf = quic.Config{
		EnableDatagrams: true,
	}

	udp_conn, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		return nil, err
	}
	new_host.RawTransport = quic.Transport{
		Conn: udp_conn,
	}

	listener, err := new_host.RawTransport.Listen(&new_host.TlsConf, &new_host.QuicConf)
	if err != nil {
		return nil, err
	}

	service_ctx, service_ctx_terminator := context.WithCancel(context.Background())
	terminated := make(chan bool)

	new_host.Terminate = func() {
		service_ctx_terminator()
		<-terminated
	}

	go AbyssHostService(
		service_ctx,
		new_host,
		func(e error) {
			select {
			case hosterr <- e:
			default:
			}
		},
		terminated,
		listener,
	)

	return new_host, nil
}

func AbyssHostService(service_ctx context.Context, host *AbyssHost, LogErr func(error), terminated chan<- bool, listener *quic.Listener) { //TODO: implement host background functionality; extend AcceptLoop.
AcceptLoop:
	for {
		_, err := listener.Accept(service_ctx)
		switch err {
		case nil:
		case context.Canceled:
			break AcceptLoop
		default:
			LogErr(errors.Join(errors.New("AbyssHost Failed: Unhandled Exception from quic.Listener.Accept()"), err))
			break AcceptLoop
		}
	}
	terminated <- true
}

//func AbyssNeighborDiscoveryService(service_ctx context.Context, host *AbyssHost, LogErr func(error), terminated chan<- bool, event_chan <-chan )

func (h *AbyssHost) Connect(ctx context.Context, addr net.Addr) error { //blocking call, returns when it accepts/connects target peer
	_, err := h.RawTransport.Dial(ctx, addr, &h.TlsConf, &h.QuicConf)
	if err != nil {
		return err
	}

	return nil
}
func (h *AbyssHost) RequestResource() error { //blocking call, returns when response arrives
	return errors.ErrUnsupported
}
func (h *AbyssHost) ReturnResource() error { //after this call, access to resource is unsafe.
	return errors.ErrUnsupported
}
func (h *AbyssHost) JoinWorld() error {
	return errors.ErrUnsupported
}
func (h *AbyssHost) CloseConnection() error { //this fails if not returned resources exist.
	return errors.ErrUnsupported
}

// func GenerateRandomString(length int) (string, error) {
// 	randomBytes := make([]byte, length)
// 	if _, err := rand.Read(randomBytes); err != nil {
// 		return "", err
// 	}
// 	return base58.Encode(randomBytes)[:length], nil
// }

// x509.Certificate{
//     SerialNumber: big.NewInt(time.Now().Unix()),
//     Subject: pkix.Name{
//         Organization: []string{"Your Organization"},
//     },
//     NotBefore:             time.Now(),
//     NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
//     KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
//     ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
//     BasicConstraintsValid: true,
//     IsCA:                  true,
// }
