package instruqt

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetInvite(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	inviteID := "invite-123"
	expectedInvite := TrackInvite{
		Id:          "invite-123",
		PublicTitle: "Test Invite",
		RuntimeParameters: struct {
			EnvironmentVariables []variable
		}{
			EnvironmentVariables: []variable{
				{Key: "ENV_VAR", Value: "value"},
			},
		},
		Claims: []TrackInviteClaim{
			{
				Id: "claim-1",
				User: struct {
					Id string
				}{
					Id: "user-1",
				},
				ClaimedAt: time.Now(),
			},
		},
	}

	queryResult := inviteQuery{
		TrackInvite: expectedInvite,
	}

	mockClient.On("Query", mock.Anything, &inviteQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*inviteQuery)
		*q = queryResult
	}).Return(nil)

	invite, err := client.GetInvite(inviteID)

	assert.NoError(t, err)
	assert.Equal(t, expectedInvite, invite)
	mockClient.AssertExpectations(t)
}

func TestGetInvites(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	expectedInvites := []TrackInvite{
		{
			Id:          "invite-123",
			PublicTitle: "Test Invite 1",
		},
		{
			Id:          "invite-456",
			PublicTitle: "Test Invite 2",
		},
	}

	queryResult := invitesQuery{
		TrackInvites: expectedInvites,
	}

	mockClient.On("Query", mock.Anything, &invitesQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*invitesQuery)
		*q = queryResult
	}).Return(nil)

	invites, err := client.GetInvites()

	assert.NoError(t, err)
	assert.Equal(t, expectedInvites, invites)
	mockClient.AssertExpectations(t)
}

func TestGetInvite_Error(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	mockClient.On("Query", mock.Anything, &inviteQuery{}, mock.Anything).Return(errors.New("graphql error"))

	invite, err := client.GetInvite("invite-123")

	assert.Error(t, err)
	assert.Empty(t, invite)
	assert.Contains(t, err.Error(), "graphql error")
	mockClient.AssertExpectations(t)
}

func TestGetInvites_Error(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	mockClient.On("Query", mock.Anything, &invitesQuery{}, mock.Anything).Return(errors.New("graphql error"))

	invites, err := client.GetInvites()

	assert.Error(t, err)
	assert.Empty(t, invites)
	assert.Contains(t, err.Error(), "graphql error")
	mockClient.AssertExpectations(t)
}
