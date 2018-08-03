/*
Package corerepo provides pinning and garbage collection for local
IPFS block services.

IPFS nodes will keep local copies of any object that have either been
added or requested locally.  Not all of these objects are worth
preserving forever though, so the node administrator can pin objects
they want to keep and unpin objects that they don't care about.

Garbage collection sweeps iterate through the local block store
removing objects that aren't pinned, which frees storage space for new
objects.
*/
package corerepo

import (
	"context"
	"fmt"

	"github.com/ipfs/go-ipfs/core"
	uio "gx/ipfs/QmQqL1Y9Wzpy8dRJmNFAwBQYaQxqaxtD93WHeYuyME6Pz1/go-unixfs/io"
	path "gx/ipfs/QmS3Pbj4ZV1VtqPQnWkmWqBn1eXt5FbxPbtUdDZMRAZZXr/go-path"
	resolver "gx/ipfs/QmS3Pbj4ZV1VtqPQnWkmWqBn1eXt5FbxPbtUdDZMRAZZXr/go-path/resolver"

	cid "gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
)

func Pin(n *core.IpfsNode, ctx context.Context, paths []string, recursive bool) ([]*cid.Cid, error) {
	out := make([]*cid.Cid, len(paths))

	r := &resolver.Resolver{
		DAG:         n.DAG,
		ResolveOnce: uio.ResolveUnixfsOnce,
	}

	for i, fpath := range paths {
		p, err := path.ParsePath(fpath)
		if err != nil {
			return nil, err
		}

		dagnode, err := core.Resolve(ctx, n.Namesys, r, p)
		if err != nil {
			return nil, fmt.Errorf("pin: %s", err)
		}
		err = n.Pinning.Pin(ctx, dagnode, recursive)
		if err != nil {
			return nil, fmt.Errorf("pin: %s", err)
		}
		out[i] = dagnode.Cid()
	}

	err := n.Pinning.Flush()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func Unpin(n *core.IpfsNode, ctx context.Context, paths []string, recursive bool) ([]*cid.Cid, error) {
	unpinned := make([]*cid.Cid, len(paths))

	r := &resolver.Resolver{
		DAG:         n.DAG,
		ResolveOnce: uio.ResolveUnixfsOnce,
	}

	for i, p := range paths {
		p, err := path.ParsePath(p)
		if err != nil {
			return nil, err
		}

		k, err := core.ResolveToCid(ctx, n.Namesys, r, p)
		if err != nil {
			return nil, err
		}

		err = n.Pinning.Unpin(ctx, k, recursive)
		if err != nil {
			return nil, err
		}
		unpinned[i] = k
	}

	err := n.Pinning.Flush()
	if err != nil {
		return nil, err
	}
	return unpinned, nil
}
