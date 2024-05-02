package openai

import (
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	oai "github.com/sashabaranov/go-openai"
)

type AssessTermExerciseResult struct {
	IsVocabularyCorrect    bool                                      `json:"isVocabularyCorrect"`
	VocabularyIssue        string                                    `json:"vocabularyIssue"`
	IsTenseCorrect         bool                                      `json:"isTenseCorrect"`
	TenseIssue             string                                    `json:"tenseIssue"`
	GrammarIssues          []AssessTermExerciseGrammarIssue          `json:"grammarIssues"`
	ImprovementSuggestions []AssessTermExerciseImprovementSuggestion `json:"improvementSuggestions"`
}

type AssessTermExerciseGrammarIssue struct {
	Issue      string `json:"issue"`
	Correction string `json:"correction"`
}

type AssessTermExerciseImprovementSuggestion struct {
	Instruction string `json:"instruction"`
	Example     string `json:"example"`
}

const assessTermExercisePrompt = `
	{{timestamp}}
	Evaluate the following sentence for {{language}} proficiency and provide responses in a structured JSON-friendly (JSON only, don't be verbose) format:

	Sentence: "{{content}}"
	Target vocabulary: "{{vocabulary}}"
	Target tense: "{{tense}}"
  
	Please focus on:
	- The correct usage of the target vocabulary '{{vocabulary}}' in any grammatical form.
	- Checking if the sentence uses the '{{tense}}' tense correctly.
	- Review and correct all aspects of grammar and spelling
	
	Required Output Format:
	{
	  "isVocabularyCorrect": [true/false],
	  "vocabularyIssue": [""/"reason if vocabulary is incorrectly used"],
	  "isTenseCorrect": [true/false],
	  "tenseIssue": [""/"reason if the incorrect tense is used"],
	  "grammarIssues": [
		{
		  "issue": "specific grammar issue",
		  "correction": "suggested correction"
		}
	  ],
	  "improvementSuggestions": [
		{
		  "instruction": "specific improvement instruction",
		  "example": "how the sentence should be corrected"
		}
	  ]
	}
`

func (o *OpenAI) AssessTermExercise(ctx *appcontext.AppContext, language, term, tense, content string) (*AssessTermExerciseResult, error) {
	prompt := strings.ReplaceAll(assessTermExercisePrompt, "{{language}}", language)
	prompt = strings.ReplaceAll(prompt, "{{vocabulary}}", term)
	prompt = strings.ReplaceAll(prompt, "{{tense}}", tense)
	prompt = strings.ReplaceAll(prompt, "{{content}}", content)
	prompt = strings.ReplaceAll(prompt, "{{timestamp}}", fmt.Sprintf("%d", time.Now().Unix()))

	resp, err := o.client.CreateChatCompletion(ctx.Context(), oai.ChatCompletionRequest{
		Model:       oai.GPT4,
		Messages:    []oai.ChatCompletionMessage{{Role: oai.ChatMessageRoleUser, Content: prompt}},
		MaxTokens:   400,
		Temperature: 0.2,
	})

	if err != nil {
		return nil, err
	}

	var result AssessTermExerciseResult
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		ctx.Logger().Print("cannot unmarshal", resp.Choices[0].Message.Content)
		return nil, err
	}

	return &result, nil
}
