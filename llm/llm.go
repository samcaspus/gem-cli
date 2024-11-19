package llm

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/google/generative-ai-go/genai"
    "github.com/samcaspus/gem-cli/iooperations"
    "github.com/samcaspus/gem-cli/utils"
    "google.golang.org/api/option"
)

// Externalize command query
const commandQuery = "Write a command for the query for mac which will be a terminal based command on execution will generate the output as expected and only send the command as response, make sure not to send anything else."

func GetGeminiModel() (context.Context, *genai.GenerativeModel, *genai.Client, error) {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(iooperations.GetApiKey()))
    if err != nil {
        return nil, nil, nil, err
    }

    model := client.GenerativeModel("gemini-1.5-flash")
    return ctx, model, client, nil
}

func ExecuteCommand(retry int, args []string) {
    if retry <= 0 {
        fmt.Println("Failed to generate content after multiple attempts")
        return
    }

    executionQuery := utils.GetMergedStringArgs(args)
    ctx, model, client, err := GetGeminiModel()
    if err != nil {
        log.Printf("Error getting Gemini model: %v", err)
        return
    }
    defer client.Close()

    resp, err := model.GenerateContent(ctx, genai.Text(commandQuery))
    if err != nil {
        time.Sleep(2 * time.Second) // Adding delay between retries
        ExecuteCommand(retry-1, args)
        return
    }

    result := GetResponse(resp)
    fmt.Print(result.Message)
    fmt.Print("\n***************\n")
    fmt.Println(result.Command)
    fmt.Print("\n***************\n")

    if confirmExecution() {
        utils.ExecuteCommand(result.Command)
    }
}

func GetResponse(resp *genai.GenerateContentResponse) LlmResponse {
    var result LlmResponse
    for _, part := range resp.Candidates[0].Content.Parts {
        if txt, ok := part.(genai.Text); ok {
            if err := json.Unmarshal([]byte(txt), &result); err != nil {
                log.Printf("Error unmarshaling response: %v", err)
                return LlmResponse{}
            }
        }
    }
    return result
}

func confirmExecution() bool {
    fmt.Println("Do you want to execute the command? (y/n)")
    var input string
    fmt.Scanln(&input)
    return input == "y"
}
