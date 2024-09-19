// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package instruqt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTrackById(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	trackID := "track-123"
	expectedTrack := Track{
		Id:          "track-123",
		Slug:        "test-slug",
		Title:       "Test Track",
		Description: "Test Description",
	}

	queryResult := trackQuery{
		Track: expectedTrack,
	}

	mockClient.On("Query", mock.Anything, &trackQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*trackQuery)
		*q = queryResult
	}).Return(nil)

	track, err := client.GetTrackById(trackID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTrack, track)
	mockClient.AssertExpectations(t)
}

func TestGetUserTrackById(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	userID := "user-123"
	trackID := "track-123"
	expectedTrack := SandboxTrack{
		Id:          "track-123",
		Slug:        "test-slug",
		Title:       "Test Track",
		Description: "Test Description",
	}

	queryResult := userTrackQueryWithChallenges{
		Track: expectedTrack,
	}

	mockClient.On("Query", mock.Anything, &userTrackQueryWithChallenges{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*userTrackQueryWithChallenges)
		*q = queryResult
	}).Return(nil)

	track, err := client.GetUserTrackById(userID, trackID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTrack, track)
	mockClient.AssertExpectations(t)
}

func TestGetTrackBySlug(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	trackSlug := "test-slug"
	expectedTrack := Track{
		Id:          "track-123",
		Slug:        "test-slug",
		Title:       "Test Track",
		Description: "Test Description",
	}

	queryResult := trackQueryBySlug{
		Track: expectedTrack,
	}

	mockClient.On("Query", mock.Anything, &trackQueryBySlug{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*trackQueryBySlug)
		*q = queryResult
	}).Return(nil)

	track, err := client.GetTrackBySlug(trackSlug)

	assert.NoError(t, err)
	assert.Equal(t, expectedTrack, track)
	mockClient.AssertExpectations(t)
}

func TestGetTrackUnlockedChallenge(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	userID := "user-123"
	trackID := "track-123"
	expectedChallenge := Challenge{
		Id:     "challenge-123",
		Slug:   "test-challenge",
		Title:  "Test Challenge",
		Status: "unlocked",
	}

	track := SandboxTrack{
		Id:          "track-123",
		Slug:        "test-slug",
		Title:       "Test Track",
		Description: "Test Description",
		Challenges: []Challenge{
			expectedChallenge,
			{Id: "challenge-456", Slug: "locked-challenge", Title: "Locked Challenge", Status: "locked"},
		},
	}

	mockClient.On("Query", mock.Anything, &userTrackQueryWithChallenges{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*userTrackQueryWithChallenges)
		q.Track = track
	}).Return(nil)

	challenge, err := client.GetTrackUnlockedChallenge(userID, trackID)

	assert.NoError(t, err)
	assert.Equal(t, expectedChallenge, challenge)
	mockClient.AssertExpectations(t)
}

func TestGetTracks(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	expectedTracks := []Track{
		{Id: "track-123", Slug: "test-slug", Title: "Test Track 1", Description: "Description 1"},
		{Id: "track-456", Slug: "another-slug", Title: "Test Track 2", Description: "Description 2"},
	}

	queryResult := tracksQuery{
		Tracks: expectedTracks,
	}

	mockClient.On("Query", mock.Anything, &tracksQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*tracksQuery)
		*q = queryResult
	}).Return(nil)

	tracks, err := client.GetTracks()

	assert.NoError(t, err)
	assert.Equal(t, expectedTracks, tracks)
	mockClient.AssertExpectations(t)
}

func TestGenerateOneTimePlayToken(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	trackID := "track-123"
	expectedToken := "one-time-play-token"

	mutationResult := struct {
		GenerateOneTimePlayToken string `graphql:"generateOneTimePlayToken(trackID: $trackID)"`
	}{
		GenerateOneTimePlayToken: expectedToken,
	}

	mockClient.On("Mutate", mock.Anything, mock.AnythingOfType("*struct { GenerateOneTimePlayToken string \"graphql:\\\"generateOneTimePlayToken(trackID: $trackID)\\\"\" }"), mock.Anything).Run(func(args mock.Arguments) {
		m := args.Get(1).(*struct {
			GenerateOneTimePlayToken string `graphql:"generateOneTimePlayToken(trackID: $trackID)"`
		})
		*m = mutationResult
	}).Return(nil)

	token, err := client.GenerateOneTimePlayToken(trackID)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockClient.AssertExpectations(t)
}
