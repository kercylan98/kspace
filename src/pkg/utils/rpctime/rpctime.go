package rpctime

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ToTimestamp(time time.Time) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{
		Seconds: int64(time.Second()),
		Nanos:   int32(time.Nanosecond()),
	}
}
