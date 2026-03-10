package builder

import (
	"time"

	builderv1 "github.com/tinywideclouds/gen-llm/go/types/builder/v1"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"
)

type CompiledCacheProvider string

type CompiledCache struct {
	ID        urn.URN               `json:"id" firestore:"-"`
	Provider  CompiledCacheProvider `json:"provider" firestore:"provider"`
	Sources   []Attachment          `json:"sources" firestore:"sources"` // Updated
	CreatedAt time.Time             `json:"createdAt" firestore:"createdAt"`
	ExpiresAt time.Time             `json:"expiresAt" firestore:"expiresAt"`
}

// --- Protobuf Converters ---

func CompiledCacheToProto(native *CompiledCache) *builderv1.CompiledCachePb {
	if native == nil {
		return nil
	}

	return &builderv1.CompiledCachePb{
		Id:        native.ID.String(),
		Provider:  string(native.Provider),
		Sources:   AttachmentsToProto(native.Sources),
		CreatedAt: native.CreatedAt.Format(time.RFC3339),
		ExpiresAt: native.ExpiresAt.Format(time.RFC3339),
	}
}

func (c CompiledCache) MarshalJSON() ([]byte, error) {
	return protojsonMarshalOptions.Marshal(CompiledCacheToProto(&c))
}

func (c *CompiledCache) UnmarshalJSON(data []byte) error {
	var pb builderv1.CompiledCachePb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &pb); err != nil {
		return err
	}

	ID, err := urn.Parse(pb.Id)
	if err != nil {
		return err
	}
	c.ID = ID
	c.Provider = CompiledCacheProvider(pb.Provider)
	c.CreatedAt, _ = time.Parse(time.RFC3339, pb.CreatedAt)
	c.ExpiresAt, _ = time.Parse(time.RFC3339, pb.ExpiresAt)

	c.Sources, err = ProtoToAttachments(pb.Sources)
	if err != nil {
		return err
	}

	return nil
}
