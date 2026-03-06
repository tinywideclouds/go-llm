package cache_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinywideclouds/go-llm/pkg/cache/v1"
)

func TestStoreCollections_JSONSerialization(t *testing.T) {
	t.Run("Marshal to camelCase", func(t *testing.T) {
		collections := cache.StoreCollections{
			BundleCollection:   "sync_bundles",
			FilesCollection:    "sync_files",
			ProfilesCollection: "sync_profiles",
		}

		data, err := json.Marshal(collections)
		require.NoError(t, err, "Failed to marshal StoreCollections")

		jsonStr := string(data)

		// Verify Protobuf forced camelCase outputs
		assert.Contains(t, jsonStr, `"bundleCollection":"sync_bundles"`)
		assert.Contains(t, jsonStr, `"filesCollection":"sync_files"`)
		assert.Contains(t, jsonStr, `"profilesCollection":"sync_profiles"`)
	})

	t.Run("Unmarshal seamlessly maps fields", func(t *testing.T) {
		inputJSON := []byte(`{
			"bundleCollection": "bundles",
			"filesCollection": "files",
			"profilesCollection": "profiles"
		}`)

		var collections cache.StoreCollections
		err := json.Unmarshal(inputJSON, &collections)
		require.NoError(t, err, "Failed to unmarshal StoreCollections")

		assert.Equal(t, "bundles", collections.BundleCollection)
		assert.Equal(t, "files", collections.FilesCollection)
		assert.Equal(t, "profiles", collections.ProfilesCollection)
	})
}
