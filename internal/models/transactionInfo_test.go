package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		transaction    TransactionInfo
		expectedResult string
	}{
		{
			name: "Transaction is sent",
			transaction: TransactionInfo{
				Username: "user1",
				Amount:   100,
				IsSent:   true,
			},
			expectedResult: `{"toUser":"user1","amount":100}`,
		},
		{
			name: "Transaction is received",
			transaction: TransactionInfo{
				Username: "user1",
				Amount:   50,
				IsSent:   false,
			},
			expectedResult: `{"fromUser":"user1","amount":50}`,
		},
		{
			name: "Transaction with zero amount",
			transaction: TransactionInfo{
				Username: "user1",
				Amount:   0,
				IsSent:   true,
			},
			expectedResult: `{"toUser":"user1","amount":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.transaction.MarshalJSON()

			assert.NoError(t, err)

			assert.JSONEq(t, tt.expectedResult, string(result))
		})
	}
}
