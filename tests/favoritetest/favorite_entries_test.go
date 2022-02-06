package favoritetest

import (
	_ "toggl_time_entry_manipulator/supports"
    "encoding/json"
	"fmt"
    "errors"
	"toggl_time_entry_manipulator/domain"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"toggl_time_entry_manipulator/command/favorite"
	"toggl_time_entry_manipulator/config"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/command"

	"github.com/jason0x43/go-toggl"
)

type FavoriteEntryTestSuite struct {
    suite.Suite
    mockedRepo *repository.MockedCachedRepository
    config *config.Config
    com *favorite.FavoriteEntryCommand
}

func TestFavoriteEntriesTestSuite(t *testing.T) {
    suite.Run(t, new(FavoriteEntryTestSuite))
}

func (suite *FavoriteEntryTestSuite) SetupTest() {
    suite.mockedRepo = &repository.MockedCachedRepository{}
    suite.config = &config.Config{}
    suite.config.WorkflowConfig.Favorites = []int{1, 3, 5}
    suite.com = &favorite.FavoriteEntryCommand{
        Repo: suite.mockedRepo,
        Config: suite.config,
    }

    suite.mockedRepo.On("GetProjects").Return([]toggl.Project{
        {ID: 0, Name: "project1"},
        {ID: 1, Name: "project2"},
        {ID: 2, Name: "project3"},
    }, nil)
}

func (suite *FavoriteEntryTestSuite) TestItems() {
    // given
    for idx, entryId := range suite.config.WorkflowConfig.Favorites {
        suite.mockedRepo.On("FindOneById", entryId).Return(domain.TimeEntryEntity{
            Entry: toggl.TimeEntry{ID: entryId, Pid: idx, Description: fmt.Sprintf("description %d", entryId), Tags: []string{fmt.Sprintf("tag %d", entryId)}},
            Estimation: domain.Estimation{Duration: 10, Memo: fmt.Sprintf("Memo %d", entryId)},
        }, nil)
    }


    // when
    items, _ := suite.com.Items("", "")

    // then
    t := suite.T()
    assert.Equal(t, 3, len(items))
    assert.Equal(t, "description 1 (project1) [tag 1]", items[0].Title)
    assert.Equal(t, "description 3 (project2) [tag 3]", items[1].Title)
    assert.Equal(t, "description 5 (project3) [tag 5]", items[2].Title)
}

func (suite *FavoriteEntryTestSuite) TestItemsWithArg() {
    // given
    for idx, entryId := range suite.config.WorkflowConfig.Favorites {
        suite.mockedRepo.On("FindOneById", entryId).Return(domain.TimeEntryEntity{
            Entry: toggl.TimeEntry{ID: entryId, Pid: idx, Description: fmt.Sprintf("description %d", entryId), Tags: []string{fmt.Sprintf("tag %d", entryId)}},
            Estimation: domain.Estimation{Duration: 10, Memo: fmt.Sprintf("Memo %d", entryId)},
        }, nil)
    }
    arg := "3"


    // when
    items, _ := suite.com.Items(arg, "")

    // then
    t := suite.T()
    assert.Equal(t, 2, len(items))
    assert.Equal(t, "description 3 (project2) [tag 3]", items[0].Title)
    assert.Equal(t, "description 5 (project3) [tag 5]", items[1].Title)
}

func (suite *FavoriteEntryTestSuite) TestItemsWhenSomeEntriesAreNotRegistered() {
    // given
    suite.mockedRepo.On("FindOneById", 1).Return(domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 1, Pid: 0, Description: "description 1", Tags: []string{"tag 1"}},
        Estimation: domain.Estimation{Duration: 10, Memo: "Memo 1"},
    }, nil)
    suite.mockedRepo.On("FindOneById", 3).Return(domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 3, Pid: 1, Description: "description 3", Tags: []string{"tag 3"}},
        Estimation: domain.Estimation{Duration: 10, Memo: "Memo 3"},
    }, nil)
    suite.mockedRepo.On("FindOneById", 5).Return(domain.TimeEntryEntity{}, errors.New("Not found"))


    // when
    items, _ := suite.com.Items("", "")

    // then
    t := suite.T()
    assert.Equal(t, 2, len(items))
    assert.Equal(t, "description 1 (project1) [tag 1]", items[0].Title)
    assert.Equal(t, "description 3 (project2) [tag 3]", items[1].Title)
}

func (suite *FavoriteEntryTestSuite) TestAddToFavorite() {
    // given
    dataStr := convertFavoriteData(command.FavoriteRefData{
        Ref: command.DetailRefData{ID: 42},
        Action: command.AddToFavorite,
    })

    // when
    out, _ := suite.com.Do(dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Entry has been added to favorite list.", out)
    assert.Contains(t, suite.config.WorkflowConfig.Favorites, 42)
}

func (suite *FavoriteEntryTestSuite) TestRemoveFromFavorite() {
    // given
    t := suite.T()
    dataStr := convertFavoriteData(command.FavoriteRefData{
        Ref: command.DetailRefData{ID: 3},
        Action: command.RemoveFromFavorite,
    })
    assert.Contains(t, suite.config.WorkflowConfig.Favorites, 3)

    // when
    out, _ := suite.com.Do(dataStr)

    // then
    assert.Equal(t, "Entry has been removed from favorite list.", out)
    assert.NotContains(t, suite.config.WorkflowConfig.Favorites, 3)
}

func (suite *FavoriteEntryTestSuite) TestRemoveWrongEntryFromFavorite() {
    // given
    t := suite.T()
    dataStr := convertFavoriteData(command.FavoriteRefData{
        Ref: command.DetailRefData{ID: 42},
        Action: command.RemoveFromFavorite,
    })
    assert.NotContains(t, suite.config.WorkflowConfig.Favorites, 42)
    assert.Equal(t, 3, len(suite.config.WorkflowConfig.Favorites))

    // when
    out, _ := suite.com.Do(dataStr)

    // then
    assert.Equal(t, "Not found 42 in favorite list.", out)
    assert.Equal(t, 3, len(suite.config.WorkflowConfig.Favorites))
    assert.NotContains(t, suite.config.WorkflowConfig.Favorites, 42)
}


func convertFavoriteData(data command.FavoriteRefData) string {
    dataBytes, _ := json.Marshal(data)
    return string(dataBytes)
}
