// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"terraform-provider-stratoid/internal/services/userflows"

	"terraform-provider-stratoid/internal/sdk"
)

//go:generate go run ../tools/generator-services/main.go -path=../../

func SupportedTypedServices() []sdk.TypedServiceRegistration {
	return []sdk.TypedServiceRegistration{}
}

func SupportedUntypedServices() []sdk.UntypedServiceRegistration {
	return []sdk.UntypedServiceRegistration{
		userflows.Registration{},
	}
}
