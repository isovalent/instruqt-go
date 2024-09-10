package instruqt

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockGraphQLClient is a mock of the GraphQLClient interface
type MockGraphQLClient struct {
	mock.Mock
}

func (m *MockGraphQLClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	args := m.Called(ctx, q, variables)
	return args.Error(0)
}

func (m *MockGraphQLClient) Mutate(ctx context.Context, mutation interface{}, variables map[string]interface{}) error {
	args := m.Called(ctx, mutation, variables)
	return args.Error(0)
}
