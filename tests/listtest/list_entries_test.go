package listtest

import (
	"toggl_time_entry_manipulator/config"
	_ "toggl_time_entry_manipulator/supports"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"encoding/json"

	"github.com/jason0x43/go-toggl"

	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/command/list"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"
)

type ListEntryTestSuite struct {
    suite.Suite
    mockedRepo *repository.MockedCachedRepository
    config *config.WorkflowConfig
    com list.IListEntryCommand
}

func TestListEntryTestSuite(t *testing.T) {
    suite.Run(t, new(ListEntryTestSuite))
}

func (suite *ListEntryTestSuite) SetupTest() {
    suite.mockedRepo = &repository.MockedCachedRepository{}
    suite.config = &config.WorkflowConfig{
        ProjectAliases: []config.AliasMap{
            {ID: 4, Alias: "hoge"},
        },
    }
    suite.com, _ = list.NewListEntryCommand(suite.mockedRepo, suite.config)
}

func (suite *ListEntryTestSuite) TestItems() {
    // given
    arg := "2"
    dataStr := ""
    suite.mockedRepo.On("Fetch").Return([]domain.TimeEntryEntity{
        {Entry: toggl.TimeEntry{ID: 1, Description: "item1"}},
        {Entry: toggl.TimeEntry{ID: 2, Description: "item2-1"}},
        {Entry: toggl.TimeEntry{ID: 3, Description: "item2-2"}},
    }, nil).Once()
    suite.mockedRepo.On("GetProjects").Return([]toggl.Project{}, nil)

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 2, len(items))
    item := items[0]
    assert.Equal(t, "item2-1 (-)", item.Title)
    var itemData command.DetailRefData
    err := json.Unmarshal([]byte(item.Arg.Data), &itemData)
    assert.Nil(t, err)
    assert.Equal(t, 2, itemData.ID)
}

func (suite *ListEntryTestSuite) TestItems_withAlias() {
    // given
    arg := "hoge"
    dataStr := ""
    suite.mockedRepo.On("Fetch").Return([]domain.TimeEntryEntity{
        {Entry: toggl.TimeEntry{ID: 1, Pid: 4, Description: "item1"}},
        {Entry: toggl.TimeEntry{ID: 2, Pid: 5, Description: "item2-1"}},
        {Entry: toggl.TimeEntry{ID: 3, Pid: 6, Description: "item2-2"}},
    }, nil).Once()
    suite.mockedRepo.On("GetProjects").Return([]toggl.Project{
        {ID: 4, Name: "project4"},
        {ID: 5, Name: "project5"},
        {ID: 6, Name: "project6"},
    }, nil)

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 1, len(items))
    item := items[0]
    assert.Equal(t, "item1 (project4)", item.Title)
    var itemData command.DetailRefData
    err := json.Unmarshal([]byte(item.Arg.Data), &itemData)
    assert.Nil(t, err)
    assert.Equal(t, 1, itemData.ID)
}
