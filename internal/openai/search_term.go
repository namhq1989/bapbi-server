package openai

import (
	"strings"

	"github.com/goccy/go-json"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	oai "github.com/sashabaranov/go-openai"
)

type SearchTermResult struct {
	From     TermByLanguage `json:"from"`
	To       TermByLanguage `json:"to"`
	Examples []TermExample  `json:"examples"`
}

type TermByLanguage struct {
	Language   string `json:"language"`
	Definition string `json:"definition"`
	Example    string `json:"example"`
}

type TermExample struct {
	PartOfSpeech string `json:"pos"`
	From         string `json:"from"`
	To           string `json:"to"`
}

const searchTermPrompt = `
	Term is "{{term}}", its language is {{fromLanguage}} and you have to translate it to {{toLanguage}}.
	Generate a detailed JSON-formatted structured as follows:
	- "from": { "language": "{{fromLanguage}}", "definition": "{{fromLanguage}} definition", "example": "{{fromLanguage}} example" }.
	- "to": {  "language": "{{toLanguage}}", "definition": "{{toLanguage}} definition", "example": "{{toLanguage}} translation" }.
	- "examples": [{ "pos": "part of speech", "from": "{{fromLanguage}} sentence", "to": "{{toLanguage}} translation"}].
	Field "examples" has up to three examples with each show a distinct usage in English. If distinct usages aren't available, provide fewer examples
`

func (o *OpenAI) SearchTerm(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*SearchTermResult, error) {
	prompt := strings.ReplaceAll(searchTermPrompt, "{{term}}", term)
	prompt = strings.ReplaceAll(prompt, "{{fromLanguage}}", fromLanguage)
	prompt = strings.ReplaceAll(prompt, "{{toLanguage}}", toLanguage)

	resp, err := o.client.CreateChatCompletion(ctx.Context(), oai.ChatCompletionRequest{
		Model:       oai.GPT3Dot5Turbo1106,
		Messages:    []oai.ChatCompletionMessage{{Role: oai.ChatMessageRoleUser, Content: prompt}},
		MaxTokens:   700,
		Temperature: 0.5,
	})

	if err != nil {
		return nil, err
	}

	var result SearchTermResult
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		ctx.Logger().Print("data", resp.Choices[0].Message.Content)
		return nil, err
	}

	return &result, nil
}
