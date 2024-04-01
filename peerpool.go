package abyssnet

type PeerPool struct {
	_Inner map[PeerID]AbyssPeer
}

func (p *PeerPool) AddAcceptedPeer() (bool, *AbyssPeer) {
	return false, nil
}

func (p *PeerPool) AddConnectedPeer() (bool, *AbyssPeer) {
	return false, nil
}
