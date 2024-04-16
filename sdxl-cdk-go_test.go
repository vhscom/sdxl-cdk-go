package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestSdxlCdkGoStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewSdxlCdkGoStack(app, "MyStack", nil)

	// THEN
	template := assertions.Template_FromStack(stack, nil)

	template.HasResourceProperties(jsii.String("AWS::Lambda::Url"), map[string]interface{}{
		"AuthType": "NONE",
		"Cors": map[string]interface{}{
			"AllowHeaders": []string{allowedHeaders},
			"AllowOrigins": []string{allowedOrigins},
		},
	})

	template.HasResourceProperties(jsii.String("AWS::S3::Bucket"), map[string]interface{}{
		"AutoDeleteObjects": nil,
		"BucketName":        imageBucketName,
	})

	template.HasResourceProperties(jsii.String("AWS::Lambda::Function"), map[string]interface{}{
		"Environment": map[string]interface{}{
			"Variables": map[string]interface{}{
				"BUCKET_NAME": imageBucketName,
			},
		},
	})

	template.HasResourceProperties(jsii.String("AWS::Lambda::Function"), map[string]interface{}{
		"Timeout": 120,
	})
}
