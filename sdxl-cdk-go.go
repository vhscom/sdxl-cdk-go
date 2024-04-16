package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"

	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	imageBucketName = "sdxl-cdk-go-generated-images"
	functionDir     = "function"
	allowedOrigins  = "*"
	allowedHeaders  = "Accept,Authorization,Content-Type,X-Amz-Date,X-Amz-Security-Token,X-Api-Key,X-Amz-Signature,X-Requested-With"
)

type SdxlCdkGoStackProps struct {
	awscdk.StackProps
}

func NewSdxlCdkGoStack(scope constructs.Construct, id string, props *SdxlCdkGoStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	awss3.NewBucket(stack, jsii.String("SdxlCdkGoBucket"), &awss3.BucketProps{
		BucketName:    jsii.String(imageBucketName),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("SdxlCdkGoLambda"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Entry:   jsii.String(functionDir),
		Environment: &map[string]*string{
			"BUCKET_NAME": jsii.String(imageBucketName),
		},
		Timeout: awscdk.Duration_Seconds(jsii.Number(120)),
	}).AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
		Cors: &awslambda.FunctionUrlCorsOptions{
			AllowedOrigins: jsii.Strings(allowedOrigins),
			AllowedHeaders: jsii.Strings(allowedHeaders),
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String("SdxlCdkGoLambdaFunctionUrl"), &awscdk.CfnOutputProps{
		Value: function.Url(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewSdxlCdkGoStack(app, "SdxlCdkGoStack", &SdxlCdkGoStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
