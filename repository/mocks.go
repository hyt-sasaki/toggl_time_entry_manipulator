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
