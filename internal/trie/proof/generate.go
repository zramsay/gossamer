// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package proof

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ChainSafe/gossamer/internal/trie/codec"
	"github.com/ChainSafe/gossamer/internal/trie/node"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/trie"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

// Database defines a key value Get method used
// for proof generation.
type Database interface {
	Get(key []byte) (value []byte, err error)
}

// Generate returns the encoded proof nodes for the trie
// corresponding to the root hash given, and for the (Little Endian)
// full key given. The database given is used to load the trie
// using the root hash given.
func Generate(rootHash common.Hash, fullKey []byte, database Database) (
	encodedProofNodes [][]byte, err error) {
	trie := trie.NewEmptyTrie()
	if err := trie.Load(database, rootHash); err != nil {
		return nil, fmt.Errorf("cannot load trie: %w", err)
	}

	rootNode := trie.RootNode()
	const isCurrentRoot = true
	fullKeyNibbles := codec.KeyLEToNibbles(fullKey)
	encodedProofNodes, err = walk(rootNode, isCurrentRoot, fullKeyNibbles)
	if err != nil {
		// Note we wrap the full key context here since find is recursive and
		// may not be aware of the initial full key.
		return nil, fmt.Errorf("cannot find node at key 0x%x in trie: %w", fullKey, err)
	}

	return encodedProofNodes, nil
}

// TODO use pointer to slice to avoid recursive appending
func walk(parent node.Node, isCurrentRoot bool, fullKey []byte) (
	encodedProofNodes [][]byte, err error) {
	if parent == nil {
		return nil, nil
	}

	encoding, _, err := parent.EncodeAndHash(isCurrentRoot)
	if err != nil {
		return nil, fmt.Errorf("cannot encode and hash node: %w", err)
	}

	encodedProofNodes = append(encodedProofNodes, encoding)

	switch parent.Type() {
	case node.BranchType, node.BranchWithValueType:
	default: // not a branch
		return encodedProofNodes, nil
	}

	b := parent.(*node.Branch)

	length := lenCommonPrefix(b.Key, fullKey)

	nodeFound := len(fullKey) == 0 || bytes.Equal(b.Key, fullKey)
	if nodeFound {
		return encodedProofNodes, nil
	}

	nodeIsDeeper := !bytes.Equal(b.Key[:length], fullKey) || len(fullKey) >= len(b.Key)
	if !nodeIsDeeper {
		return nil, ErrKeyNotFound
	}

	isCurrentRoot = false
	deeperEncodedProofNodes, err := walk(
		b.Children[fullKey[length]], isCurrentRoot,
		fullKey[length+1:])
	if err != nil {
		return nil, err // note: do not wrap since this is recursive
	}

	encodedProofNodes = append(encodedProofNodes, deeperEncodedProofNodes...)
	return encodedProofNodes, nil
}

// lenCommonPrefix returns the length of the
// common prefix between two byte slices.
func lenCommonPrefix(a, b []byte) (length int) {
	min := len(a)
	if len(b) < min {
		min = len(b)
	}

	for length = 0; length < min; length++ {
		if a[length] != b[length] {
			break
		}
	}

	return length
}
