package builder_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinywideclouds/go-llm/pkg/builder/v1"
)

func TestSession_JSONSerialization(t *testing.T) {
	// Use a fixed time to avoid monotonic clock issues during deep equality checks
	fixedTime := time.Date(2026, 2, 27, 12, 0, 0, 0, time.UTC)

	original := builder.Session{
		ID:              "sess-123",
		CompiledCacheID: "cc-456",
		UpdatedAt:       fixedTime,
		AcceptedOverlays: map[string]builder.FileState{
			"main.go": {
				Content:   "package main",
				IsDeleted: false,
			},
		},
		PendingProposals: map[string]builder.ChangeProposal{
			"prop-1": {
				ID:         "prop-1",
				FilePath:   "auth.go",
				NewContent: "package auth",
				Reasoning:  "Updated auth logic",
				Status:     builder.StatusPending,
				CreatedAt:  fixedTime,
			},
		},
	}

	// 1. Marshal to JSON
	data, err := json.Marshal(original)
	require.NoError(t, err, "Failed to marshal Session")

	jsonStr := string(data)

	// Verify Protobuf camelCase naming was enforced
	assert.Contains(t, jsonStr, `"compiledCacheId":"cc-456"`)
	assert.Contains(t, jsonStr, `"acceptedOverlays":`)
	assert.Contains(t, jsonStr, `"pendingProposals":`)
	assert.Contains(t, jsonStr, `"filePath":"auth.go"`)
	assert.Contains(t, jsonStr, `"isDeleted":false`)

	// 2. Unmarshal back to a new struct
	var parsed builder.Session
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err, "Failed to unmarshal Session")

	// 3. Verify exact equality
	assert.Equal(t, original.ID, parsed.ID)
	assert.Equal(t, original.CompiledCacheID, parsed.CompiledCacheID)
	assert.True(t, original.UpdatedAt.Equal(parsed.UpdatedAt), "UpdatedAt time mismatch")

	// Verify Maps
	require.Contains(t, parsed.AcceptedOverlays, "main.go")
	assert.Equal(t, original.AcceptedOverlays["main.go"].Content, parsed.AcceptedOverlays["main.go"].Content)
	assert.Equal(t, original.AcceptedOverlays["main.go"].IsDeleted, parsed.AcceptedOverlays["main.go"].IsDeleted)

	require.Contains(t, parsed.PendingProposals, "prop-1")
	prop := parsed.PendingProposals["prop-1"]
	assert.Equal(t, "auth.go", prop.FilePath)
	assert.Equal(t, "Updated auth logic", prop.Reasoning)
	assert.Equal(t, builder.StatusPending, prop.Status)
	assert.True(t, original.PendingProposals["prop-1"].CreatedAt.Equal(prop.CreatedAt), "Proposal CreatedAt mismatch")
}
