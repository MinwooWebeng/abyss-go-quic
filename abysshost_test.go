package abyssnet

import (
	"context"
	"net"
	"testing"
)

func TestNewAbyssHost(t *testing.T) {
	host_masterkey, _ := GenerateRSAKeypairPKCS8()
	host, err := NewAbyssHost(host_masterkey)
	if err != nil {
		t.Errorf(err.Error())
	}
	if host == nil {
		t.Errorf("returned nil host")
		return
	}
	host.Terminate()
}
func TestLocalConnect(t *testing.T) {
	hostA_masterkey, _ := GenerateRSAKeypairPKCS8()
	hostA, _ := NewAbyssHost(hostA_masterkey)
	defer hostA.Terminate()
	hostB_masterkey, _ := GenerateRSAKeypairPKCS8()
	hostB, _ := NewAbyssHost(hostB_masterkey)
	defer hostB.Terminate()

	if err := hostA.Connect(context.TODO(), &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: hostB.RawTransport.Conn.LocalAddr().(*net.UDPAddr).Port,
	}); err != nil {
		t.Errorf(err.Error())
	}
}
func TestRequestResource(t *testing.T) { //blocking call

}
func TestReturnResource(t *testing.T) {

}
func TestCloseConnection(t *testing.T) {

}
