package tswatcher

import (
	"context"
	"fmt"
	"strings"

	"tailscale.com/client/tailscale"
)

type TailscaleLocal struct {
	client tailscale.LocalClient
}

func NewLocalClient() *TailscaleLocal {
	return &TailscaleLocal{
		client: tailscale.LocalClient{},
	}
}

func (t *TailscaleLocal) GetDevices(ctx context.Context) ([]*TailnetDevice, error) {
	status, err := t.client.Status(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*TailnetDevice, 0, len(status.Peers()))
	for _, v := range status.Peer {
		result = append(result, FromPeerStatus(v))
	}
	return result, nil
}

func (t *TailscaleLocal) GetDevice(ctx context.Context, id string) (*TailnetDevice, error) {
	status, err := t.client.Status(ctx)
	if err != nil {
		return nil, err
	}

	if strings.EqualFold(string(status.Self.ID), id) {
		// it's me
		return FromPeerStatus(status.Self), nil
	}

	for _, p := range status.Peer {
		if strings.EqualFold(string(p.ID), id) {
			return FromPeerStatus(p), nil
		}
	}

	return nil, fmt.Errorf("%s not found in tailnet", id)
}
