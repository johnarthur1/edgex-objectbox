//
// Copyright (c) 2018
// Cavium
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"testing"

	"github.com/objectbox/edgex-objectbox/internal/core/command"
	"github.com/objectbox/edgex-objectbox/internal/pkg/config"
)

func TestToml(t *testing.T) {
	configuration := &command.ConfigurationStruct{}
	if err := config.VerifyTomlFiles(configuration); err != nil {
		t.Fatalf("%v", err)
	}
	if configuration.Service.StartupMsg == "" {
		t.Errorf("configuration.StartupMsg is zero length.")
	}
}
