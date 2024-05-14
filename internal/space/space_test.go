package space_test

import (
	"testing"

	"github.com/omc/terraform-provider-bonsai/internal/test"
	"github.com/stretchr/testify/suite"
)

type SpaceTestSuite struct {
	*test.ProviderTestSuite
}

func TestSpaceTestSuite(t *testing.T) {
	suite.Run(t, &SpaceTestSuite{ProviderTestSuite: &test.ProviderTestSuite{}})
}

func (s *SpaceTestSuite) SetupSuite() {
	suite.SetupAllSuite(s.ProviderTestSuite).SetupSuite()
}
