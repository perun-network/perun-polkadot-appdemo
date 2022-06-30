package client

import "perun.network/go-perun/wire"

func (c *Client) RegisterPeer(addr wire.Address, host string) {
	c.dialer.Register(addr, host)
}
