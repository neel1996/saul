package migrations

import (
	"context"
	"core/log"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sirupsen/logrus"
)

type Migration interface {
	Setup()
	CreateUserTable(ctx context.Context)
	CreateDocumentHistoryTable(ctx context.Context)
}

type migration struct {
	client   *dynamodb.Client
	tableMap map[string]bool
}

var logger *logrus.Entry

func (m migration) Setup() {
	ctx := context.Background()
	logger = log.NewLogger(ctx)

	logger.Info("Setting up DB migration")

	m.CreateUserTable(ctx)
	m.CreateDocumentHistoryTable(ctx)

	logger.Info("Migration setup complete")
}

func NewMigration(client *dynamodb.Client) Migration {
	ctx := context.Background()
	tables, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		logger.Panic(err)
	}

	var tableMap = make(map[string]bool)
	for _, table := range tables.TableNames {
		tableMap[table] = true
	}

	return migration{
		client:   client,
		tableMap: tableMap,
	}
}
