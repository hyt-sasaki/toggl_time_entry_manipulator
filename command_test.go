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
	"toggl_time_entry_manipulator/command/add"
	"toggl_time_entry_manipulator/command/list"
	"toggl_time_entry_manipulator/command/get"
	"toggl_time_entry_manipulator/command/stop"
	"toggl_time_entry_manipulator/command/modify"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"
)

type AddEntryTestSuite struct {
    suite.Suite
    mockedRepo *repository.MockedCachedRepository
    com *add.AddEntryCommand
}

func TestSuite(t *testing.T) {
    suite.Run(t, new(AddEntryTestSuite))
}

func (suite *AddEntryTestSuite) SetupTest() {
    suite.mockedRepo = &repository.MockedCachedRepository{}
    suite.com = &add.AddEntryCommand{
        Repo: suite.mockedRepo,
    }
}

func (suite *AddEntryTestSuite) TestItems1() {
    // given
    type Data struct {
    }
    data := Data{}
    dataBytes, _ := json.Marshal(data)
    dataStr := string(dataBytes)

    arg := "this is description"

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 1, len(items))
    item := items[0]
    assert.Equal(t, fmt.Sprintf("New description: %s", arg), item.Title)
    itemArg := item.Arg
    assert.Equal(t, alfred.ModeTell, itemArg.Mode)
    assert.Equal(t, fmt.Sprintf("{\"Current\":1,\"Args\":{\"Description\":\"%s\",\"Project\":0,\"Tag\":\"\",\"TimeEstimation\":0}}", arg), itemArg.Data)
}

func (suite *AddEntryTestSuite) TestItems2() {
    // given
    dataStr := `{"Current":1,"Args":{"Description":"arg","Project":0,"Tag":"","TimeEstimation":0}}`
    arg := "ho"
    projects := []toggl.Project{
        { ID: 1, Name: "hoge", }, 
        { ID: 2, Name: "fuga", },
        { ID: 3, Name: "hoo", },
    }
    suite.mockedRepo.On("GetProjects").Return(projects, nil)

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 2, len(items))
    assert.Equal(t, fmt.Sprintf("Project: %s", projects[0].Name), items[0].Title)
    assert.Equal(t, fmt.Sprintf("Project: %s", projects[2].Name), items[1].Title)
    assert.Equal(t, alfred.ModeTell, items[0].Arg.Mode)
    assert.Equal(t, fmt.Sprintf("{\"Current\":2,\"Args\":{\"Description\":\"arg\",\"Project\":%d,\"Tag\":\"\",\"TimeEstimation\":0}}", projects[0].ID), items[0].Arg.Data)
}

func (suite *AddEntryTestSuite) TestItems3() {
    // given
    dataStr := `{"Current":2,"Args":{"Description":"arg","Project":1,"Tag":"","TimeEstimation":0}}`
    arg := "ho"
    tags := []toggl.Tag{
        { ID: 1, Name: "hoge", }, 
        { ID: 2, Name: "fuga", },
        { ID: 3, Name: "hoo", },
    }
    suite.mockedRepo.On("GetTags").Return(tags, nil)

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 2, len(items))
    assert.Equal(t, fmt.Sprintf("Tag: %s", tags[0].Name), items[0].Title)
    assert.Equal(t, fmt.Sprintf("Tag: %s", tags[2].Name), items[1].Title)
    assert.Equal(t, alfred.ModeTell, items[0].Arg.Mode)
    assert.Equal(t, fmt.Sprintf("{\"Current\":3,\"Args\":{\"Description\":\"arg\",\"Project\":1,\"Tag\":\"%s\",\"TimeEstimation\":0}}", tags[0].Name), items[0].Arg.Data)
}

func (suite *AddEntryTestSuite) TestItems4_Normal() {
    // given
    dataStr := `{"Current":3,"Args":{"Description":"arg","Project":1,"Tag":"hoge","TimeEstimation":0}}`

    arg := "31"

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 1, len(items))
    item := items[0]
    assert.Equal(t, fmt.Sprintf("Time estimation [min]: %s", arg), item.Title)
    itemArg := item.Arg
    assert.Equal(t, alfred.ModeDo, itemArg.Mode)
    assert.Equal(t, fmt.Sprintf("{\"Current\":4,\"Args\":{\"Description\":\"arg\",\"Project\":1,\"Tag\":\"hoge\",\"TimeEstimation\":%s}}", arg), item.Arg.Data)
}

func (suite *AddEntryTestSuite) TestItems4_Invalid() {
    // given
    dataStr := `{"Current":3,"Args":{"Description":"arg","Project":1,"Tag":"hoge","TimeEstimation":0}}`

    arg := "aa"

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 1, len(items))
    item := items[0]
    assert.Equal(t, "Time estimation [min]: 30", item.Title)
    itemArg := item.Arg
    assert.Equal(t, alfred.ModeDo, itemArg.Mode)
    assert.Equal(t, "{\"Current\":4,\"Args\":{\"Description\":\"arg\",\"Project\":1,\"Tag\":\"hoge\",\"TimeEstimation\":30}}", item.Arg.Data)
}

func (suite *AddEntryTestSuite) TestDo() {
    // given
    dataStr := "{\"Current\":4,\"Args\":{\"Description\":\"arg\",\"Project\":1,\"Tag\":\"hoge\",\"TimeEstimation\":30}}"
    suite.mockedRepo.On("Insert", mock.Anything).Return(nil).Once()

    // when
    suite.com.Do(dataStr)

    // then
    t := suite.T()
    suite.mockedRepo.AssertExpectations(t)
    entity := suite.mockedRepo.Calls[0].Arguments[0].(*domain.TimeEntryEntity)
    assert.Equal(t, 1, entity.Entry.Pid)
    assert.Equal(t, "arg", entity.Entry.Description)
    assert.Equal(t, []string{"hoge"}, entity.Entry.Tags)
    assert.Equal(t, 30, entity.Estimation.Duration)
}


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
    assert.Equal(t, "actual duration: 0 [min], estimation: -", item.Subtitle)
    var itemData command.DetailRefData
    err := json.Unmarshal([]byte(item.Arg.Data), &itemData)
    assert.Nil(t, err)
    assert.Equal(t, 2, itemData.ID)
}

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
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start, Stop: &stop, Duration: duration},
        Estimation: domain.Estimation{
            Duration: 66,
            Memo: "memo test",
        },
    }, nil).Once()


    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 5, len(items))
    assert.Equal(t, "Description: item42", items[0].Title)
    assert.Equal(t, "Estimated duration: 66 [min]", items[1].Title)
    assert.Equal(t, "Start: 2022/01/24 13:50", items[2].Title)
    assert.Equal(t, "Stop: 2022/01/24 15:53", items[3].Title)
    assert.Equal(t, "Memo: memo test", items[4].Title)
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
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start, Duration: -1},
        Estimation: domain.Estimation{
            Duration: 66,
            Memo: "memo test",
        },
    }, nil).Once()


    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 5, len(items))
    assert.Equal(t, "Description: item42", items[0].Title)
    assert.Equal(t, "Estimated duration: 66 [min]", items[1].Title)
    assert.Equal(t, "Start: 2022/01/24 13:50", items[2].Title)
    assert.Equal(t, "Memo: memo test", items[3].Title)
    assert.Equal(t, "Stop this entry", items[4].Title)
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
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Start: &start, Stop: &stop, Duration: duration},
    }, nil).Once()


    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 3, len(items))
    assert.Equal(t, "Description: item42", items[0].Title)
    assert.Equal(t, "Start: 2022/01/24 13:50", items[1].Title)
    assert.Equal(t, "Stop: 2022/01/24 15:53", items[2].Title)
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
