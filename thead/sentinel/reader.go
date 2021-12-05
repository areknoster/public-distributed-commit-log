package sentinel

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"

	"github.com/areknoster/public-distributed-commit-log/pdcl"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
)

// Reader implements head.Reader based on sentinel client. It should be initialized with New function
type Reader struct {
	client sentinelpb.SentinelClient
}

// New initializes sentinel head.Reader
func New(client sentinelpb.SentinelClient) *Reader {
	return &Reader{client: client}
}

// ReadHead fetches head cid from sentinel
func (s *Reader) ReadHead(ctx context.Context) (cid.Cid, error) {
	resp, err := s.client.GetHeadCID(ctx, &sentinelpb.GetHeadCIDRequest{})
	if err != nil {
		return cid.Cid{}, fmt.Errorf("get head cid from sentinel: %w", err)
	}

	decodedCID, err := pdcl.ParseCID(resp.Cid)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("decode CID: %w", err)
	}
	return decodedCID, nil
}
