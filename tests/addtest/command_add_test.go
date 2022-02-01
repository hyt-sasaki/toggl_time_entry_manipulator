package addtest

import (
	_ "toggl_time_entry_manipulator/supports"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/mock"

	"fmt"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"

	"toggl_time_entry_manipulator/command/add"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/domain"
)

type AddEntryTestSuite struct {
    suite.Suite
    mockedRepo *repository.MockedCachedRepository
    projects []toggl.Project
    tags []toggl.Tag
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

    suite.projects = []toggl.Project{
        { ID: 1, Name: "hoge", }, 
        { ID: 2, Name: "fuga", },
        { ID: 3, Name: "hoo", },
    }
    suite.tags = []toggl.Tag{
        { ID: 1, Name: "hoge", }, 
        { ID: 2, Name: "fuga", },
        { ID: 3, Name: "hoo", },
    }
    suite.mockedRepo.On("GetProjects").Return(suite.projects, nil)
    suite.mockedRepo.On("GetTags").Return(suite.tags, nil)
}

func (suite *AddEntryTestSuite) TestItems_ProjectEdit() {
    // given
    dataStr := ""

    arg := "ho"

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 2, len(items))
    assert.Equal(t, fmt.Sprintf("Project: %s", suite.projects[0].Name), items[0].Title)
    assert.Equal(t, fmt.Sprintf("Project: %s", suite.projects[2].Name), items[1].Title)
    assertAddItemArg(t, items[0], add.StateData{
        Current: add.TagEdit, Entity: domain.TimeEntryEntity{Entry: toggl.TimeEntry{Pid: 1}}}, alfred.ModeTell)
}

func (suite *AddEntryTestSuite) TestItems_TagEdit() {
    // given
    dataStr := convertAddStateData(add.StateData{
            Current: add.TagEdit,
            Entity: domain.TimeEntryEntity{Entry: toggl.TimeEntry{Pid: 1}}})
    arg := "ho"

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 3, len(items))
    assert.Equal(t, fmt.Sprintf("Tag: %s", suite.tags[0].Name), items[0].Title)
    assert.Equal(t, fmt.Sprintf("Tag: %s", suite.tags[2].Name), items[1].Title)
    assertAddItemArg(t, items[0], add.StateData{
        Current: add.DescriptionEdit,
        Entity: domain.TimeEntryEntity{Entry: toggl.TimeEntry{Pid: 1, Tags:[]string{"hoge"}}}}, alfred.ModeTell)
}

func (suite *AddEntryTestSuite) TestItems_TagEditNoInput() {
    // given
    dataStr := convertAddStateData(add.StateData{
            Current: add.TagEdit,
            Entity: domain.TimeEntryEntity{Entry: toggl.TimeEntry{Pid: 1}}})
    arg := ""

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 5, len(items))
    assert.Equal(t, "No tag", items[0].Title)
    assertAddItemArg(t, items[0], add.StateData{
        Current: add.DescriptionEdit,
        Entity: domain.TimeEntryEntity{Entry: toggl.TimeEntry{Pid: 1, Tags: []string{}}}}, alfred.ModeTell)
}

func (suite *AddEntryTestSuite) TestItems_DescriptionEdit() {
    // given
    dataStr := convertAddStateData(add.StateData{
            Current: add.DescriptionEdit,
            Entity: domain.TimeEntryEntity{Entry: toggl.TimeEntry{Pid: 1, Tags:[]string{"hoge"}}}})
    arg := "new description"

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 2, len(items))
    assert.Equal(t, fmt.Sprintf("New description: %s", arg), items[0].Title)

    assertAddItemArg(t, items[0], add.StateData{
        Current: add.TimeEstimationEdit,
        Entity: domain.TimeEntryEntity{Entry: toggl.TimeEntry{Pid: 1, Tags:[]string{"hoge"}, Description: arg}}}, alfred.ModeTell)
}

func (suite *AddEntryTestSuite) TestItems_TimeEstimationEdit() {
    // given
    dataStr := convertAddStateData(add.StateData{
            Current: add.TimeEstimationEdit,
            Entity: domain.TimeEntryEntity{Entry: toggl.TimeEntry{Pid: 1, Tags:[]string{"hoge"}, Description: "new description"}}})
    arg := "31"

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 2, len(items))
    assert.Equal(t, fmt.Sprintf("Time estimation [min]: %s", arg), items[0].Title)

    assertAddItemArg(t, items[0], add.StateData{
        Current: add.EndEdit,
        Entity: domain.TimeEntryEntity{
            Entry: toggl.TimeEntry{Pid: 1, Tags:[]string{"hoge"}, Description: "new description"},
            Estimation: domain.Estimation{Duration: 31},}}, alfred.ModeDo)
}

func (suite *AddEntryTestSuite) TestItems_TimeEstimationEdit_Invalid() {
    // given
    dataStr := convertAddStateData(add.StateData{
            Current: add.TimeEstimationEdit,
            Entity: domain.TimeEntryEntity{Entry: toggl.TimeEntry{Pid: 1, Tags:[]string{"hoge"}, Description: "new description"}}})
    arg := "aa"

    // when
    items, _ := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Equal(t, 2, len(items))
    assert.Equal(t, "Time estimation [min]: 30", items[0].Title)

    assertAddItemArg(t, items[0], add.StateData{
        Current: add.EndEdit,
        Entity: domain.TimeEntryEntity{
            Entry: toggl.TimeEntry{Pid: 1, Tags:[]string{"hoge"}, Description: "new description"},
            Estimation: domain.Estimation{Duration: 30},}}, alfred.ModeDo)
}


func (suite *AddEntryTestSuite) TestDo() {
    // given
    dataStr := convertAddStateData(add.StateData{
        Current: add.EndEdit,
        Entity: domain.TimeEntryEntity{
            Entry: toggl.TimeEntry{Pid: 1, Tags:[]string{"hoge"}, Description: "new description"},
            Estimation: domain.Estimation{Duration: 31}}})
    suite.mockedRepo.On("Insert", mock.Anything).Return(nil).Once()

    // when
    out, _ := suite.com.Do(dataStr)

    // then
    t := suite.T()
    entity := suite.mockedRepo.Calls[0].Arguments[0].(*domain.TimeEntryEntity)
    assert.Equal(t, 1, entity.Entry.Pid)
    assert.Equal(t, "new description", entity.Entry.Description)
    assert.Equal(t, []string{"hoge"}, entity.Entry.Tags)
    assert.Equal(t, 31, entity.Estimation.Duration)
    assert.Equal(t, "Time entry [new description] has started", out)
}
