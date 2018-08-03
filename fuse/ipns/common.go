package ipns

import (
	"context"

	"github.com/ipfs/go-ipfs/core"
	nsys "github.com/ipfs/go-ipfs/namesys"
	ft "gx/ipfs/QmQqL1Y9Wzpy8dRJmNFAwBQYaQxqaxtD93WHeYuyME6Pz1/go-unixfs"
	path "gx/ipfs/QmS3Pbj4ZV1VtqPQnWkmWqBn1eXt5FbxPbtUdDZMRAZZXr/go-path"
	ci "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

// InitializeKeyspace sets the ipns record for the given key to
// point to an empty directory.
func InitializeKeyspace(n *core.IpfsNode, key ci.PrivKey) error {
	ctx, cancel := context.WithCancel(n.Context())
	defer cancel()

	emptyDir := ft.EmptyDirNode()

	err := n.Pinning.Pin(ctx, emptyDir, false)
	if err != nil {
		return err
	}

	err = n.Pinning.Flush()
	if err != nil {
		return err
	}

	pub := nsys.NewIpnsPublisher(n.Routing, n.Repo.Datastore())

	return pub.Publish(ctx, key, path.FromCid(emptyDir.Cid()))
}
