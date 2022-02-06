// Package ipns provides implementation for pdcl head management based on ipns protocol.
package ipns

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multibase"

	"github.com/areknoster/public-distributed-commit-log/thead"
)

var _ thead.Reader = (*BasicHeadReader)(nil)

// BasicHeadReader get's topic head from given IPNS address
type BasicHeadReader struct {
	sh       *shell.Shell
	resolver resolver
	ipnsAddr string
}

func NewBasicHeadReader(sh *shell.Shell, ipnsAddr string) *BasicHeadReader {
	bhr := &BasicHeadReader{
		sh:       sh,
		resolver: newShellResolver(sh),
		ipnsAddr: ipnsAddr,
	}
	bhr.ipnsAddr = ipnsAddr
	return bhr
}

func (hr *BasicHeadReader) ReadHead(ctx context.Context) (cid.Cid, error) {
	id, err := hr.resolver.resolveIPNS(hr.ipnsAddr)
	if err != nil {
		return cid.Cid{}, err
	}
	return id, nil
}

// BasicHeadSetter sets default IPFS daemon's IPNS key address to point to given CID
type BasicHeadSetter struct {
	sh *shell.Shell
}

func NewBasicHeadSetter(sh *shell.Shell) *BasicHeadSetter {
	return &BasicHeadSetter{sh: sh}
}

func (bhs *BasicHeadSetter) SetHead(ctx context.Context, cid cid.Cid) error {
	ipfsAddr := path.Join("/ipfs/", cid.String())
	_, err := bhs.sh.PublishWithDetails(ipfsAddr, "", 24*time.Hour, 10*time.Minute, false)
	if err != nil {
		return fmt.Errorf("publishing ipns update to ipfs daemon: %v", err)
	}
	return nil
}

// BasicHeadManager can be used by Sentinel to manage topic head
// Since consumers don't set topic's head, it should not be used by them.
// Use HeadReader implementation to get topic head for consumer
// and some other (e.g. memory or disk) implementation to store internal consumer offset
type BasicHeadManager struct {
	*BasicHeadReader
	*BasicHeadSetter
}

// NewBasicHeadManager initializes BasicHeadManager with default daemon's key used as PK for topic's head
func NewBasicHeadManager(sh *shell.Shell) (BasicHeadManager, error) {
	ipnsAddr, err := GetDaemonIPNSAddress(sh)
	if err != nil {
		return BasicHeadManager{}, fmt.Errorf("get daemon IPNS address: %w", err)
	}
	return BasicHeadManager{
		BasicHeadReader: NewBasicHeadReader(sh, ipnsAddr),
		BasicHeadSetter: NewBasicHeadSetter(sh),
	}, nil
}

// GetDaemonIPNSAddress gets IPNS address attached to daemon's default key.
//	In most scenarios it's to be used when initializing
//  sentinel with some existing daemon to use its default
//  ipns address as IPNS head.
func GetDaemonIPNSAddress(sh *shell.Shell) (string, error) {
	// this implementation is extremely non-obvious
	// because IPFS doesn't normally allow for finding IPNS address
	// of given key unless some file is added to it.
	resp, err := sh.ID()
	if err != nil {
		return "", fmt.Errorf("get IPFS ID: %w", err)
	}
	pid, err := peer.Decode(resp.ID)
	if err != nil {
		return "", fmt.Errorf("decode peer ID: %w", err)
	}
	ipnsAddr, err := peer.ToCid(pid).StringOfBase(multibase.Base36)
	if err != nil {
		return "", fmt.Errorf("encode ipns address: %w", err)
	}
	return ipnsAddr, nil
}
