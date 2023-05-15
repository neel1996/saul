package repository

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/neel1996/saul/constants"
	"github.com/neel1996/saul/mocks"
	"github.com/neel1996/saul/model/db"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	context            context.Context
	mockController     *gomock.Controller
	mockDynamodbClient *mocks.MockDynamoDBClient
	repository         UserRepository
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	suite.context = context.Background()
	suite.mockController = gomock.NewController(suite.T())
	suite.mockDynamodbClient = mocks.NewMockDynamoDBClient(suite.mockController)

	suite.repository = NewUserRepository(suite.mockDynamodbClient)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *UserRepositoryTestSuite) TestGetUser_ShouldReturnUserWithEmailID() {
	email := "test@test.com"

	query := &dynamodb.QueryInput{
		TableName: aws.String(constants.UserTableName),
		ExpressionAttributeNames: map[string]string{
			"#email": "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{
				Value: email,
			},
		},
		FilterExpression:       aws.String("#email = :email"),
		Limit:                  aws.Int32(1),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		Select:                 types.SelectAllAttributes,
	}

	output := &dynamodb.QueryOutput{
		ConsumedCapacity: nil,
		Count:            1,
		Items: []map[string]types.AttributeValue{
			{
				"email": &types.AttributeValueMemberS{
					Value: email,
				},
				"name": &types.AttributeValueMemberS{
					Value: "Test User",
				},
			},
		},
		ScannedCount: 1,
	}

	expected := db.User{
		Email: email,
		Name:  "Test User",
	}

	suite.mockDynamodbClient.EXPECT().Query(suite.context, query).Return(output, nil).Times(1)

	user, err := suite.repository.GetUser(suite.context, email)

	suite.Nil(err)
	suite.Equal(expected, user)
}

func (suite *UserRepositoryTestSuite) TestGetUser_WhenQueryExecutionFails_ShouldReturnError() {
	email := "test@test.com"

	query := &dynamodb.QueryInput{
		TableName: aws.String(constants.UserTableName),
		ExpressionAttributeNames: map[string]string{
			"#email": "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{
				Value: email,
			},
		},
		FilterExpression:       aws.String("#email = :email"),
		Limit:                  aws.Int32(1),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		Select:                 types.SelectAllAttributes,
	}

	suite.mockDynamodbClient.EXPECT().Query(suite.context, query).
		Return(nil, errors.New("failed to returned data")).
		Times(1)

	_, err := suite.repository.GetUser(suite.context, email)

	suite.NotNil(err)
}

func (suite *UserRepositoryTestSuite) TestGetUser_WhenNoUsersMatchEmailID_ShouldReturnError() {
	email := "test@test.com"

	query := &dynamodb.QueryInput{
		TableName: aws.String(constants.UserTableName),
		ExpressionAttributeNames: map[string]string{
			"#email": "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{
				Value: email,
			},
		},
		FilterExpression:       aws.String("#email = :email"),
		Limit:                  aws.Int32(1),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		Select:                 types.SelectAllAttributes,
	}

	output := &dynamodb.QueryOutput{
		ConsumedCapacity: nil,
		Count:            1,
		Items:            []map[string]types.AttributeValue{},
		ScannedCount:     1,
	}

	suite.mockDynamodbClient.EXPECT().Query(suite.context, query).Return(output, nil).Times(1)

	_, err := suite.repository.GetUser(suite.context, email)

	suite.NotNil(err)
	suite.Equal(constants.UserNotFoundError, err)
}
