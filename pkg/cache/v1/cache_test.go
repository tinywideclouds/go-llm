package cache_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinywideclouds/go-llm/pkg/cache/v1"
	"github.com/tinywideclouds/go-llm/pkg/yaml/filter"
)

func TestCreateCacheRequest_JSONSerialization(t *testing.T) {
	t.Run("Marshal to camelCase", func(t *testing.T) {
		req := cache.CreateCacheRequest{
			Repo:   "tinywideclouds/go-llm",
			Branch: "main",
		}

		data, err := json.Marshal(req)
		require.NoError(t, err, "Failed to marshal CreateCacheRequest")

		jsonStr := string(data)
		assert.Contains(t, jsonStr, `"repo":"tinywideclouds/go-llm"`)
		assert.Contains(t, jsonStr, `"branch":"main"`)
	})

	t.Run("Unmarshal accurately mapping fields", func(t *testing.T) {
		inputJSON := []byte(`{
			"repo": "tinywideclouds/go-llm",
			"branch": "develop"
		}`)

		var req cache.CreateCacheRequest
		err := json.Unmarshal(inputJSON, &req)
		require.NoError(t, err, "Failed to unmarshal CreateCacheRequest")

		assert.Equal(t, "tinywideclouds/go-llm", req.Repo)
		assert.Equal(t, "develop", req.Branch)
	})
}

func TestSyncRequest_JSONSerialization(t *testing.T) {
	t.Run("Marshal to camelCase with nested rules", func(t *testing.T) {
		req := cache.SyncRequest{
			IngestionRules: filter.FilterRules{
				Include: []string{"**/*.go", "**/*.ts"},
				Exclude: []string{"vendor/**", "node_modules/**"},
			},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err, "Failed to marshal SyncRequest")

		jsonStr := string(data)
		// Protobuf JSON marshaler forces camelCase for nested structs
		assert.Contains(t, jsonStr, `"ingestionRules":`)
		assert.Contains(t, jsonStr, `"include":["**/*.go","**/*.ts"]`)
		assert.Contains(t, jsonStr, `"exclude":["vendor/**","node_modules/**"]`)
	})

	t.Run("Unmarshal snake_case and camelCase seamlessly", func(t *testing.T) {
		// Testing resilience: simulating an incoming payload that might mix cases
		inputJSON := []byte(`{
			"ingestion_rules": {
				"include": ["src/**/*.ts"],
				"exclude": ["dist/**"]
			}
		}`)

		var req cache.SyncRequest
		err := json.Unmarshal(inputJSON, &req)
		require.NoError(t, err, "Failed to unmarshal SyncRequest")

		require.Len(t, req.IngestionRules.Include, 1)
		assert.Equal(t, "src/**/*.ts", req.IngestionRules.Include[0])

		require.Len(t, req.IngestionRules.Exclude, 1)
		assert.Equal(t, "dist/**", req.IngestionRules.Exclude[0])
	})
}

func TestProfileRequest_JSONSerialization(t *testing.T) {
	t.Run("Marshal to camelCase", func(t *testing.T) {
		req := cache.ProfileRequest{
			Name:      "Backend Go Files",
			RulesYaml: "include:\n  - \"**/*.go\"\n",
		}

		data, err := json.Marshal(req)
		require.NoError(t, err, "Failed to marshal ProfileRequest")

		jsonStr := string(data)
		assert.Contains(t, jsonStr, `"name":"Backend Go Files"`)
		// Protobuf converts `rules_yaml` to `rulesYaml`
		assert.Contains(t, jsonStr, `"rulesYaml":"include:\n  - \"**/*.go\"\n"`)
	})

	t.Run("Unmarshal accurately mapping fields", func(t *testing.T) {
		inputJSON := []byte(`{
			"name": "Frontend TS",
			"rulesYaml": "include:\n  - \"**/*.ts\"\n"
		}`)

		var req cache.ProfileRequest
		err := json.Unmarshal(inputJSON, &req)
		require.NoError(t, err, "Failed to unmarshal ProfileRequest")

		assert.Equal(t, "Frontend TS", req.Name)
		assert.Equal(t, "include:\n  - \"**/*.ts\"\n", req.RulesYaml)
	})
}
