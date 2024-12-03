package model

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &IdentityUserFlowId{}

// IdentityUserFlowAttributeId is a struct representing the Resource ID for a Identity User Flow Attribute
type IdentityUserFlowId struct {
	IdentityUserFlowId string
}

// NewIdentityUserFlowAttributeID returns a new IdentityUserFlowAttributeId struct
func NewIdentityUserFlowID(identityUserFlowId string) IdentityUserFlowId {
	return IdentityUserFlowId{
		IdentityUserFlowId: identityUserFlowId,
	}
}

func ParseIdentityUserFlowID(input string) (*IdentityUserFlowId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IdentityUserFlowId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IdentityUserFlowId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func ParseIdentityUserFlowAttributeIDInsensitively(input string) (*IdentityUserFlowId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IdentityUserFlowId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IdentityUserFlowId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *IdentityUserFlowId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.IdentityUserFlowId, ok = input.Parsed["identityUserFlowId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "identityUserFlowId", input)
	}

	return nil
}

// ValidateIdentityUserFlowAttributeID checks that 'input' can be parsed as a Identity User Flow Attribute ID
func ValidateIdentityUserFlowID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIdentityUserFlowID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

func (id IdentityUserFlowId) ID() string {
	fmtString := "/identity/AuthenticationEventsFlows/%s"
	return fmt.Sprintf(fmtString, id.IdentityUserFlowId)
}

func (id IdentityUserFlowId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("identity", "identity", "identity"),
		resourceids.StaticSegment("AuthenticationEventsFlows", "AuthenticationEventsFlows", "AuthenticationEventsFlows"),
		resourceids.UserSpecifiedSegment("identityUserFlowId", "identityUserFlowId"),
	}
}

// String returns a human-readable description of this Identity User Flow Attribute ID
func (id IdentityUserFlowId) String() string {
	components := []string{
		fmt.Sprintf("Identity User Flow Attribute: %q", id.IdentityUserFlowId),
	}
	return fmt.Sprintf("Identity User Flow Attribute (%s)", strings.Join(components, "\n"))
}
