package proof

import (
	"testing"

	"github.com/ChainSafe/gossamer/internal/trie/node"
	"github.com/stretchr/testify/assert"
)

func Test_walk(t *testing.T) {
	t.Parallel()

	// bigBranch is a branch resulting in an encoding larger than 32 bytes.
	// bigBranch := &node.Branch{
	// 	Key:   []byte{1, 2, 3, 4, 5},
	// 	Value: []byte{3},
	// 	Children: [16]node.Node{
	// 		&node.Leaf{
	// 			Key:   []byte{4},
	// 			Value: []byte{5},
	// 		},
	// 		&node.Leaf{
	// 			Key:   []byte{4},
	// 			Value: []byte{5},
	// 		},
	// 	},
	// }

	testCases := map[string]struct {
		parent            node.Node
		isCurrentRoot     bool
		fullKey           []byte // nibbles
		encodedProofNodes [][]byte
		errWrapped        error
		errMessage        string
	}{
		"nil parent": {},
		"parent encode and hash error": {
			parent: &node.Leaf{
				Key: make([]byte, int(^uint16(0))+63),
			},
			errWrapped: node.ErrPartialKeyTooBig,
			errMessage: "cannot encode and hash node: cannot encode header: " +
				"partial key length cannot be larger than or equal to 2^16: 65535",
		},
		"parent leaf": {
			parent: &node.Leaf{
				Key: []byte{1, 2},
			},
			encodedProofNodes: [][]byte{{0b0100_0000 | 2, 0x12, 0x00}},
		},
		"branch and empty search key": {
			parent: &node.Branch{
				Key:   []byte{1, 2},
				Value: []byte{3},
				Children: [16]node.Node{
					&node.Leaf{
						Key:   []byte{4},
						Value: []byte{5},
					},
				},
			},
			encodedProofNodes: [][]byte{
				{
					0b1100_0000 | 2,          // Node variant and partial key length
					0x12,                     // partial key
					0b0000_0001, 0b0000_0000, // children bitmap
					0x4, 0x3, // scale encoded value
					0x10,
					0b0100_0000 | 1, // Node variant and partial key length
					0x4, 0x4,        // partial key
					0x5, // leaf value
				}},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			encodedProofNodes, err := walk(testCase.parent, testCase.isCurrentRoot,
				testCase.fullKey)

			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
			assert.Equal(t, testCase.encodedProofNodes, encodedProofNodes)
		})
	}
}

func Test_lenCommonPrefix(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		a      []byte
		b      []byte
		length int
	}{
		"nil slices": {},
		"empty slices": {
			a: []byte{},
			b: []byte{},
		},
		"fully different": {
			a: []byte{1, 2, 3},
			b: []byte{4, 5, 6},
		},
		"fully same": {
			a:      []byte{1, 2, 3},
			b:      []byte{1, 2, 3},
			length: 3,
		},
		"different and common prefix": {
			a:      []byte{1, 2, 3, 4},
			b:      []byte{1, 2, 4, 4},
			length: 2,
		},
		"first bigger than second": {
			a:      []byte{1, 2, 3},
			b:      []byte{1, 2},
			length: 2,
		},
		"first smaller than second": {
			a:      []byte{1, 2},
			b:      []byte{1, 2, 3},
			length: 2,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			length := lenCommonPrefix(testCase.a, testCase.b)

			assert.Equal(t, testCase.length, length)
		})
	}
}
