package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"toggl_time_entry_manipulator/command/add"
	"toggl_time_entry_manipulator/command/list"
	"toggl_time_entry_manipulator/command/get"
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

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 2, len(items))
    item := items[0]
    assert.Equal(t, "Description: item2-1", item.Title)
    var itemData list.ItemData
    err := json.Unmarshal([]byte(item.Arg.Data), &itemData)
    assert.Nil(t, err)
    assert.Equal(t, 1, itemData.ID)
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
