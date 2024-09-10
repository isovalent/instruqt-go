package instruqt

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSandboxVariable(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	playID := "sandbox-123"
	key := "MY_VAR"
	expectedValue := "value123"

	queryResult := sandboxVarQuery{
		GetSandboxVariable: SandboxVar{
			Key:   key,
			Value: expectedValue,
		},
	}

	mockClient.On("Query", mock.Anything, &sandboxVarQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*sandboxVarQuery)
		*q = queryResult
	}).Return(nil)

	value, err := client.GetSandboxVariable(playID, key)

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
	mockClient.AssertExpectations(t)
}

func TestGetSandboxVariable_Error(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	playID := "sandbox-123"
	key := "MY_VAR"

	mockClient.On("Query", mock.Anything, &sandboxVarQuery{}, mock.Anything).Return(errors.New("graphql error"))

	value, err := client.GetSandboxVariable(playID, key)

	assert.Error(t, err)
	assert.Empty(t, value)
	assert.Contains(t, err.Error(), "graphql error")
	mockClient.AssertExpectations(t)
}

func TestGetSandboxes(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	expectedSandboxes := []Sandbox{
		{
			Last_Activity_At: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			State:            "running",
			Track: SandboxTrack{
				Id:    "track-123",
				Title: "Track 1",
			},
			Invite: TrackInvite{
				Id:          "invite-123",
				PublicTitle: "Invite 1",
			},
		},
		{
			Last_Activity_At: time.Date(2023, 1, 2, 14, 30, 0, 0, time.UTC),
			State:            "completed",
			Track: SandboxTrack{
				Id:    "track-456",
				Title: "Track 2",
			},
			Invite: TrackInvite{
				Id:          "invite-456",
				PublicTitle: "Invite 2",
			},
		},
	}

	queryResult := sandboxesQuery{
		Sandboxes: struct {
			Nodes []Sandbox
		}{
			Nodes: expectedSandboxes,
		},
	}

	mockClient.On("Query", mock.Anything, &sandboxesQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*sandboxesQuery)
		*q = queryResult
	}).Return(nil)

	sandboxes, err := client.GetSandboxes()

	assert.NoError(t, err)
	assert.Equal(t, expectedSandboxes, sandboxes)
	mockClient.AssertExpectations(t)
}

func TestGetSandboxes_Error(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	mockClient.On("Query", mock.Anything, &sandboxesQuery{}, mock.Anything).Return(errors.New("graphql error"))

	sandboxes, err := client.GetSandboxes()

	assert.Error(t, err)
	assert.Empty(t, sandboxes)
	assert.Contains(t, err.Error(), "graphql error")
	mockClient.AssertExpectations(t)
}
