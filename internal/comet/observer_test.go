package comet

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type JSONObserverSuite struct {
	suite.Suite
}

func TestJSONFilter(t *testing.T) {
	suite.Run(t, new(JSONObserverSuite))
}


