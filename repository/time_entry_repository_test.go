package repository

import (
	"testing"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/estimation_client"

	"github.com/jason0x43/go-toggl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
    suite.Suite
    mockedToggleClient *estimation_client.MockedToggleClient
    mockedEstimationClient *estimation_client.MockedEstimationClient
    repo *TimeEntryRepository
}

func TestSuite(t *testing.T) {
    suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) SetupTest() {
    suite.mockedEstimationClient = &estimation_client.MockedEstimationClient{}
    suite.mockedToggleClient = &estimation_client.MockedToggleClient{}
    suite.repo = &TimeEntryRepository{
        config: &Config{
            TogglAPIKey: "test",
        },
        estimationClient: suite.mockedEstimationClient,
        togglClient: suite.mockedToggleClient,
    }
}
func (suite *RepositoryTestSuite) TestFetchAccount() {
    // given
    mockedAccount := toggl.Account{}
    mockedAccount.Data.TimeEntries = []toggl.TimeEntry{
        { ID: 1, }, 
        { ID: 2, },
    }
    mockedAccount.Data.Projects = []toggl.Project{
        { ID: 3, }, 
        { ID: 4, },
    }
    mockedAccount.Data.Tags = []toggl.Tag{
        { ID: 5, Name: "tag1", }, 
        { ID: 6, Name: "tag2", },
    }
    suite.mockedToggleClient.On("GetAccount").Return(mockedAccount, nil).Once()

    // when
    account, _ := suite.repo.FetchTogglAccount()

    // then
    t := suite.T()
    assert.Equal(t, mockedAccount, account)
    suite.mockedToggleClient.AssertExpectations(t)
}

func (suite *RepositoryTestSuite) TestFetch() {
    // given
    entryIds := []int64{1, 2}

    // mock
    account := toggl.Account{}
    account.Data.TimeEntries = []toggl.TimeEntry{
        { ID: 1, }, 
        { ID: 2, },
    }
    mockedEstimations := []estimation_client.Estimation{
        { Duration: 30, Memo: "memo1", },
        { Duration: 40, Memo: "memo2", },
    }
    suite.mockedEstimationClient.On("Fetch", entryIds).Return(mockedEstimations, nil).Once()


    // when
    entities, _ := suite.repo.Fetch(account)

    // then
    t := suite.T()
    assert.Equal(t, []domain.TimeEntryEntity{
        {
            Entry: account.Data.TimeEntries[0],
            Estimation: mockedEstimations[0],
        },
        {
            Entry: account.Data.TimeEntries[1],
            Estimation: mockedEstimations[1],
        },
    }, entities)
    assert.Equal(t, "1", entities[0].GetId())
    assert.Equal(t, "2", entities[1].GetId())
    suite.mockedToggleClient.AssertExpectations(t)
}

func (suite *RepositoryTestSuite) TestInsert() {
    // given
    description := "description"
    pid := 1
    tag := "tag"
    duration := 33
    entity := domain.Create(description, pid, tag, duration)
    timeEntry := toggl.TimeEntry{
        ID: 10,
        Pid: pid,
        Tags: []string{tag},
        Duration: -1,
    }
    suite.mockedToggleClient.On("StartTimeEntry", description, pid, []string{tag}).Return(timeEntry, nil).Once()
    suite.mockedEstimationClient.On("Insert", "10", entity.Estimation).Return(nil).Once()

    // when
    suite.repo.Insert(&entity)

    // then
    t := suite.T()
    suite.mockedToggleClient.AssertExpectations(t)
    suite.mockedEstimationClient.AssertExpectations(t)
    assert.Equal(t, entity.Entry, timeEntry)
}