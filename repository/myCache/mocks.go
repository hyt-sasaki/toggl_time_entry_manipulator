package myCache

import (
    "github.com/stretchr/testify/mock"
)

type MockedCache struct {
    mock.Mock
}

func (m *MockedCache) Save() {
    m.Called()
}

func (m *MockedCache) GetData() *Data {
    args := m.Called()
    return args.Get(0).(*Data)
}
