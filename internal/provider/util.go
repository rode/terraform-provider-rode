package provider

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func formatProtoTimestamp(timestamp *timestamppb.Timestamp) string {
	return timestamp.AsTime().Format(time.RFC3339Nano)
}
