package test

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/omc/bonsai-api-go/v2/bonsai"
	"github.com/omc/terraform-provider-bonsai/internal/provider"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"bonsai": providerserver.NewProtocol6WithError(provider.New("test")()),
}

type ClientTestSuite struct {
	// Assertions embedded here allows all tests to reach through the suite to access assertion methods
	*require.Assertions
	// Suite is the testify/suite used for all HTTP request tests
	suite.Suite

	// client allows each test to have a reachable *bonsai.Client for testing
	client *bonsai.Client
}

// ProviderTestSuite is used for all provider acceptance tests.
type ProviderTestSuite struct {
	ClientTestSuite
}

func (s *ProviderTestSuite) SetupSuite() {
	s.client = bonsai.NewClient(
		bonsai.WithApplication(
			bonsai.Application{
				Name:    "terraform-provider-bonsai",
				Version: "0.1.0-dev",
			},
		),
		bonsai.WithCredentialPair(
			bonsai.CredentialPair{
				AccessKey:   bonsai.AccessKey("TerraformTestKey"),
				AccessToken: bonsai.AccessToken("TerraformTestToken"),
			},
		),
	)
	// configure testify
	s.Assertions = require.New(s.T())
}
