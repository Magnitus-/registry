package providercache

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/opentofu/registry/internal/providers"
	"golang.org/x/exp/slog"
)

type Handler struct {
	TableName *string
	Client    *dynamodb.Client
}

func NewHandler(awsConfig aws.Config, tableName string) *Handler {
	ddbClient := dynamodb.NewFromConfig(awsConfig)

	return &Handler{
		TableName: aws.String(tableName),
		Client:    ddbClient,
	}
}

const AllowedAge = (1 * time.Hour) - (5 * time.Minute) //nolint:gomnd // 55 minutes

type VersionListingItem struct {
	Provider    string              `dynamodbav:"provider"`
	Versions    []providers.Version `dynamodbav:"versions"`
	LastUpdated time.Time           `dynamodbav:"last_updated"`
}

func (p *Handler) Store(ctx context.Context, key string, versions []providers.Version) error {
	item := VersionListingItem{
		Provider:    key,
		Versions:    versions,
		LastUpdated: time.Now(),
	}

	marshalledItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		slog.Error("got error marshalling dynamodb item", "error", err)
		return fmt.Errorf("got error marshalling dynamodb item: %w", err)
	}

	putItemInput := &dynamodb.PutItemInput{
		Item:      marshalledItem,
		TableName: p.TableName,
	}

	slog.Info("Storing provider versions", "key", key, "versions", len(versions))
	_, err = p.Client.PutItem(ctx, putItemInput)
	if err != nil {
		slog.Error("got error calling PutItem", "error", err)
		return fmt.Errorf("got error calling PutItem: %w", err)
	}

	slog.Info("Successfully stored provider versions", "key", key, "versions", len(versions))
	return nil
}
