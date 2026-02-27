package builder_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinywideclouds/go-llm/pkg/builder/v1"
)

func TestCompiledCache_JSONSerialization(t *testing.T) {
	fixedTime := time.Date(2026, 2, 27, 10, 30, 0, 0, time.UTC)

	original := builder.CompiledCache{
		ID:         "internal-cc-1",
		ExternalID: "cachedContents/gemini-789",
		Provider:   "gemini",
		CreatedAt:  fixedTime,
		AttachmentsUsed: []builder.Attachment{
			{
				ID:        "att-1",
				CacheID:   "bundle-1",
				ProfileID: "prof-xyz",
			},
			{
				ID:      "att-2",
				CacheID: "bundle-2",
				// No ProfileID to test optional pointer serialization
			},
		},
	}

	// 1. Marshal to JSON
	data, err := json.Marshal(original)
	require.NoError(t, err, "Failed to marshal CompiledCache")

	jsonStr := string(data)

	// Verify Protobuf camelCase naming was enforced
	assert.Contains(t, jsonStr, `"externalId":"cachedContents/gemini-789"`)
	assert.Contains(t, jsonStr, `"attachmentsUsed":`)
	assert.Contains(t, jsonStr, `"profileId":"prof-xyz"`)

	// 2. Unmarshal back to a new struct
	var parsed builder.CompiledCache
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err, "Failed to unmarshal CompiledCache")

	// 3. Verify exact equality
	assert.Equal(t, original.ID, parsed.ID)
	assert.Equal(t, original.ExternalID, parsed.ExternalID)
	assert.Equal(t, original.Provider, parsed.Provider)
	assert.True(t, original.CreatedAt.Equal(parsed.CreatedAt), "CreatedAt time mismatch")

	// Verify Attachments Array
	require.Len(t, parsed.AttachmentsUsed, 2)

	assert.Equal(t, "att-1", parsed.AttachmentsUsed[0].ID)
	assert.Equal(t, "bundle-1", parsed.AttachmentsUsed[0].CacheID)
	assert.Equal(t, "prof-xyz", parsed.AttachmentsUsed[0].ProfileID)

	assert.Equal(t, "att-2", parsed.AttachmentsUsed[1].ID)
	assert.Equal(t, "bundle-2", parsed.AttachmentsUsed[1].CacheID)
	assert.Equal(t, "", parsed.AttachmentsUsed[1].ProfileID) // Ensure missing profile is empty string
}
