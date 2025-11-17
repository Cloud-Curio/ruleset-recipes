/*
 * Data Ingestion TUI - AI Integration Client
 * 
 * This package provides AI integration capabilities using OpenRouter and other
 * AI providers. It includes RAG (Retrieval-Augmented Generation) support,
 * MCP (Model Context Protocol) server integration, and chat functionality.
 * 
 * Features:
 * - OpenRouter API integration
 * - Multiple AI model support
 * - RAG with vector embeddings
 * - MCP server communication
 * - Context-aware assistance
 * - Tool integration (Cline, Continue.dev, etc.)
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"data-ingestion-tui/internal/config"
	"data-ingestion-tui/internal/models"
)

// @decorator: Client
// @description: AI client for interacting with various AI providers
// @version: 1.0.0
// @author: Codegen AI Assistant

// Client represents the AI integration client
type Client struct {
	config     *config.AIConfig
	logger     *logrus.Logger
	openaiClient *openai.Client
	ragEngine  *RAGEngine
	mcpManager *MCPManager
	
	// Context and conversation management
	conversationHistory []models.ChatMessage
	systemPrompt        string
	
	// Tool integrations
	tools map[string]Tool
}

// @decorator: Tool
// @description: Interface for AI tools and integrations
type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
	Schema() map[string]interface{}
}

// @decorator: ChatMessage
// @description: Represents a chat message in the conversation
type ChatMessage struct {
	Role      string                 `json:"role"`
	Content   string                 `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// @decorator: ChatResponse
// @description: Response from AI chat completion
type ChatResponse struct {
	Message       string                 `json:"message"`
	Model         string                 `json:"model"`
	TokensUsed    int                    `json:"tokens_used"`
	ResponseTime  time.Duration          `json:"response_time"`
	Suggestions   []string               `json:"suggestions,omitempty"`
	ToolCalls     []ToolCall             `json:"tool_calls,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// @decorator: ToolCall
// @description: Represents a tool call from the AI
type ToolCall struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Args     map[string]interface{} `json:"args"`
	Response interface{}            `json:"response,omitempty"`
}

// @decorator: New
// @description: Create a new AI client instance
// @param config: AI configuration
// @param logger: Logger instance
// @return *Client: New AI client
// @return error: Any error that occurred during creation
func New(config *config.AIConfig, logger *logrus.Logger) (*Client, error) {
	if config == nil {
		return nil, errors.New("AI configuration cannot be nil")
	}
	
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	
	// Initialize OpenAI client for OpenRouter
	clientConfig := openai.DefaultConfig(config.APIKey)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}
	
	openaiClient := openai.NewClientWithConfig(clientConfig)
	
	// Create AI client
	client := &Client{
		config:              config,
		logger:              logger,
		openaiClient:        openaiClient,
		conversationHistory: make([]models.ChatMessage, 0),
		tools:               make(map[string]Tool),
		systemPrompt:        getDefaultSystemPrompt(),
	}
	
	// Initialize RAG engine if enabled
	if config.RAG.Enabled {
		ragEngine, err := NewRAGEngine(&config.RAG, logger)
		if err != nil {
			logger.WithError(err).Warn("Failed to initialize RAG engine, continuing without RAG")
		} else {
			client.ragEngine = ragEngine
		}
	}
	
	// Initialize MCP manager if enabled
	if config.MCP.Enabled {
		mcpManager, err := NewMCPManager(&config.MCP, logger)
		if err != nil {
			logger.WithError(err).Warn("Failed to initialize MCP manager, continuing without MCP")
		} else {
			client.mcpManager = mcpManager
		}
	}
	
	// Initialize built-in tools
	if err := client.initializeTools(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize AI tools")
	}
	
	logger.Info("AI client initialized successfully")
	return client, nil
}

// @decorator: Chat
// @description: Send a chat message and get AI response
// @param ctx: Context for cancellation
// @param message: User message
// @param options: Chat options
// @return *ChatResponse: AI response
// @return error: Any error that occurred during chat
func (c *Client) Chat(ctx context.Context, message string, options *ChatOptions) (*ChatResponse, error) {
	startTime := time.Now()
	
	// Validate input
	if strings.TrimSpace(message) == "" {
		return nil, errors.New("message cannot be empty")
	}
	
	// Set default options
	if options == nil {
		options = &ChatOptions{
			Model:       c.config.Model,
			Temperature: c.config.Temperature,
			MaxTokens:   c.config.MaxTokens,
		}
	}
	
	// Add user message to history
	userMsg := models.ChatMessage{
		Role:      "user",
		Content:   message,
		Timestamp: time.Now(),
	}
	c.conversationHistory = append(c.conversationHistory, userMsg)
	
	// Enhance message with RAG if available
	enhancedMessage := message
	if c.ragEngine != nil && options.UseRAG {
		ragContext, err := c.ragEngine.Query(ctx, message, 5)
		if err != nil {
			c.logger.WithError(err).Warn("Failed to get RAG context")
		} else if len(ragContext) > 0 {
			enhancedMessage = fmt.Sprintf("Context:\n%s\n\nQuestion: %s", 
				strings.Join(ragContext, "\n"), message)
		}
	}
	
	// Prepare messages for OpenAI API
	messages := c.buildMessages(enhancedMessage, options)
	
	// Create chat completion request
	req := openai.ChatCompletionRequest{
		Model:       options.Model,
		Messages:    messages,
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
		Stream:      options.Stream,
	}
	
	// Add tools if available
	if len(c.tools) > 0 && options.UseTools {
		req.Tools = c.buildToolDefinitions()
	}
	
	// Make API call
	resp, err := c.openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create chat completion")
	}
	
	// Process response
	if len(resp.Choices) == 0 {
		return nil, errors.New("no response choices returned")
	}
	
	choice := resp.Choices[0]
	responseTime := time.Since(startTime)
	
	// Handle tool calls
	var toolCalls []ToolCall
	if len(choice.Message.ToolCalls) > 0 {
		toolCalls, err = c.executeToolCalls(ctx, choice.Message.ToolCalls)
		if err != nil {
			c.logger.WithError(err).Warn("Failed to execute tool calls")
		}
	}
	
	// Add assistant message to history
	assistantMsg := models.ChatMessage{
		Role:      "assistant",
		Content:   choice.Message.Content,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"model":         resp.Model,
			"tokens_used":   resp.Usage.TotalTokens,
			"response_time": responseTime,
		},
	}
	c.conversationHistory = append(c.conversationHistory, assistantMsg)
	
	// Build response
	response := &ChatResponse{
		Message:      choice.Message.Content,
		Model:        resp.Model,
		TokensUsed:   resp.Usage.TotalTokens,
		ResponseTime: responseTime,
		ToolCalls:    toolCalls,
		Metadata: map[string]interface{}{
			"finish_reason": choice.FinishReason,
			"prompt_tokens": resp.Usage.PromptTokens,
			"completion_tokens": resp.Usage.CompletionTokens,
		},
	}
	
	// Generate suggestions if requested
	if options.GenerateSuggestions {
		response.Suggestions = c.generateSuggestions(message, choice.Message.Content)
	}
	
	return response, nil
}

// @decorator: ChatOptions
// @description: Options for chat completion
type ChatOptions struct {
	Model               string  `json:"model"`
	Temperature         float32 `json:"temperature"`
	MaxTokens           int     `json:"max_tokens"`
	Stream              bool    `json:"stream"`
	UseRAG              bool    `json:"use_rag"`
	UseTools            bool    `json:"use_tools"`
	GenerateSuggestions bool    `json:"generate_suggestions"`
	SystemPrompt        string  `json:"system_prompt,omitempty"`
}

// @decorator: StreamChat
// @description: Stream chat responses for real-time interaction
// @param ctx: Context for cancellation
// @param message: User message
// @param options: Chat options
// @param callback: Callback function for streaming responses
// @return error: Any error that occurred during streaming
func (c *Client) StreamChat(ctx context.Context, message string, options *ChatOptions, 
	callback func(chunk string) error) error {
	
	if callback == nil {
		return errors.New("callback function cannot be nil")
	}
	
	// Set streaming option
	if options == nil {
		options = &ChatOptions{}
	}
	options.Stream = true
	
	// Prepare messages
	enhancedMessage := message
	if c.ragEngine != nil && options.UseRAG {
		ragContext, err := c.ragEngine.Query(ctx, message, 5)
		if err != nil {
			c.logger.WithError(err).Warn("Failed to get RAG context")
		} else if len(ragContext) > 0 {
			enhancedMessage = fmt.Sprintf("Context:\n%s\n\nQuestion: %s", 
				strings.Join(ragContext, "\n"), message)
		}
	}
	
	messages := c.buildMessages(enhancedMessage, options)
	
	// Create streaming request
	req := openai.ChatCompletionRequest{
		Model:       options.Model,
		Messages:    messages,
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
		Stream:      true,
	}
	
	// Create stream
	stream, err := c.openaiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to create chat completion stream")
	}
	defer stream.Close()
	
	// Process stream
	var fullResponse strings.Builder
	for {
		response, err := stream.Recv()
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "stream receive error")
		}
		
		if len(response.Choices) > 0 {
			chunk := response.Choices[0].Delta.Content
			if chunk != "" {
				fullResponse.WriteString(chunk)
				if err := callback(chunk); err != nil {
					return errors.Wrap(err, "callback error")
				}
			}
		}
	}
}

// @decorator: AddDocument
// @description: Add a document to the RAG knowledge base
// @param ctx: Context for cancellation
// @param document: Document to add
// @return error: Any error that occurred during addition
func (c *Client) AddDocument(ctx context.Context, document *models.Document) error {
	if c.ragEngine == nil {
		return errors.New("RAG engine not initialized")
	}
	
	return c.ragEngine.AddDocument(ctx, document)
}

// @decorator: SearchDocuments
// @description: Search documents in the RAG knowledge base
// @param ctx: Context for cancellation
// @param query: Search query
// @param limit: Maximum number of results
// @return []string: Search results
// @return error: Any error that occurred during search
func (c *Client) SearchDocuments(ctx context.Context, query string, limit int) ([]string, error) {
	if c.ragEngine == nil {
		return nil, errors.New("RAG engine not initialized")
	}
	
	return c.ragEngine.Query(ctx, query, limit)
}

// @decorator: GetConversationHistory
// @description: Get the conversation history
// @return []models.ChatMessage: Conversation history
func (c *Client) GetConversationHistory() []models.ChatMessage {
	return c.conversationHistory
}

// @decorator: ClearConversationHistory
// @description: Clear the conversation history
func (c *Client) ClearConversationHistory() {
	c.conversationHistory = make([]models.ChatMessage, 0)
	c.logger.Info("Conversation history cleared")
}

// @decorator: SetSystemPrompt
// @description: Set the system prompt for AI interactions
// @param prompt: System prompt
func (c *Client) SetSystemPrompt(prompt string) {
	c.systemPrompt = prompt
	c.logger.Info("System prompt updated")
}

// @decorator: RegisterTool
// @description: Register a new tool for AI use
// @param tool: Tool to register
// @return error: Any error that occurred during registration
func (c *Client) RegisterTool(tool Tool) error {
	if tool == nil {
		return errors.New("tool cannot be nil")
	}
	
	name := tool.Name()
	if name == "" {
		return errors.New("tool name cannot be empty")
	}
	
	c.tools[name] = tool
	c.logger.WithField("tool", name).Info("Tool registered")
	return nil
}

// @decorator: GetAvailableModels
// @description: Get list of available AI models
// @param ctx: Context for cancellation
// @return []string: Available models
// @return error: Any error that occurred during retrieval
func (c *Client) GetAvailableModels(ctx context.Context) ([]string, error) {
	// This would typically call the OpenRouter API to get available models
	// For now, return a predefined list
	models := []string{
		"anthropic/claude-3-haiku",
		"anthropic/claude-3-sonnet",
		"anthropic/claude-3-opus",
		"openai/gpt-4-turbo",
		"openai/gpt-3.5-turbo",
		"google/gemini-pro",
		"meta-llama/llama-2-70b-chat",
		"mistralai/mixtral-8x7b-instruct",
		"qwen/qwen-72b-chat",
	}
	
	// Add custom models from configuration
	for name := range c.config.CustomModels {
		models = append(models, name)
	}
	
	return models, nil
}

// @decorator: buildMessages
// @description: Build messages array for OpenAI API
// @param message: User message
// @param options: Chat options
// @return []openai.ChatCompletionMessage: Messages array
func (c *Client) buildMessages(message string, options *ChatOptions) []openai.ChatCompletionMessage {
	var messages []openai.ChatCompletionMessage
	
	// Add system prompt
	systemPrompt := c.systemPrompt
	if options.SystemPrompt != "" {
		systemPrompt = options.SystemPrompt
	}
	
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
	})
	
	// Add conversation history (last 10 messages to avoid token limits)
	historyStart := 0
	if len(c.conversationHistory) > 10 {
		historyStart = len(c.conversationHistory) - 10
	}
	
	for i := historyStart; i < len(c.conversationHistory); i++ {
		msg := c.conversationHistory[i]
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	
	// Add current message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})
	
	return messages
}

// @decorator: buildToolDefinitions
// @description: Build tool definitions for OpenAI API
// @return []openai.Tool: Tool definitions
func (c *Client) buildToolDefinitions() []openai.Tool {
	var tools []openai.Tool
	
	for _, tool := range c.tools {
		toolDef := openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        tool.Name(),
				Description: tool.Description(),
				Parameters:  tool.Schema(),
			},
		}
		tools = append(tools, toolDef)
	}
	
	return tools
}

// @decorator: executeToolCalls
// @description: Execute tool calls from AI response
// @param ctx: Context for cancellation
// @param toolCalls: Tool calls to execute
// @return []ToolCall: Executed tool calls with responses
// @return error: Any error that occurred during execution
func (c *Client) executeToolCalls(ctx context.Context, toolCalls []openai.ToolCall) ([]ToolCall, error) {
	var results []ToolCall
	
	for _, call := range toolCalls {
		// Parse arguments
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
			c.logger.WithError(err).WithField("tool", call.Function.Name).Warn("Failed to parse tool arguments")
			continue
		}
		
		// Find and execute tool
		tool, exists := c.tools[call.Function.Name]
		if !exists {
			c.logger.WithField("tool", call.Function.Name).Warn("Tool not found")
			continue
		}
		
		response, err := tool.Execute(ctx, args)
		if err != nil {
			c.logger.WithError(err).WithField("tool", call.Function.Name).Warn("Tool execution failed")
			continue
		}
		
		results = append(results, ToolCall{
			ID:       call.ID,
			Name:     call.Function.Name,
			Args:     args,
			Response: response,
		})
	}
	
	return results, nil
}

// @decorator: initializeTools
// @description: Initialize built-in AI tools
// @return error: Any error that occurred during initialization
func (c *Client) initializeTools() error {
	// Register built-in tools
	tools := []Tool{
		NewAPIEndpointTool(),
		NewDataIngestionTool(),
		NewSQLQueryTool(),
		NewFileOperationTool(),
		NewSchedulerTool(),
	}
	
	for _, tool := range tools {
		if err := c.RegisterTool(tool); err != nil {
			return errors.Wrapf(err, "failed to register tool %s", tool.Name())
		}
	}
	
	return nil
}

// @decorator: generateSuggestions
// @description: Generate follow-up suggestions based on conversation
// @param userMessage: User's message
// @param aiResponse: AI's response
// @return []string: Generated suggestions
func (c *Client) generateSuggestions(userMessage, aiResponse string) []string {
	// This is a simplified implementation
	// In a real system, you might use another AI call or rule-based logic
	suggestions := []string{
		"Can you explain this in more detail?",
		"What are the next steps?",
		"Are there any alternatives?",
		"How can I implement this?",
		"What are the potential issues?",
	}
	
	return suggestions
}

// @decorator: getDefaultSystemPrompt
// @description: Get the default system prompt for AI interactions
// @return string: Default system prompt
func getDefaultSystemPrompt() string {
	return `You are an AI assistant specialized in data ingestion and management. You help users with:

1. API endpoint configuration and testing
2. Data ingestion script development
3. SQL query writing and optimization
4. Job scheduling and monitoring
5. Troubleshooting data pipeline issues
6. Best practices for data management

You have access to various tools and can help with:
- Analyzing API responses and schemas
- Writing and debugging ingestion scripts
- Creating SQL migrations
- Setting up scheduled jobs
- Monitoring data quality

Always provide practical, actionable advice and ask clarifying questions when needed.
Be concise but thorough in your explanations.`
}
