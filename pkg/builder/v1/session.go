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

// ChangeProposal now represents a document in the Global Diff Registry
type ChangeProposal struct {
	ID         string         `json:"id" firestore:"-"`                // Stored as Firestore Doc ID
	SessionID  string         `json:"sessionId" firestore:"sessionId"` // NEW: Ties proposal to the chat session
	FilePath   string         `json:"filePath" firestore:"filePath"`
	NewContent string         `json:"newContent" firestore:"newContent"`
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
	return &builderv1.ChangeProposalPb{
		Id:         native.ID,
		SessionId:  native.SessionID,
		FilePath:   native.FilePath,
		NewContent: native.NewContent,
		Reasoning:  native.Reasoning,
		Status:     string(native.Status),
		CreatedAt:  native.CreatedAt.Format(time.RFC3339),
	}
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
	p.NewContent = pb.NewContent
	p.Reasoning = pb.Reasoning
	p.Status = ProposalStatus(pb.Status)
	p.CreatedAt, _ = time.Parse(time.RFC3339, pb.CreatedAt)

	return nil
}
