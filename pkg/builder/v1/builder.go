package builder

import (
	"time"

	builderv1 "github.com/tinywideclouds/gen-llm/go/types/builder/v1"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	protojsonMarshalOptions = &protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: false,
	}
	protojsonUnmarshalOptions = &protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

// --- BUILD CACHE RESPONSE ---

type BuildCacheResponse struct {
	CompiledCacheId urn.URN   `json:"compiledCacheId"`
	ExpiresAt       time.Time `json:"expiresAt"`
}

func (pk BuildCacheResponse) MarshalJSON() ([]byte, error) {
	expires := pk.ExpiresAt.Format(time.RFC3339)
	protoPb := &builderv1.BuildCacheResponsePb{CompiledCacheId: pk.CompiledCacheId.String(), ExpiresAt: expires}
	return protojsonMarshalOptions.Marshal(protoPb)
}

func (pk *BuildCacheResponse) UnmarshalJSON(data []byte) error {
	var protoPb builderv1.BuildCacheResponsePb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}
	var err error
	pk.CompiledCacheId, err = urn.Parse(protoPb.CompiledCacheId)
	if err != nil {
		return err
	}
	expiresAt, err := time.Parse(time.RFC3339, protoPb.ExpiresAt)
	if err != nil {
		expiresAt = time.Now().Add(time.Hour)
	}
	pk.ExpiresAt = expiresAt
	return nil
}

// --- BUILD CACHE REQUEST ---

type Attachment struct {
	ID        urn.URN  `json:"id"`
	CacheID   urn.URN  `json:"cacheId"`
	ProfileID *urn.URN `json:"profileId,omitempty"`
}

type BuildCacheRequest struct {
	SessionID     urn.URN      `json:"sessionId"`
	Model         string       `json:"model"`
	Attachments   []Attachment `json:"attachments"`
	ExpiresAtHint *time.Time   `json:"expiresAtHint,omitempty"`
}

func CacheRequestToProto(native *BuildCacheRequest) *builderv1.BuildCacheRequestPb {
	if native == nil {
		return nil
	}

	pb := &builderv1.BuildCacheRequestPb{
		SessionId:   native.SessionID.String(),
		Model:       native.Model,
		Attachments: AttachmentsToProto(native.Attachments),
	}

	if native.ExpiresAtHint != nil {
		expiresHint := native.ExpiresAtHint.Format(time.RFC3339)
		pb.ExpiresAtHint = &expiresHint
	}

	return pb
}

func (pk BuildCacheRequest) MarshalJSON() ([]byte, error) {
	return protojsonMarshalOptions.Marshal(CacheRequestToProto(&pk))
}

func (pk *BuildCacheRequest) UnmarshalJSON(data []byte) error {
	var protoPb builderv1.BuildCacheRequestPb
	err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb)
	if err != nil {
		return err
	}

	sessionID, err := urn.Parse(protoPb.SessionId)
	if err != nil {
		return err
	}

	pk.SessionID = sessionID
	pk.Model = protoPb.Model

	pk.Attachments, err = ProtoToAttachments(protoPb.Attachments)
	if err != nil {
		return err
	}

	if protoPb.ExpiresAtHint != nil {
		expires, err := time.Parse(time.RFC3339, *protoPb.ExpiresAtHint)
		if err != nil {
			return err
		}
		pk.ExpiresAtHint = &expires
	}

	return nil
}

// --- GENERATE STREAM REQUEST ---

type Message struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type GenerateStreamRequest struct {
	SessionID         urn.URN      `json:"sessionId"`
	Model             string       `json:"model"`
	History           []Message    `json:"history"`
	CompiledCacheID   *urn.URN     `json:"compiledCacheId,omitempty"`
	InlineAttachments []Attachment `json:"inlineAttachments,omitempty"`
}

func ToStreamProto(native *GenerateStreamRequest) *builderv1.GenerateStreamRequestPb {
	if native == nil {
		return nil
	}

	history := make([]*builderv1.NetworkMessagePb, 0, len(native.History))
	for _, m := range native.History {
		history = append(history, &builderv1.NetworkMessagePb{
			Id:        m.ID,
			Role:      m.Role,
			Content:   m.Content,
			Timestamp: m.Timestamp,
		})
	}

	pb := &builderv1.GenerateStreamRequestPb{
		SessionId:         native.SessionID.String(),
		Model:             native.Model,
		History:           history,
		InlineAttachments: AttachmentsToProto(native.InlineAttachments),
	}

	if native.CompiledCacheID != nil {
		c := native.CompiledCacheID.String()
		pb.CompiledCacheId = &c
	}

	return pb
}

func (pk GenerateStreamRequest) MarshalJSON() ([]byte, error) {
	return protojsonMarshalOptions.Marshal(ToStreamProto(&pk))
}

func (pk *GenerateStreamRequest) UnmarshalJSON(data []byte) error {
	var protoPb builderv1.GenerateStreamRequestPb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}

	sessionID, err := urn.Parse(protoPb.SessionId)
	if err != nil {
		return err
	}
	pk.SessionID = sessionID
	pk.Model = protoPb.Model

	if protoPb.CompiledCacheId != nil {
		c, err := urn.Parse(*protoPb.CompiledCacheId)
		if err != nil {
			return err
		}
		pk.CompiledCacheID = &c
	}

	pk.History = make([]Message, 0, len(protoPb.History))
	for _, m := range protoPb.History {
		pk.History = append(pk.History, Message{
			ID:        m.Id,
			Role:      m.Role,
			Content:   m.Content,
			Timestamp: m.Timestamp,
		})
	}

	pk.InlineAttachments, err = ProtoToAttachments(protoPb.InlineAttachments)
	if err != nil {
		return err
	}

	return nil
}

// --- HELPERS ---

func AttachmentsToProto(native []Attachment) []*builderv1.NetworkAttachmentPb {
	attachments := make([]*builderv1.NetworkAttachmentPb, 0, len(native))
	for _, m := range native {
		pbAtt := &builderv1.NetworkAttachmentPb{
			Id:      m.ID.String(),
			CacheId: m.CacheID.String(),
		}
		if m.ProfileID != nil {
			pid := m.ProfileID.String()
			pbAtt.ProfileId = &pid
		}
		attachments = append(attachments, pbAtt)
	}
	return attachments
}

func ProtoToAttachments(pb []*builderv1.NetworkAttachmentPb) ([]Attachment, error) {
	attachments := make([]Attachment, 0, len(pb))
	for _, a := range pb {
		id, err := urn.Parse(a.Id)
		if err != nil {
			return nil, err
		}
		cacheID, err := urn.Parse(a.CacheId)
		if err != nil {
			return nil, err
		}
		att := Attachment{
			ID:      id,
			CacheID: cacheID,
		}
		if a.ProfileId != nil {
			profileID, err := urn.Parse(*a.ProfileId)
			if err != nil {
				return nil, err
			}
			att.ProfileID = &profileID
		}
		attachments = append(attachments, att)
	}
	return attachments, nil
}
