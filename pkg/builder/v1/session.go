package builder

import (
	"time"

	builderv1 "github.com/tinywideclouds/gen-llm/go/types/builder/v1"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"
)

type FileState struct {
	Content   string `json:"content" firestore:"content"`
	IsDeleted bool   `json:"isDeleted" firestore:"isDeleted"`
}

// ChangeProposal represents a document in the Global Diff Registry
type ChangeProposal struct {
	ID         string    `json:"id" firestore:"-"`
	SessionID  urn.URN   `json:"sessionId" firestore:"sessionId"`
	FilePath   string    `json:"filePath" firestore:"filePath"`
	Patch      string    `json:"patch,omitempty" firestore:"patch,omitempty"`
	NewContent string    `json:"newContent,omitempty" firestore:"newContent,omitempty"`
	Reasoning  string    `json:"reasoning" firestore:"reasoning"`
	CreatedAt  time.Time `json:"createdAt" firestore:"createdAt"`
}

type Session struct {
	ID              urn.URN   `json:"id" firestore:"-"`
	CompiledCacheID urn.URN   `json:"compiledCacheId" firestore:"compiledCacheId"`
	UpdatedAt       time.Time `json:"updatedAt" firestore:"updatedAt"`

	// ACCEPTED_OVERLAYS AND PENDING_PROPOSALS DELETED
}

// --- Protobuf Converters for Session ---

func SessionToProto(native *Session) *builderv1.SessionPb {
	if native == nil {
		return nil
	}

	return &builderv1.SessionPb{
		Id:              native.ID.String(),
		CompiledCacheId: native.CompiledCacheID.String(),
		UpdatedAt:       native.UpdatedAt.Format(time.RFC3339),
	}
}

func (s Session) MarshalJSON() ([]byte, error) {
	return protojsonMarshalOptions.Marshal(SessionToProto(&s))
}

func (s *Session) UnmarshalJSON(data []byte) error {
	var pb builderv1.SessionPb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &pb); err != nil {
		return err
	}

	id, err := urn.Parse(pb.Id)
	c, err := urn.Parse(pb.CompiledCacheId)

	if err != nil {
		return err
	}
	s.ID = id
	s.CompiledCacheID = c
	s.UpdatedAt, _ = time.Parse(time.RFC3339, pb.UpdatedAt)

	return nil
}

// --- Protobuf Converters for ChangeProposal ---

func ChangeProposalToProto(native *ChangeProposal) *builderv1.ChangeProposalPb {
	if native == nil {
		return nil
	}

	pb := &builderv1.ChangeProposalPb{
		Id:        native.ID,
		SessionId: native.SessionID.String(),
		FilePath:  native.FilePath,
		Reasoning: native.Reasoning,
		CreatedAt: native.CreatedAt.Format(time.RFC3339),
	}

	if native.Patch != "" {
		pb.Patch = &native.Patch
	}
	if native.NewContent != "" {
		pb.NewContent = &native.NewContent
	}

	return pb
}

func (p ChangeProposal) MarshalJSON() ([]byte, error) {
	return protojsonMarshalOptions.Marshal(ChangeProposalToProto(&p))
}

func (p *ChangeProposal) UnmarshalJSON(data []byte) error {
	var pb builderv1.ChangeProposalPb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &pb); err != nil {
		return err
	}

	s, err := urn.Parse(pb.SessionId)
	if err != nil {
		return err
	}

	p.ID = pb.Id
	p.SessionID = s
	p.FilePath = pb.FilePath
	p.Reasoning = pb.Reasoning
	p.CreatedAt, _ = time.Parse(time.RFC3339, pb.CreatedAt)

	if pb.Patch != nil {
		p.Patch = *pb.Patch
	}
	if pb.NewContent != nil {
		p.NewContent = *pb.NewContent
	}

	return nil
}
