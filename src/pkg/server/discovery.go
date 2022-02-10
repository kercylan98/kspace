package server

import "context"

// Discovery 服务发现
type Discovery interface {

	// Release 服务发布
	Release(ctx context.Context, state *State) (<-chan error, error)
}
