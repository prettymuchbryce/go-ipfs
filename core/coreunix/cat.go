package coreunix

import (
	"context"

	core "github.com/ipfs/go-ipfs/core"
	uio "gx/ipfs/QmQqL1Y9Wzpy8dRJmNFAwBQYaQxqaxtD93WHeYuyME6Pz1/go-unixfs/io"
	path "gx/ipfs/QmS3Pbj4ZV1VtqPQnWkmWqBn1eXt5FbxPbtUdDZMRAZZXr/go-path"
	resolver "gx/ipfs/QmS3Pbj4ZV1VtqPQnWkmWqBn1eXt5FbxPbtUdDZMRAZZXr/go-path/resolver"
)

func Cat(ctx context.Context, n *core.IpfsNode, pstr string) (uio.DagReader, error) {
	r := &resolver.Resolver{
		DAG:         n.DAG,
		ResolveOnce: uio.ResolveUnixfsOnce,
	}

	dagNode, err := core.Resolve(ctx, n.Namesys, r, path.Path(pstr))
	if err != nil {
		return nil, err
	}

	return uio.NewDagReader(ctx, dagNode, n.DAG)
}
