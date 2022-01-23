package client

import (
	"github.com/jason0x43/go-toggl"
    "github.com/stretchr/testify/mock"
    "toggl_time_entry_manipulator/domain"
)

type MockedEstimationClient struct {
    mock.Mock
}

func (m *MockedEstimationClient) Fetch(entryIds []int64) ([]domain.Estimation, error) {
    args := m.Called(entryIds)
    return args.Get(0).([]domain.Estimation), args.Error(1)
}

func (m *MockedEstimationClient) Insert(id string, estimation domain.Estimation) (error) {
    args := m.Called(id, estimation)
    return args.Error(0)
}

func (m *MockedEstimationClient) Close() {
    m.Called()
}

type MockedToggleClient struct {
    mock.Mock
}

func (m *MockedToggleClient) GetAccount() (toggl.Account, error) {
    args := m.Called()
    return args.Get(0).(toggl.Account), args.Error(1)
}

func (m *MockedToggleClient) StartTimeEntry(description string, pid int, tags []string) (toggl.TimeEntry, error) {
    args := m.Called(description, pid, tags)
    return args.Get(0).(toggl.TimeEntry), args.Error(1)
}

func (m *MockedToggleClient) StopTimeEntry(entry toggl.TimeEntry) (toggl.TimeEntry, error) {
    args := m.Called(entry)
    return args.Get(0).(toggl.TimeEntry), args.Error(1)
}
