package daemonstorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	"github.com/areknoster/public-distributed-commit-log/pdcl"
	"github.com/areknoster/public-distributed-commit-log/storage"
)

type Config struct {
	IPFSDaemonPort string `envconfig:"IPFS_DAEMON_PORT" default:"5001"`
	IPFSDaemonHost string `envconfig:"IPFS_DAEMON_HOST" required:"true"`
}

func NewShell(config Config) *shell.Shell {
	return shell.NewShell(net.JoinHostPort(config.IPFSDaemonHost, config.IPFSDaemonPort))
}

// Storage is IPFS-based storage.Storage interface implementation
type Storage struct {
	shell       *shell.Shell
	marshalOpts proto.MarshalOptions
}

func NewStorage(sh *shell.Shell) *Storage {
	return &Storage{
		shell: sh,
		marshalOpts: proto.MarshalOptions{
			Deterministic: true,
		},
	}
}

func (s *Storage) Read(ctx context.Context, cid cid.Cid) (storage.ProtoUnmarshallable, error) {
	rc, err := s.shell.Cat(cid.String())
	if err != nil {
		return nil, fmt.Errorf("cat %s from IPFS: %w", cid.String(), err)
	}

	content, err := io.ReadAll(rc)
	if err != nil {
		if err := rc.Close(); err != nil {
			log.Ctx(ctx).Error().Err(err).Stringer("cid", cid).Msg("close message reader")
		}
		return nil, fmt.Errorf("read message content: %w", err)
	}
	return storage.ProtoDecode(content), nil
}

func (s *Storage) Write(ctx context.Context, message proto.Message) (cid.Cid, error) {
	encoded, err := s.marshalOpts.Marshal(message)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("marshall message: %w", err)
	}

	cidStr, err := s.shell.Add(bytes.NewReader(encoded), shell.CidVersion(1), shell.Pin(true))
	if err != nil {
		return cid.Cid{}, fmt.Errorf("add marshaled message to IPFS: %w", err)
	}

	pdclCid, err := pdcl.ParseCID(cidStr)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("parse cid from added message: %w", err)
	}
	return pdclCid, nil
}
