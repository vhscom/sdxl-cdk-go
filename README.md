# SDXL CDK Go

This is a CDK project for Lambda development with Go.

The `cdk.json` file tells the CDK toolkit how to execute your app.

## Useful commands

Accounts and regions

- `aws sso login` login to AWS using Identity Center
- `aws sts get-caller-identity` get your AWS identity
- `aws configure sso` configure your AWS credentials
- `aws configure list` list your AWS credentials

Deployment and testing

- `cdk deploy` deploy this stack to your default AWS account/region
- `cdk diff` compare deployed stack with current state
- `cdk synth` emits the synthesized CloudFormation template
- `go test` run unit tests

Development

- `go mod download` download dependencies
- `go run sdxl-cdk-go.go` bundle handler
- `go mod tidy` update dependencies
- `go mod verify` verify dependencies
