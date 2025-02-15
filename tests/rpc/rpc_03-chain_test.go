// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ChainSafe/gossamer/dot/rpc/modules"
	"github.com/ChainSafe/gossamer/dot/rpc/subscription"
	"github.com/ChainSafe/gossamer/lib/common"
	libutils "github.com/ChainSafe/gossamer/lib/utils"
	"github.com/ChainSafe/gossamer/tests/utils"
	"github.com/ChainSafe/gossamer/tests/utils/config"
	"github.com/ChainSafe/gossamer/tests/utils/node"
	"github.com/ChainSafe/gossamer/tests/utils/retry"
	"github.com/ChainSafe/gossamer/tests/utils/rpc"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	regex32BytesHex      = `^0x[0-9a-f]{64}$`
	regexBytesHex        = `^0x[0-9a-f]{2}[0-9a-f]*$`
	regexBytesHexOrEmpty = `^0x[0-9a-f]*$`
)

func TestChainRPC(t *testing.T) {
	if utils.MODE != rpcSuite {
		t.Log("Going to skip RPC suite tests")
		return
	}

	genesisPath := libutils.GetDevGenesisSpecPathTest(t)
	tomlConfig := config.Default()
	tomlConfig.Init.Genesis = genesisPath
	tomlConfig.Core.BABELead = true
	node := node.New(t, tomlConfig)
	ctx, cancel := context.WithCancel(context.Background())
	node.InitAndStartTest(ctx, t, cancel)

	// Wait for Gossamer to produce block 2
	errBlockNumberTooHigh := errors.New("block number is too high")
	const retryWaitDuration = 200 * time.Millisecond
	err := retry.UntilOK(ctx, retryWaitDuration, func() (ok bool, err error) {
		var header modules.ChainBlockHeaderResponse
		fetchWithTimeout(ctx, t, "chain_getHeader", "[]", &header)
		number, err := common.HexToUint(header.Number)
		if err != nil {
			return false, fmt.Errorf("cannot convert header number to uint: %w", err)
		}

		switch number {
		case 0, 1:
			return false, nil
		case 2:
			return true, nil
		default:
			return false, fmt.Errorf("%w: %d", errBlockNumberTooHigh, number)
		}
	})
	require.NoError(t, err)

	var finalizedHead string
	fetchWithTimeout(ctx, t, "chain_getFinalizedHead", "[]", &finalizedHead)
	assert.Regexp(t, regex32BytesHex, finalizedHead)

	var header modules.ChainBlockHeaderResponse
	fetchWithTimeout(ctx, t, "chain_getHeader", "[]", &header)

	// Check and clear unpredictable fields
	assert.Regexp(t, regex32BytesHex, header.StateRoot)
	header.StateRoot = ""
	assert.Regexp(t, regex32BytesHex, header.ExtrinsicsRoot)
	header.ExtrinsicsRoot = ""
	assert.Len(t, header.Digest.Logs, 2)
	for _, digestLog := range header.Digest.Logs {
		assert.Regexp(t, regexBytesHex, digestLog)
	}
	header.Digest.Logs = nil

	// Assert remaining struct with predictable fields
	expectedHeader := modules.ChainBlockHeaderResponse{
		ParentHash: finalizedHead,
		Number:     "0x02",
	}
	assert.Equal(t, expectedHeader, header)

	var block modules.ChainBlockResponse
	fetchWithTimeout(ctx, t, "chain_getBlock", fmt.Sprintf(`["`+header.ParentHash+`"]`), &block)

	// Check and clear unpredictable fields
	assert.Regexp(t, regex32BytesHex, block.Block.Header.ParentHash)
	block.Block.Header.ParentHash = ""
	assert.Regexp(t, regex32BytesHex, block.Block.Header.StateRoot)
	block.Block.Header.StateRoot = ""
	assert.Regexp(t, regex32BytesHex, block.Block.Header.ExtrinsicsRoot)
	block.Block.Header.ExtrinsicsRoot = ""
	assert.Len(t, block.Block.Header.Digest.Logs, 3)
	for _, digestLog := range block.Block.Header.Digest.Logs {
		assert.Regexp(t, regexBytesHex, digestLog)
	}
	block.Block.Header.Digest.Logs = nil
	assert.Len(t, block.Block.Body, 1)
	const bodyRegex = `^0x280403000b[0-9a-z]{8}8101$`
	assert.Regexp(t, bodyRegex, block.Block.Body[0])
	block.Block.Body = nil

	// Assert remaining struct with predictable fields
	expectedBlock := modules.ChainBlockResponse{
		Block: modules.ChainBlock{
			Header: modules.ChainBlockHeaderResponse{
				Number: "0x01",
			},
		},
	}
	assert.Equal(t, expectedBlock, block)

	var blockHash string
	fetchWithTimeout(ctx, t, "chain_getBlockHash", "[]", &blockHash)
	assert.Regexp(t, regex32BytesHex, blockHash)
	assert.NotEqual(t, finalizedHead, blockHash)
}

func TestChainSubscriptionRPC(t *testing.T) {
	if utils.MODE != rpcSuite {
		t.Log("Going to skip RPC suite tests")
		return
	}

	genesisPath := libutils.GetDevGenesisSpecPathTest(t)
	tomlConfig := config.Default()
	tomlConfig.Init.Genesis = genesisPath
	tomlConfig.Core.BABELead = true
	tomlConfig.RPC.WS = true // WS port is set in the node.New constructor
	node := node.New(t, tomlConfig)
	ctx, cancel := context.WithCancel(context.Background())
	node.InitAndStartTest(ctx, t, cancel)

	const endpoint = "ws://localhost:8546/"

	t.Run("chain_subscribeNewHeads", func(t *testing.T) {
		t.Parallel()

		const numberOfMesages = 2
		messages := callAndSubscribeWebsocket(ctx, t, endpoint, "chain_subscribeNewHeads", "[]", numberOfMesages)

		allParams := make([]subscription.Params, numberOfMesages)
		for i, message := range messages {
			err := rpc.Decode(message, &allParams[i])
			require.NoError(t, err, "cannot decode websocket message for message index %d", i)
		}

		for i, params := range allParams {
			result := getResultMapFromParams(t, params)

			number := getResultNumber(t, result)
			assert.Equal(t, uint(i+1), number)

			assertResultRegex(t, result, "parentHash", regex32BytesHex)
			assertResultRegex(t, result, "stateRoot", regex32BytesHex)
			assertResultRegex(t, result, "extrinsicsRoot", regex32BytesHex)
			assertResultDigest(t, result)

			remainingExpected := subscription.Params{
				Result:         map[string]interface{}{},
				SubscriptionID: 1,
			}
			assert.Equal(t, remainingExpected, params)
		}
	})

	t.Run("state_subscribeStorage", func(t *testing.T) {
		t.Parallel()

		const numberOfMesages = 2
		messages := callAndSubscribeWebsocket(ctx, t, endpoint, "state_subscribeStorage", "[]", numberOfMesages)

		allParams := make([]subscription.Params, numberOfMesages)
		for i := range allParams {
			message := messages[i]
			err := rpc.Decode(message, &allParams[i])
			require.NoError(t, err, "cannot decode websocket message for message index %d", i)
		}

		for i, params := range allParams {
			errorContext := fmt.Sprintf("for response at index %d", i)

			result := getResultMapFromParams(t, params)

			blockHex, ok := result["block"].(string)
			require.True(t, ok, errorContext)
			assert.Regexp(t, regex32BytesHex, blockHex, errorContext)
			delete(result, "block")

			changes, ok := result["changes"].([]interface{})
			require.True(t, ok, errorContext)

			for _, change := range changes {
				fromTo, ok := change.([]interface{})
				require.Truef(t, ok, "%s and change: %v", errorContext, change)
				from, ok := fromTo[0].(string)
				require.Truef(t, ok, "%s and from: %v", errorContext, fromTo[0])
				to, ok := fromTo[1].(string)
				require.Truef(t, ok, "%s and to: %v", errorContext, fromTo[1])
				assert.Regexp(t, regexBytesHexOrEmpty, from, errorContext)
				assert.Regexp(t, regexBytesHexOrEmpty, to, errorContext)
			}
			delete(result, "changes")

			remainingExpected := map[string]interface{}{}
			assert.Equal(t, remainingExpected, result, errorContext)
		}
	})

	t.Run("chain_subscribeFinalizedHeads", func(t *testing.T) {
		t.Parallel()

		const numberOfMesages = 4
		messages := callAndSubscribeWebsocket(ctx, t, endpoint, "chain_subscribeFinalizedHeads", "[]", numberOfMesages)

		allParams := make([]subscription.Params, numberOfMesages)
		for i, message := range messages {
			err := rpc.Decode(message, &allParams[i])
			require.NoError(t, err, "cannot decode websocket message for message index %d", i)
		}

		var blockNumbers []uint
		for _, params := range allParams {
			result := getResultMapFromParams(t, params)

			number := getResultNumber(t, result)
			blockNumbers = append(blockNumbers, number)

			assertResultRegex(t, result, "parentHash", regex32BytesHex)
			assertResultRegex(t, result, "stateRoot", regex32BytesHex)
			assertResultRegex(t, result, "extrinsicsRoot", regex32BytesHex)
			assertResultDigest(t, result)

			remainingExpected := subscription.Params{
				Result:         map[string]interface{}{},
				SubscriptionID: 1,
			}
			assert.Equal(t, remainingExpected, params)
		}

		// Check block numbers grow by zero or one in order of responses.
		for i, blockNumber := range blockNumbers {
			if i == 0 {
				assert.Equal(t, uint(1), blockNumber)
				continue
			}
			assert.GreaterOrEqual(t, blockNumber, blockNumbers[i-1])
		}
	})
}

func getResultMapFromParams(t *testing.T, params subscription.Params) (
	resultMap map[string]interface{}) {
	t.Helper()

	resultMap, ok := params.Result.(map[string]interface{})
	require.True(t, ok)

	return resultMap
}

// getResultNumber returns the number value from the result map
// and deletes the "number" key from the map.
func getResultNumber(t *testing.T, result map[string]interface{}) uint {
	t.Helper()

	hexNumber, ok := result["number"].(string)
	require.True(t, ok)

	number, err := common.HexToUint(hexNumber)
	require.NoError(t, err)
	delete(result, "number")

	return number
}

// assertResultRegex gets the value from the map and asserts that it matches the regex.
// It then removes the key from the map.
func assertResultRegex(t *testing.T, result map[string]interface{}, key, regex string) {
	t.Helper()

	value, ok := result[key]
	require.True(t, ok, "cannot find key %q in result", key)
	assert.Regexp(t, regex, value, "at result key %q", key)
	delete(result, key)
}

func assertResultDigest(t *testing.T, result map[string]interface{}) {
	t.Helper()

	digest, ok := result["digest"].(map[string]interface{})
	require.True(t, ok)

	logs, ok := digest["logs"].([]interface{})
	require.True(t, ok)

	assert.NotEmpty(t, logs)
	for _, log := range logs {
		assert.Regexp(t, regexBytesHex, log)
	}

	delete(result, "digest")
}

func callAndSubscribeWebsocket(ctx context.Context, t *testing.T,
	endpoint, method, params string, numberOfMesages uint) (
	messages [][]byte) {
	t.Helper()

	connection, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	require.NoError(t, err, "cannot dial websocket")
	defer connection.Close() // in case of failed required assertion

	const maxid = 100000 // otherwise it becomes a float64
	id := rand.Intn(maxid)
	messageData := fmt.Sprintf(`{
    "jsonrpc": "2.0",
    "method": %q,
    "params": [%s],
    "id": %d
}`, method, params, id)
	err = connection.WriteMessage(websocket.TextMessage, []byte(messageData))
	require.NoError(t, err, "cannot write websocket message")

	// Read subscription id result
	var target subscription.ResponseJSON
	err = connection.ReadJSON(&target)
	require.NoError(t, err, "cannot read websocket message")
	assert.Equal(t, float64(id), target.ID, "request id mismatch")
	assert.NotZero(t, target.Result, "subscription id is 0")

	for i := uint(0); i < numberOfMesages; i++ {
		_, data, err := connection.ReadMessage()
		require.NoError(t, err, "cannot read websocket message")

		messages = append(messages, data)
	}

	// Close connection
	const messageType = websocket.CloseMessage
	data := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	err = connection.WriteMessage(messageType, data)
	assert.NoError(t, err, "cannot write close websocket message")
	err = connection.Close()
	assert.NoError(t, err, "cannot close websocket connection")

	return messages
}
