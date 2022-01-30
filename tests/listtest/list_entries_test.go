package listtest

import (
	_ "toggl_time_entry_manipulator/supports"

	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"

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
    com *list.ListEntryCommand
}

func TestListEntryTestSuite(t *testing.T) {
    suite.Run(t, new(ListEntryTestSuite))
}

func (suite *ListEntryTestSuite) SetupTest() {
    suite.mockedRepo = &repository.MockedCachedRepository{}
    suite.com = &list.ListEntryCommand{
        Repo: suite.mockedRepo,
    }
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
