package service

import (
	"context"

	"github.com/areknoster/public-distributed-commit-log/sentinel"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/ipfs/go-cid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	sentinelpb.UnimplementedSentinelServer
	validator sentinel.Validator
	pinner    sentinel.Pinner
	commiter  sentinel.Commiter
}

func New(validator sentinel.Validator, pinner sentinel.Pinner, commiter sentinel.Commiter) *Service {
	return &Service{validator: validator, pinner: pinner, commiter: commiter}
}

func (s *Service) Publish(ctx context.Context, req *sentinelpb.PublishRequest) (*sentinelpb.PublishResponse, error) {
	cid, err := cid.Decode(req.Cid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "can't parse CID")
	}

	if err := s.validator.Validate(ctx, cid); err != nil {
		return nil, validationErrorToProtoStatus(err)
	}
	if err := s.pinner.Pin(ctx, cid); err != nil {
		return nil, status.Error(codes.Internal, "can't pin message")
	}
	if err := s.commiter.Add(ctx, cid); err != nil {
		return nil, status.Error(codes.Internal, "can't add message to next commit")
	}
	return &sentinelpb.PublishResponse{}, nil
}

func validationErrorToProtoStatus(err error) error {
	// todo: map validation errors to correct statuses
	return status.Error(codes.InvalidArgument, err.Error())
}
