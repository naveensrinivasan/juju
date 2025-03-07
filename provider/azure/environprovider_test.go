// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package azure_test

import (
	stdcontext "context"
	"net/http"
	"time"

	"github.com/juju/clock/testclock"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/cloud"
	"github.com/juju/juju/environs"
	environscloudspec "github.com/juju/juju/environs/cloudspec"
	"github.com/juju/juju/provider/azure"
	"github.com/juju/juju/provider/azure/internal/azureauth"
	"github.com/juju/juju/provider/azure/internal/azurecli"
	"github.com/juju/juju/provider/azure/internal/azuretesting"
	"github.com/juju/juju/testing"
)

type environProviderSuite struct {
	testing.BaseSuite
	provider environs.EnvironProvider
	spec     environscloudspec.CloudSpec
	requests []*http.Request
	sender   azuretesting.Senders
}

var _ = gc.Suite(&environProviderSuite{})

func (s *environProviderSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)
	s.provider = newProvider(c, azure.ProviderConfig{
		Sender:                     &s.sender,
		RequestInspector:           azuretesting.RequestRecorder(&s.requests),
		RandomWindowsAdminPassword: func() string { return "sorandom" },
	})
	s.spec = environscloudspec.CloudSpec{
		Type:             "azure",
		Name:             "azure",
		Region:           "westus",
		Endpoint:         "https://api.azurestack.local",
		IdentityEndpoint: "https://login.azurestack.local",
		StorageEndpoint:  "https://storage.azurestack.local",
		Credential:       fakeServicePrincipalCredential(),
	}
	s.sender = nil
}

func fakeServicePrincipalCredential() *cloud.Credential {
	cred := cloud.NewCredential(
		"service-principal-secret",
		map[string]string{
			"application-id":       fakeApplicationId,
			"subscription-id":      fakeSubscriptionId,
			"application-password": "opensezme",
		},
	)
	return &cred
}

func (s *environProviderSuite) TestPrepareConfig(c *gc.C) {
	cfg := makeTestModelConfig(c)
	cfg, err := s.provider.PrepareConfig(environs.PrepareConfigParams{
		Cloud:  s.spec,
		Config: cfg,
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Check(cfg, gc.NotNil)
}

func (s *environProviderSuite) TestOpen(c *gc.C) {
	s.sender = azuretesting.Senders{
		makeResourceGroupNotFoundSender(".*/resourcegroups/juju-testmodel-model-deadbeef-.*"),
		makeSender(".*/resourcegroups/juju-testmodel-.*", makeResourceGroupResult()),
	}
	env, err := environs.Open(stdcontext.TODO(), s.provider, environs.OpenParams{
		Cloud:  s.spec,
		Config: makeTestModelConfig(c),
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(env, gc.NotNil)
}

func (s *environProviderSuite) TestOpenMissingCredential(c *gc.C) {
	s.spec.Credential = nil
	s.testOpenError(c, s.spec, `validating cloud spec: missing credential not valid`)
}

func (s *environProviderSuite) TestOpenUnsupportedCredential(c *gc.C) {
	credential := cloud.NewCredential(cloud.OAuth1AuthType, map[string]string{})
	s.spec.Credential = &credential
	s.testOpenError(c, s.spec, `validating cloud spec: "oauth1" auth-type not supported`)
}

func (s *environProviderSuite) testOpenError(c *gc.C, spec environscloudspec.CloudSpec, expect string) {
	s.sender = azuretesting.Senders{
		makeResourceGroupNotFoundSender(".*/resourcegroups/juju-testmodel-model-deadbeef-.*"),
		makeSender(".*/resourcegroups/juju-testmodel-.*", makeResourceGroupResult()),
	}
	_, err := environs.Open(stdcontext.TODO(), s.provider, environs.OpenParams{
		Cloud:  spec,
		Config: makeTestModelConfig(c),
	})
	c.Assert(err, gc.ErrorMatches, expect)
}

func newProvider(c *gc.C, config azure.ProviderConfig) environs.EnvironProvider {
	if config.NewStorageClient == nil {
		var storage azuretesting.MockStorageClient
		config.NewStorageClient = storage.NewClient
	}
	if config.RetryClock == nil {
		config.RetryClock = testclock.NewClock(time.Time{})
	}
	if config.ServicePrincipalCreator == nil {
		config.ServicePrincipalCreator = &azureauth.ServicePrincipalCreator{}
	}
	if config.AzureCLI == nil {
		config.AzureCLI = azurecli.AzureCLI{}
	}
	config.RandomWindowsAdminPassword = func() string { return "sorandom" }
	config.GenerateSSHKey = func(string) (string, string, error) {
		return "private", "public", nil
	}
	environProvider, err := azure.NewProvider(config)
	c.Assert(err, jc.ErrorIsNil)
	return environProvider
}
