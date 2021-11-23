package service

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/areknoster/public-distributed-commit-log/head"
	"github.com/areknoster/public-distributed-commit-log/sentinel"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
)

type Service struct {
	sentinelpb.UnimplementedSentinelServer
	validator  sentinel.Validator
	pinner     sentinel.Pinner
	commiter   sentinel.Commiter
	headReader head.Reader
}

func New(validator sentinel.Validator, pinner sentinel.Pinner, commiter sentinel.Commiter, headReader head.Reader) *Service {
	return &Service{
		validator:  validator,
		pinner:     pinner,
		commiter:   commiter,
		headReader: headReader,
	}
}

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

func (s *Service) GetHeadCID(ctx context.Context, req *sentinelpb.GetHeadCIDRequest) (*sentinelpb.GetHeadCIDResponse, error) {
	headCID, err := s.headReader.ReadHead(ctx)
	if err != nil {
		return nil, fmt.Errorf("get head CID: %w", err)
	}
	return &sentinelpb.GetHeadCIDResponse{Cid: headCID.String()}, nil
}
