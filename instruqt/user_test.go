package instruqt

import (
	"errors"
	"log"
	"testing"

	"github.com/shurcooL/graphql"
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
			Details: struct {
				FirstName graphql.String
				LastName  graphql.String
				Email     graphql.String
			}{
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
			Profile: struct {
				Display_Name graphql.String
				Email        graphql.String
			}{
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
