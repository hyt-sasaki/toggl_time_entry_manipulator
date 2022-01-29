package command_test

import (
	_ "toggl_time_entry_manipulator/supports"

    "testing"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"

	"toggl_time_entry_manipulator/command"
)


type MatchTestSuite struct {
    suite.Suite
}

func TestSuite(t *testing.T) {
    suite.Run(t, new(MatchTestSuite))
}

func (suite *MatchTestSuite) TestMatches() {
    t := suite.T()
    assert.True(t, command.Match("[PBI]アーカイブされないプロジェクト", ""))
    assert.True(t, command.Match("[PBI]アーカイブされないプロジェクト", "PBI アーカイブ"))
    assert.True(t, command.Match("[PBI]アーカイブされないプロジェクト", "PBI アーカイブ プロジェクト"))
    assert.True(t, command.Match("[PBI]アーカイブされないプロジェクト", "[PBI]アーカイブされないプロジェクト"))
    assert.False(t, command.Match("[PBI]アーカイブされないプロジェクト", "PBI hoge"))
    assert.False(t, command.Match("[PBI]アーカイブされないプロジェクト", "hoge"))
}
