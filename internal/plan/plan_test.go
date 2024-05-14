package plan_test

import (
	"testing"

	"github.com/omc/terraform-provider-bonsai/internal/test"
	"github.com/stretchr/testify/suite"
)

type PlanTestSuite struct {
	*test.ProviderTestSuite
}

func TestPlanTestSuite(t *testing.T) {
	suite.Run(t, &PlanTestSuite{ProviderTestSuite: &test.ProviderTestSuite{}})
}

func (s *PlanTestSuite) SetupSuite() {
	suite.SetupAllSuite(s.ProviderTestSuite).SetupSuite()
}
