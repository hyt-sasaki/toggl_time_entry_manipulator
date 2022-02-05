package continuetest

import (
	_ "toggl_time_entry_manipulator/supports"

	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

    "encoding/json"

	"github.com/jason0x43/go-toggl"
	"github.com/jason0x43/go-alfred"

	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/command/continue_entry"
	"toggl_time_entry_manipulator/tests"
)

type ContinueEntryTestSuite struct {
    suite.Suite
    mockedRepo *repository.MockedCachedRepository
    com *continue_entry.ContinueEntryCommand
}

func TestContinueEntryTestSuite(t *testing.T) {
    suite.Run(t, new(ContinueEntryTestSuite))
}

func (suite *ContinueEntryTestSuite) SetupTest() {
    suite.mockedRepo = &repository.MockedCachedRepository{}
    suite.com = &continue_entry.ContinueEntryCommand{
        Repo: suite.mockedRepo,
    }
}

func (suite *ContinueEntryTestSuite) TestItems() {
    // given
    dataStr := tests.StringifyDetailRefData(command.DetailRefData{ID: 42})
    arg := "31"
    originalEntity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Duration: 1200},
        Estimation: domain.Estimation{Duration: 10, Memo: "memo"},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(originalEntity, nil).Once()

    // when
    items, err := suite.com.Items(arg, dataStr)

    // then
    t := suite.T()
    assert.Nil(t, err)
    assert.Equal(t, 2, len(items))

    item := items[0]
    assert.Equal(t, "Duration: 31", item.Title)
    assertItemArg(t, item, domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Duration: 1200},
        Estimation: domain.Estimation{Duration: 31},
    })
}

func (suite *ContinueEntryTestSuite) TestDo() {
    // given
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Duration: 1200},
        Estimation: domain.Estimation{Duration: 31},
    }
    dataStr := tests.StringifyEntity(entity)
    suite.mockedRepo.On("Continue", &entity).Return(domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 43, Description: "item42", Duration: 1200},
        Estimation: domain.Estimation{Duration: 31},
    }, nil).Once()

    // when
    out, _ := suite.com.Do(dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Entry has been copied. Description: item42", out)
    suite.mockedRepo.AssertExpectations(t)
}

func (suite *ContinueEntryTestSuite) TestDo_byId() {
    // given
    dataStr := tests.StringifyDetailRefData(command.DetailRefData{ID: 42})
    entity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Duration: 1200},
        Estimation: domain.Estimation{Duration: 31},
    }
    originalEntity := domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 42, Description: "item42", Duration: 1200},
        Estimation: domain.Estimation{Duration: 31, Memo: "memo"},
    }
    suite.mockedRepo.On("FindOneById", 42).Return(originalEntity, nil).Once()
    suite.mockedRepo.On("Continue", &entity).Return(domain.TimeEntryEntity{
        Entry: toggl.TimeEntry{ID: 43, Description: "item42", Duration: 1200},
        Estimation: domain.Estimation{Duration: 31},
    }, nil).Once()

    // when
    out, _ := suite.com.Do(dataStr)

    // then
    t := suite.T()
    assert.Equal(t, "Entry has been copied. Description: item42", out)
    suite.mockedRepo.AssertExpectations(t)
}

func assertItemArg(t *testing.T, actualItem alfred.Item, expected domain.TimeEntryEntity) {
    actualItemArg := actualItem.Arg
    assert.Equal(t, alfred.ModeDo, actualItemArg.Mode)
    var actualEntity domain.TimeEntryEntity
    json.Unmarshal([]byte(actualItemArg.Data), &actualEntity)
    assert.Equal(t, expected, actualEntity)
}
