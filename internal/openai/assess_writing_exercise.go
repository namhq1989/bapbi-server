package openai

import (
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	oai "github.com/sashabaranov/go-openai"
)

type AssessWritingExerciseResult struct {
	IsTopicRelevance bool     `json:"isTopicRelevance"`
	Score            int      `json:"score"`
	Improvement      []string `json:"improvement"`
	Comment          string   `json:"comment"`
}

const assessWritingExercisePrompt = `
	{{timestamp}}
	Language: {{language}}
	Topic: {{topic}}
	Level: {{level}}
	Paragraph: {{content}}
	
	Please assess the following:
	1. Topic Relevance: Determine if the paragraph is relevant to the given topic.
	2. Quality Score: Based on the specified level, evaluate the writing quality of the paragraph and assign a score from 0 to 10.
	3. Suggestions for Improvement: Provide up to 3 specific recommendations to enhance the clarity, coherence, grammar, style, and relevance of the paragraph.
       Let it empty if the score is larger than 7.
	4. Critical Assessment: Provide a short critical evaluation of the paragraph, discussing its strengths and weaknesses
	
	Output the assessment in JSON format as follows:
	{
	  "isTopicRelevance": [true/false based on assessment],
	  "score": [0 to 10 based on quality and level],
	  "improvement": [list of suggestions],
	  "comment": [critical evaluation]
	}
`

func (o *OpenAI) AssessWritingExercise(ctx *appcontext.AppContext, language, topic, level, content string) (*AssessWritingExerciseResult, error) {
	prompt := strings.ReplaceAll(assessWritingExercisePrompt, "{{language}}", language)
	prompt = strings.ReplaceAll(prompt, "{{topic}}", topic)
	prompt = strings.ReplaceAll(prompt, "{{level}}", level)
	prompt = strings.ReplaceAll(prompt, "{{content}}", content)
	prompt = strings.ReplaceAll(prompt, "{{timestamp}}", fmt.Sprintf("%d", time.Now().Unix()))

	resp, err := o.client.CreateChatCompletion(ctx.Context(), oai.ChatCompletionRequest{
		Model:       oai.GPT3Dot5Turbo1106,
		Messages:    []oai.ChatCompletionMessage{{Role: oai.ChatMessageRoleUser, Content: prompt}},
		MaxTokens:   500,
		Temperature: 0.3,
	})

	if err != nil {
		return nil, err
	}

	var result AssessWritingExerciseResult
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		ctx.Logger().Print("data", resp.Choices[0].Message.Content)
		return nil, err
	}

	return &result, nil
}
