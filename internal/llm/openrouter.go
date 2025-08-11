package llm

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
)

type Message struct {
    Role       string `json:"role"`
    Content    string `json:"content,omitempty"`
    ToolCallID string `json:"tool_call_id,omitempty"`
    Name       string `json:"name,omitempty"`
}

type Tool struct {
    Type     string        `json:"type"` // "function"
    Function ToolFunction  `json:"function"`
}

type ToolFunction struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description,omitempty"`
    Parameters  map[string]any         `json:"parameters"`
}

type ChatRequest struct {
    Model       string    `json:"model"`
    Temperature float64   `json:"temperature,omitempty"`
    Messages    []Message `json:"messages"`
    Tools       []Tool    `json:"tools,omitempty"`
    ToolChoice  string    `json:"tool_choice,omitempty"` // "auto"
}

type ChatResponse struct {
    Choices []struct {
        Message struct {
            Role      string `json:"role"`
            Content   string `json:"content"`
            ToolCalls []struct {
                ID       string `json:"id"`
                Type     string `json:"type"`
                Function struct {
                    Name      string `json:"name"`
                    Arguments string `json:"arguments"`
                } `json:"function"`
            } `json:"tool_calls"`
        } `json:"message"`
        FinishReason string `json:"finish_reason"`
    } `json:"choices"`
}

type Client struct {
    httpClient *http.Client
    baseURL    string
    apiKey     string
    model      string
}

func NewClient() (*Client, error) {
    key := os.Getenv("OPENROUTER_API_KEY")
    if key == "" {
        return nil, fmt.Errorf("OPENROUTER_API_KEY not set")
    }
    base := os.Getenv("OPENROUTER_BASE_URL")
    if base == "" { base = "https://openrouter.ai/api/v1" }
    model := os.Getenv("OPENROUTER_MODEL")
    if model == "" { model = "anthropic/claude-3-opus" }
    return &Client{
        httpClient: &http.Client{Timeout: 30 * time.Second},
        baseURL:    base,
        apiKey:     key,
        model:      model,
    }, nil
}

func (c *Client) Chat(ctx context.Context, messages []Message, tools []Tool) (ChatResponse, error) {
    reqBody := ChatRequest{
        Model:       c.model,
        Temperature: 0.1,
        Messages:    messages,
        Tools:       tools,
        ToolChoice:  "auto",
    }
    b, err := json.Marshal(reqBody)
    if err != nil { return ChatResponse{}, err }
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(b))
    if err != nil { return ChatResponse{}, err }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("X-Title", "Skybit CLI")

    resp, err := c.httpClient.Do(req)
    if err != nil { return ChatResponse{}, err }
    defer resp.Body.Close()
    if resp.StatusCode/100 != 2 {
        return ChatResponse{}, fmt.Errorf("openrouter http %d", resp.StatusCode)
    }
    var out ChatResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return ChatResponse{}, err }
    return out, nil
}


