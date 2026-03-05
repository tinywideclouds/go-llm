package builder

import (
	"time"

	builderv1 "github.com/tinywideclouds/gen-llm/go/types/builder/v1"
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
	GeminiCacheId string    `json:"geminiCacheId"`
	ExpiresAt     time.Time `json:"expiresAt"`
}

func (pk BuildCacheResponse) MarshalJSON() ([]byte, error) {
	expires := pk.ExpiresAt.Format(time.RFC1123)
	protoPb := &builderv1.BuildCacheResponsePb{GeminiCacheId: pk.GeminiCacheId, ExpiresAt: expires}
	return protojsonMarshalOptions.Marshal(protoPb)
}

func (pk *BuildCacheResponse) UnmarshalJSON(data []byte) error {
	var protoPb builderv1.BuildCacheResponsePb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}
	pk.GeminiCacheId = protoPb.GeminiCacheId
	expiresAt, err := time.Parse(protoPb.ExpiresAt, time.RFC3339)
	if err != nil {
		expiresAt = time.Now().Add(time.Hour)
	}
	pk.ExpiresAt = expiresAt
	return nil
}

// --- BUILD CACHE REQUEST ---

type Attachment struct {
	ID        string `json:"id"`
	CacheID   string `json:"cacheId"`
	ProfileID string `json:"profileId,omitempty"`
}

type BuildCacheRequest struct {
	SessionID     string       `json:"sessionId"`
	Model         string       `json:"model"`
	Attachments   []Attachment `json:"attachments"`
	ExpiresAtHint time.Time    `json:"expiresAtHint"`
}

func CacheRequestToProto(native *BuildCacheRequest) *builderv1.BuildCacheRequestPb {
	if native == nil {
		return nil
	}

	// FIX: Pre-allocate capacity, but start length at 0, then append correctly.
	attachments := make([]*builderv1.NetworkAttachmentPb, 0, len(native.Attachments))

	for _, m := range native.Attachments {
		pbAtt := &builderv1.NetworkAttachmentPb{
			Id:      m.ID,
			CacheId: m.CacheID,
		}
		if m.ProfileID != "" {
			pid := m.ProfileID // Need a local variable to take the pointer
			pbAtt.ProfileId = &pid
		}
		attachments = append(attachments, pbAtt)
	}

	expiresHint := native.ExpiresAtHint.Format(time.RFC1123)

	return &builderv1.BuildCacheRequestPb{
		SessionId:     native.SessionID,
		Model:         native.Model,
		Attachments:   attachments,
		ExpiresAtHint: &expiresHint,
	}
}

func (pk BuildCacheRequest) MarshalJSON() ([]byte, error) {
	return protojsonMarshalOptions.Marshal(CacheRequestToProto(&pk))
}

func (pk *BuildCacheRequest) UnmarshalJSON(data []byte) error {
	var protoPb builderv1.BuildCacheRequestPb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}

	pk.SessionID = protoPb.SessionId
	pk.Model = protoPb.Model

	pk.Attachments = make([]Attachment, 0, len(protoPb.Attachments))
	for _, a := range protoPb.Attachments {
		att := Attachment{
			ID:      a.Id,
			CacheID: a.CacheId,
		}
		if a.ProfileId != nil {
			att.ProfileID = *a.ProfileId
		}
		pk.Attachments = append(pk.Attachments, att)
	}

	if protoPb.ExpiresAtHint != nil {
		expires, err := time.Parse(*protoPb.ExpiresAtHint, time.RFC3339)
		if err != nil {
			return err
		}
		pk.ExpiresAtHint = expires
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
	SessionID         string       `json:"sessionId"`
	Model             string       `json:"model"`
	History           []Message    `json:"history"`
	GeminiCacheID     string       `json:"geminiCacheId,omitempty"`
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

	inlines := make([]*builderv1.NetworkAttachmentPb, 0, len(native.InlineAttachments))
	for _, a := range native.InlineAttachments {
		pbAtt := &builderv1.NetworkAttachmentPb{
			Id:      a.ID,
			CacheId: a.CacheID,
		}
		if a.ProfileID != "" {
			pid := a.ProfileID
			pbAtt.ProfileId = &pid
		}
		inlines = append(inlines, pbAtt)
	}

	pb := &builderv1.GenerateStreamRequestPb{
		SessionId:         native.SessionID,
		Model:             native.Model,
		History:           history,
		InlineAttachments: inlines,
	}

	if native.GeminiCacheID != "" {
		cid := native.GeminiCacheID
		pb.GeminiCacheId = &cid
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

	pk.SessionID = protoPb.SessionId
	pk.Model = protoPb.Model

	if protoPb.GeminiCacheId != nil {
		pk.GeminiCacheID = *protoPb.GeminiCacheId
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

	pk.InlineAttachments = make([]Attachment, 0, len(protoPb.InlineAttachments))
	for _, a := range protoPb.InlineAttachments {
		att := Attachment{
			ID:      a.Id,
			CacheID: a.CacheId,
		}
		if a.ProfileId != nil {
			att.ProfileID = *a.ProfileId
		}
		pk.InlineAttachments = append(pk.InlineAttachments, att)
	}

	return nil
}
