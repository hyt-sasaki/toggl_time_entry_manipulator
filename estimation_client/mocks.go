package estimation_client

import (
	"github.com/jason0x43/go-toggl"
    "github.com/stretchr/testify/mock"
)

type MockedEstimationClient struct {
    mock.Mock
}

func (m *MockedEstimationClient) Fetch(entryIds []int64) ([]Estimation, error) {
    args := m.Called(entryIds)
    return args.Get(0).([]Estimation), args.Error(1)
}

func (m *MockedEstimationClient) Insert(id string, estimation Estimation) (err error) {
    args := m.Called(id, estimation)
    return args.Error(0)
}

func (m *MockedEstimationClient) Close() {
    m.Called()
}

type MockedToggleClient struct {
    mock.Mock
}

func (m *MockedToggleClient) GetAccount() (account toggl.Account, err error) {
    args := m.Called()
    return args.Get(0).(toggl.Account), args.Error(1)
}

func (m *MockedToggleClient) StartTimeEntry(description string, pid int, tags []string) (err error) {
    args := m.Called(description, pid, tags)
    return args.Error(0)
}

