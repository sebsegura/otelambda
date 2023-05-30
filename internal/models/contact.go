package models

type Contact struct {
	ID        string `dynamodbav:"id" json:"id"`
	FirstName string `dynamodbav:"firstName" json:"first_name"`
}
