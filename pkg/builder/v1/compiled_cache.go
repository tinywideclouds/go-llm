package builder

import (
	"time"

	builderv1 "github.com/tinywideclouds/gen-llm/go/types/builder/v1"
)

type CompiledCache struct {
	ID              string       `json:"id" firestore:"-"`
	ExternalID      string       `json:"externalId" firestore:"externalId"`
	Provider        string       `json:"provider" firestore:"provider"`
	AttachmentsUsed []Attachment `json:"attachmentsUsed" firestore:"attachmentsUsed"`
	CreatedAt       time.Time    `json:"createdAt" firestore:"createdAt"`
	ExpiresAt       time.Time    `json:"expiresAt" firestore:"expiresAt"`
}

// --- Protobuf Converters ---

func CompiledCacheToProto(native *CompiledCache) *builderv1.CompiledCachePb {
	if native == nil {
		return nil
	}

	attachments := make([]*builderv1.NetworkAttachmentPb, 0, len(native.AttachmentsUsed))
	for _, att := range native.AttachmentsUsed {
		pbAtt := &builderv1.NetworkAttachmentPb{
			Id:      att.ID,
			CacheId: att.CacheID,
		}
		if att.ProfileID != "" {
			pid := att.ProfileID
			pbAtt.ProfileId = &pid
		}
		attachments = append(attachments, pbAtt)
	}

	return &builderv1.CompiledCachePb{
		Id:              native.ID,
		ExternalId:      native.ExternalID,
		Provider:        native.Provider,
		AttachmentsUsed: attachments,
		CreatedAt:       native.CreatedAt.Format(time.RFC3339),
		ExpiresAt:       native.ExpiresAt.Format(time.RFC3339),
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

	c.ID = pb.Id
	c.ExternalID = pb.ExternalId
	c.Provider = pb.Provider
	c.CreatedAt, _ = time.Parse(time.RFC3339, pb.CreatedAt)
	c.ExpiresAt, _ = time.Parse(time.RFC3339, pb.ExpiresAt)

	c.AttachmentsUsed = make([]Attachment, 0, len(pb.AttachmentsUsed))
	for _, a := range pb.AttachmentsUsed {
		att := Attachment{
			ID:      a.Id,
			CacheID: a.CacheId,
		}
		if a.ProfileId != nil {
			att.ProfileID = *a.ProfileId
		}
		c.AttachmentsUsed = append(c.AttachmentsUsed, att)
	}

	return nil
}
