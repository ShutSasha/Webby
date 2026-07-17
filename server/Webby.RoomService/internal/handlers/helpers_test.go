package handlers_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertErrorResponse(t *testing.T, body string) {
	t.Helper()
	var resp map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	assert.False(t, resp["success"].(bool))
}

func assertSuccessResponse(t *testing.T, body string) {
	t.Helper()
	var resp map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	assert.True(t, resp["success"].(bool))
}

func assertSuccessWithData(t *testing.T, body string) map[string]any {
	t.Helper()
	var resp map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	require.True(t, resp["success"].(bool))
	require.NotNil(t, resp["data"])
	return resp["data"].(map[string]any)
}
