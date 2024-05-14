package release_test

import (
	"testing"

	"github.com/omc/terraform-provider-bonsai/internal/test"
	"github.com/stretchr/testify/suite"
)

type ReleaseTestSuite struct {
	*test.ProviderTestSuite
}

func TestReleaseTestSuite(t *testing.T) {
	suite.Run(t, &ReleaseTestSuite{ProviderTestSuite: &test.ProviderTestSuite{}})
}

func (s *ReleaseTestSuite) SetupSuite() {
	suite.SetupAllSuite(s.ProviderTestSuite).SetupSuite()
}
