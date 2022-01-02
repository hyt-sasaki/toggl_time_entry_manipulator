package repository

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/jason0x43/go-toggl"
    "toggl_time_entry_manipulator/estimation_client"
    "github.com/stretchr/testify/assert"
)

type RepositoryTestSuite struct {
    suite.Suite
    mockedToggleClient estimation_client.MockedToggleClient
    mockedEstimationClient estimation_client.MockedEstimationClient
}

func TestSuite(t *testing.T) {
    suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) SetupTest() {
    suite.mockedEstimationClient = estimation_client.MockedEstimationClient{}
    suite.mockedToggleClient = estimation_client.MockedToggleClient{}
}

func (suite *RepositoryTestSuite) TestFetch() {
    // given
    entryIds := []int64{1, 2}

    // mock
    mockedAccount := toggl.Account{}
    mockedAccount.Data.TimeEntries = []toggl.TimeEntry{
        { ID: 1, }, 
        { ID: 2, },
    }
    mockedEstimations := []estimation_client.Estimation{
        { Duration: 30, Memo: "memo1", },
        { Duration: 40, Memo: "memo2", },
    }
    suite.mockedToggleClient.On("GetAccount").Return(mockedAccount, nil)
    suite.mockedEstimationClient.On("Fetch", entryIds).Return(mockedEstimations, nil)

    repo := TimeEntryRepository{
        config: &Config{
            TogglAPIKey: "test",
        },
        estimationClient: &suite.mockedEstimationClient,
        togglClient: &suite.mockedToggleClient,
    }

    // when
    account, _ := repo.FetchTogglAccount()
    entities, _ := repo.Fetch(account)

    // then
    assert.Equal(suite.T(), 2, len(entities))
}

