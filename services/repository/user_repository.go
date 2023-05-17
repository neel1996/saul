package repository

//go:generate mockgen -destination=../mocks/mock_user_repository.go -package=mocks -source=user_repository.go

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	dynamodbPkg "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/neel1996/saul/constants"
	"github.com/neel1996/saul/dynamodb"
	"github.com/neel1996/saul/log"
	"github.com/neel1996/saul/model/db"
	"github.com/neel1996/saul/model/request"
	"time"
)

type UserRepository interface {
	GetUser(ctx context.Context, email string) (db.User, error)
	CreateUser(ctx context.Context, request request.UserLoginRequest) error
	DoesUserExist(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	dynamodbClient dynamodb.DynamoDBClient
}

func (repository userRepository) DoesUserExist(ctx context.Context, email string) (bool, error) {
	logger := log.NewLogger(ctx)
	logger.Infof("Checking if user with email: %s exists", email)

	query := dynamodbPkg.QueryInput{
		TableName: aws.String(constants.UserTableName),
		ExpressionAttributeNames: map[string]string{
			"#email": "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{
				Value: email,
			},
		},
		KeyConditionExpression: aws.String("#email = :email"),
		Limit:                  aws.Int32(1),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		Select:                 types.SelectCount,
	}

	output, err := repository.dynamodbClient.Query(ctx, &query)
	if err != nil {
		logger.Errorf("Error occurred while querying for user with email: %s, error: %v", email, err)
		return false, err
	}

	if output.Count == 0 {
		logger.Infof("User with email: %s does not exist", email)
		return false, nil
	}

	return true, nil
}

func (repository userRepository) CreateUser(ctx context.Context, userRequest request.UserLoginRequest) error {
	logger := log.NewLogger(ctx)
	logger.Infof("Creating new user with email: %s", userRequest.Email)

	input := &dynamodbPkg.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName: aws.String(constants.UserTableName),
					Item: map[string]types.AttributeValue{
						"email": &types.AttributeValueMemberS{
							Value: userRequest.Email,
						},
						"name": &types.AttributeValueMemberS{
							Value: userRequest.Name,
						},
						"avatar": &types.AttributeValueMemberS{
							Value: userRequest.Avatar,
						},
						"user_id": &types.AttributeValueMemberS{
							Value: userRequest.UserId,
						},
						"created_at": &types.AttributeValueMemberS{
							Value: time.Now().String(),
						},
					},
				},
			},
		},
		ClientRequestToken:          aws.String(uuid.NewString()),
		ReturnConsumedCapacity:      types.ReturnConsumedCapacityTotal,
		ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
	}

	_, err := repository.dynamodbClient.TransactWriteItems(ctx, input)
	if err != nil {
		logger.Errorf("Error creating user in DB %v", err)
		return err
	}

	logger.Infof("Successfully created user with email: %s", userRequest.Email)
	return nil
}

func (repository userRepository) GetUser(ctx context.Context, email string) (db.User, error) {
	logger := log.NewLogger(ctx)
	logger.Info("Getting user from DB")

	query := &dynamodbPkg.QueryInput{
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

	output, err := repository.dynamodbClient.Query(ctx, query)
	if err != nil {
		logger.Errorf("Error getting user from DB %v", err)
		return db.User{}, err
	}

	var user []db.User
	err = attributevalue.UnmarshalListOfMaps(output.Items, &user)
	if err != nil {
		logger.Errorf("Error unmarshalling user from DB %v", err)
		return db.User{}, err
	}

	if len(user) == 0 || user[0].Email != email {
		logger.Errorf("User with email %s not found in DB", email)
		return db.User{}, constants.UserNotFoundError
	}

	return user[0], nil
}

func NewUserRepository(dynamodbClient dynamodb.DynamoDBClient) UserRepository {
	return userRepository{
		dynamodbClient,
	}
}
