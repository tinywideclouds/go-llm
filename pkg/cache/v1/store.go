package cache

import (
	cachev1 "github.com/tinywideclouds/gen-llm/go/types/cache/v1"
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

type StoreCollections struct {
	BundleCollection   string
	FilesCollection    string
	ProfilesCollection string
}

// ToProto converts the idiomatic Go struct into its Protobuf representation.
func ToProto(native *StoreCollections) *cachev1.StoreCollectionsPb {
	if native == nil {
		return nil
	}
	return &cachev1.StoreCollectionsPb{
		BundleCollection:   native.BundleCollection,
		FilesCollection:    native.FilesCollection,
		ProfilesCollection: native.ProfilesCollection,
	}
}

// FromProto converts the Protobuf representation into the idiomatic Go struct.
func FromProto(proto *cachev1.StoreCollectionsPb) (*StoreCollections, error) {
	if proto == nil {
		return nil, nil
	}
	return &StoreCollections{
		BundleCollection:   proto.BundleCollection,
		FilesCollection:    proto.FilesCollection,
		ProfilesCollection: proto.ProfilesCollection,
	}, nil
}

// --- JSON METHODS ---

// MarshalJSON implements the json.Marshaler interface.
func (pk StoreCollections) MarshalJSON() ([]byte, error) {
	protoPb := ToProto(&pk)
	return protojsonMarshalOptions.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (pk *StoreCollections) UnmarshalJSON(data []byte) error {
	var protoPb cachev1.StoreCollectionsPb

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
		*pk = StoreCollections{}
	}
	return nil
}
