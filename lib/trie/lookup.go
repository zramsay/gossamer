// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package trie

import (
	"bytes"

	"github.com/ChainSafe/gossamer/internal/trie/node"
	"github.com/ChainSafe/gossamer/internal/trie/record"
)

var _ recorder = (*record.Recorder)(nil)

type recorder interface {
	Record(hash, rawData []byte)
}

// findAndRecord search for a desired key recording all the nodes in the path including the desired node
func findAndRecord(t *Trie, key []byte, recorder recorder) error {
	return find(t.root, key, recorder, true)
}

func find(parent *Node, key []byte, recorder recorder, isCurrentRoot bool) error {
	enc, hash, err := parent.EncodeAndHash(isCurrentRoot)
	if err != nil {
		return err
	}

	recorder.Record(hash, enc)

	if parent.Type() != node.Branch {
		return nil
	}

	branch := parent
	length := lenCommonPrefix(branch.Key, key)

	// found the value at this node
	if bytes.Equal(branch.Key, key) || len(key) == 0 {
		return nil
	}

	// did not find value
	if bytes.Equal(branch.Key[:length], key) && len(key) < len(branch.Key) {
		return nil
	}

	return find(branch.Children[key[length]], key[length+1:], recorder, false)
}
