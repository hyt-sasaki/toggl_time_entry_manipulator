package repository

import (
    "fmt"
	"testing"
	"time"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository/myCache"

	"github.com/jason0x43/go-toggl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CacheRepositoryTestSuite struct {
    suite.Suite
    mockedTimeEntryRepository *MockedTimeEntryRepository
    mockedCache *myCache.MockedCache
    repo ICachedRepository
}

func TestCacheRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new (CacheRepositoryTestSuite))
}

func (suite *CacheRepositoryTestSuite) SetupTest() {
    suite.mockedTimeEntryRepository = &MockedTimeEntryRepository{}
    suite.mockedCache = &myCache.MockedCache{}
    suite.repo = NewCachedRepository(suite.mockedCache, suite.mockedTimeEntryRepository)
}

func (suite *CacheRepositoryTestSuite) TestFetchOneById() {
    // given
    entryId := 42
    mockedData := &myCache.Data{
        Workspace: 1,
        Projects: []toggl.Project{},
        Tags: []toggl.Tag{},
        Entities: []domain.TimeEntryEntity{
            { Entry: toggl.TimeEntry{ ID: entryId } },
        },
        Time: time.Now(),
    }
    suite.mockedCache.On("GetData").Return(mockedData)

    // when
    entity, _ := suite.repo.FindOneById(entryId)

    // then
    t := suite.T()
    assert.Equal(t, mockedData.Entities[0], entity)
}

func (suite *CacheRepositoryTestSuite) TestFetchOneById_whenNoMatchedEntity() {
    // given
    entryId := 42
    mockedData := &myCache.Data{
        Workspace: 1,
        Entities: []domain.TimeEntryEntity{
            { Entry: toggl.TimeEntry{ ID: 1 } },
        },
        Time: time.Now(),
    }
    suite.mockedCache.On("GetData").Return(mockedData)

    // when
    entity, err := suite.repo.FindOneById(entryId)

    // then
    t := suite.T()
    assert.Equal(t, entity, domain.TimeEntryEntity{})
    assert.Equal(t, err.Error(), fmt.Sprintf("Resource not found: %d", entryId))
}
