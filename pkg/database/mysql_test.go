package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMySqlConnection(t *testing.T) {
	db, err := NewMySQLClient()
	require.NotNil(t, db)
	require.Nil(t, err)
}
