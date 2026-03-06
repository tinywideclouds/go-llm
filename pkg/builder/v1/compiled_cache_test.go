package builder_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinywideclouds/go-llm/pkg/builder/v1"
)

// mustURN and urnPtr are reused or redefined here depending on your test suite setup
func TestCompiledCache_JSONSerialization(t *testing.T) {
	fixedTime := time.Date(2026, 2, 27, 10, 30, 0, 0, time.UTC)

	original := builder.CompiledCache{
		ID:        mustURN("urn:llm:compiled-cache:124"),
		Provider:  "gemini",
		CreatedAt: fixedTime,
		ExpiresAt: fixedTime.Add(time.Hour),
		AttachmentsUsed: []builder.Attachment{
			{
				ID:        mustURN("urn:llm:attachment:1"),
				CacheID:   mustURN("urn:llm:cache:1"),
				ProfileID: urnPtr("urn:llm:profile:xyz"),
			},
			{
				ID:      mustURN("urn:llm:attachment:2"),
				CacheID: mustURN("urn:llm:cache:2"),
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
	assert.Contains(t, jsonStr, `"profileId":"urn:llm:profile:xyz"`)

	// 2. Unmarshal back to a new struct
	var parsed builder.CompiledCache
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err, "Failed to unmarshal CompiledCache")

	// 3. Verify exact equality
	assert.Equal(t, original.ID, parsed.ID)
	assert.Equal(t, original.Provider, parsed.Provider)
	assert.True(t, original.CreatedAt.Equal(parsed.CreatedAt), "CreatedAt time mismatch")

	// Verify Attachments Array
	require.Len(t, parsed.AttachmentsUsed, 2)

	assert.Equal(t, mustURN("urn:llm:attachment:1"), parsed.AttachmentsUsed[0].ID)
	assert.Equal(t, mustURN("urn:llm:cache:1"), parsed.AttachmentsUsed[0].CacheID)
	assert.Equal(t, urnPtr("urn:llm:profile:xyz"), parsed.AttachmentsUsed[0].ProfileID)

	assert.Equal(t, mustURN("urn:llm:attachment:2"), parsed.AttachmentsUsed[1].ID)
	assert.Equal(t, mustURN("urn:llm:cache:2"), parsed.AttachmentsUsed[1].CacheID)
	assert.Nil(t, parsed.AttachmentsUsed[1].ProfileID)
}
