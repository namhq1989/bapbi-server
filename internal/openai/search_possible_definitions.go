package openai

import (
	"strings"

	"github.com/goccy/go-json"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	oai "github.com/sashabaranov/go-openai"
)

type SearchPossibleDefinitionsResult struct {
	List []SearchPossibleDefinitionsResultItem `json:"list"`
}

type SearchPossibleDefinitionsResultItem struct {
	Definition string `json:"definition"`
	Pos        string `json:"pos"`
}

const searchPossibleDefinitionsPrompt = `
	Provide all possible translations for the term '{{term}}' from {{fromLanguage}} to {{toLanguage}}. Generate a detailed JSON-formatted response including:
	- "list": An array of translations. Each translation item should contain:
		- "definition": The translation to {{toLanguage}}.
		- "pos": The part of speech (in {{fromLanguage}}) of this translation.
`

func (o *OpenAI) SearchPossibleDefinitions(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*SearchPossibleDefinitionsResult, error) {
	prompt := strings.ReplaceAll(searchPossibleDefinitionsPrompt, "{{term}}", term)
	prompt = strings.ReplaceAll(prompt, "{{fromLanguage}}", fromLanguage)
	prompt = strings.ReplaceAll(prompt, "{{toLanguage}}", toLanguage)

	resp, err := o.client.CreateChatCompletion(ctx.Context(), oai.ChatCompletionRequest{
		Model:       oai.GPT3Dot5Turbo1106,
		Messages:    []oai.ChatCompletionMessage{{Role: oai.ChatMessageRoleUser, Content: prompt}},
		MaxTokens:   500,
		Temperature: 0.3,
	})

	if err != nil {
		return nil, err
	}

	// cleaning the input JSON string
	jsonStr := resp.Choices[0].Message.Content
	cleanJsonStr := strings.Replace(jsonStr, "`json\n", "", 1)
	cleanJsonStr = strings.Replace(cleanJsonStr, "\n", "", -1)
	cleanJsonStr = strings.Trim(cleanJsonStr, "`")

	var result SearchPossibleDefinitionsResult
	if err = json.Unmarshal([]byte(cleanJsonStr), &result); err != nil {
		ctx.Logger().Print("data", string(resp.Choices[0].Message.Content))
		return nil, err
	}

	return &result, nil
}
