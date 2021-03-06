package repository

import (
	"toggl_time_entry_manipulator/domain"

	"github.com/jason0x43/go-toggl"
	"github.com/stretchr/testify/mock"
)

type MockedCachedRepository struct {
    mock.Mock
}

func (m *MockedCachedRepository) Fetch() ([]domain.TimeEntryEntity, error) {
    args := m.Called()
    return args.Get(0).([]domain.TimeEntryEntity), args.Error(1)
}

func (m *MockedCachedRepository) FindOneById(entryId int) (domain.TimeEntryEntity, error) {
    args := m.Called(entryId)
    return args.Get(0).(domain.TimeEntryEntity), args.Error(1)
}

func (m *MockedCachedRepository) GetProjects() ([]toggl.Project, error) {
    args := m.Called()
    return args.Get(0).([]toggl.Project), args.Error(1)
}

func (m *MockedCachedRepository) GetTags() ([]toggl.Tag, error) {
    args := m.Called()
    return args.Get(0).([]toggl.Tag), args.Error(1)
}

func (m *MockedCachedRepository) Insert(entity *domain.TimeEntryEntity) error {
    args := m.Called(entity)
    return args.Error(0)
}

func (m *MockedCachedRepository) Update(entity *domain.TimeEntryEntity) error {
    args := m.Called(entity)
    return args.Error(0)
}

func (m *MockedCachedRepository) Stop(entity *domain.TimeEntryEntity) error {
    args := m.Called(entity)
    return args.Error(0)
}

func (m *MockedCachedRepository) Delete(entity *domain.TimeEntryEntity) error {
    args := m.Called(entity)
    return args.Error(0)
}

func (m *MockedCachedRepository) Continue(entity *domain.TimeEntryEntity) (newEntity domain.TimeEntryEntity, err error) {
    args := m.Called(entity)
    return args.Get(0).(domain.TimeEntryEntity), args.Error(1)
}

type MockedTimeEntryRepository struct {
    mock.Mock
}

func (m *MockedTimeEntryRepository) Fetch(account toggl.Account) ([]domain.TimeEntryEntity, error) {
    args := m.Called(account)
    return args.Get(0).([]domain.TimeEntryEntity), args.Error(1)
}

func (m *MockedTimeEntryRepository) FetchTogglAccount() (toggl.Account, error) {
    args := m.Called()
    return args.Get(0).(toggl.Account), args.Error(1)
}

func (m *MockedTimeEntryRepository) Insert(entity *domain.TimeEntryEntity, tags []toggl.Tag) error {
    args := m.Called(entity, tags)
    return args.Error(0)
}

func (m *MockedTimeEntryRepository) Update(entity *domain.TimeEntryEntity, tags []toggl.Tag) error {
    args := m.Called(entity, tags)
    return args.Error(0)
}

func (m *MockedTimeEntryRepository) Stop(entity *domain.TimeEntryEntity) error {
    args := m.Called(entity)
    return args.Error(0)
}

func (m *MockedTimeEntryRepository) Continue(entity *domain.TimeEntryEntity) (domain.TimeEntryEntity, error) {
    args := m.Called(entity)
    return args.Get(0).(domain.TimeEntryEntity), args.Error(1)
}

func (m *MockedTimeEntryRepository) Delete(entity *domain.TimeEntryEntity) error {
    args := m.Called(entity)
    return args.Error(0)
}
