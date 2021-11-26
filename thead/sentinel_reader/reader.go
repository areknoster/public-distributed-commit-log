package sentinel_reader

import (
	"context"
	"fmt"
	"github.com/areknoster/public-distributed-commit-log/pdcl"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/ipfs/go-cid"
)

type SentinelHeadReader struct{
	client sentinelpb.SentinelClient
}

func NewSentinelHeadReader(client sentinelpb.SentinelClient) *SentinelHeadReader {
	return &SentinelHeadReader{client: client}
}

func (s *SentinelHeadReader) ReadHead(ctx context.Context) (cid.Cid, error) {
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

