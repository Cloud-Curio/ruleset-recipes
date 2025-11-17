/*
 * AI Assistant Helper Methods - Supporting Functions
 * 
 * This file contains helper methods and supporting functions for the AI assistant,
 * including message processing, tool execution, and utility functions.
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
)

// @decorator: buildMessageHistory
// @description: Build message history for AI context
func (a *AIAssistant) buildMessageHistory(conversation *Conversation) []*ChatMessage {
	// Return recent messages within context window
	messages := conversation.Messages
	
	// If we have too many messages, keep only the most recent ones
	maxMessages := 20 // Reasonable limit for context
	if len(messages) > maxMessages {
		messages = messages[len(messages)-maxMessages:]
	}
	
	return messages
}

// @decorator: getAvailableTools
// @description: Get list of available AI tools
func (a *AIAssistant) getAvailableTools() []*AITool {
	tools := make([]*AITool, 0, len(a.tools))
	for _, tool := range a.tools {
		if tool.IsEnabled {
			tools = append(tools, tool)
		}
	}
	return tools
}

// @decorator: sendChatRequest
// @description: Send chat request to AI model (mock implementation)
func (a *AIAssistant) sendChatRequest(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
	// This is a mock implementation for demo purposes
	// In production, this would make actual API calls to OpenRouter or other providers
	
	a.logger.WithField("model", request.Model).Info("🤖 Sending request to AI model")
	
	// Simulate processing time
	time.Sleep(500 * time.Millisecond)
	
	// Generate mock response based on the user's message
	lastMessage := ""
	if len(request.Messages) > 0 {
		lastMessage = request.Messages[len(request.Messages)-1].Content
	}
	
	response := &ChatResponse{
		ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
		Role:      "assistant",
		Model:     request.Model,
		CreatedAt: time.Now(),
		Usage: &TokenUsage{
			PromptTokens:     len(lastMessage) / 4, // Rough estimate
			CompletionTokens: 100,                  // Mock value
			TotalTokens:      len(lastMessage)/4 + 100,
		},
		Metadata: make(map[string]interface{}),
	}
	
	// Generate contextual response based on message content
	response.Content = a.generateMockResponse(lastMessage, request.Tools)
	
	// Check if we should call tools
	if request.Tools != nil && len(request.Tools) > 0 {
		toolCalls := a.detectToolCalls(lastMessage)
		if len(toolCalls) > 0 {
			response.ToolCalls = toolCalls
		}
	}
	
	return response, nil
}

// @decorator: generateMockResponse
// @description: Generate mock AI response based on input
func (a *AIAssistant) generateMockResponse(message string, tools []*AITool) string {
	messageLower := strings.ToLower(message)
	
	// Handle specific commands
	if strings.HasPrefix(messageLower, "/analyze") {
		return "🔍 I'll analyze that API endpoint for you. Let me examine its structure, authentication requirements, rate limits, and generate recommendations for optimal data ingestion."
	}
	
	if strings.HasPrefix(messageLower, "/generate") {
		return "🔧 I'll generate a custom ingestion script for that API. The script will include proper error handling, rate limiting, authentication, and data validation."
	}
	
	if strings.HasPrefix(messageLower, "/optimize") {
		return "⚡ I'll optimize that database query for better performance. Let me analyze the query structure and suggest improvements like proper indexing, query rewriting, and execution plan optimization."
	}
	
	if strings.HasPrefix(messageLower, "/debug") {
		return "🐛 I'll help debug that error. Let me analyze the error message, identify the root cause, and provide specific solutions with code examples."
	}
	
	if strings.HasPrefix(messageLower, "/schema") {
		return "🏗️ I'll reverse engineer the API schema from those responses. This will help create proper data models and validation rules for your ingestion pipeline."
	}
	
	if strings.HasPrefix(messageLower, "/quality") {
		return "📊 I'll analyze the data quality and provide detailed metrics on completeness, consistency, accuracy, and duplicates, along with improvement suggestions."
	}
	
	// Handle general topics
	if strings.Contains(messageLower, "congress") || strings.Contains(messageLower, "bill") {
		return "🏛️ I can help you with Congress.gov data ingestion! This includes fetching congressional bills, sponsors, committees, actions, and votes. I can generate scripts with proper rate limiting (100 req/min) and handle the complex nested data structures."
	}
	
	if strings.Contains(messageLower, "fred") || strings.Contains(messageLower, "economic") {
		return "📈 For FRED economic data, I can help you ingest GDP, unemployment, inflation, and other economic indicators. The API supports time series data with various frequencies and I can handle the data transformations needed."
	}
	
	if strings.Contains(messageLower, "sec") || strings.Contains(messageLower, "edgar") {
		return "📋 SEC EDGAR data ingestion involves company filings, financial facts, and ownership data. I can help navigate the complex filing structures and extract key financial metrics with proper data validation."
	}
	
	if strings.Contains(messageLower, "database") || strings.Contains(messageLower, "sql") {
		return "💾 I can help optimize your database operations! This includes query optimization, index recommendations, schema design, and performance tuning for both PostgreSQL and SQLite."
	}
	
	if strings.Contains(messageLower, "error") || strings.Contains(messageLower, "debug") {
		return "🔧 I'm here to help debug issues! I can analyze error messages, stack traces, and provide specific solutions for common ingestion problems like rate limiting, authentication, parsing errors, and network issues."
	}
	
	if strings.Contains(messageLower, "api") {
		return "🌐 I can analyze any API endpoint for you! I'll examine the response structure, authentication requirements, rate limits, and generate optimized ingestion code with proper error handling."
	}
	
	// Default helpful response
	return "🤖 I'm your AI data ingestion assistant! I can help with:\n\n" +
		"• API analysis and reverse engineering\n" +
		"• Code generation for data ingestion\n" +
		"• Database query optimization\n" +
		"• Error debugging and troubleshooting\n" +
		"• Data quality analysis\n" +
		"• Schema design and validation\n\n" +
		"Try using commands like /analyze, /generate, /optimize, /debug, /schema, or /quality, or just ask me about any data ingestion challenge you're facing!"
}

// @decorator: detectToolCalls
// @description: Detect if message should trigger tool calls
func (a *AIAssistant) detectToolCalls(message string) []*ToolCall {
	var toolCalls []*ToolCall
	messageLower := strings.ToLower(message)
	
	// Detect analyze command
	if strings.HasPrefix(messageLower, "/analyze") {
		parts := strings.Fields(message)
		if len(parts) > 1 {
			toolCall := &ToolCall{
				ID:   fmt.Sprintf("call_%d", time.Now().UnixNano()),
				Type: "function",
				Function: &FunctionCall{
					Name: "analyze_api",
					Arguments: map[string]interface{}{
						"url":    parts[1],
						"method": "GET",
					},
				},
			}
			toolCalls = append(toolCalls, toolCall)
		}
	}
	
	// Detect generate command
	if strings.HasPrefix(messageLower, "/generate") {
		parts := strings.Fields(message)
		if len(parts) > 1 {
			toolCall := &ToolCall{
				ID:   fmt.Sprintf("call_%d", time.Now().UnixNano()),
				Type: "function",
				Function: &FunctionCall{
					Name: "generate_script",
					Arguments: map[string]interface{}{
						"api_type":      parts[1],
						"output_format": "json",
						"rate_limit":    100,
					},
				},
			}
			toolCalls = append(toolCalls, toolCall)
		}
	}
	
	// Detect optimize command
	if strings.HasPrefix(messageLower, "/optimize") {
		query := strings.TrimPrefix(message, "/optimize ")
		if query != message {
			toolCall := &ToolCall{
				ID:   fmt.Sprintf("call_%d", time.Now().UnixNano()),
				Type: "function",
				Function: &FunctionCall{
					Name: "optimize_query",
					Arguments: map[string]interface{}{
						"query":    query,
						"database": "postgresql",
					},
				},
			}
			toolCalls = append(toolCalls, toolCall)
		}
	}
	
	// Detect debug command
	if strings.HasPrefix(messageLower, "/debug") {
		errorMsg := strings.TrimPrefix(message, "/debug ")
		if errorMsg != message {
			toolCall := &ToolCall{
				ID:   fmt.Sprintf("call_%d", time.Now().UnixNano()),
				Type: "function",
				Function: &FunctionCall{
					Name: "debug_error",
					Arguments: map[string]interface{}{
						"error_message": errorMsg,
					},
				},
			}
			toolCalls = append(toolCalls, toolCall)
		}
	}
	
	return toolCalls
}

// @decorator: executeToolCalls
// @description: Execute AI tool calls
func (a *AIAssistant) executeToolCalls(ctx context.Context, toolCalls []*ToolCall) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	
	for _, toolCall := range toolCalls {
		if toolCall.Function == nil {
			continue
		}
		
		toolName := toolCall.Function.Name
		tool, exists := a.tools[toolName]
		if !exists || !tool.IsEnabled {
			a.logger.WithField("tool", toolName).Warn("Tool not found or disabled")
			continue
		}
		
		a.logger.WithField("tool", toolName).Info("🔧 Executing AI tool")
		
		// Execute the tool
		result, err := tool.Handler(ctx, toolCall.Function.Arguments)
		if err != nil {
			a.logger.WithError(err).WithField("tool", toolName).Error("Tool execution failed")
			results[toolName] = map[string]interface{}{
				"error": err.Error(),
			}
		} else {
			results[toolName] = result
			toolCall.Result = map[string]interface{}{
				"success": true,
				"data":    result,
			}
		}
	}
	
	return results, nil
}

// @decorator: updateConversationMemory
// @description: Update conversation memory with new messages
func (a *AIAssistant) updateConversationMemory(conversation *Conversation, userMsg, assistantMsg *ChatMessage) {
	a.memory.mutex.Lock()
	defer a.memory.mutex.Unlock()
	
	// Update short-term memory with recent context
	a.memory.shortTerm["last_user_message"] = userMsg.Content
	a.memory.shortTerm["last_assistant_message"] = assistantMsg.Content
	a.memory.shortTerm["conversation_id"] = conversation.ID
	a.memory.shortTerm["message_count"] = len(conversation.Messages)
	
	// Analyze patterns in user messages
	userContent := strings.ToLower(userMsg.Content)
	keywords := []string{"api", "database", "error", "optimize", "generate", "analyze", "debug"}
	
	for _, keyword := range keywords {
		if strings.Contains(userContent, keyword) {
			if a.memory.patterns[keyword] == 0 {
				a.memory.patterns[keyword] = 1
			} else {
				a.memory.patterns[keyword]++
			}
		}
	}
	
	// Update preferences based on usage patterns
	if strings.Contains(userContent, "claude") || strings.Contains(userContent, "anthropic") {
		a.memory.preferences["preferred_model"] = "anthropic/claude-3.5-sonnet"
	}
	if strings.Contains(userContent, "gpt") || strings.Contains(userContent, "openai") {
		a.memory.preferences["preferred_model"] = "openai/gpt-4-turbo"
	}
}

// Additional helper methods for query analysis and optimization

// @decorator: analyzeQueryStructure
// @description: Analyze SQL query structure
func (a *AIAssistant) analyzeQueryStructure(query string) map[string]interface{} {
	analysis := make(map[string]interface{})
	queryLower := strings.ToLower(query)
	
	// Basic query type detection
	if strings.Contains(queryLower, "select") {
		analysis["type"] = "SELECT"
	} else if strings.Contains(queryLower, "insert") {
		analysis["type"] = "INSERT"
	} else if strings.Contains(queryLower, "update") {
		analysis["type"] = "UPDATE"
	} else if strings.Contains(queryLower, "delete") {
		analysis["type"] = "DELETE"
	}
	
	// Complexity indicators
	analysis["has_joins"] = strings.Contains(queryLower, "join")
	analysis["has_subqueries"] = strings.Contains(queryLower, "(select")
	analysis["has_aggregates"] = strings.Contains(queryLower, "group by") || 
		strings.Contains(queryLower, "count(") || 
		strings.Contains(queryLower, "sum(")
	analysis["has_order_by"] = strings.Contains(queryLower, "order by")
	analysis["has_limit"] = strings.Contains(queryLower, "limit")
	
	// Estimate complexity score
	complexity := 1
	if analysis["has_joins"].(bool) {
		complexity += 2
	}
	if analysis["has_subqueries"].(bool) {
		complexity += 3
	}
	if analysis["has_aggregates"].(bool) {
		complexity += 2
	}
	analysis["complexity_score"] = complexity
	
	return analysis
}

// @decorator: optimizeQueryForDatabase
// @description: Generate optimized query for specific database
func (a *AIAssistant) optimizeQueryForDatabase(query, database string) string {
	// This is a simplified optimization for demo purposes
	// In production, this would use actual query optimization techniques
	
	optimized := query
	
	// Add basic optimizations
	if !strings.Contains(strings.ToLower(query), "limit") && strings.Contains(strings.ToLower(query), "select") {
		optimized += " LIMIT 1000" // Add reasonable limit
	}
	
	// Database-specific optimizations
	switch database {
	case "postgresql":
		// PostgreSQL-specific optimizations
		if strings.Contains(strings.ToLower(query), "like") {
			optimized = strings.ReplaceAll(optimized, "LIKE", "ILIKE") // Case-insensitive
		}
	case "sqlite":
		// SQLite-specific optimizations
		optimized = strings.ReplaceAll(optimized, "ILIKE", "LIKE") // SQLite doesn't have ILIKE
	}
	
	return optimized
}

// @decorator: calculateQueryImprovements
// @description: Calculate potential query improvements
func (a *AIAssistant) calculateQueryImprovements(original, optimized, database string) []string {
	improvements := make([]string, 0)
	
	if !strings.Contains(strings.ToLower(original), "limit") && strings.Contains(strings.ToLower(optimized), "limit") {
		improvements = append(improvements, "Added LIMIT clause to prevent excessive result sets")
	}
	
	if strings.Contains(strings.ToLower(optimized), "index") {
		improvements = append(improvements, "Suggested index creation for better performance")
	}
	
	improvements = append(improvements, "Query structure analyzed and optimized for "+database)
	
	return improvements
}

// @decorator: generateQueryRecommendations
// @description: Generate query optimization recommendations
func (a *AIAssistant) generateQueryRecommendations(analysis map[string]interface{}, database string) []string {
	recommendations := make([]string, 0)
	
	if hasJoins, ok := analysis["has_joins"].(bool); ok && hasJoins {
		recommendations = append(recommendations, "Consider adding indexes on join columns")
		recommendations = append(recommendations, "Ensure join conditions use appropriate data types")
	}
	
	if hasSubqueries, ok := analysis["has_subqueries"].(bool); ok && hasSubqueries {
		recommendations = append(recommendations, "Consider rewriting subqueries as JOINs for better performance")
	}
	
	if hasAggregates, ok := analysis["has_aggregates"].(bool); ok && hasAggregates {
		recommendations = append(recommendations, "Ensure GROUP BY columns are indexed")
		recommendations = append(recommendations, "Consider using covering indexes for aggregate queries")
	}
	
	if complexity, ok := analysis["complexity_score"].(int); ok && complexity > 5 {
		recommendations = append(recommendations, "Consider breaking complex query into smaller parts")
		recommendations = append(recommendations, "Use EXPLAIN ANALYZE to identify bottlenecks")
	}
	
	// Database-specific recommendations
	switch database {
	case "postgresql":
		recommendations = append(recommendations, "Consider using PostgreSQL-specific features like partial indexes")
		recommendations = append(recommendations, "Use VACUUM ANALYZE regularly for optimal performance")
	case "sqlite":
		recommendations = append(recommendations, "Consider using SQLite's query planner with ANALYZE")
		recommendations = append(recommendations, "Use appropriate SQLite pragmas for performance tuning")
	}
	
	return recommendations
}

// Additional helper methods for error analysis, data quality, etc. would be implemented here...
