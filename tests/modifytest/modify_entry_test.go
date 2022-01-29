package modifytest

import (
	_ "toggl_time_entry_manipulator/supports"
	"encoding/json"
	"fmt"
    "time"
	"testing"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/command/modify"
)
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

