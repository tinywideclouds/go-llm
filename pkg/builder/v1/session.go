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

type ChangeProposal struct {
	ID         string         `json:"id" firestore:"id"`
	FilePath   string         `json:"filePath" firestore:"filePath"`
	NewContent string         `json:"newContent" firestore:"newContent"`
	Reasoning  string         `json:"reasoning" firestore:"reasoning"`
	Status     ProposalStatus `json:"status" firestore:"status"`
	CreatedAt  time.Time      `json:"createdAt" firestore:"createdAt"`
}

type Session struct {
	ID               string                    `json:"id" firestore:"-"`
	CompiledCacheID  string                    `json:"compiledCacheId" firestore:"compiledCacheId"`
	AcceptedOverlays map[string]FileState      `json:"acceptedOverlays" firestore:"acceptedOverlays"`
	PendingProposals map[string]ChangeProposal `json:"pendingProposals" firestore:"pendingProposals"`
	UpdatedAt        time.Time                 `json:"updatedAt" firestore:"updatedAt"`
}

// --- Protobuf Converters ---

func SessionToProto(native *Session) *builderv1.SessionPb {
	if native == nil {
		return nil
	}

	overlays := make(map[string]*builderv1.FileStatePb)
	for k, v := range native.AcceptedOverlays {
		overlays[k] = &builderv1.FileStatePb{
			Content:   v.Content,
			IsDeleted: v.IsDeleted,
		}
	}

	proposals := make(map[string]*builderv1.ChangeProposalPb)
	for k, v := range native.PendingProposals {
		proposals[k] = &builderv1.ChangeProposalPb{
			Id:         v.ID,
			FilePath:   v.FilePath,
			NewContent: v.NewContent,
			Reasoning:  v.Reasoning,
			Status:     string(v.Status),
			CreatedAt:  v.CreatedAt.Format(time.RFC3339),
		}
	}

	return &builderv1.SessionPb{
		Id:               native.ID,
		CompiledCacheId:  native.CompiledCacheID,
		AcceptedOverlays: overlays,
		PendingProposals: proposals,
		UpdatedAt:        native.UpdatedAt.Format(time.RFC3339),
	}
}

// MarshalJSON uses protojson to ensure exact gRPC/Protobuf JSON mapping
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

	s.AcceptedOverlays = make(map[string]FileState)
	for k, v := range pb.AcceptedOverlays {
		s.AcceptedOverlays[k] = FileState{
			Content:   v.Content,
			IsDeleted: v.IsDeleted,
		}
	}

	s.PendingProposals = make(map[string]ChangeProposal)
	for k, v := range pb.PendingProposals {
		createdAt, _ := time.Parse(time.RFC3339, v.CreatedAt)
		s.PendingProposals[k] = ChangeProposal{
			ID:         v.Id,
			FilePath:   v.FilePath,
			NewContent: v.NewContent,
			Reasoning:  v.Reasoning,
			Status:     ProposalStatus(v.Status),
			CreatedAt:  createdAt,
		}
	}

	return nil
}
