package builder

import (
	"time"

	builderv1 "github.com/tinywideclouds/gen-llm/go/types/builder/v1"
)

type ProposalStatus string

const (
	StatusPending  ProposalStatus = "pending"
	StatusAccepted ProposalStatus = "accepted"
	StatusRejected ProposalStatus = "rejected"
)

type FileState struct {
	Content   string `json:"content" firestore:"content"`
	IsDeleted bool   `json:"isDeleted" firestore:"isDeleted"`
}

// ChangeProposal represents a document in the Global Diff Registry
type ChangeProposal struct {
	ID         string         `json:"id" firestore:"-"`
	SessionID  string         `json:"sessionId" firestore:"sessionId"`
	FilePath   string         `json:"filePath" firestore:"filePath"`
	Patch      string         `json:"patch,omitempty" firestore:"patch,omitempty"`           // NEW
	NewContent string         `json:"newContent,omitempty" firestore:"newContent,omitempty"` // UPDATED
	Reasoning  string         `json:"reasoning" firestore:"reasoning"`
	Status     ProposalStatus `json:"status" firestore:"status"`
	CreatedAt  time.Time      `json:"createdAt" firestore:"createdAt"`
}

type Session struct {
	ID              string    `json:"id" firestore:"-"`
	CompiledCacheID string    `json:"compiledCacheId" firestore:"compiledCacheId"`
	UpdatedAt       time.Time `json:"updatedAt" firestore:"updatedAt"`

	// ACCEPTED_OVERLAYS AND PENDING_PROPOSALS DELETED
}

// --- Protobuf Converters for Session ---

func SessionToProto(native *Session) *builderv1.SessionPb {
	if native == nil {
		return nil
	}

	return &builderv1.SessionPb{
		Id:              native.ID,
		CompiledCacheId: native.CompiledCacheID,
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

	s.ID = pb.Id
	s.CompiledCacheID = pb.CompiledCacheId
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
		SessionId: native.SessionID,
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

	p.ID = pb.Id
	p.SessionID = pb.SessionId
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
