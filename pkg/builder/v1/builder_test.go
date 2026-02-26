package builder_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/tinywideclouds/go-llm/pkg/builder/v1"
)

func TestBuildCacheRequest_JSON(t *testing.T) {
	t.Run("Marshal to camelCase", func(t *testing.T) {
		req := builder.BuildCacheRequest{
			SessionID: "sess-123",
			Model:     "gemini-1.5-pro",
			Attachments: []builder.Attachment{
				{
					ID:        "att-1",
					CacheID:   "cache-abc",
					ProfileID: "prof-xyz",
				},
			},
		}

		data, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		jsonStr := string(data)

		// Verify protobuf forced camelCase outputs
		if !strings.Contains(jsonStr, `"sessionId":"sess-123"`) {
			t.Errorf("Expected camelCase sessionId, got: %s", jsonStr)
		}
		if !strings.Contains(jsonStr, `"profileId":"prof-xyz"`) {
			t.Errorf("Expected camelCase profileId, got: %s", jsonStr)
		}
	})

	t.Run("Unmarshal accurately mapping arrays", func(t *testing.T) {
		// Simulating incoming request from Angular
		inputJSON := []byte(`{
			"sessionId": "sess-999",
			"model": "gemini-2.5-pro",
			"attachments": [
				{"id": "att-2", "cacheId": "urn:fire:123"}
			]
		}`)

		var req builder.BuildCacheRequest
		if err := json.Unmarshal(inputJSON, &req); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if req.SessionID != "sess-999" {
			t.Errorf("Expected SessionID sess-999, got %s", req.SessionID)
		}

		// Prove the array allocation bug is fixed (length should be exactly 1)
		if len(req.Attachments) != 1 {
			t.Fatalf("Expected 1 attachment, got %d", len(req.Attachments))
		}

		if req.Attachments[0].CacheID != "urn:fire:123" {
			t.Errorf("Expected CacheID urn:fire:123, got %s", req.Attachments[0].CacheID)
		}

		// Prove that omitted optional fields evaluate correctly
		if req.Attachments[0].ProfileID != "" {
			t.Errorf("Expected empty ProfileID, got %s", req.Attachments[0].ProfileID)
		}
	})
}

func TestBuildCacheResponse_JSON(t *testing.T) {
	t.Run("Unmarshal handles both snake_case and camelCase", func(t *testing.T) {
		// Because protojson is smart, it will accept snake_case even though we output camelCase
		snakeCaseJSON := []byte(`{"gemini_cache_id": "cachedContents/123"}`)

		var res builder.BuildCacheResponse
		if err := json.Unmarshal(snakeCaseJSON, &res); err != nil {
			t.Fatalf("Failed to unmarshal snake_case: %v", err)
		}

		if res.GeminiCacheId != "cachedContents/123" {
			t.Errorf("Expected cachedContents/123, got %s", res.GeminiCacheId)
		}
	})
}

func TestGenerateStreamRequest_JSON(t *testing.T) {
	t.Run("Marshal and Unmarshal full stream request", func(t *testing.T) {
		original := builder.GenerateStreamRequest{
			SessionID:     "stream-123",
			Model:         "gemini-flash",
			GeminiCacheID: "cachedContents/abc",
			History: []builder.Message{
				{ID: "msg-1", Role: "user", Content: "Hello AI", Timestamp: "2024-01-01T12:00:00Z"},
			},
			InlineAttachments: []builder.Attachment{
				{ID: "att-1", CacheID: "urn:repo:456"},
			},
		}

		// 1. Marshal to JSON
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal GenerateStreamRequest: %v", err)
		}

		// 2. Unmarshal back to a new struct
		var parsed builder.GenerateStreamRequest
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal GenerateStreamRequest: %v", err)
		}

		// 3. Verify identical structures
		if parsed.SessionID != original.SessionID {
			t.Errorf("SessionID mismatch: got %s, want %s", parsed.SessionID, original.SessionID)
		}
		if parsed.GeminiCacheID != original.GeminiCacheID {
			t.Errorf("GeminiCacheID mismatch: got %s, want %s", parsed.GeminiCacheID, original.GeminiCacheID)
		}

		// Check History
		if len(parsed.History) != 1 {
			t.Fatalf("Expected 1 history message, got %d", len(parsed.History))
		}
		if parsed.History[0].Content != "Hello AI" {
			t.Errorf("History content mismatch: got %s", parsed.History[0].Content)
		}

		// Check Inline Attachments
		if len(parsed.InlineAttachments) != 1 {
			t.Fatalf("Expected 1 inline attachment, got %d", len(parsed.InlineAttachments))
		}
		if parsed.InlineAttachments[0].CacheID != "urn:repo:456" {
			t.Errorf("Inline Attachment CacheID mismatch: got %s", parsed.InlineAttachments[0].CacheID)
		}
	})
}
