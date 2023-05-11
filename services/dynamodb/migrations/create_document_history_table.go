package migrations

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (m migration) CreateDocumentHistoryTable(ctx context.Context) {
	if m.tableMap[DocumentHistory] {
		logger.Infof("%s table already exists", DocumentHistory)
		return
	}

	table, err := m.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName:                 aws.String(DocumentHistory),
		DeletionProtectionEnabled: aws.Bool(true),
		BillingMode:               types.BillingModePayPerRequest,
		TableClass:                types.TableClassStandard,
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("document_id"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("created_at"),
				KeyType:       types.KeyTypeRange,
			},
		},
		LocalSecondaryIndexes: []types.LocalSecondaryIndex{
			{
				IndexName: aws.String("document_id-created_at-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("document_id"),
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
				AttributeName: aws.String("document_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("created_at"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
	})

	if err != nil {
		logger.Errorf("Error creating table %s : %v", DocumentHistory, err)
		return
	}

	logger.Infof("Created table %s successfully", *table.TableDescription.TableName)
}
