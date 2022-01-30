package gettest

import (
    _ "toggl_time_entry_manipulator/supports"
    "time"
	"testing"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/command/get"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/tests"
)

type GetEntryTestSuite struct {
    suite.Suite
    mockedRepo *repository.MockedCachedRepository
    com *get.GetEntryCommand
}

func TestGetEntryTestSuite(t *testing.T) {
    suite.Run(t, new(GetEntryTestSuite))
}

func (suite *GetEntryTestSuite) SetupTest() {
    suite.mockedRepo = &repository.MockedCachedRepository{}
    suite.com = &get.GetEntryCommand{
        Repo: suite.mockedRepo,
    }
}

func (suite *GetEntryTestSuite) TestItems() {
    // given
    arg := ""
    dataStr := tests.StringifyDetailRefData(command.DetailRefData{ID: 42})
    loc, _ := time.LoadLocation("Asia/Tokyo")
    timeLayout := "2006-01-02 15:04:05"
    start, _ := time.ParseInLocation(timeLayout, "2022-01-24 13:50:31", loc)
    stop, _ := time.ParseInLocation(timeLayout, "2022-01-24 15:53:01", loc)
    duration := int64(stop.Sub(start).Seconds())
    suite.mockedRepo.On("FindOneById", 42).Return(domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{
            ID: 42,
            Pid: 3,
            Description: "item42",
            Tags: []string{"tag2"},
            Start: &start,
            Stop: &stop,
            Duration: duration,
        },
        Estimation: domain.Estimation{
            Duration: 66,
            Memo: "memo test",
        },
    }, nil).Once()
    suite.mockedRepo.On("GetProjects").Return([]toggl.Project{
        { ID: 1, Name: "project1", }, 
        { ID: 2, Name: "project2", },
        { ID: 3, Name: "project3", },
    }, nil).Once()


    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 10, len(items))
    assert.Equal(t, "Description: item42", items[0].Title)
    assert.Equal(t, "Project: project3", items[1].Title)
    assert.Equal(t, "Tag: [tag2]", items[2].Title)
    assert.Equal(t, "Estimated duration: 66 [min]", items[3].Title)
    assert.Equal(t, "Start: 22/01/24 13:50", items[4].Title)
    assert.Equal(t, "Stop: 22/01/24 15:53", items[5].Title)
    assert.Equal(t, "Memo: memo test", items[6].Title)
    assert.Equal(t, "Delete this entry (Press Cmd+Enter)", items[7].Title)
    assert.Equal(t, "Continue this entry", items[8].Title)
    assert.Equal(t, "Back", items[9].Title)
    normalItemIds := []int{0, 1, 2, 3, 4, 5, 6}
    for _, i := range normalItemIds {
        item := items[i]
        assert.Equal(t, command.ModifyEntryKeyword, item.Arg.Keyword)
        assert.Equal(t, alfred.ModeTell, item.Arg.Mode)
    }
    assert.Nil(t, items[7].Arg)
    assert.Equal(t, command.ContinueEntryKeyword, items[8].Arg.Keyword)
    assert.Equal(t, alfred.ModeTell, items[8].Arg.Mode)
    assert.Equal(t, command.ListEntryKeyword, items[9].Arg.Keyword)
    assert.Equal(t, alfred.ModeTell, items[9].Arg.Mode)
}

func (suite *GetEntryTestSuite) TestItems_whenEntryIsRunning() {
    // given
    arg := ""
    dataStr := tests.StringifyDetailRefData(command.DetailRefData{ID: 42})
    loc, _ := time.LoadLocation("Asia/Tokyo")
    timeLayout := "2006-01-02 15:04:05"
    start, _ := time.ParseInLocation(timeLayout, "2022-01-24 13:50:31", loc)
    suite.mockedRepo.On("FindOneById", 42).Return(domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{
            ID: 42,
            Pid: 3,
            Description: "item42",
            Tags: []string{"tag2"},
            Start: &start,
            Duration: -1,
        },
        Estimation: domain.Estimation{
            Duration: 66,
            Memo: "memo test",
        },
    }, nil).Once()
    suite.mockedRepo.On("GetProjects").Return([]toggl.Project{
        { ID: 1, Name: "project1", }, 
        { ID: 2, Name: "project2", },
        { ID: 3, Name: "project3", },
    }, nil).Once()


    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 9, len(items))
    assert.Equal(t, "Description: item42", items[0].Title)
    assert.Equal(t, "Project: project3", items[1].Title)
    assert.Equal(t, "Tag: [tag2]", items[2].Title)
    assert.Equal(t, "Estimated duration: 66 [min]", items[3].Title)
    assert.Equal(t, "Start: 22/01/24 13:50", items[4].Title)
    assert.Equal(t, "Memo: memo test", items[5].Title)
    assert.Equal(t, "Delete this entry (Press Cmd+Enter)", items[6].Title)
    assert.Equal(t, "Stop this entry", items[7].Title)
    assert.Equal(t, "Back", items[8].Title)
    normalItemIds := []int{0, 1, 2, 3, 4, 5}
    for _, i := range normalItemIds {
        item := items[i]
        assert.Equal(t, command.ModifyEntryKeyword, item.Arg.Keyword)
        assert.Equal(t, alfred.ModeTell, item.Arg.Mode)
    }
    assert.Nil(t, items[6].Arg)
    assert.Equal(t, command.StopEntryKeyword, items[7].Arg.Keyword)
    assert.Equal(t, alfred.ModeDo, items[7].Arg.Mode)
    assert.Equal(t, command.ListEntryKeyword, items[8].Arg.Keyword)
    assert.Equal(t, alfred.ModeTell, items[8].Arg.Mode)
}

func (suite *GetEntryTestSuite) TestItems_whenNoEstimation() {
    // given
    arg := ""
    dataStr := tests.StringifyDetailRefData(command.DetailRefData{ID: 42})
    loc, _ := time.LoadLocation("Asia/Tokyo")
    timeLayout := "2006-01-02 15:04:05"
    start, _ := time.ParseInLocation(timeLayout, "2022-01-24 13:50:31", loc)
    stop, _ := time.ParseInLocation(timeLayout, "2022-01-24 15:53:01", loc)
    duration := int64(stop.Sub(start).Seconds())
    suite.mockedRepo.On("FindOneById", 42).Return(domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{
            ID: 42,
            Pid: 3,
            Description: "item42",
            Tags: []string{"tag2"},
            Start: &start,
            Stop: &stop,
            Duration: duration,
        },

    }, nil).Once()
    suite.mockedRepo.On("GetProjects").Return([]toggl.Project{
        { ID: 1, Name: "project1", }, 
        { ID: 2, Name: "project2", },
        { ID: 3, Name: "project3", },
    }, nil).Once()


    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 8, len(items))
    assert.Equal(t, "Description: item42", items[0].Title)
    assert.Equal(t, "Project: project3", items[1].Title)
    assert.Equal(t, "Tag: [tag2]", items[2].Title)
    assert.Equal(t, "Start: 22/01/24 13:50", items[3].Title)
    assert.Equal(t, "Stop: 22/01/24 15:53", items[4].Title)
    assert.Equal(t, "Delete this entry (Press Cmd+Enter)", items[5].Title)
    assert.Equal(t, "Continue this entry", items[6].Title)
    assert.Equal(t, "Back", items[7].Title)
}
