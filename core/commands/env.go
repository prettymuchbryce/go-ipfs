package commands

import (
	"fmt"

	"github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"

	config "gx/ipfs/QmYyFh6g1C9uieTpH8CR8PpWBUQjvMDJTsRhJWx5qkXy39/go-ipfs-config"
)

// GetNode extracts the node from the environment.
func GetNode(env interface{}) (*core.IpfsNode, error) {
	ctx, ok := env.(*commands.Context)
	if !ok {
		return nil, fmt.Errorf("expected env to be of type %T, got %T", ctx, env)
	}

	return ctx.GetNode()
}

// GetConfig extracts the config from the environment.
func GetConfig(env interface{}) (*config.Config, error) {
	ctx, ok := env.(*commands.Context)
	if !ok {
		return nil, fmt.Errorf("expected env to be of type %T, got %T", ctx, env)
	}

	return ctx.GetConfig()
}
