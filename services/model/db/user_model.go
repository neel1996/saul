package db

type User struct {
	Id        string `dynamodbav:"id" validate:"required"`
	Name      string `dynamodbav:"name" validate:"required"`
	Username  string `dynamodbav:"username" validate:"required"`
	Email     string `dynamodbav:"email" validate:"required"`
	CreatedAt string `dynamodbav:"createdAt,omitempty"`
	UpdatedAt string `dynamodbav:"updatedAt,omitempty"`
}
