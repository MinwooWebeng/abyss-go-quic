package abyssresource

import (
	"time"
)

type DatagramResource struct {
	ResourceBase
	TxBuffer chan []byte
	RxBuffer chan []byte
}

func CreateDatagramResource(max_tx_count int, max_rx_count int) DatagramResource {
	return DatagramResource{
		ResourceBase{last_modified: time.Time{}},
		make(chan []byte, max_tx_count),
		make(chan []byte, max_rx_count)}
}

func (r *DatagramResource) GetMIME() string {
	return "transport/datagram"
}
func (r *DatagramResource) GetMTU() int {
	return 1200 //TODO: check for quic spec and modify this to better reflect reality
}
func (r *DatagramResource) SendTo(payload []byte) {
	r.TxBuffer <- payload
}
func (r *DatagramResource) NonblockSendTo(payload []byte) bool {
	select {
	case r.TxBuffer <- payload: //buffer has space
		return true
	default: //buffer full
		return false
	}
}
func (r *DatagramResource) RecvFrom() (bool, []byte) {
	result, more := <-r.RxBuffer
	if more {
		return true, result
	} else {
		return false, nil
	}
}
func (r *DatagramResource) NonblockRecvFrom() (bool, []byte, bool) {
	select {
	case payload, more := <-r.RxBuffer:
		if more { //payload available
			return true, payload, more
		} else { //channel closed
			return false, nil, more
		}
	default: //buffer empty, channel not closed
		return false, nil, true
	}
}
func (r *DatagramResource) Close() { //gracefull termination, try to send leftover datagrams

}
func (r *DatagramResource) Abandon() {
	// r.rx_done <- false
}
