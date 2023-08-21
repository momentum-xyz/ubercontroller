package open_ai

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type OpenAI struct {
	client         *openai.Client
	tiktokenClient *tiktoken.Tiktoken

	logger      *zap.SugaredLogger
	apiKey      string
	temperature float32
	maxTokens   int
	model       string
}

func NewOpenAI(logger *zap.SugaredLogger, apiKey string, temperature float32, maxTokens int) *OpenAI {
	model := openai.GPT4

	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		log.Fatal(err)
	}

	client := openai.NewClient(apiKey)

	return &OpenAI{
		client:         client,
		tiktokenClient: tkm,
		apiKey:         apiKey,
		temperature:    temperature,
		maxTokens:      maxTokens,
		model:          model,
	}
}

func (o *OpenAI) CreateChatCompletion(text string) (string, error) {
	token := o.tiktokenClient.Encode(text, nil, nil)
	if len(token) > o.maxTokens {
		return "", errors.New("Max tokens exceeded")
	}

	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: o.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
		},
	)

	if err != nil {
		return "", errors.WithMessagef(err, "ChatCompletion error: %v\n")
	}

	return resp.Choices[0].Message.Content, nil
}
