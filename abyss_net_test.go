package abyss_net

import (
	"testing"
)

func TestNewAbyssHost(t *testing.T) {
	host, err := NewAbyssHost()
	if err != nil {
		t.Errorf(err.Error())
	}
	if host == nil {
		t.Errorf("returned nil host")
	}
}
func TestInit(t *testing.T) {
	host, _ := NewAbyssHost()

	err := host.Init()
	if err != nil {
		t.Errorf(err.Error())
	}
}
func TestConnect(t *testing.T) { //blocking call

}
func TestRequestResource(t *testing.T) { //blocking call

}
func TestReturnResource(t *testing.T) {

}
func TestCloseConnection(t *testing.T) {

}
