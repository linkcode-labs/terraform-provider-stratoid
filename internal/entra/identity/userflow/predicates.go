package userflow

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"terraform-provider-stratoid/internal/entra/identity/userflow/model"
)

type IdentityUserFlowOperationPredicate struct {
}

func (p IdentityUserFlowOperationPredicate) Matches(input model.IdentityUserFlow) bool {

	return true
}
