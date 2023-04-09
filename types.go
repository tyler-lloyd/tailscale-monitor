package tswatcher

import (
	"time"

	"tailscale.com/ipn/ipnstate"
)

type TailnetDevice struct {
	Addresses                 []string  `json:"addresses"`
	Authorized                bool      `json:"authorized"`
	BlocksIncomingConnections bool      `json:"blocksIncomingConnections"`
	ClientVersion             string    `json:"clientVersion"`
	Created                   time.Time `json:"created"`
	Expires                   time.Time `json:"expires"`
	Hostname                  string    `json:"hostname"`
	ID                        string    `json:"id"`
	IsExternal                bool      `json:"isExternal"`
	KeyExpiryDisabled         bool      `json:"keyExpiryDisabled"`
	LastSeen                  time.Time `json:"lastSeen"`
	MachineKey                string    `json:"machineKey"`
	Name                      string    `json:"name"`
	NodeID                    string    `json:"nodeId"`
	NodeKey                   string    `json:"nodeKey"`
	Os                        string    `json:"os"`
	TailnetLockError          string    `json:"tailnetLockError"`
	TailnetLockKey            string    `json:"tailnetLockKey"`
	UpdateAvailable           bool      `json:"updateAvailable"`
	User                      string    `json:"user"`
	Online                    bool      `json:"online"`
	DNSName                   string    `json:"dnsName"`
}

func FromPeerStatus(p *ipnstate.PeerStatus) *TailnetDevice {
	return &TailnetDevice{
		Addresses: p.Addrs,
		Created:   p.Created,
		Hostname:  p.HostName,
		DNSName:   p.DNSName,
		ID:        string(p.ID),
		Os:        p.OS,
		Name:      p.HostName,
		NodeID:    string(p.ID),
		NodeKey:   p.PublicKey.String(),
		Online:    p.Online,
		LastSeen:  p.LastSeen,
	}
}
