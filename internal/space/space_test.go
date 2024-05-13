package space_test

import (
	"testing"

	"github.com/omc/terraform-provider-bonsai/internal/test"
	"github.com/stretchr/testify/suite"
)

type SpaceTestSuite struct {
	test.ProviderTestSuite
}

func TestSpaceTestSuite(t *testing.T) {
	suite.Run(t, &SpaceTestSuite{})
}

func (s *SpaceTestSuite) SetupTest() {
	suite.SetupAllSuite(s).SetupSuite()
}
