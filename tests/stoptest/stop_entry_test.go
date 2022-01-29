package stoptest

import (
	_ "toggl_time_entry_manipulator/supports"

	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"encoding/json"
    "time"

	"github.com/jason0x43/go-toggl"

	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/command/stop"
)

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
    assert.Equal(t, "Entry has stopped. Description: item42", out)
}
