package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

var bedrockSvc *bedrockruntime.Client

const (
	stableDiffusionXLModelID = "stability.stable-diffusion-xl-v1"
	defaultRegion            = "us-east-1"
)

func init() {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultRegion
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	bedrockSvc = bedrockruntime.NewFromConfig(cfg)
}

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", req.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(req.Body))

	var bucketname = os.Getenv("BUCKET_NAME")
	fmt.Printf("Bucket name = %s.\n", bucketname)

	fmt.Println("Headers:")
	for key, value := range req.Headers {
		fmt.Printf(" ++ %s: %s\n", key, value)
	}

	fmt.Println("Query string parameters:")
	for key, value := range req.QueryStringParameters {
		fmt.Printf(" && %s: %s\n", key, value)
	}

	prompt := req.Body
	cfgScaleF, _ := strconv.ParseFloat(req.QueryStringParameters["cfg_scale"], 64)
	seed, _ := strconv.Atoi(req.QueryStringParameters["seed"])
	steps, _ := strconv.Atoi(req.QueryStringParameters["steps"])

	payload := BedrockRequestPayload{
		TextPrompts: []TextPrompt{{Text: prompt}},
		CfgScale:    cfgScaleF,
		Steps:       steps,
	}

	if seed > 0 {
		payload.Seed = seed
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("failed to marshall json\n", err)
	}

	fmt.Printf("Payload = %s.\n", string(payloadBytes))

	result, err := bedrockSvc.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(stableDiffusionXLModelID),
		Body:        payloadBytes,
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		log.Fatal("failed to invoke model\n", err)
	}

	var resp BedrockResponseBody
	err = json.Unmarshal(result.Body, &resp)
	if err != nil {
		log.Fatal("failed to unmarshal json\n", err)
	}

	image := resp.Artifacts[0].Base64

	return events.LambdaFunctionURLResponse{Body: image, StatusCode: 200}, nil
}

func main() {
	lambda.Start(handler)
}

type BedrockRequestPayload struct {
	TextPrompts []TextPrompt `json:"text_prompts"`
	CfgScale    float64      `json:"cfg_scale"`
	Steps       int          `json:"steps"`
	Seed        int          `json:"seed"`
}

type BedrockResponseBody struct {
	Result    string     `json:"result"`
	Artifacts []Artifact `json:"artifacts"`
}

type TextPrompt struct {
	Text string `json:"text"`
}

type Artifact struct {
	Base64       string `json:"base64"`
	FinishReason string `json:"finishReason"`
}
