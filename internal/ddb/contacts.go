package ddb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"os"
	"sebsegura/otelambda/internal/models"
	"sebsegura/otelambda/internal/session"
	"sebsegura/otelambda/pkg/caller"
	"sebsegura/otelambda/pkg/instr"
	"sebsegura/otelambda/pkg/logger"
)

type Repository interface {
	Create(ctx context.Context, contact *models.Contact) error
	Update(ctx context.Context, contact *models.Contact) error
}

type repository struct {
	ddb       *dynamodb.Client
	tableName string
}

func New(ctx context.Context) Repository {
	cfg := session.GetConfig(ctx)
	ddb := dynamodb.NewFromConfig(cfg)
	tableName := os.Getenv("TABLE_NAME")

	return &repository{
		ddb:       ddb,
		tableName: tableName,
	}
}

func (r *repository) Create(ctx context.Context, contact *models.Contact) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(instr.GetTracerName()).Start(ctx, "PutItem")
	defer span.End()

	log := logger.GetLogger(ctx)
	item, err := attributevalue.MarshalMap(contact)
	if err != nil {
		return err
	}

	_, err = r.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String(r.tableName),
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	})
	log.WithField("item", contact).Debug("inserted new item")

	if err != nil {
		span.RecordError(err)
	}

	span.SetAttributes(
		attribute.String("Table", r.tableName))

	return err
}

func (r *repository) Update(ctx context.Context, contact *models.Contact) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(instr.GetTracerName()).Start(ctx, caller.GetRunTimeCaller(2))
	defer span.End()

	log := logger.GetLogger(ctx)

	upd := expression.Set(expression.Name("lastName"), expression.Value("processed"))
	expr, err := expression.NewBuilder().WithUpdate(upd).Build()
	_, err = r.ddb.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: contact.ID},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		span.RecordError(err)
	}
	log.WithField("id", contact.ID).Debug("contact updated")

	return err
}
