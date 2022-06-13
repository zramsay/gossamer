package wasmer

import (
	"fmt"
	"github.com/ChainSafe/gossamer/lib/transaction"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplyExtrinsicErrors(t *testing.T) {
	testCases := []struct {
		name        string
		test        []byte
		expErr      error
		expvalidity *transaction.Validity
	}{
		{
			name:   "Valid validity",
			test:   []byte{0, 1, 1, 0},
			expErr: &TransactionValidityError{errLookupFailed},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			validity, err := decodeValidity(c.test)
			//if c.expErr == nil {
			//	require.NoError(t, err)
			//	return
			//}

			//if c.test[0] == 0 {
			//	_, ok := err.(*DispatchOutcomeError)
			//	require.True(t, ok)
			//} else {
			//	_, ok := err.(*TransactionValidityError)
			//	require.True(t, ok)
			//}
			val, ok := err.(*TransactionValidityError)
			fmt.Println(val)
			require.True(t, ok)
			require.Equal(t, c.expErr, err)
			require.Equal(t, c.expvalidity, validity)
		})
	}
}
