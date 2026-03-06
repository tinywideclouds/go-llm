package builder_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinywideclouds/go-llm/pkg/builder/v1"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"
)

func TestSession_JSONSerialization(t *testing.T) {
	fixedTime := time.Date(2026, 2, 27, 12, 0, 0, 0, time.UTC)

	id, _ := urn.Parse("urn:llm:sess:123")
	cID, _ := urn.Parse("urn:gemini:compiled_cache:124/cc-eddg")

	original := builder.Session{
		ID:              id,
		CompiledCacheID: cID,
		UpdatedAt:       fixedTime,
	}

	data, err := json.Marshal(original)
	require.NoError(t, err, "Failed to marshal Session")

	jsonStr := string(data)
	assert.Contains(t, jsonStr, `"compiledCacheId":"urn:cc:456"`)

	var parsed builder.Session
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err, "Failed to unmarshal Session")

	assert.Equal(t, original.ID, parsed.ID)
	assert.Equal(t, original.CompiledCacheID, parsed.CompiledCacheID)
	assert.True(t, original.UpdatedAt.Equal(parsed.UpdatedAt), "UpdatedAt time mismatch")
}

func TestChangeProposal_JSONSerialization(t *testing.T) {
	fixedTime := time.Date(2026, 2, 27, 13, 0, 0, 0, time.UTC)

	// Note: We intentionally omit the 'Status' field here because it is a
	// Firestore-only property and does not exist in the protobuf definition.
	s, _ := urn.Parse("urn:llm:session:123")
	original := builder.ChangeProposal{
		ID:         "prop-999",
		SessionID:  s,
		FilePath:   "src/main.go",
		NewContent: "package main",
		Reasoning:  "Refactoring",
		CreatedAt:  fixedTime,
	}

	data, err := json.Marshal(original)
	require.NoError(t, err, "Failed to marshal ChangeProposal")

	jsonStr := string(data)

	// Verify Protobuf camelCase enforcement
	assert.Contains(t, jsonStr, `"sessionId":"urn:sess:123"`)
	assert.Contains(t, jsonStr, `"filePath":"src/main.go"`)
	assert.Contains(t, jsonStr, `"newContent":"package main"`)

	var parsed builder.ChangeProposal
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err, "Failed to unmarshal ChangeProposal")

	assert.Equal(t, original.ID, parsed.ID)
	assert.Equal(t, original.SessionID, parsed.SessionID)
	assert.Equal(t, original.FilePath, parsed.FilePath)
	assert.Equal(t, original.NewContent, parsed.NewContent)
	assert.True(t, original.CreatedAt.Equal(parsed.CreatedAt), "CreatedAt mismatch")
}
