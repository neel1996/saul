package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	dynamodbPkg "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/neel1996/saul/constants"
	"github.com/neel1996/saul/dynamodb"
	"github.com/neel1996/saul/log"
	"github.com/neel1996/saul/model/db"
)

type UserRepository interface {
	GetUser(ctx context.Context, email string) (db.User, error)
}

type userRepository struct {
	dynamodbClient dynamodb.DynamoDBClient
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
