// Package importer implements utilities used to create IPFS DAGs from files
// and readers.
package importer

import (
	chunker "gx/ipfs/QmVDjhUMtkRskBFAVNwyXuLSKbeAya7JKPnzAxMKDaK4x4/go-ipfs-chunker"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"

	bal "gx/ipfs/QmQqL1Y9Wzpy8dRJmNFAwBQYaQxqaxtD93WHeYuyME6Pz1/go-unixfs/importer/balanced"
	h "gx/ipfs/QmQqL1Y9Wzpy8dRJmNFAwBQYaQxqaxtD93WHeYuyME6Pz1/go-unixfs/importer/helpers"
	trickle "gx/ipfs/QmQqL1Y9Wzpy8dRJmNFAwBQYaQxqaxtD93WHeYuyME6Pz1/go-unixfs/importer/trickle"
)

// BuildDagFromReader creates a DAG given a DAGService and a Splitter
// implementation (Splitters are io.Readers), using a Balanced layout.
func BuildDagFromReader(ds ipld.DAGService, spl chunker.Splitter) (ipld.Node, error) {
	dbp := h.DagBuilderParams{
		Dagserv:  ds,
		Maxlinks: h.DefaultLinksPerBlock,
	}

	return bal.Layout(dbp.New(spl))
}

// BuildTrickleDagFromReader creates a DAG given a DAGService and a Splitter
// implementation (Splitters are io.Readers), using a Trickle Layout.
func BuildTrickleDagFromReader(ds ipld.DAGService, spl chunker.Splitter) (ipld.Node, error) {
	dbp := h.DagBuilderParams{
		Dagserv:  ds,
		Maxlinks: h.DefaultLinksPerBlock,
	}

	return trickle.Layout(dbp.New(spl))
}
