package cluster_test

import (
	"testing"

	"github.com/omc/terraform-provider-bonsai/internal/test"
	"github.com/stretchr/testify/suite"
)

type ClusterTestSuite struct {
	*test.ProviderTestSuite
}

func TestClusterTestSuite(t *testing.T) {
	suite.Run(t, &ClusterTestSuite{ProviderTestSuite: &test.ProviderTestSuite{}})
}

func (s *ClusterTestSuite) SetupSuite() {
	suite.SetupAllSuite(s.ProviderTestSuite).SetupSuite()
}
