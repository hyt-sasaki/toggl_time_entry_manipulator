package repository

import (
	"testing"
    "errors"
	"time"
	"toggl_time_entry_manipulator/client"
	"toggl_time_entry_manipulator/domain"

	"github.com/jason0x43/go-toggl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
    suite.Suite
    mockedToggleClient *client.MockedToggleClient
    mockedEstimationClient *client.MockedEstimationClient
    mockedAccount toggl.Account
    repo *timeEntryRepository
}

func TestRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) SetupTest() {
    suite.mockedEstimationClient = &client.MockedEstimationClient{}
    suite.mockedToggleClient = &client.MockedToggleClient{}
    suite.mockedAccount = toggl.Account{}
    suite.mockedAccount.Data.TimeEntries = []toggl.TimeEntry{
        { ID: 1, }, 
        { ID: 2, },
    }
    suite.mockedAccount.Data.Projects = []toggl.Project{
        { ID: 3, }, 
        { ID: 4, },
    }
    suite.mockedAccount.Data.Tags = []toggl.Tag{
        { ID: 5, Name: "tag1", }, 
        { ID: 6, Name: "tag2", },
    }
    suite.repo = &timeEntryRepository{
        estimationClient: suite.mockedEstimationClient,
        togglClient: suite.mockedToggleClient,
    }
}
func (suite *RepositoryTestSuite) TestFetchAccount() {
    // given
    suite.mockedToggleClient.On("GetAccount").Return(suite.mockedAccount, nil).Once()

    // when
    account, _ := suite.repo.FetchTogglAccount()

    // then
    t := suite.T()
    assert.Equal(t, suite.mockedAccount, account)
    suite.mockedToggleClient.AssertExpectations(t)
}

func (suite *RepositoryTestSuite) TestFetch() {
    // given
    entryIds := []int64{1, 2}

    // mock
    mockedEstimations := []*domain.Estimation{
        { Duration: 30, Memo: "memo1", },
        { Duration: 40, Memo: "memo2", },
    }
    suite.mockedEstimationClient.On("Fetch", entryIds).Return(mockedEstimations, nil).Once()


    // when
    entities, _ := suite.repo.Fetch(suite.mockedAccount)

    // then
    t := suite.T()
    assert.Equal(t, []domain.TimeEntryEntity{
        {
            Entry: suite.mockedAccount.Data.TimeEntries[0],
            Estimation: *mockedEstimations[0],
        },
        {
            Entry: suite.mockedAccount.Data.TimeEntries[1],
            Estimation: *mockedEstimations[1],
        },
    }, entities)
    assert.Equal(t, "1", entities[0].GetId())
    assert.Equal(t, "2", entities[1].GetId())
    suite.mockedToggleClient.AssertExpectations(t)
}

func (suite *RepositoryTestSuite) TestFetch_whenSomeEstimationsDoNotExist() {

    // given
    entryIds := []int64{1, 2}

    // mock
    account := toggl.Account{}
    account.Data.TimeEntries = []toggl.TimeEntry{
        { ID: 1, }, 
        { ID: 2, },
    }
    mockedEstimations := []*domain.Estimation{
        nil,
        { Duration: 40, Memo: "memo2", },
    }

    suite.mockedEstimationClient.On("Fetch", entryIds).Return(mockedEstimations, nil).Once()

    // when
    entities, _ := suite.repo.Fetch(account)

    // then
    t := suite.T()
    assert.Equal(t, "1", entities[0].GetId())
    assert.Equal(t, "2", entities[1].GetId())
    assert.Equal(t, domain.Estimation{}, entities[0].Estimation)
    assert.Equal(t, *mockedEstimations[1], entities[1].Estimation)
}

func (suite *RepositoryTestSuite) TestInsert() {
    // given
    description := "description"
    pid := 1
    tag := "tag2"
    duration := 33
    entity := domain.Create(description, pid, tag, duration)
    timeEntry := toggl.TimeEntry{
        ID: 10,
        Pid: pid,
        Tags: []string{tag},
        Duration: -1,
    }
    suite.mockedToggleClient.On("StartTimeEntry", description, pid, []string{tag}).Return(timeEntry, nil).Once()
    suite.mockedEstimationClient.On("Insert", "10", mock.Anything).Return(nil).Once()

    // when
    suite.repo.Insert(&entity, suite.mockedAccount.Data.Tags)

    // then
    t := suite.T()
    suite.mockedToggleClient.AssertExpectations(t)
    capturedId := suite.mockedEstimationClient.Calls[0].Arguments[0]
    capturedEstimation := suite.mockedEstimationClient.Calls[0].Arguments[1].(domain.Estimation)
    assert.Equal(t, "10", capturedId)
    assert.Equal(t, 33, capturedEstimation.Duration)
}

func (suite *RepositoryTestSuite) TestUpdate() {
    // given
    description := "description"
    pid := 1
    tag := "tag2"
    duration := 33
    entity := domain.Create(description, pid, tag, duration)
    suite.mockedToggleClient.On("UpdateTimeEntry", entity.Entry).Return(toggl.TimeEntry{
        ID: 10,
        Pid: pid,
        Tags: []string{tag},
        Duration: -1,
    }, nil).Once()
    suite.mockedEstimationClient.On("Update", "10", entity.Estimation).Return(nil).Once()

    // when
    suite.repo.Update(&entity, suite.mockedAccount.Data.Tags)

    // then
    t := suite.T()
    suite.mockedToggleClient.AssertExpectations(t)
    suite.mockedEstimationClient.AssertExpectations(t)
}

func (suite *RepositoryTestSuite) TestStop() {
    // given
    start := time.Now().Add(-time.Hour)
    stop := time.Now()
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ ID: 10, Start: &start },
    }
    suite.mockedToggleClient.On("StopTimeEntry", entity.Entry).Return(toggl.TimeEntry{
        ID: 10, Start: &start, Stop: &stop,
    }, nil).Once()

    // when
    suite.repo.Stop(&entity)

    // then
    t := suite.T()
    assert.Equal(t, entity.Entry.Stop.After(*entity.Entry.Start), true)
}

func (suite *RepositoryTestSuite) TestDelete() {
    // given
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ ID: 10 },
    }
    suite.mockedEstimationClient.On("Delete", "10").Return(nil).Once()
    suite.mockedToggleClient.On("DeleteTimeEntry", entity.Entry).Return(nil).Once()
    
    // when
    suite.repo.Delete(&entity)

    // then
    t := suite.T()
    suite.mockedToggleClient.AssertExpectations(t)
    suite.mockedEstimationClient.AssertExpectations(t)
}

func (suite *RepositoryTestSuite) TestDelete_estimationDeletionFailed() {
    // given
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ ID: 10 },
    }
    suite.mockedEstimationClient.On("Delete", "10").Return(errors.New("deletion failed")).Once()
    
    // when
    suite.repo.Delete(&entity)

    // then
    t := suite.T()
    suite.mockedEstimationClient.AssertExpectations(t)
    suite.mockedToggleClient.AssertNotCalled(t, "DeleteTimeEntry")
}

func (suite *RepositoryTestSuite) TestDelete_rollback() {
    // given
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ ID: 10 },
    }
    suite.mockedEstimationClient.On("Delete", "10").Return(nil).Once()
    suite.mockedToggleClient.On("DeleteTimeEntry", entity.Entry).Return(errors.New("deletion faield")).Once()
    suite.mockedEstimationClient.On("Insert", "10", entity.Estimation).Return(nil).Once()
    
    // when
    suite.repo.Delete(&entity)

    // then
    t := suite.T()
    suite.mockedToggleClient.AssertExpectations(t)
    suite.mockedEstimationClient.AssertExpectations(t)
}

func (suite *RepositoryTestSuite) TestContinue() {
    // given
    start, _ := time.ParseInLocation("06-01-02 15:04", "21-01-27 13:25", time.Local)
    stop, _ := time.ParseInLocation("06-01-02 15:04", "21-01-27 13:45", time.Local)
    now, _ := time.ParseInLocation("06-01-02 15:04", "21-01-27 18:00", time.Local)
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Pid: 1, Description: "item42", Start: &start, Stop: &stop, Duration: 1200},
        Estimation: domain.Estimation{Duration: 100, Memo: "old memo", CreatedTm: start, UpdatedTm: stop},
    }
    suite.mockedToggleClient.On("ContinueTimeEntry", entity.Entry).Return(toggl.TimeEntry{
        ID: 43, Pid: 1, Description: "item42", Start: &now, Duration: -1,
    }, nil).Once()
    suite.mockedEstimationClient.On("Insert", "43", entity.Estimation).Return(nil).Once()

    // when
    newEntity, _ := suite.repo.Continue(&entity)

    // then
    t := suite.T()
    assert.Equal(t, 43, newEntity.Entry.ID)
    assert.Equal(t, 1, newEntity.Entry.Pid)
    assert.Equal(t, "item42", newEntity.Entry.Description)
    assert.Equal(t, &now, newEntity.Entry.Start)
    assert.Nil(t, newEntity.Entry.Stop)
    assert.True(t, newEntity.IsRunning())
    assert.Equal(t, entity.Estimation.Duration, newEntity.Estimation.Duration)
    assert.Equal(t, entity.Estimation.Memo, newEntity.Estimation.Memo)
}
