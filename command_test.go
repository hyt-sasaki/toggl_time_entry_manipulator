package main

import (
	"encoding/json"
	"fmt"
    "time"
	"testing"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/command/get"
	"toggl_time_entry_manipulator/command/stop"
	"toggl_time_entry_manipulator/command/modify"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"
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
    data := command.DetailRefData{ID: 42}
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
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
    assert.Equal(t, 7, len(items))
    assert.Equal(t, "Description: item42", items[0].Title)
    assert.Equal(t, "Project: project3", items[1].Title)
    assert.Equal(t, "Tag: [tag2]", items[2].Title)
    assert.Equal(t, "Estimated duration: 66 [min]", items[3].Title)
    assert.Equal(t, "Start: 22/01/24 13:50", items[4].Title)
    assert.Equal(t, "Stop: 22/01/24 15:53", items[5].Title)
    assert.Equal(t, "Memo: memo test", items[6].Title)
    for _, item := range items {
        assert.Equal(t, command.ModifyEntryKeyword, item.Arg.Keyword)
        assert.Equal(t, alfred.ModeTell, item.Arg.Mode)
    }
}

func (suite *GetEntryTestSuite) TestItems_whenEntryIsRunning() {
    // given
    arg := ""
    data := command.DetailRefData{ID: 42}
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
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
    assert.Equal(t, 7, len(items))
    assert.Equal(t, "Description: item42", items[0].Title)
    assert.Equal(t, "Project: project3", items[1].Title)
    assert.Equal(t, "Tag: [tag2]", items[2].Title)
    assert.Equal(t, "Estimated duration: 66 [min]", items[3].Title)
    assert.Equal(t, "Start: 22/01/24 13:50", items[4].Title)
    assert.Equal(t, "Memo: memo test", items[5].Title)
    assert.Equal(t, "Stop this entry", items[6].Title)
}

func (suite *GetEntryTestSuite) TestItems_whenNoEstimation() {
    // given
    arg := ""
    data := command.DetailRefData{ID: 42}
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
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
    assert.Equal(t, 5, len(items))
    assert.Equal(t, "Description: item42", items[0].Title)
    assert.Equal(t, "Project: project3", items[1].Title)
    assert.Equal(t, "Tag: [tag2]", items[2].Title)
    assert.Equal(t, "Start: 22/01/24 13:50", items[3].Title)
    assert.Equal(t, "Stop: 22/01/24 15:53", items[4].Title)
}

type StopEntryTestSuite struct {
    suite.Suite
    mockedRepo *repository.MockedCachedRepository
    com *stop.StopEntryCommand
}

func TestStopEntryTestSuite(t *testing.T) {
    suite.Run(t, new(StopEntryTestSuite))
}

func (suite *StopEntryTestSuite) SetupTest() {
    suite.mockedRepo = &repository.MockedCachedRepository{}
    suite.com = &stop.StopEntryCommand{
        Repo: suite.mockedRepo,
    }
}

func (suite *StopEntryTestSuite) TestDo() {
    // given
    // given
    data := command.DetailRefData{ID: 42}
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
    start := time.Now().Add(-time.Hour)
    runningEntity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(runningEntity, nil).Once()
    suite.mockedRepo.On("Stop", &runningEntity).Return(nil).Once()

    // when
    out, _ := suite.com.Do(dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Entry has stopped. Description: item42", out)     // TODO WIP
}

type ModifyEntryTestSuite struct {
    suite.Suite
    mockedRepo *repository.MockedCachedRepository
    com *modify.ModifyEntryCommand
}

func TestModifyEntryTestSuite(t *testing.T) {
    suite.Run(t, new(ModifyEntryTestSuite))
}

func (suite *ModifyEntryTestSuite) SetupTest() {
    suite.mockedRepo = &repository.MockedCachedRepository{}
    suite.com = &modify.ModifyEntryCommand{
        Repo: suite.mockedRepo,
    }
}

func (suite *ModifyEntryTestSuite) TestItems_whenTargetIsDescription() {
    // given
    data := command.ModifyData{
        Ref: command.DetailRefData{ID: 42},
        Target: command.ModifyDescription,
    }
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
    arg := "new description"
    start := time.Now().Add(-time.Hour)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(entity, nil).Once()

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, fmt.Sprintf("Description: %s", arg), items[0].Title)
    assert.Equal(t, command.ModifyEntryKeyword, items[0].Arg.Keyword)
    assert.Equal(t, alfred.ModeDo, items[0].Arg.Mode)
}

func (suite *ModifyEntryTestSuite) TestItems_whenTargetIsDuration() {
    // given
    data := command.ModifyData{
        Ref: command.DetailRefData{ID: 42},
        Target: command.ModifyDuration,
    }
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
    arg := "50"
    start := time.Now().Add(-time.Hour)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(entity, nil).Once()

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, fmt.Sprintf("Duration: %s", arg), items[0].Title)
    assert.Equal(t, command.ModifyEntryKeyword, items[0].Arg.Keyword)
    assert.Equal(t, alfred.ModeDo, items[0].Arg.Mode)
}

func (suite *ModifyEntryTestSuite) TestItems_whenTargetIsDurationAndArgIsNotNumber() {
    // given
    data := command.ModifyData{
        Ref: command.DetailRefData{ID: 42},
        Target: command.ModifyDuration,
    }
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
    arg := "test"
    start := time.Now().Add(-time.Hour)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start},
        Estimation: domain.Estimation{Duration: 100},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(entity, nil).Once()

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Duration: 100", items[0].Title)
    assert.Nil(t, items[0].Arg)
}

func (suite *ModifyEntryTestSuite) TestItems_whenTargetIsMemo() {
    // given
    data := command.ModifyData{
        Ref: command.DetailRefData{ID: 42},
        Target: command.ModifyMemo,
    }
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
    arg := "test"
    start := time.Now().Add(-time.Hour)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start},
        Estimation: domain.Estimation{Duration: 100, Memo: "old memo"},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(entity, nil).Once()

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, fmt.Sprintf("Memo: %s", arg), items[0].Title)
    assert.Equal(t, command.ModifyEntryKeyword, items[0].Arg.Keyword)
    assert.Equal(t, alfred.ModeDo, items[0].Arg.Mode)
}

func (suite *ModifyEntryTestSuite) TestItems_whenTargetIsStartAndHMFormat() {
    // given
    data := command.ModifyData{
        Ref: command.DetailRefData{ID: 42},
        Target: command.ModifyStart,
    }
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
    arg := "20:31"
    start, _ := time.ParseInLocation("06-01-02 15:04", "21-01-27 13:45", time.Local)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start},
        Estimation: domain.Estimation{Duration: 100, Memo: "old memo"},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(entity, nil).Once()

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Start: 21/01/27 20:31", items[0].Title)
    assert.Equal(t, "Modify start time (21/01/27 13:45)", items[0].Subtitle)
    assert.Equal(t, command.ModifyEntryKeyword, items[0].Arg.Keyword)
    assert.Equal(t, alfred.ModeDo, items[0].Arg.Mode)
}

func (suite *ModifyEntryTestSuite) TestItems_whenTargetIsStartAndFullFormat() {
    // given
    data := command.ModifyData{
        Ref: command.DetailRefData{ID: 42},
        Target: command.ModifyStart,
    }
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
    arg := "21/01/28 20:31"
    start, _ := time.ParseInLocation("06-01-02 15:04", "21-01-27 13:45", time.Local)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start},
        Estimation: domain.Estimation{Duration: 100, Memo: "old memo"},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(entity, nil).Once()

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Start: 21/01/28 20:31", items[0].Title)
    assert.Equal(t, "Modify start time (21/01/27 13:45)", items[0].Subtitle)
    assert.Equal(t, command.ModifyEntryKeyword, items[0].Arg.Keyword)
    assert.Equal(t, alfred.ModeDo, items[0].Arg.Mode)
}

func (suite *ModifyEntryTestSuite) TestItems_whenTargetIsStartAndInvalidFormat() {
    // given
    data := command.ModifyData{
        Ref: command.DetailRefData{ID: 42},
        Target: command.ModifyStart,
    }
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
    arg := "21/01/282031"
    start, _ := time.ParseInLocation("06-01-02 15:04", "21-01-27 13:45", time.Local)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start},
        Estimation: domain.Estimation{Duration: 100, Memo: "old memo"},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(entity, nil).Once()

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Start: -", items[0].Title)
    assert.Equal(t, "Modify start time (21/01/27 13:45)", items[0].Subtitle)
    assert.Nil(t, items[0].Arg)
}

func (suite *ModifyEntryTestSuite) TestItems_whenTargetIsStopAndHMFormat() {
    // given
    data := command.ModifyData{
        Ref: command.DetailRefData{ID: 42},
        Target: command.ModifyStop,
    }
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
    arg := "20:31"
    stop, _ := time.ParseInLocation("06-01-02 15:04", "21-01-27 13:45", time.Local)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Stop: &stop},
        Estimation: domain.Estimation{Duration: 100, Memo: "old memo"},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(entity, nil).Once()

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Stop: 21/01/27 20:31", items[0].Title)
    assert.Equal(t, "Modify stop time (21/01/27 13:45)", items[0].Subtitle)
    assert.Equal(t, command.ModifyEntryKeyword, items[0].Arg.Keyword)
    assert.Equal(t, alfred.ModeDo, items[0].Arg.Mode)
}

func (suite *ModifyEntryTestSuite) TestItems_whenTargetIsStopAndFullFormat() {
    // given
    data := command.ModifyData{
        Ref: command.DetailRefData{ID: 42},
        Target: command.ModifyStop,
    }
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
    arg := "21/01/28 20:31"
    stop, _ := time.ParseInLocation("06-01-02 15:04", "21-01-27 13:45", time.Local)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Stop: &stop},
        Estimation: domain.Estimation{Duration: 100, Memo: "old memo"},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(entity, nil).Once()

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Stop: 21/01/28 20:31", items[0].Title)
    assert.Equal(t, "Modify stop time (21/01/27 13:45)", items[0].Subtitle)
    assert.Equal(t, command.ModifyEntryKeyword, items[0].Arg.Keyword)
    assert.Equal(t, alfred.ModeDo, items[0].Arg.Mode)
}

func (suite *ModifyEntryTestSuite) TestItems_whenTargetIsStopAndInvalidFormat() {
    // given
    data := command.ModifyData{
        Ref: command.DetailRefData{ID: 42},
        Target: command.ModifyStop,
    }
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)
    arg := "21/01/282031"
    stop, _ := time.ParseInLocation("06-01-02 15:04", "21-01-27 13:45", time.Local)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Stop: &stop},
        Estimation: domain.Estimation{Duration: 100, Memo: "old memo"},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(entity, nil).Once()

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Stop: -", items[0].Title)
    assert.Equal(t, "Modify stop time (21/01/27 13:45)", items[0].Subtitle)
    assert.Nil(t, items[0].Arg)
}

func (suite *ModifyEntryTestSuite) TestDo() {
    // given
    start := time.Now().Add(-time.Hour)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start},
        Estimation: domain.Estimation{Duration: 100, Memo: "memo"},
    }
    suite.mockedRepo.On("Update", mock.Anything).Return(nil).Once()
    dataBytes, _ := json.Marshal(entity)
    dataStr := string(dataBytes)

    // when
    out, _ := suite.com.Do(dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Time entry has been updated successfully", out)
}
