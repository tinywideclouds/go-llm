package builder_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/tinywideclouds/go-llm/pkg/builder/v1"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"
)

// Test helpers to cleanly initialize strict URN types
func mustURN(s string) urn.URN {
	u, err := urn.Parse(s)
	if err != nil {
		panic("invalid test URN: " + s)
	}
	return u
}

func urnPtr(s string) *urn.URN {
	u := mustURN(s)
	return &u
}

func TestBuildCacheRequest_JSON(t *testing.T) {
	t.Run("Marshal to camelCase", func(t *testing.T) {
		req := builder.BuildCacheRequest{
			Model: "gemini-1.5-pro",
			Sources: []builder.Attachment{
				{
					ID:           mustURN("urn:llm:attachment:1"),
					DataSourceID: mustURN("urn:llm:cache:abc"),
					ProfileID:    urnPtr("urn:llm:profile:xyz"),
				},
			},
		}

		data, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		jsonStr := string(data)

		// Verify protobuf forced camelCase outputs
		if !strings.Contains(jsonStr, `"dataSourceId":"urn:llm:cache:abc"`) {
			t.Errorf("Expected camelCase dataSourceId, got: %s", jsonStr)
		}
		if !strings.Contains(jsonStr, `"profileId":"urn:llm:profile:xyz"`) {
			t.Errorf("Expected camelCase profileId, got: %s", jsonStr)
		}
	})

	t.Run("Unmarshal accurately mapping arrays", func(t *testing.T) {
		inputJSON := []byte(`{
			"model": "gemini-2.5-pro",
			"sources": [
				{"id": "urn:llm:attachment:2", "dataSourceId": "urn:llm:cache:123"}
			]
		}`)

		var req builder.BuildCacheRequest
		if err := json.Unmarshal(inputJSON, &req); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if len(req.Sources) != 1 {
			t.Fatalf("Expected 1 attachment, got %d", len(req.Sources))
		}

		if req.Sources[0].DataSourceID.String() != "urn:llm:cache:123" {
			t.Errorf("Expected DataSourceID urn:llm:cache:123, got %s", req.Sources[0].DataSourceID)
		}

		if req.Sources[0].ProfileID != nil {
			t.Errorf("Expected nil ProfileID, got %s", req.Sources[0].ProfileID)
		}
	})
}

func TestBuildCacheResponse_JSON(t *testing.T) {
	t.Run("Unmarshal handles both snake_case and camelCase", func(t *testing.T) {
		snakeCaseJSON := []byte(`{"compiled_cache_id": "urn:llm:compiled-cache:123"}`)

		var res builder.BuildCacheResponse
		if err := json.Unmarshal(snakeCaseJSON, &res); err != nil {
			t.Fatalf("Failed to unmarshal snake_case: %v", err)
		}

		if res.CompiledCacheId.String() != "urn:llm:compiled-cache:123" {
			t.Errorf("Expected urn:llm:compiled-cache:123, got %s", res.CompiledCacheId)
		}
	})
}

func TestGenerateStreamRequest_JSON(t *testing.T) {
	t.Run("Marshal and Unmarshal full stream request", func(t *testing.T) {
		original := builder.GenerateStreamRequest{
			SessionID:       mustURN("urn:llm:session:stream-123"),
			Model:           "gemini-flash",
			CompiledCacheID: urnPtr("urn:llm:compiled-cache:abc"),
			History: []builder.Message{
				{ID: "msg-1", Role: "user", Content: "Hello AI", Timestamp: "2024-01-01T12:00:00Z"},
			},
			InlineAttachments: []builder.Attachment{
				{ID: mustURN("urn:llm:attachment:1"), DataSourceID: mustURN("urn:llm:cache:456")},
			},
		}

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal GenerateStreamRequest: %v", err)
		}

		var parsed builder.GenerateStreamRequest
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal GenerateStreamRequest: %v", err)
		}

		if parsed.SessionID != original.SessionID {
			t.Errorf("SessionID mismatch: got %s, want %s", parsed.SessionID, original.SessionID)
		}
		if parsed.CompiledCacheID == nil || *parsed.CompiledCacheID != *original.CompiledCacheID {
			t.Errorf("CompiledCacheID mismatch: got %v, want %v", parsed.CompiledCacheID, original.CompiledCacheID)
		}

		if len(parsed.History) != 1 {
			t.Fatalf("Expected 1 history message, got %d", len(parsed.History))
		}
		if parsed.History[0].Content != "Hello AI" {
			t.Errorf("History content mismatch: got %s", parsed.History[0].Content)
		}

		if len(parsed.InlineAttachments) != 1 {
			t.Fatalf("Expected 1 inline attachment, got %d", len(parsed.InlineAttachments))
		}
		if parsed.InlineAttachments[0].DataSourceID != original.InlineAttachments[0].DataSourceID {
			t.Errorf("Inline Attachment DataSourceID mismatch: got %s", parsed.InlineAttachments[0].DataSourceID)
		}
	})
}
