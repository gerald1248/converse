package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
	systemPrompt = "You are a coding assistant. Output only code without explanations, comments, or any other text. Do not wrap the code in markdown code blocks."
	modelID	  = "anthropic.claude-3-sonnet-20240229-v1:0"
)

func main() {
	// Parse command line arguments
	flag.Parse()
	prompt := strings.Join(flag.Args(), " ")

	if prompt == "" {
		fmt.Fprintln(os.Stderr, "Error: No prompt provided")
		fmt.Fprintln(os.Stderr, "Usage: converse <prompt>")
		os.Exit(1)
	}

	// Configure AWS SDK
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading AWS configuration: %v\n", err)
		os.Exit(1)
	}

	// Create Bedrock client
	client := bedrockruntime.NewFromConfig(cfg)

	// Call Claude via AWS Bedrock
	response, err := callClaude(client, prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling Claude: %v\n", err)
		os.Exit(1)
	}

	// Write response to stdout
	fmt.Print(response)
}

func callClaude(client *bedrockruntime.Client, prompt string) (string, error) {
	// Create the request payload
	payload := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":		4096 * 16,
		"messages": []map[string]interface{}{
			{
				"role":	"user",
				"content": []map[string]string{
					{
						"type": "text",
						"text": systemPrompt + "\n\n" + prompt,
					},
				},
			},
		},
	}
	// Convert payload to JSON
	payloadBytes, err := marshalJSON(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	// Create the API request
	input := &bedrockruntime.InvokeModelInput{
		ModelId:	 aws.String(modelID),
		ContentType: aws.String("application/json"),
		Body:		payloadBytes,
	}

	// Call the API
	resp, err := client.InvokeModel(context.Background(), input)
	if err != nil {
		return "", fmt.Errorf("error invoking model: %w", err)
	}

	// Parse the response
	var responseBody map[string]interface{}
	if err := unmarshalJSON(bytes.NewReader(resp.Body), &responseBody); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Extract the content from the response
	messages, ok := responseBody["content"].([]interface{})
	if !ok || len(messages) == 0 {
		return "", fmt.Errorf("unexpected response format")
	}

	var result strings.Builder
	for _, msg := range messages {
		msgMap, ok := msg.(map[string]interface{})
		if !ok {
			continue
		}

		if text, ok := msgMap["text"].(string); ok {
			result.WriteString(text)
		}
	}

	return result.String(), nil
}

// Helper functions for JSON marshaling/unmarshaling
func marshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func unmarshalJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
