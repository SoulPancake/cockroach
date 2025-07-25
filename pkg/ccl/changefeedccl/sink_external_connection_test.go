// Copyright 2024 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package changefeedccl

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/cloud/externalconn"
	"github.com/cockroachdb/cockroach/pkg/security/username"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfra"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

// TestValidateExternalConnectionSinkURICallsDial tests that the validation
// function properly calls Dial() on the created sink.
func TestValidateExternalConnectionSinkURICallsDial(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	// Test with a fake URI that would normally succeed in creation but fail on dial
	ctx := context.Background()
	
	// Create a minimal server config for testing
	serverCfg := &execinfra.ServerConfig{}
	
	env := externalconn.ExternalConnEnv{
		Username:  username.RootUserName(),
		ServerCfg: serverCfg,
	}

	// Test with a kafka URI that will fail to dial
	kafkaURI := "kafka://nonexistent-host:9092?topic_name=test"
	
	err := validateExternalConnectionSinkURI(ctx, env, kafkaURI)
	
	// We expect this to fail with a dial error, not a creation error
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to dial changefeed sink")
}