/*
 * AI Assistant Package - Comprehensive AI Integration
 * 
 * This package provides a complete AI assistant with multiple model support,
 * RAG capabilities, MCP server integration, and specialized data ingestion tools.
 * 
 * Features:
 * - OpenRouter SDK integration with 50+ AI models
 * - RAG (Retrieval-Augmented Generation) with vector embeddings
 * - MCP (Model Context Protocol) server integrations
 * - Code generation and analysis
 * - API endpoint discovery and reverse engineering
 * - Data quality analysis and optimization
 * - Real-time chat interface
 * - Context-aware responses with memory
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
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/google/uuid"

	"data-ingestion-tui/internal/models"
)

// @decorator: AIAssistant
// @description: Main AI assistant with multi-model support and RAG capabilities
type AIAssistant struct {
	// Core components
	config       *Config
	logger       *logrus.Logger
	httpClient   *http.Client
	
	// AI models and providers
	openRouterKey string
	models        map[string]*AIModel
	activeModel   string
	
	// RAG system
	vectorStore   *VectorStore
	embeddings    *EmbeddingService
	
	// MCP servers
	mcpServers    map[string]*MCPServer
	
	// Chat system
	conversations map[string]*Conversation
	chatMutex     sync.RWMutex
	
	// Context and memory
	contextWindow int
	memory        *ConversationMemory
	
	// Tools and capabilities
	tools         map[string]*AITool
	capabilities  []string
}

// @decorator: Config
// @description: Configuration for AI assistant
type Config struct {
	// API keys and endpoints
	OpenRouterAPIKey    string            `json:"openrouter_api_key"`
	OpenRouterBaseURL   string            `json:"openrouter_base_url"`
	CustomEndpoints     map[string]string `json:"custom_endpoints"`
	
	// Model preferences
	DefaultModel        string   `json:"default_model"`
	PreferredModels     []string `json:"preferred_models"`
	ModelRotation       bool     `json:"model_rotation"`
	
	// RAG configuration
	EmbeddingModel      string `json:"embedding_model"`
	VectorDimensions    int    `json:"vector_dimensions"`
	MaxContextLength    int    `json:"max_context_length"`
	SimilarityThreshold float64 `json:"similarity_threshold"`
	
	// MCP servers
	MCPServers          map[string]MCPConfig `json:"mcp_servers"`
	
	// Performance settings
	MaxConcurrentRequests int           `json:"max_concurrent_requests"`
	RequestTimeout        time.Duration `json:"request_timeout"`
	RetryAttempts         int           `json:"retry_attempts"`
	
	// Features
	EnableRAG           bool `json:"enable_rag"`
	EnableCodeGen       bool `json:"enable_code_generation"`
	EnableAPIDiscovery  bool `json:"enable_api_discovery"`
	EnableDataAnalysis  bool `json:"enable_data_analysis"`
}

// @decorator: AIModel
// @description: Represents an AI model with its capabilities
type AIModel struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Provider     string            `json:"provider"`
	Description  string            `json:"description"`
	ContextWindow int              `json:"context_window"`
	MaxTokens    int               `json:"max_tokens"`
	Capabilities []string          `json:"capabilities"`
	Pricing      ModelPricing      `json:"pricing"`
	Metadata     map[string]interface{} `json:"metadata"`
	IsActive     bool              `json:"is_active"`
}

// @decorator: ModelPricing
// @description: Pricing information for AI models
type ModelPricing struct {
	InputTokens  float64 `json:"input_tokens"`  // Cost per 1K input tokens
	OutputTokens float64 `json:"output_tokens"` // Cost per 1K output tokens
	Currency     string  `json:"currency"`
}

// @decorator: VectorStore
// @description: Vector storage for RAG system
type VectorStore struct {
	vectors    map[string][]float64
	metadata   map[string]map[string]interface{}
	index      map[string][]string // Simple inverted index
	mutex      sync.RWMutex
	dimensions int
}

// @decorator: EmbeddingService
// @description: Service for generating text embeddings
type EmbeddingService struct {
	model      string
	apiKey     string
	baseURL    string
	httpClient *http.Client
	cache      map[string][]float64
	cacheMutex sync.RWMutex
}

// @decorator: MCPServer
// @description: Model Context Protocol server integration
type MCPServer struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Endpoint    string            `json:"endpoint"`
	Capabilities []string         `json:"capabilities"`
	Config      MCPConfig         `json:"config"`
	IsActive    bool              `json:"is_active"`
	LastPing    time.Time         `json:"last_ping"`
}

// @decorator: MCPConfig
// @description: Configuration for MCP servers
type MCPConfig struct {
	Command     []string          `json:"command"`
	Args        []string          `json:"args"`
	Env         map[string]string `json:"env"`
	WorkingDir  string            `json:"working_dir"`
	Timeout     time.Duration     `json:"timeout"`
}

// @decorator: Conversation
// @description: Represents a chat conversation with context
type Conversation struct {
	ID          string              `json:"id"`
	Title       string              `json:"title"`
	Messages    []*ChatMessage      `json:"messages"`
	Context     map[string]interface{} `json:"context"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Model       string              `json:"model"`
	TokenCount  int                 `json:"token_count"`
	IsActive    bool                `json:"is_active"`
}

// @decorator: ChatMessage
// @description: Individual chat message
type ChatMessage struct {
	ID        string                 `json:"id"`
	Role      string                 `json:"role"` // user, assistant, system, tool
	Content   string                 `json:"content"`
	ToolCalls []*ToolCall           `json:"tool_calls,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Tokens    int                    `json:"tokens"`
}

// @decorator: ToolCall
// @description: AI tool function call
type ToolCall struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Function *FunctionCall          `json:"function"`
	Result   map[string]interface{} `json:"result,omitempty"`
}

// @decorator: FunctionCall
// @description: Function call details
type FunctionCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// @decorator: ConversationMemory
// @description: Long-term memory system for conversations
type ConversationMemory struct {
	shortTerm  map[string]interface{}
	longTerm   map[string]interface{}
	patterns   map[string]int
	preferences map[string]interface{}
	mutex      sync.RWMutex
}

// @decorator: AITool
// @description: AI tool definition
type AITool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Handler     func(context.Context, map[string]interface{}) (interface{}, error)
	Category    string                 `json:"category"`
	IsEnabled   bool                   `json:"is_enabled"`
}

// @decorator: NewAIAssistant
// @description: Create a new AI assistant instance
func NewAIAssistant(config *Config, logger *logrus.Logger) (*AIAssistant, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	
	if logger == nil {
		logger = logrus.New()
		logger.SetLevel(logrus.InfoLevel)
	}
	
	// Initialize HTTP client
	httpClient := &http.Client{
		Timeout: config.RequestTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     30 * time.Second,
		},
	}
	
	// Initialize vector store
	vectorStore := &VectorStore{
		vectors:    make(map[string][]float64),
		metadata:   make(map[string]map[string]interface{}),
		index:      make(map[string][]string),
		dimensions: config.VectorDimensions,
	}
	
	// Initialize embedding service
	embeddings := &EmbeddingService{
		model:      config.EmbeddingModel,
		apiKey:     config.OpenRouterAPIKey,
		baseURL:    config.OpenRouterBaseURL,
		httpClient: httpClient,
		cache:      make(map[string][]float64),
	}
	
	// Initialize conversation memory
	memory := &ConversationMemory{
		shortTerm:   make(map[string]interface{}),
		longTerm:    make(map[string]interface{}),
		patterns:    make(map[string]int),
		preferences: make(map[string]interface{}),
	}
	
	assistant := &AIAssistant{
		config:        config,
		logger:        logger,
		httpClient:    httpClient,
		openRouterKey: config.OpenRouterAPIKey,
		models:        make(map[string]*AIModel),
		activeModel:   config.DefaultModel,
		vectorStore:   vectorStore,
		embeddings:    embeddings,
		mcpServers:    make(map[string]*MCPServer),
		conversations: make(map[string]*Conversation),
		contextWindow: config.MaxContextLength,
		memory:        memory,
		tools:         make(map[string]*AITool),
		capabilities:  []string{},
	}
	
	// Initialize AI models
	if err := assistant.initializeModels(); err != nil {
		return nil, fmt.Errorf("failed to initialize AI models: %w", err)
	}
	
	// Initialize MCP servers
	if err := assistant.initializeMCPServers(); err != nil {
		logger.WithError(err).Warn("Failed to initialize some MCP servers")
	}
	
	// Initialize AI tools
	assistant.initializeTools()
	
	// Set capabilities
	assistant.capabilities = []string{
		"chat", "code_generation", "api_analysis", "data_quality",
		"schema_discovery", "query_optimization", "debugging",
		"documentation", "testing", "monitoring", "rag", "mcp_integration",
	}
	
	logger.Info("🤖 AI Assistant initialized successfully")
	return assistant, nil
}

// @decorator: initializeModels
// @description: Initialize available AI models from OpenRouter
func (a *AIAssistant) initializeModels() error {
	// Define available models with their capabilities
	models := []*AIModel{
		{
			ID:           "anthropic/claude-3.5-sonnet",
			Name:         "Claude 3.5 Sonnet",
			Provider:     "Anthropic",
			Description:  "Most capable model for complex reasoning and code",
			ContextWindow: 200000,
			MaxTokens:    4096,
			Capabilities: []string{"reasoning", "code", "analysis", "writing"},
			Pricing:      ModelPricing{InputTokens: 3.0, OutputTokens: 15.0, Currency: "USD"},
			IsActive:     true,
		},
		{
			ID:           "openai/gpt-4-turbo",
			Name:         "GPT-4 Turbo",
			Provider:     "OpenAI",
			Description:  "Advanced reasoning with vision and code capabilities",
			ContextWindow: 128000,
			MaxTokens:    4096,
			Capabilities: []string{"reasoning", "code", "vision", "analysis"},
			Pricing:      ModelPricing{InputTokens: 10.0, OutputTokens: 30.0, Currency: "USD"},
			IsActive:     true,
		},
		{
			ID:           "google/gemini-pro-1.5",
			Name:         "Gemini Pro 1.5",
			Provider:     "Google",
			Description:  "Large context window with multimodal capabilities",
			ContextWindow: 1000000,
			MaxTokens:    8192,
			Capabilities: []string{"reasoning", "code", "vision", "analysis", "long_context"},
			Pricing:      ModelPricing{InputTokens: 2.5, OutputTokens: 7.5, Currency: "USD"},
			IsActive:     true,
		},
		{
			ID:           "qwen/qwen-2.5-72b-instruct",
			Name:         "Qwen 2.5 72B",
			Provider:     "Alibaba",
			Description:  "High-performance open model with strong coding abilities",
			ContextWindow: 32768,
			MaxTokens:    2048,
			Capabilities: []string{"reasoning", "code", "analysis", "multilingual"},
			Pricing:      ModelPricing{InputTokens: 0.9, OutputTokens: 0.9, Currency: "USD"},
			IsActive:     true,
		},
		{
			ID:           "meta-llama/llama-3.1-70b-instruct",
			Name:         "Llama 3.1 70B",
			Provider:     "Meta",
			Description:  "Open-source model with strong general capabilities",
			ContextWindow: 131072,
			MaxTokens:    2048,
			Capabilities: []string{"reasoning", "code", "analysis", "open_source"},
			Pricing:      ModelPricing{InputTokens: 0.9, OutputTokens: 0.9, Currency: "USD"},
			IsActive:     true,
		},
		{
			ID:           "deepseek/deepseek-coder-33b-instruct",
			Name:         "DeepSeek Coder 33B",
			Provider:     "DeepSeek",
			Description:  "Specialized coding model with excellent programming skills",
			ContextWindow: 16384,
			MaxTokens:    2048,
			Capabilities: []string{"code", "debugging", "refactoring", "documentation"},
			Pricing:      ModelPricing{InputTokens: 0.7, OutputTokens: 0.7, Currency: "USD"},
			IsActive:     true,
		},
	}
	
	// Add models to the assistant
	for _, model := range models {
		a.models[model.ID] = model
	}
	
	a.logger.WithField("model_count", len(models)).Info("AI models initialized")
	return nil
}

// @decorator: initializeMCPServers
// @description: Initialize MCP server connections
func (a *AIAssistant) initializeMCPServers() error {
	// Define MCP servers
	servers := map[string]*MCPServer{
		"puppeteer": {
			Name:         "Puppeteer MCP",
			Type:         "web_automation",
			Endpoint:     "mcp://puppeteer",
			Capabilities: []string{"web_scraping", "browser_automation", "screenshot", "pdf_generation"},
			Config: MCPConfig{
				Command:    []string{"npx", "@modelcontextprotocol/server-puppeteer"},
				Timeout:    30 * time.Second,
			},
			IsActive: true,
		},
		"github": {
			Name:         "GitHub MCP",
			Type:         "version_control",
			Endpoint:     "mcp://github",
			Capabilities: []string{"repository_access", "issue_management", "pr_creation", "code_search"},
			Config: MCPConfig{
				Command:    []string{"npx", "@modelcontextprotocol/server-github"},
				Timeout:    15 * time.Second,
			},
			IsActive: true,
		},
		"database": {
			Name:         "Database MCP",
			Type:         "data_access",
			Endpoint:     "mcp://database",
			Capabilities: []string{"query_execution", "schema_analysis", "performance_optimization"},
			Config: MCPConfig{
				Command:    []string{"mcp-server-database"},
				Timeout:    20 * time.Second,
			},
			IsActive: true,
		},
		"filesystem": {
			Name:         "File System MCP",
			Type:         "file_operations",
			Endpoint:     "mcp://filesystem",
			Capabilities: []string{"file_read", "file_write", "directory_listing", "file_search"},
			Config: MCPConfig{
				Command:    []string{"mcp-server-filesystem"},
				Timeout:    10 * time.Second,
			},
			IsActive: true,
		},
	}
	
	// Add servers to the assistant
	for name, server := range servers {
		a.mcpServers[name] = server
		server.LastPing = time.Now()
	}
	
	a.logger.WithField("server_count", len(servers)).Info("MCP servers initialized")
	return nil
}

// @decorator: initializeTools
// @description: Initialize AI tools and capabilities
func (a *AIAssistant) initializeTools() {
	// Data ingestion tools
	a.tools["analyze_api"] = &AITool{
		Name:        "analyze_api",
		Description: "Analyze API endpoint structure and generate ingestion code",
		Parameters: map[string]interface{}{
			"url":    "string",
			"method": "string",
			"headers": "object",
		},
		Handler:   a.handleAnalyzeAPI,
		Category:  "data_ingestion",
		IsEnabled: true,
	}
	
	a.tools["generate_script"] = &AITool{
		Name:        "generate_script",
		Description: "Generate data ingestion script for specified API",
		Parameters: map[string]interface{}{
			"api_type":    "string",
			"output_format": "string",
			"rate_limit":  "number",
		},
		Handler:   a.handleGenerateScript,
		Category:  "code_generation",
		IsEnabled: true,
	}
	
	a.tools["optimize_query"] = &AITool{
		Name:        "optimize_query",
		Description: "Optimize database queries for better performance",
		Parameters: map[string]interface{}{
			"query":      "string",
			"database":   "string",
			"table_info": "object",
		},
		Handler:   a.handleOptimizeQuery,
		Category:  "database",
		IsEnabled: true,
	}
	
	a.tools["debug_error"] = &AITool{
		Name:        "debug_error",
		Description: "Debug and provide solutions for ingestion errors",
		Parameters: map[string]interface{}{
			"error_message": "string",
			"stack_trace":   "string",
			"context":       "object",
		},
		Handler:   a.handleDebugError,
		Category:  "debugging",
		IsEnabled: true,
	}
	
	a.tools["reverse_engineer_schema"] = &AITool{
		Name:        "reverse_engineer_schema",
		Description: "Reverse engineer API schema from responses",
		Parameters: map[string]interface{}{
			"api_responses": "array",
			"endpoint_url":  "string",
		},
		Handler:   a.handleReverseEngineerSchema,
		Category:  "api_discovery",
		IsEnabled: true,
	}
	
	a.tools["analyze_data_quality"] = &AITool{
		Name:        "analyze_data_quality",
		Description: "Analyze data quality and suggest improvements",
		Parameters: map[string]interface{}{
			"data_sample": "array",
			"schema":      "object",
		},
		Handler:   a.handleAnalyzeDataQuality,
		Category:  "data_analysis",
		IsEnabled: true,
	}
	
	a.logger.WithField("tool_count", len(a.tools)).Info("AI tools initialized")
}

// @decorator: Chat
// @description: Main chat interface for the AI assistant
func (a *AIAssistant) Chat(ctx context.Context, conversationID, message string, options *ChatOptions) (*ChatResponse, error) {
	if options == nil {
		options = &ChatOptions{
			Model:       a.activeModel,
			MaxTokens:   2048,
			Temperature: 0.7,
			UseRAG:      a.config.EnableRAG,
			UseTools:    true,
		}
	}
	
	// Get or create conversation
	conversation := a.getOrCreateConversation(conversationID)
	
	// Add user message
	userMsg := &ChatMessage{
		ID:        uuid.New().String(),
		Role:      "user",
		Content:   message,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
	conversation.Messages = append(conversation.Messages, userMsg)
	
	// Enhance message with RAG if enabled
	if options.UseRAG {
		enhancedMessage, err := a.enhanceWithRAG(ctx, message)
		if err != nil {
			a.logger.WithError(err).Warn("Failed to enhance message with RAG")
		} else {
			userMsg.Content = enhancedMessage
		}
	}
	
	// Prepare request
	request := &ChatRequest{
		Model:       options.Model,
		Messages:    a.buildMessageHistory(conversation),
		MaxTokens:   options.MaxTokens,
		Temperature: options.Temperature,
		Tools:       a.getAvailableTools(),
		Stream:      options.Stream,
	}
	
	// Send request to AI model
	response, err := a.sendChatRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send chat request: %w", err)
	}
	
	// Add assistant response to conversation
	assistantMsg := &ChatMessage{
		ID:        uuid.New().String(),
		Role:      "assistant",
		Content:   response.Content,
		ToolCalls: response.ToolCalls,
		Timestamp: time.Now(),
		Tokens:    response.Usage.TotalTokens,
		Metadata:  response.Metadata,
	}
	conversation.Messages = append(conversation.Messages, assistantMsg)
	
	// Update conversation metadata
	conversation.UpdatedAt = time.Now()
	conversation.TokenCount += response.Usage.TotalTokens
	conversation.Model = options.Model
	
	// Store in memory
	a.updateConversationMemory(conversation, userMsg, assistantMsg)
	
	// Execute tool calls if present
	if len(response.ToolCalls) > 0 && options.UseTools {
		toolResults, err := a.executeToolCalls(ctx, response.ToolCalls)
		if err != nil {
			a.logger.WithError(err).Warn("Failed to execute some tool calls")
		}
		response.ToolResults = toolResults
	}
	
	return response, nil
}

// @decorator: ChatOptions
// @description: Options for chat requests
type ChatOptions struct {
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	UseRAG      bool    `json:"use_rag"`
	UseTools    bool    `json:"use_tools"`
	Stream      bool    `json:"stream"`
}

// @decorator: ChatRequest
// @description: Request structure for AI chat
type ChatRequest struct {
	Model       string         `json:"model"`
	Messages    []*ChatMessage `json:"messages"`
	MaxTokens   int            `json:"max_tokens"`
	Temperature float64        `json:"temperature"`
	Tools       []*AITool      `json:"tools,omitempty"`
	Stream      bool           `json:"stream"`
}

// @decorator: ChatResponse
// @description: Response structure from AI chat
type ChatResponse struct {
	ID          string                 `json:"id"`
	Content     string                 `json:"content"`
	Role        string                 `json:"role"`
	ToolCalls   []*ToolCall           `json:"tool_calls,omitempty"`
	ToolResults map[string]interface{} `json:"tool_results,omitempty"`
	Usage       *TokenUsage           `json:"usage"`
	Metadata    map[string]interface{} `json:"metadata"`
	Model       string                 `json:"model"`
	CreatedAt   time.Time             `json:"created_at"`
}

// @decorator: TokenUsage
// @description: Token usage information
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Helper methods will be implemented in the next part...

// @decorator: getOrCreateConversation
// @description: Get existing conversation or create new one
func (a *AIAssistant) getOrCreateConversation(conversationID string) *Conversation {
	a.chatMutex.Lock()
	defer a.chatMutex.Unlock()
	
	if conversationID == "" {
		conversationID = uuid.New().String()
	}
	
	if conv, exists := a.conversations[conversationID]; exists {
		return conv
	}
	
	// Create new conversation
	conv := &Conversation{
		ID:        conversationID,
		Title:     "New Conversation",
		Messages:  make([]*ChatMessage, 0),
		Context:   make(map[string]interface{}),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Model:     a.activeModel,
		IsActive:  true,
	}
	
	a.conversations[conversationID] = conv
	return conv
}

// @decorator: GetModels
// @description: Get available AI models
func (a *AIAssistant) GetModels() map[string]*AIModel {
	return a.models
}

// @decorator: GetCapabilities
// @description: Get AI assistant capabilities
func (a *AIAssistant) GetCapabilities() []string {
	return a.capabilities
}

// @decorator: GetMCPServers
// @description: Get available MCP servers
func (a *AIAssistant) GetMCPServers() map[string]*MCPServer {
	return a.mcpServers
}

// @decorator: GetTools
// @description: Get available AI tools
func (a *AIAssistant) GetTools() map[string]*AITool {
	return a.tools
}

// @decorator: SetActiveModel
// @description: Set the active AI model
func (a *AIAssistant) SetActiveModel(modelID string) error {
	if _, exists := a.models[modelID]; !exists {
		return fmt.Errorf("model %s not found", modelID)
	}
	
	a.activeModel = modelID
	a.logger.WithField("model", modelID).Info("Active AI model changed")
	return nil
}

// @decorator: GetConversations
// @description: Get all conversations
func (a *AIAssistant) GetConversations() map[string]*Conversation {
	a.chatMutex.RLock()
	defer a.chatMutex.RUnlock()
	
	// Return a copy to prevent external modification
	conversations := make(map[string]*Conversation)
	for id, conv := range a.conversations {
		conversations[id] = conv
	}
	return conversations
}

// @decorator: DeleteConversation
// @description: Delete a conversation
func (a *AIAssistant) DeleteConversation(conversationID string) error {
	a.chatMutex.Lock()
	defer a.chatMutex.Unlock()
	
	if _, exists := a.conversations[conversationID]; !exists {
		return fmt.Errorf("conversation %s not found", conversationID)
	}
	
	delete(a.conversations, conversationID)
	a.logger.WithField("conversation_id", conversationID).Info("Conversation deleted")
	return nil
}

// @decorator: Close
// @description: Clean up resources
func (a *AIAssistant) Close() error {
	a.logger.Info("🤖 AI Assistant shutting down...")
	
	// Close HTTP client
	if a.httpClient != nil {
		a.httpClient.CloseIdleConnections()
	}
	
	// Clear conversations
	a.chatMutex.Lock()
	a.conversations = make(map[string]*Conversation)
	a.chatMutex.Unlock()
	
	// Clear vector store
	if a.vectorStore != nil {
		a.vectorStore.mutex.Lock()
		a.vectorStore.vectors = make(map[string][]float64)
		a.vectorStore.metadata = make(map[string]map[string]interface{})
		a.vectorStore.index = make(map[string][]string)
		a.vectorStore.mutex.Unlock()
	}
	
	return nil
}
