package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/solisamicus/occstructor/internal/config"
)

type LLMClient struct {
	client *openai.Client
	config *config.Config
}

func NewLLMClient(cfg *config.Config) *LLMClient {
	if !cfg.AI.Enabled {
		return nil
	}

	client := openai.NewClient(
		option.WithAPIKey(cfg.GetAPIKey()),
		option.WithBaseURL(cfg.AI.BaseURL),
	)

	return &LLMClient{
		client: client,
		config: cfg,
	}
}

func (l *LLMClient) MergeNamesWithLLM(namesText string) []string {
	if l == nil || l.client == nil {
		return l.fallbackProcessing(namesText)
	}

	pureChineseText := KeepOnlyChinese(namesText)
	if pureChineseText == "" {
		fmt.Println("Warning: input is empty after filtering non-Chinese characters")
		return nil
	}

	prompt := fmt.Sprintf(`Please standardize these Chinese job titles into a clean JSON array:
Rules:
1. Merge fragmented names into complete job titles
2. Split combined titles if they contain multiple independent jobs  
3. Each entry should be a complete, standalone job title
4. Keep only Chinese characters
5. Output valid JSON format: ["职业名称1", "职业名称2", ...]

Input text:
%s

Expected output: JSON array of standardized Chinese job titles`, pureChineseText)

	resp, err := l.client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are a job title processing expert, specializing in merging split Chinese job titles."),
				openai.UserMessage(prompt),
			}),
			Model:       openai.F(l.config.AI.Model),
			Temperature: openai.F(l.config.AI.Temperature),
		},
	)
	if err != nil {
		fmt.Println("API call failed:", err)
		return l.fallbackProcessing(namesText)
	}

	var result []string
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		fmt.Println("JSON parsing failed:", err)
		return l.fallbackProcessing(namesText)
	}

	var finalResults []string
	for _, name := range result {
		cleaned := KeepOnlyChinese(name)
		if cleaned != "" {
			finalResults = append(finalResults, cleaned)
		}
	}

	return finalResults
}

func (l *LLMClient) fallbackProcessing(namesText string) []string {
	lines := strings.Split(namesText, "\n")
	var names []string
	for _, line := range lines {
		line = KeepOnlyChinese(strings.TrimSpace(line))
		if line != "" {
			names = append(names, line)
		}
	}
	return names
}
