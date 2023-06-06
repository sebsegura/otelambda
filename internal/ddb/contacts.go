package ddb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"os"
	"sebsegura/otelambda/internal/models"
	"sebsegura/otelambda/internal/session"
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
	ctx, span := instr.NewSpan(ctx, "CreateContact")
	defer span.End()

	log := logger.GetLogger(ctx)
	item, err := attributevalue.MarshalMap(contact)
	if err != nil {
		return err
	}

	span.SetAttributes(
		attribute.String("item.pk", contact.ID))

	_, err = r.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String(r.tableName),
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	})
	log.
		WithField("item", contact).
		WithField("span.id", span.SpanContext().SpanID()).
		Debug("inserted new item")

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	return err
}

func (r *repository) Update(ctx context.Context, contact *models.Contact) error {
	ctx, span := instr.NewSpan(ctx, "UpdateContact")
	defer span.End()

	log := logger.GetLogger(ctx)

	upd := expression.Set(expression.Name("lastName"), expression.Value("processed"))
	expr, err := expression.NewBuilder().WithUpdate(upd).Build()

	span.SetAttributes(
		attribute.String("expression", *expr.Update()),
		attribute.String("pk", contact.ID))

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
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	log.
		WithField("id", contact.ID).
		WithField("span.id", span.SpanContext().SpanID()).
		Debug("contact updated")

	return err
}
