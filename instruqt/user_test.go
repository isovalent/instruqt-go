// Copyright 2024 Cisco Systems, Inc. and its affiliates

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
	"errors"
	"log"
	"testing"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserInfo_Details(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		InfoLogger:    log.New(log.Writer(), "INFO: ", log.LstdFlags), // Initialize InfoLogger
	}

	userID := "12345"
	expectedUserInfo := UserInfo{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	queryResult := userInfoQuery{
		User: User{
			Details: &UserDetails{
				FirstName: graphql.String("John"),
				LastName:  graphql.String("Doe"),
				Email:     graphql.String("john.doe@example.com"),
			},
		},
	}

	mockClient.On("Query", mock.Anything, &userInfoQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*userInfoQuery)
		*q = queryResult
	}).Return(nil)

	userInfo, err := client.GetUserInfo(userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserInfo, userInfo)
	mockClient.AssertExpectations(t)
}

func TestGetUserInfo_Profile(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		InfoLogger:    log.New(log.Writer(), "INFO: ", log.LstdFlags), // Initialize InfoLogger
	}

	userID := "12345"
	expectedUserInfo := UserInfo{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.smith@example.com",
	}

	queryResult := userInfoQuery{
		User: User{
			Profile: &UserProfile{
				Display_Name: graphql.String("Jane Smith"),
				Email:        graphql.String("jane.smith@example.com"),
			},
		},
	}

	mockClient.On("Query", mock.Anything, &userInfoQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*userInfoQuery)
		*q = queryResult
	}).Return(nil)

	userInfo, err := client.GetUserInfo(userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserInfo, userInfo)
	mockClient.AssertExpectations(t)
}

func TestGetUserInfo_Error(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		InfoLogger:    log.New(log.Writer(), "INFO: ", log.LstdFlags), // Initialize InfoLogger
	}

	userID := "12345"

	mockClient.On("Query", mock.Anything, &userInfoQuery{}, mock.Anything).Return(errors.New("graphql error"))

	userInfo, err := client.GetUserInfo(userID)

	assert.Error(t, err)
	assert.Equal(t, UserInfo{}, userInfo)
	assert.Contains(t, err.Error(), "Failed to retrieve user info")
	mockClient.AssertExpectations(t)
}
