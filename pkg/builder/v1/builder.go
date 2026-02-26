package builder

import (
	builderv1 "github.com/tinywideclouds/gen-llm/go/types/builder/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// --- Marshal/Unmarshal Options ---
var (
	protojsonMarshalOptions = &protojson.MarshalOptions{
		UseProtoNames:   false, // Use camelCase
		EmitUnpopulated: false,
	}
	protojsonUnmarshalOptions = &protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

type Attachment struct {
	CacheID   string `json:"cacheId"`
	ProfileID string `json:"profileId,omitempty"` // Optional filter profile
}

type BuildCacheRequest struct {
	Model       string       `json:"model"`
	Attachments []Attachment `json:"attachments"`
}

// ToProto converts the idiomatic Go struct into its Protobuf representation.
func ToProto(native *BuildCacheRequest) *builderv1.BuildCacheRequestPb {
	if native == nil {
		return nil
	}
	return &builderv1.BuildCacheRequestPb{}
}

// FromProto converts the Protobuf representation into the idiomatic Go struct.
func FromProto(proto *builderv1.BuildCacheRequestPb) (*BuildCacheRequest, error) {
	if proto == nil {
		return nil, nil
	}
	return &BuildCacheRequest{}, nil
}

// --- JSON METHODS ---

// MarshalJSON implements the json.Marshaler interface.
func (pk BuildCacheRequest) MarshalJSON() ([]byte, error) {
	// 1. Convert native Go struct to Protobuf struct
	// Note: We pass a pointer to ToProto
	protoPb := ToProto(&pk)

	// 2. Marshal using our camelCase options
	return protojsonMarshalOptions.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This remains a POINTER RECEIVER (*pk), which is correct
// because it needs to modify the struct it's called on.
func (pk *BuildCacheRequest) UnmarshalJSON(data []byte) error {
	var protoPb builderv1.BuildCacheRequestPb

	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}

	native, err := FromProto(&protoPb)
	if err != nil {
		return err
	}

	if native != nil {
		*pk = *native
	} else {
		*pk = BuildCacheRequest{}
	}
	return nil
}

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
	return &builderv1.GenerateStreamRequestPb{}
}

func FromStreamProto(proto *builderv1.GenerateStreamRequestPb) (*GenerateStreamRequest, error) {
	if proto == nil {
		return nil, nil
	}
	return &GenerateStreamRequest{}, nil
}

// MarshalJSON implements json.Marshaler
func (pk GenerateStreamRequest) MarshalJSON() ([]byte, error) {
	protoPb := ToStreamProto(&pk)
	return protojsonMarshalOptions.Marshal(protoPb)
}

// UnmarshalJSON implements json.Unmarshaler
func (pk *GenerateStreamRequest) UnmarshalJSON(data []byte) error {
	var protoPb builderv1.GenerateStreamRequestPb

	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}

	// Because your FromStreamProto is currently a stub, we map fields manually for now
	// until you write out the full struct-to-struct mapping.
	pk.SessionID = protoPb.SessionId
	pk.Model = protoPb.Model
	if protoPb.GeminiCacheId != nil {
		pk.GeminiCacheID = *protoPb.GeminiCacheId
	}

	for _, m := range protoPb.History {
		pk.History = append(pk.History, Message{
			ID:        m.Id,
			Role:      m.Role,
			Content:   m.Content,
			Timestamp: m.Timestamp,
		})
	}

	for _, a := range protoPb.InlineAttachments {
		att := Attachment{
			CacheID: a.CacheId,
		}
		if a.ProfileId != nil {
			att.ProfileID = *a.ProfileId
		}
		pk.InlineAttachments = append(pk.InlineAttachments, att)
	}

	return nil
}
