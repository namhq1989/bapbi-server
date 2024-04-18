package openai

import (
	"strings"

	"github.com/goccy/go-json"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	oai "github.com/sashabaranov/go-openai"
)

type SearchTermResult struct {
	IsValid  bool           `json:"isValid"`
	Term     string         `json:"term"`
	From     TermByLanguage `json:"from"`
	To       TermByLanguage `json:"to"`
	Synonyms []string       `json:"synonyms"`
	Antonyms []string       `json:"antonyms"`
}

type TermByLanguage struct {
	Language   string `json:"language"`
	Definition string `json:"definition"`
	Example    string `json:"example"`
}

const searchTermPrompt = `
	Term is "{{term}}", its language is {{fromLanguage}} (source language) and you have to translate it to {{toLanguage}} (target language).
	Check if the term is a valid {{fromLanguage}} word or phrase. If valid, generate a detailed JSON-formatted (only data, no redundant information) response including:
	- "isValid": True
	- "term": The input word or phrase.
	- "from": Object with "language" is the source language, "definition" is the definition in the source language and "example" is an example sentence in the source language. All fields are mandatory.
	- "to": Object with "language" is the target language, "definition" is the definition in the target language and "example" is a translated sentence of source language's example. All fields are mandatory.
	- "synonyms": List max 3 synonyms of the term in source language, or provide an empty list if none are applicable.
	- "antonyms": List max 3 antonyms of the term in source language, or provide an empty list if none are applicable.
	If the term is not valid, return only: { "isValid": false }
`

func (o *OpenAI) SearchTerm(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*SearchTermResult, error) {
	prompt := strings.ReplaceAll(searchTermPrompt, "{{term}}", term)
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

	// Cleaning the input JSON string
	jsonStr := resp.Choices[0].Message.Content
	cleanJsonStr := strings.Replace(jsonStr, "`json\n", "", 1) // Removes the initial backtick, `json, and newline
	cleanJsonStr = strings.Replace(cleanJsonStr, "\n", "", -1) // Removes all newline characters
	cleanJsonStr = strings.Trim(cleanJsonStr, "`")             // Removes any remaining backticks

	var result SearchTermResult
	if err = json.Unmarshal([]byte(cleanJsonStr), &result); err != nil {
		ctx.Logger().Print("data", string(resp.Choices[0].Message.Content))
		return nil, err
	}

	return &result, nil
}
