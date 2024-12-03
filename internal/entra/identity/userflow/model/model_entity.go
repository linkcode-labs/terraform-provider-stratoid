package model

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Entity interface {
	Entity() BaseEntityImpl
}

var _ Entity = BaseEntityImpl{}

type BaseEntityImpl struct {
	// The unique identifier for an entity. Read-only.
	Id *string `json:"id,omitempty"`

	// The OData ID of this entity
	ODataId *string `json:"@odata.id,omitempty"`

	// The OData Type of this entity
	ODataType *string `json:"@odata.type,omitempty"`

	// Model Behaviors
	OmitDiscriminatedValue bool `json:"-"`
}

func (s BaseEntityImpl) Entity() BaseEntityImpl {
	return s
}

var _ Entity = RawEntityImpl{}

// RawEntityImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEntityImpl struct {
	entity BaseEntityImpl
	Type   string
	Values map[string]interface{}
}

func (s RawEntityImpl) Entity() BaseEntityImpl {
	return s.entity
}

var _ json.Marshaler = BaseEntityImpl{}

func (s BaseEntityImpl) MarshalJSON() ([]byte, error) {
	type wrapper BaseEntityImpl
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BaseEntityImpl: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BaseEntityImpl: %+v", err)
	}

	delete(decoded, "id")

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BaseEntityImpl: %+v", err)
	}

	return encoded, nil
}
