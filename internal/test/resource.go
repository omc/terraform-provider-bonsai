package test

import (
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/omc/bonsai-api-go/v2/bonsai"
	"github.com/omc/terraform-provider-bonsai/internal/provider"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	// Assertions embedded here allows all tests to reach through the suite to access assertion methods
	*require.Assertions
	// Suite is the testify/suite used for all HTTP request tests
	suite.Suite

	// client allows each test to have a reachable *bonsai.Client for testing
	client *bonsai.Client

	ProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
}

// ProviderTestSuite is used for all provider acceptance tests.
type ProviderTestSuite struct {
	ClientTestSuite
}

func (s *ProviderTestSuite) SetupSuite() {
	version := "0.1.0-test"

	// configure terraform provider factory
	s.ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"bonsai": providerserver.NewProtocol6WithError(
			provider.New(
				provider.WithVersion(version),
			)(),
		),
	}

	// configure testify
	s.Assertions = require.New(s.T())
}

// ProviderMockRequestTestSuite is used for tests where the target endpoints
// are mocked; allowing for isolating the terraform provider functionality,
// without requiring live responses from the production Bonsai API.
type ProviderMockRequestTestSuite struct {
	ClientTestSuite
	serveMux *chi.Mux
	server   *httptest.Server
}

func (s *ProviderMockRequestTestSuite) SetupSuite() {
	version := "0.1.0-test"

	// Configure http client and other miscellany
	s.serveMux = chi.NewRouter()
	s.server = httptest.NewServer(s.serveMux)
	s.client = bonsai.NewClient(
		bonsai.WithEndpoint(s.server.URL),
		bonsai.WithApplication(
			bonsai.Application{
				Name:    "terraform-provider-bonsai",
				Version: version,
			},
		),
		bonsai.WithCredentialPair(
			bonsai.CredentialPair{
				AccessKey:   bonsai.AccessKey("TerraformTestKey"),
				AccessToken: bonsai.AccessToken("TerraformTestToken"),
			},
		),
	)

	// configure terraform provider factory
	s.ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"bonsai": providerserver.NewProtocol6WithError(
			provider.New(
				provider.WithAPIClient(
					s.client,
				),
				provider.WithVersion(version),
			)(),
		),
	}

	// configure testify
	s.Assertions = require.New(s.T())
}
