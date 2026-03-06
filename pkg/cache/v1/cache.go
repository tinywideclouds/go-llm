package cache

import (
	cachev1 "github.com/tinywideclouds/gen-llm/go/types/cache/v1"
	"github.com/tinywideclouds/go-llm/pkg/yaml/filter"
)

// --- API Request Types ---

type CreateCacheRequest struct {
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
}

func (req CreateCacheRequest) MarshalJSON() ([]byte, error) {
	pb := &cachev1.CreateCacheRequestPb{
		Repo:   req.Repo,
		Branch: req.Branch,
	}
	return protojsonMarshalOptions.Marshal(pb)
}

func (req *CreateCacheRequest) UnmarshalJSON(data []byte) error {
	var pb cachev1.CreateCacheRequestPb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &pb); err != nil {
		return err
	}
	req.Repo = pb.Repo
	req.Branch = pb.Branch
	return nil
}

type SyncRequest struct {
	IngestionRules filter.FilterRules `json:"ingestionRules"`
}

func (req SyncRequest) MarshalJSON() ([]byte, error) {
	pb := &cachev1.SyncRequestPb{
		IngestionRules: &cachev1.FilterRulesPb{
			Include: req.IngestionRules.Include,
			Exclude: req.IngestionRules.Exclude,
		},
	}
	return protojsonMarshalOptions.Marshal(pb)
}

func (req *SyncRequest) UnmarshalJSON(data []byte) error {
	var pb cachev1.SyncRequestPb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &pb); err != nil {
		return err
	}
	if pb.IngestionRules != nil {
		req.IngestionRules = filter.FilterRules{
			Include: pb.IngestionRules.Include,
			Exclude: pb.IngestionRules.Exclude,
		}
	}
	return nil
}

type ProfileRequest struct {
	Name      string `json:"name"`
	RulesYaml string `json:"rulesYaml"`
}

func (req ProfileRequest) MarshalJSON() ([]byte, error) {
	pb := &cachev1.ProfileRequestPb{
		Name:      req.Name,
		RulesYaml: req.RulesYaml,
	}
	return protojsonMarshalOptions.Marshal(pb)
}

func (req *ProfileRequest) UnmarshalJSON(data []byte) error {
	var pb cachev1.ProfileRequestPb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &pb); err != nil {
		return err
	}
	req.Name = pb.Name
	req.RulesYaml = pb.RulesYaml
	return nil
}
