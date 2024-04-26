package openai

import (
	"fmt"
	"strings"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"

	"github.com/goccy/go-json"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	oai "github.com/sashabaranov/go-openai"
)

var writingAnalysisTopics = []string{
	"Weather Trends",
	"Population Growth",
	"Economic Indicators",
	"Healthcare Statistics",
	"Environmental Conservation",
	"Sports Performance",
	"Educational Outcomes",
	"Technology Adoption",
	"Travel and Tourism",
	"Consumer Behavior",
	"Public Transportation Usage",
	"Real Estate Markets",
	"Energy Consumption",
	"Cultural Participation",
	"Social Media Trends",
}

type WritingExerciseResult struct {
	Topic      string              `json:"topic"`
	Question   string              `json:"question"`
	Vocabulary []string            `json:"vocabulary"`
	Data       WritingExerciseData `json:"data,omitempty"`
}

type WritingExerciseData struct {
	TableHeader []string   `json:"tableHeader,omitempty"`
	TableData   [][]string `json:"tableData,omitempty"`
}

var writingExercisePrompt = map[string]string{
	domain.WritingExerciseTypeBasic.String(): `
		{{timestamp}}
		Generate a writing topic for {{level}} (user's level) {{language}} learners that includes a JSON-formatted data on topic "{{topic}}".
		Suggest some vocabulary words relevant to the analysis of this data. Provide a specific question for analyzing this table.
		Format the complete response in JSON: { 'topic': '<brief topic description>', 'vocabulary': [random 3 to 6 related words and user's level], 'question': '<question related to the topic>'}.
		Ensure the task is suitable for enhancing critical thinking and language skills.
	`,
	domain.WritingExerciseTypeAnalyze.String(): `
		{{timestamp}}
		Generate a writing topic for {{level}} (user's level) {{language}} learners that includes a JSON-formatted hypothetical table of data on topic "{{topic}}".
		The table should include a 'tableHeader' representing column names and 'tableData' consisting of 5 to 10 rows, each represented as an array of values (display as strings) corresponding to the headers.
		Suggest some vocabulary words relevant to the analysis of this data. Provide a specific question for analyzing this table.
		Format the complete response in JSON: { 'topic': '<brief topic description>', 'data': {'tableHeader': [list of column names], 'tableData': [random 5 to 10 values]}, 'vocabulary': [random 3 to 6 related words and user's level], 'question': '<analysis question>'}.
		Ensure the task is suitable for enhancing critical thinking and language skills.
	`,
}

func (o *OpenAI) WritingExercise(ctx *appcontext.AppContext, language, exType, level string) (*WritingExerciseResult, error) {
	prompt := writingExercisePrompt[exType]
	if prompt == "" {
		return nil, apperrors.Language.InvalidWritingExerciseData
	}
	prompt = strings.ReplaceAll(prompt, "{{timestamp}}", fmt.Sprintf("%d", time.Now().Unix()))

	// random a topic
	topic := writingAnalysisTopics[manipulation.RandomIntInRange(0, len(writingAnalysisTopics))]

	// prepare prompt
	if exType == domain.WritingExerciseTypeBasic.String() {
		prompt = strings.ReplaceAll(prompt, "{{language}}", language)
		prompt = strings.ReplaceAll(prompt, "{{level}}", level)
		prompt = strings.ReplaceAll(prompt, "{{topic}}", topic)
	} else if exType == domain.WritingExerciseTypeAnalyze.String() {
		prompt = strings.ReplaceAll(prompt, "{{language}}", language)
		prompt = strings.ReplaceAll(prompt, "{{level}}", level)
		prompt = strings.ReplaceAll(prompt, "{{topic}}", topic)
	}

	resp, err := o.client.CreateChatCompletion(ctx.Context(), oai.ChatCompletionRequest{
		Model:       oai.GPT3Dot5Turbo1106,
		Messages:    []oai.ChatCompletionMessage{{Role: oai.ChatMessageRoleUser, Content: prompt}},
		MaxTokens:   500,
		Temperature: 0.9,
	})

	if err != nil {
		return nil, err
	}

	var result WritingExerciseResult
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		ctx.Logger().Print("data", resp.Choices[0].Message.Content)
		return nil, err
	}

	return &result, nil
}
