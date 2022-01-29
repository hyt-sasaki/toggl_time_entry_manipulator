package deletetest

import (
	_ "toggl_time_entry_manipulator/supports"

	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jason0x43/go-toggl"

	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/command/delete"
	"toggl_time_entry_manipulator/tests"
)


type DeleteEntryTestSuite struct {
    suite.Suite
    mockedRepo *repository.MockedCachedRepository
    com *delete.DeleteEntryCommand
}

func TestDeleteEntryTestSuite(t *testing.T) {
    suite.Run(t, new(DeleteEntryTestSuite))
}

func (suite *DeleteEntryTestSuite) SetupTest() {
    suite.mockedRepo = &repository.MockedCachedRepository{}
    suite.com = &delete.DeleteEntryCommand{
        Repo: suite.mockedRepo,
    }
}

func (suite *DeleteEntryTestSuite) TestDo() {
    // given
    dataStr := tests.StringifyDetailRefData(command.DetailRefData{ID: 42})

    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42"},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(entity, nil).Once()
    suite.mockedRepo.On("Delete", &entity).Return(nil).Once()

    // when
    out, _ := suite.com.Do(dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Entry has been deleted. Description: item42", out)
    suite.mockedRepo.AssertExpectations(t)
}
