// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package node

import (
	"bytes"
	"fmt"

	"github.com/ChainSafe/gossamer/internal/trie/pools"
	"github.com/ChainSafe/gossamer/lib/common"
)

// MerkleValue produces the Merkle value from the encoding of a node.
// For root nodes, the Merkle value is always the Blak2b hash of the encoding.
// For other nodes, the Merkle value is either:
// - the encoding if it is less than 32 bytes
// - the Blake2b hash of the encoding
func MerkleValue(encoding []byte, isRoot bool) (merkleValue []byte, err error) {
	if !isRoot && len(encoding) < 32 {
		merkleValue = make([]byte, len(encoding))
		copy(merkleValue, encoding)
		return merkleValue, nil
	}

	hashDigest, err := common.Blake2bHash(encoding)
	if err != nil {
		return nil, err
	}

	merkleValue = hashDigest[:]
	return merkleValue, nil
}

// MerkleValue returns the encoding of the node and
// the blake2b hash digest of the encoding of the node.
// If the encoding is less than 32 bytes, the hash returned
// is the encoding and not the hash of the encoding.
func (n *Node) MerkleValue(isRoot bool) (merkleValue []byte, err error) {
	if !n.Dirty && n.HashDigest != nil {
		return n.HashDigest, nil
	}

	_, merkleValue, err = n.EncodeAndHash(isRoot)
	return merkleValue, err
}

// EncodeAndHash returns the encoding of the node and the
// Merkle value of the node. See the `MerkleValue`
// method for more details on the value of the Merkle value.
func (n *Node) EncodeAndHash(isRoot bool) (encoding, merkleValue []byte, err error) {
	if !n.Dirty && n.Encoding != nil && n.HashDigest != nil {
		return n.Encoding, n.HashDigest, nil
	}

	if n.Dirty || n.Encoding == nil {
		buffer := pools.EncodingBuffers.Get().(*bytes.Buffer)
		buffer.Reset()
		defer pools.EncodingBuffers.Put(buffer)

		err = n.Encode(buffer)
		if err != nil {
			return nil, nil, fmt.Errorf("encode node: %w", err)
		}

		bufferBytes := buffer.Bytes()

		// TODO remove this copying since it defeats the purpose of `buffer`
		// and the sync.Pool.
		n.Encoding = make([]byte, len(bufferBytes))
		copy(n.Encoding, bufferBytes)
	}
	encoding = n.Encoding // no need to copy

	merkleValue, err = MerkleValue(encoding, isRoot)
	if err != nil {
		return nil, nil, fmt.Errorf("merkle value: %w", err)
	}
	n.HashDigest = merkleValue // no need to copy

	return encoding, merkleValue, nil
}
