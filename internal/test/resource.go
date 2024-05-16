package test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"os"

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

	// Client allows each test to have a reachable *bonsai.Client for testing
	Client *bonsai.Client

	ProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
}

// ProviderTestSuite is used for all provider acceptance tests.
type ProviderTestSuite struct {
	ClientTestSuite
}

const version = "0.1.0-test"

var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"bonsai": providerserver.NewProtocol6WithError(
		provider.New(
			provider.WithVersion(version),
		)(),
	),
}

func NewApiClient() *bonsai.Client {
	envKey := os.Getenv("BONSAI_API_KEY")
	envToken := os.Getenv("BONSAI_API_TOKEN")

	accessKey, err := bonsai.NewAccessKey(envKey)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid user received: %w", err))
	}

	accessToken, err := bonsai.NewAccessToken(envToken)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid token/password received: %w", err))
	}

	return bonsai.NewClient(
		bonsai.WithApplication(
			bonsai.Application{
				Name:    "terraform-provider-bonsai",
				Version: version,
			},
		),
		bonsai.WithCredentialPair(
			bonsai.CredentialPair{
				AccessKey:   accessKey,
				AccessToken: accessToken,
			},
		),
	)
}

func (s *ProviderTestSuite) SetupSuite() {

	// configure terraform provider factory
	s.ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"bonsai": providerserver.NewProtocol6WithError(
			provider.New(
				provider.WithVersion(version),
			)(),
		),
	}

	s.Client = NewApiClient()

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
	s.Client = bonsai.NewClient(
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
					s.Client,
				),
				provider.WithVersion(version),
			)(),
		),
	}

	// configure testify
	s.Assertions = require.New(s.T())
}
