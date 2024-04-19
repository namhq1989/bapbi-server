package openai

import (
	"strings"

	"github.com/goccy/go-json"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	oai "github.com/sashabaranov/go-openai"
)

type SearchSemanticRelationsResult struct {
	Synonyms []string `json:"synonyms"`
	Antonyms []string `json:"antonyms"`
}

const searchSemanticRelationsPrompt = `
	Term is "{{term}}", its language is {{language}}.
	Generate a detailed JSON-formatted (only data, no redundant information) response including:
	- "synonyms": List of 3 synonyms of the term, or provide an empty list if none are applicable.
	- "antonyms": List of 3 antonyms of the term, or provide an empty list if none are applicable.
	If the term is not valid, return only: { }
`

func (o *OpenAI) SearchSemanticRelations(ctx *appcontext.AppContext, term, language string) (*SearchSemanticRelationsResult, error) {
	prompt := strings.ReplaceAll(searchSemanticRelationsPrompt, "{{term}}", term)
	prompt = strings.ReplaceAll(prompt, "{{language}}", language)

	resp, err := o.client.CreateChatCompletion(ctx.Context(), oai.ChatCompletionRequest{
		Model:       oai.GPT3Dot5Turbo1106,
		Messages:    []oai.ChatCompletionMessage{{Role: oai.ChatMessageRoleUser, Content: prompt}},
		MaxTokens:   300,
		Temperature: 0.8,
	})

	if err != nil {
		return nil, err
	}

	// cleaning the input JSON string
	jsonStr := resp.Choices[0].Message.Content
	cleanJsonStr := strings.Replace(jsonStr, "`json\n", "", 1)
	cleanJsonStr = strings.Replace(cleanJsonStr, "\n", "", -1)
	cleanJsonStr = strings.Trim(cleanJsonStr, "`")

	var result SearchSemanticRelationsResult
	if err = json.Unmarshal([]byte(cleanJsonStr), &result); err != nil {
		ctx.Logger().Print("data", string(resp.Choices[0].Message.Content))
		return nil, err
	}

	return &result, nil
}
