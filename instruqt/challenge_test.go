package instruqt

import (
	"errors"
	"testing"

	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetChallenge(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	challengeID := "challenge-123"
	expectedChallenge := Challenge{
		Id:     "challenge-123",
		Slug:   "test-slug",
		Title:  "Test Challenge",
		Index:  1,
		Status: "active",
	}

	queryResult := challengeQuery{
		Challenge: expectedChallenge,
	}

	mockClient.On("Query", mock.Anything, &challengeQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*challengeQuery)
		*q = queryResult
	}).Return(nil)

	challenge, err := client.GetChallenge(challengeID)

	assert.NoError(t, err)
	assert.Equal(t, expectedChallenge, challenge)
	mockClient.AssertExpectations(t)
}

func TestGetUserChallenge(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	userID := "user-123"
	challengeID := "challenge-123"
	expectedChallenge := Challenge{
		Id:     "challenge-123",
		Slug:   "test-slug",
		Title:  "Test Challenge",
		Index:  1,
		Status: "completed",
	}

	queryResult := userChallengeQuery{
		Challenge: expectedChallenge,
	}

	mockClient.On("Query", mock.Anything, &userChallengeQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*userChallengeQuery)
		*q = queryResult
	}).Return(nil)

	challenge, err := client.GetUserChallenge(userID, challengeID)

	assert.NoError(t, err)
	assert.Equal(t, expectedChallenge, challenge)
	mockClient.AssertExpectations(t)
}
func TestSkipToChallenge(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	userID := "user-123"
	trackID := "track-123"
	challengeID := "challenge-123"

	mutationResult := struct {
		SkipToChallenge struct {
			Id     graphql.String
			Status graphql.String
		} `graphql:"skipToChallenge(trackID: $trackID, challengeID: $challengeID, userID: $userID)"`
	}{
		SkipToChallenge: struct {
			Id     graphql.String
			Status graphql.String
		}{
			Id:     graphql.String(challengeID),
			Status: graphql.String("skipped"),
		},
	}

	mockClient.On("Mutate", mock.Anything, mock.AnythingOfType("*struct { SkipToChallenge struct { Id graphql.String; Status graphql.String } \"graphql:\\\"skipToChallenge(trackID: $trackID, challengeID: $challengeID, userID: $userID)\\\"\" }"), mock.Anything).
		Run(func(args mock.Arguments) {
			m := args.Get(1).(*struct {
				SkipToChallenge struct {
					Id     graphql.String
					Status graphql.String
				} `graphql:"skipToChallenge(trackID: $trackID, challengeID: $challengeID, userID: $userID)"`
			})
			*m = mutationResult
		}).
		Return(nil)

	err := client.SkipToChallenge(userID, trackID, challengeID)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestGetChallenge_Error(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	mockClient.On("Query", mock.Anything, &challengeQuery{}, mock.Anything).Return(errors.New("graphql error"))

	challenge, err := client.GetChallenge("challenge-123")

	assert.Error(t, err)
	assert.Empty(t, challenge)
	assert.Contains(t, err.Error(), "graphql error")
	mockClient.AssertExpectations(t)
}

func TestSkipToChallenge_Error(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	mockClient.On("Mutate", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("graphql mutation error"))

	err := client.SkipToChallenge("user-123", "track-123", "challenge-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "graphql mutation error")
	mockClient.AssertExpectations(t)
}
