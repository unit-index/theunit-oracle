package transport

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/libp2p/go-libp2p-core/crypto"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/messages"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/p2p"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/p2p/crypto/ethkey"
)

var p2pTransportFactory = func(ctx context.Context, cfg p2p.Config) (transport.Transport, error) {
	return p2p.New(ctx, cfg)
}

type Transport struct {
	P2P P2P `json:"p2p"`
}

type P2P struct {
	PrivKeySeed      string   `json:"privKeySeed"`
	ListenAddrs      []string `json:"listenAddrs"`
	BootstrapAddrs   []string `json:"bootstrapAddrs"`
	DirectPeersAddrs []string `json:"directPeersAddrs"`
	BlockedAddrs     []string `json:"blockedAddrs"`
	DisableDiscovery bool     `json:"disableDiscovery"`
}

type Dependencies struct {
	Context context.Context
	Signer  ethereum.Signer
	Feeds   []ethereum.Address
	Logger  log.Logger
}

type BootstrapDependencies struct {
	Context context.Context
	Logger  log.Logger
}

func (c *Transport) Configure(d Dependencies) (transport.Transport, error) {
	peerPrivKey, err := c.generatePrivKey()
	if err != nil {
		return nil, err
	}
	cfg := p2p.Config{
		Mode:             p2p.ClientMode,
		PeerPrivKey:      peerPrivKey,
		Topics:           map[string]transport.Message{messages.PriceMessageName: (*messages.Price)(nil)},
		MessagePrivKey:   ethkey.NewPrivKey(d.Signer),
		ListenAddrs:      c.P2P.ListenAddrs,
		BootstrapAddrs:   c.P2P.BootstrapAddrs,
		DirectPeersAddrs: c.P2P.DirectPeersAddrs,
		BlockedAddrs:     c.P2P.BlockedAddrs,
		FeedersAddrs:     d.Feeds,
		Discovery:        !c.P2P.DisableDiscovery,
		Signer:           d.Signer,
		Logger:           d.Logger,
		AppName:          "spire",
		AppVersion:       "1",
	}
	p, err := p2pTransportFactory(d.Context, cfg)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (c *Transport) ConfigureP2PBoostrap(d BootstrapDependencies) (transport.Transport, error) {
	peerPrivKey, err := c.generatePrivKey()
	if err != nil {
		return nil, err
	}
	cfg := p2p.Config{
		Mode:             p2p.BootstrapMode,
		PeerPrivKey:      peerPrivKey,
		ListenAddrs:      c.P2P.ListenAddrs,
		BootstrapAddrs:   c.P2P.BootstrapAddrs,
		DirectPeersAddrs: c.P2P.DirectPeersAddrs,
		BlockedAddrs:     c.P2P.BlockedAddrs,
		Logger:           d.Logger,
		AppName:          "bootstrap",
		AppVersion:       "1",
	}
	p, err := p2pTransportFactory(d.Context, cfg)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (c *Transport) generatePrivKey() (crypto.PrivKey, error) {
	seedReader := rand.Reader
	if len(c.P2P.PrivKeySeed) != 0 {
		seed, err := hex.DecodeString(c.P2P.PrivKeySeed)
		if err != nil {
			return nil, fmt.Errorf("invalid privKeySeed value, failed to decode hex data: %w", err)
		}
		if len(seed) != ed25519.SeedSize {
			return nil, fmt.Errorf("invalid privKeySeed value, 32 bytes expected")
		}
		seedReader = bytes.NewReader(seed)
	}
	privKey, _, err := crypto.GenerateEd25519Key(seedReader)
	if err != nil {
		return nil, fmt.Errorf("invalid privKeySeed value, failed to generate key: %w", err)
	}
	return privKey, nil
}
