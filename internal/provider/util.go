package provider

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func formatProtoTimestamp(timestamp *timestamppb.Timestamp) string {
	return timestamp.AsTime().Format(time.RFC3339Nano)
}

func asProtoTimestamp(val interface{}) (*timestamppb.Timestamp, error) {
	t, err := time.Parse(val.(string), time.RFC3339Nano)

	if err != nil {
		return nil, err
	}

	return timestamppb.New(t), nil
}

func isDeletionError(err error) bool {
	if err == nil {
		return false
	}

	// ignore NOT_FOUND for the purpose of delete
	return status.Code(err) != codes.NotFound
}