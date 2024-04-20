package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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
	timeoutSeconds           = 10
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

	// Cancel the context after timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer func() {
		fmt.Println("Canceling context")
		cancel()
	}()
	go func() {
		<-ctx.Done()
		fmt.Println("Context cancelled")
	}()

	// Print the request data
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

	// Get the request body
	prompt := req.Body

	// Parse the query string
	cfgScaleF, _ := strconv.ParseFloat(req.QueryStringParameters["cfg_scale"], 64)
	seed, _ := strconv.Atoi(req.QueryStringParameters["seed"])
	steps, _ := strconv.Atoi(req.QueryStringParameters["steps"])
	width, _ := strconv.Atoi(req.QueryStringParameters["width"])
	height, _ := strconv.Atoi(req.QueryStringParameters["height"])

	// Create the payload
	payload := BedrockRequestPayload{
		TextPrompts: []TextPrompt{{Text: prompt}},
		CfgScale:    cfgScaleF,
		Seed:        seed,
		Steps:       steps,
		Width:       width,
		Height:      height,
	}

	// Initialize zero values
	payload.Init()

	// Validate the payload
	if err := payload.Validate(); err != nil {
		fmt.Println("failed to validate payload\n", err)

		return events.LambdaFunctionURLResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	// Marshall the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("failed to marshall json\n", err)
	}

	// Print the payload for debugging
	fmt.Printf("Inference payload = %s.\n", string(payloadBytes))

	// Invoke the model
	result, err := bedrockSvc.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(stableDiffusionXLModelID),
		Body:        payloadBytes,
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		log.Fatal("failed to invoke model\n", err)
	}

	// Unmarshall the response
	var resp BedrockResponseBody
	err = json.Unmarshal(result.Body, &resp)
	if err != nil {
		log.Fatal("failed to unmarshal json\n", err)
	}

	// Get and return the image
	image := resp.Artifacts[0].Base64
	return events.LambdaFunctionURLResponse{Body: image, StatusCode: 200}, nil
}

// Validate the payload
func (p *BedrockRequestPayload) Validate() error {
	if p.CfgScale < 0 || p.CfgScale > 35 {
		return fmt.Errorf("cfg_scale must be between 0 and 35")
	}

	if p.Steps < 10 || p.Steps > 50 {
		return fmt.Errorf("steps must be between 10 and 50")
	}

	if p.Seed < 0 || p.Seed > 4294967295 {
		return fmt.Errorf("seed must be between 0 and 4294967295")
	}

	switch p.Height*1000 + p.Width {
	case 1025024, 1152806, 1216832, 1344768, 1536640, 641536, 769344, 833216, 897152:
		// do nothing
	default:
		return fmt.Errorf("width and height must be one of 1024x1024, 1152x896, 1216x832, 1344x768, 1536x640, 640x1536, 768x1344, 832x1216, 896x1152")
	}

	return nil
}

// Initialize payload as needed
func (p *BedrockRequestPayload) Init() BedrockRequestPayload {
	// set values if not already set
	if p.TextPrompts == nil {
		p.TextPrompts = make([]TextPrompt, 0)
	}
	if p.CfgScale == 0 {
		p.CfgScale = 7.0
	}
	if p.Steps == 0 {
		p.Steps = 20
	}
	if p.Seed == 0 {
		p.Seed = 0
	}
	if p.Width == 0 {
		p.Width = 1024
	}
	if p.Height == 0 {
		p.Height = 1024
	}

	return *p
}
func main() {
	lambda.Start(handler)
}

type BedrockRequestPayload struct {
	TextPrompts []TextPrompt `json:"text_prompts"`
	CfgScale    float64      `json:"cfg_scale"`
	Steps       int          `json:"steps"`
	Seed        int          `json:"seed"`
	Width       int          `json:"width"`
	Height      int          `json:"height"`
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
