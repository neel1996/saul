package migrations

import (
	"context"
	"core/constants"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (m migration) CreateUserTable(ctx context.Context) {
	if m.tableMap[constants.UserTableName] {
		logger.Infof("%s table already exists", constants.UserTableName)
		return
	}

	table, err := m.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName:                 aws.String(constants.UserTableName),
		DeletionProtectionEnabled: aws.Bool(true),
		BillingMode:               types.BillingModePayPerRequest,
		TableClass:                types.TableClassStandard,
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("email"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("created_at"),
				KeyType:       types.KeyTypeRange,
			},
		},
		LocalSecondaryIndexes: []types.LocalSecondaryIndex{
			{
				IndexName: aws.String("email-created_at-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("email"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("created_at"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("email"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("created_at"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
	})

	if err != nil {
		logger.Errorf("Error creating table %s : %v", constants.UserTableName, err)
		return
	}

	logger.Infof("Created table %s successfully", *table.TableDescription.TableName)
}
