package service

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/sentinel"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/thead"
)

// Service implements sentinel GRPC server
type Service struct {
	sentinelpb.UnimplementedSentinelServer
	validator   sentinel.Validator
	pinner      sentinel.Pinner
	commiter    sentinel.Commiter
	headReader  thead.Reader
	ipnsManager ipns.Manager
}

// New initializes sentinel service
func New(validator sentinel.Validator, pinner sentinel.Pinner, commiter sentinel.Commiter, headReader thead.Reader,
	ipnsManager ipns.Manager) *Service {
	return &Service{
		validator:   validator,
		pinner:      pinner,
		commiter:    commiter,
		headReader:  headReader,
		ipnsManager: ipnsManager,
	}
}

// Publish is GRPC method that producers can use to add message to topic
func (s *Service) Publish(ctx context.Context, req *sentinelpb.PublishRequest) (*sentinelpb.PublishResponse, error) {
	logger := log.With().Str("cid", req.Cid).Logger()
	cid, err := cid.Decode(req.Cid)
	if err != nil {
		logger.Info().Err(err).Msg("error decoding CID")
		return nil, status.Error(codes.InvalidArgument, "can't parse CID")
	}

	if err := s.validator.Validate(ctx, cid); err != nil {
		logger.Info().Err(err).Msg("error validating message")
		return nil, validationErrorToProtoStatus(err)
	}
	if err := s.pinner.Pin(ctx, cid); err != nil {
		logger.Error().Err(err).Msg("error pinning CID")
		return nil, status.Error(codes.Internal, "can't pin message")
	}
	if err := s.commiter.Add(ctx, cid); err != nil {
		logger.Error().Err(err).Msg("error pinning CID")
		return nil, status.Error(codes.Internal, "can't add message to next commit")
	}
	return &sentinelpb.PublishResponse{}, nil
}

func validationErrorToProtoStatus(err error) error {
	// todo: map validation errors to correct statuses
	return status.Error(codes.InvalidArgument, err.Error())
}

// GetHeadCID can be used by consumers to fetch head from sentinel
func (s *Service) GetHeadCID(ctx context.Context, req *sentinelpb.GetHeadCIDRequest) (*sentinelpb.GetHeadCIDResponse, error) {
	headCID, err := s.headReader.ReadHead(ctx)
	if err != nil {
		return nil, fmt.Errorf("get head CID: %w", err)
	}
	return &sentinelpb.GetHeadCIDResponse{Cid: headCID.String()}, nil
}

// GetHeadCID can be used by consumers to fetch head from sentinel
func (s *Service) GetHeadIPNS(ctx context.Context, req *sentinelpb.GetHeadIPNSRequest) (*sentinelpb.GetHeadIPNSResponse, error) {
	return &sentinelpb.GetHeadIPNSResponse{IpnsAddr: s.ipnsManager.GetIPNSAddr()}, nil
}
