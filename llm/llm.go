package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/samcaspus/gem-cli/iooperations"
	"github.com/samcaspus/gem-cli/utils"
	"google.golang.org/api/option"
)

func GetGeminiModel() (context.Context, *genai.GenerativeModel, *genai.Client) {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(iooperations.GetApiKey()))
	if err != nil {
		log.Fatal(err)
	}

	// The Gemini 1.5 models are versatile and work with most use cases
	model := client.GenerativeModel("gemini-1.5-flash")
	return ctx, model, client
}

func ExecuteCommand(retry int, args []string) {
	if retry == 0 {
		fmt.Println("Failed to generate content")
		return
	}
	executionQuery := utils.GetMergedStringArgs(args)
	ctx, model, client := GetGeminiModel()
	defer client.Close()

	// prompt to generate and get response as follows {'message':'<ai generated message on what the command will do>', 'command': 'actual command here'}
	prompt := "generate a command and fill it inside response for this json {'message':'<ai generated message on what the command will do>', 'command': 'actual command here'} and return only json back no additonal data. The query is " + executionQuery

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		ExecuteCommand(retry-1, args)
	}
	result := GetResponse(resp)
	fmt.Print(result.Message)
	fmt.Print("\n***************\n")
	fmt.Println(result.Command)
	fmt.Print("\n***************\n")

	fmt.Println("Do you want to execute the command? (y/n)")
	var input string
	fmt.Scanln(&input)
	if input == "y" {
		utils.ExecuteCommand(result.Command)
	}

}

func GetResponse(resp *genai.GenerateContentResponse) LlmResponse {
	var result LlmResponse
	fmt.Println("response ")
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			response_text := string(txt)
			response_text = strings.ReplaceAll(response_text, "```", "")
			response_text = strings.ReplaceAll(response_text, "json", "")
			response_text = strings.TrimSpace(response_text)
			if err := json.Unmarshal([]byte(response_text), &result); err != nil {
				log.Fatal(err)
			}
		}
	}

	return result
}
