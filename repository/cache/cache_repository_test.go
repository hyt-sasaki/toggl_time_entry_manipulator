package cache

import (
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/stretchr/testify/assert"
)

type CachedRepositoryTestSuite struct {
    suite.Suite
}

func TestSuite(t *testing.T) {
    suite.Run(t, new(CachedRepositoryTestSuite))
}

func (suite *CachedRepositoryTestSuite) TestA() {
    // given

    // when

    // then
    assert.Equal(suite.T(), 1, 1)
}
