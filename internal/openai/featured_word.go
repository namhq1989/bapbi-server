package openai

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	oai "github.com/sashabaranov/go-openai"
)

type FeaturedWordResult struct {
	Word string `json:"word"`
}

const featuredWordPrompt = `
	Using a random factor of {{timestamp}}, generate a random {{language}} vocabulary word.
	Generate a detailed JSON-formatted (only data, no redundant information) response including: { 'word': random word }.
`

func (o *OpenAI) FeaturedWord(ctx *appcontext.AppContext, language string) (*FeaturedWordResult, error) {
	prompt := strings.ReplaceAll(featuredWordPrompt, "{{language}}", language)
	prompt = strings.ReplaceAll(prompt, "{{timestamp}}", fmt.Sprintf("%d", time.Now().Unix()))

	// random int number
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := randSource.Intn(10000)

	resp, err := o.client.CreateChatCompletion(ctx.Context(), oai.ChatCompletionRequest{
		Model:       oai.GPT3Dot5Turbo1106,
		Messages:    []oai.ChatCompletionMessage{{Role: oai.ChatMessageRoleUser, Content: prompt}},
		MaxTokens:   100,
		Temperature: 0.9,
		Seed:        &seed,
	})

	if err != nil {
		return nil, err
	}

	var result FeaturedWordResult
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		ctx.Logger().Print("data", resp.Choices[0].Message.Content)
		return nil, err
	}

	return &result, nil
}
