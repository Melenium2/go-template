package tx

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestExtractTx_Should_extract_provided_in_the_context_tx(t *testing.T) {
	expected := &sqlx.Tx{}

	ctx := context.WithValue(context.Background(), txKey, expected)

	tx, err := extractTx(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, tx)
}

func TestExtractTx_Should_return_error_if_tx_by_txKey_will_be_nil(t *testing.T) {
	ctx := context.WithValue(context.Background(), txKey, nil)

	tx, err := extractTx(ctx)
	assert.Error(t, err)
	assert.Nil(t, tx)
}

func TestExtractTx_Should_return_error_if_nothing_found_by_txKey(t *testing.T) {
	_, err := extractTx(context.Background())
	assert.Error(t, err)
}

func TestInjectTx_Should_inject_provided_tx_to_context(t *testing.T) {
	expected := &sqlx.Tx{}

	ctx := injectTx(context.Background(), expected)

	tx, ok := ctx.Value(txKey).(*sqlx.Tx)
	assert.True(t, ok)
	assert.Equal(t, expected, tx)
}
