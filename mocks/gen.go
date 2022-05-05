package mocks

//go:generate go run github.com/golang/mock/mockgen -destination awsmocks/s3.go -package awsmocks github.com/aws/aws-sdk-go/service/s3/s3iface S3API
//go:generate go run github.com/golang/mock/mockgen -destination awsmocks/ddb.go -package awsmocks github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface DynamoDBAPI
