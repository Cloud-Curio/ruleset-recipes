/*
 * RAG System Implementation - Retrieval-Augmented Generation
 * 
 * This file implements the RAG (Retrieval-Augmented Generation) system
 * for the AI assistant, providing context-aware responses with knowledge
 * of user files, configurations, and data ingestion patterns.
 * 
 * Features:
 * - Vector embeddings for semantic search
 * - Document indexing and retrieval
 * - Context enhancement for AI responses
 * - Knowledge base management
 * - Similarity search and ranking
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package ai

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// @decorator: enhanceWithRAG
// @description: Enhance user message with relevant context from RAG system
func (a *AIAssistant) enhanceWithRAG(ctx context.Context, message string) (string, error) {
	if !a.config.EnableRAG {
		return message, nil
	}
	
	a.logger.WithField("message_length", len(message)).Info("🧠 Enhancing message with RAG")
	
	// Generate embedding for the user message
	messageEmbedding, err := a.embeddings.GenerateEmbedding(ctx, message)
	if err != nil {
		return message, fmt.Errorf("failed to generate embedding: %w", err)
	}
	
	// Search for relevant documents
	relevantDocs, err := a.vectorStore.Search(messageEmbedding, a.config.SimilarityThreshold, 5)
	if err != nil {
		return message, fmt.Errorf("failed to search vector store: %w", err)
	}
	
	if len(relevantDocs) == 0 {
		a.logger.Info("No relevant documents found in RAG system")
		return message, nil
	}
	
	// Build enhanced message with context
	enhancedMessage := a.buildEnhancedMessage(message, relevantDocs)
	
	a.logger.WithField("relevant_docs", len(relevantDocs)).Info("✅ Message enhanced with RAG context")
	
	return enhancedMessage, nil
}

// @decorator: IndexDocument
// @description: Index a document in the RAG system
func (a *AIAssistant) IndexDocument(ctx context.Context, docID, content string, metadata map[string]interface{}) error {
	a.logger.WithField("doc_id", docID).Info("📚 Indexing document in RAG system")
	
	// Generate embedding for the document
	embedding, err := a.embeddings.GenerateEmbedding(ctx, content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding for document: %w", err)
	}
	
	// Store in vector store
	err = a.vectorStore.Store(docID, embedding, metadata)
	if err != nil {
		return fmt.Errorf("failed to store document in vector store: %w", err)
	}
	
	// Update inverted index for keyword search
	a.vectorStore.UpdateIndex(docID, content)
	
	a.logger.WithField("doc_id", docID).Info("✅ Document indexed successfully")
	
	return nil
}

// @decorator: IndexDirectory
// @description: Index all files in a directory
func (a *AIAssistant) IndexDirectory(ctx context.Context, dirPath string) error {
	a.logger.WithField("directory", dirPath).Info("📁 Indexing directory")
	
	var indexedCount int
	
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories and non-text files
		if info.IsDir() || !a.isTextFile(path) {
			return nil
		}
		
		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			a.logger.WithError(err).WithField("file", path).Warn("Failed to read file")
			return nil // Continue with other files
		}
		
		// Create document ID from file path
		docID := a.generateDocumentID(path)
		
		// Create metadata
		metadata := map[string]interface{}{
			"file_path":    path,
			"file_name":    info.Name(),
			"file_size":    info.Size(),
			"modified_at":  info.ModTime(),
			"content_type": a.detectContentType(path),
		}
		
		// Index the document
		err = a.IndexDocument(ctx, docID, string(content), metadata)
		if err != nil {
			a.logger.WithError(err).WithField("file", path).Warn("Failed to index file")
			return nil // Continue with other files
		}
		
		indexedCount++
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}
	
	a.logger.WithField("indexed_files", indexedCount).Info("✅ Directory indexing completed")
	
	return nil
}

// @decorator: SearchDocuments
// @description: Search for documents using semantic similarity
func (a *AIAssistant) SearchDocuments(ctx context.Context, query string, limit int) ([]*DocumentResult, error) {
	a.logger.WithField("query", query).Info("🔍 Searching documents with RAG")
	
	// Generate embedding for the query
	queryEmbedding, err := a.embeddings.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}
	
	// Search vector store
	results, err := a.vectorStore.Search(queryEmbedding, a.config.SimilarityThreshold, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector store: %w", err)
	}
	
	a.logger.WithField("results", len(results)).Info("✅ Document search completed")
	
	return results, nil
}

// @decorator: GenerateEmbedding
// @description: Generate embedding for text using the embedding service
func (e *EmbeddingService) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	// Check cache first
	e.cacheMutex.RLock()
	if embedding, exists := e.cache[text]; exists {
		e.cacheMutex.RUnlock()
		return embedding, nil
	}
	e.cacheMutex.RUnlock()
	
	// For demo purposes, generate a simple embedding based on text characteristics
	// In production, this would call an actual embedding API
	embedding := e.generateSimpleEmbedding(text)
	
	// Cache the result
	e.cacheMutex.Lock()
	e.cache[text] = embedding
	e.cacheMutex.Unlock()
	
	return embedding, nil
}

// @decorator: generateSimpleEmbedding
// @description: Generate a simple embedding based on text characteristics
func (e *EmbeddingService) generateSimpleEmbedding(text string) []float64 {
	// This is a simplified embedding for demo purposes
	// In production, use actual embedding models like OpenAI's text-embedding-ada-002
	
	dimensions := 384 // Common embedding dimension
	embedding := make([]float64, dimensions)
	
	// Simple hash-based embedding
	hash := md5.Sum([]byte(text))
	
	// Convert hash to float values
	for i := 0; i < dimensions; i++ {
		byteIndex := i % len(hash)
		embedding[i] = float64(hash[byteIndex]) / 255.0
	}
	
	// Add some text-based features
	textLower := strings.ToLower(text)
	
	// Length feature
	if len(embedding) > 0 {
		embedding[0] += float64(len(text)) / 1000.0
	}
	
	// Word count feature
	if len(embedding) > 1 {
		wordCount := len(strings.Fields(text))
		embedding[1] += float64(wordCount) / 100.0
	}
	
	// Keyword features
	keywords := []string{"api", "data", "ingestion", "database", "query", "error", "config", "script"}
	for i, keyword := range keywords {
		if i+2 < len(embedding) && strings.Contains(textLower, keyword) {
			embedding[i+2] += 0.5
		}
	}
	
	// Normalize the embedding
	norm := 0.0
	for _, val := range embedding {
		norm += val * val
	}
	norm = math.Sqrt(norm)
	
	if norm > 0 {
		for i := range embedding {
			embedding[i] /= norm
		}
	}
	
	return embedding
}

// @decorator: Store
// @description: Store a document vector in the vector store
func (vs *VectorStore) Store(docID string, vector []float64, metadata map[string]interface{}) error {
	vs.mutex.Lock()
	defer vs.mutex.Unlock()
	
	if len(vector) != vs.dimensions {
		return fmt.Errorf("vector dimension mismatch: expected %d, got %d", vs.dimensions, len(vector))
	}
	
	vs.vectors[docID] = vector
	vs.metadata[docID] = metadata
	
	return nil
}

// @decorator: Search
// @description: Search for similar documents using cosine similarity
func (vs *VectorStore) Search(queryVector []float64, threshold float64, limit int) ([]*DocumentResult, error) {
	vs.mutex.RLock()
	defer vs.mutex.RUnlock()
	
	if len(queryVector) != vs.dimensions {
		return nil, fmt.Errorf("query vector dimension mismatch: expected %d, got %d", vs.dimensions, len(queryVector))
	}
	
	var results []*DocumentResult
	
	// Calculate similarity with all stored vectors
	for docID, vector := range vs.vectors {
		similarity := vs.cosineSimilarity(queryVector, vector)
		
		if similarity >= threshold {
			result := &DocumentResult{
				DocumentID: docID,
				Similarity: similarity,
				Metadata:   vs.metadata[docID],
			}
			results = append(results, result)
		}
	}
	
	// Sort by similarity (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})
	
	// Limit results
	if len(results) > limit {
		results = results[:limit]
	}
	
	return results, nil
}

// @decorator: UpdateIndex
// @description: Update the inverted index for keyword search
func (vs *VectorStore) UpdateIndex(docID, content string) {
	vs.mutex.Lock()
	defer vs.mutex.Unlock()
	
	// Simple tokenization
	words := strings.Fields(strings.ToLower(content))
	
	// Add document to index for each word
	for _, word := range words {
		// Clean the word
		word = strings.Trim(word, ".,!?;:\"'()[]{}*")
		if len(word) < 3 { // Skip short words
			continue
		}
		
		if vs.index[word] == nil {
			vs.index[word] = make([]string, 0)
		}
		
		// Check if document is already in the list
		found := false
		for _, existingDocID := range vs.index[word] {
			if existingDocID == docID {
				found = true
				break
			}
		}
		
		if !found {
			vs.index[word] = append(vs.index[word], docID)
		}
	}
}

// @decorator: cosineSimilarity
// @description: Calculate cosine similarity between two vectors
func (vs *VectorStore) cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}
	
	var dotProduct, normA, normB float64
	
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	
	if normA == 0 || normB == 0 {
		return 0.0
	}
	
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// @decorator: DocumentResult
// @description: Result from document search
type DocumentResult struct {
	DocumentID string                 `json:"document_id"`
	Similarity float64                `json:"similarity"`
	Metadata   map[string]interface{} `json:"metadata"`
	Content    string                 `json:"content,omitempty"`
}

// @decorator: buildEnhancedMessage
// @description: Build enhanced message with RAG context
func (a *AIAssistant) buildEnhancedMessage(originalMessage string, relevantDocs []*DocumentResult) string {
	if len(relevantDocs) == 0 {
		return originalMessage
	}
	
	var contextBuilder strings.Builder
	contextBuilder.WriteString("Context from relevant documents:\n\n")
	
	for i, doc := range relevantDocs {
		contextBuilder.WriteString(fmt.Sprintf("Document %d (similarity: %.3f):\n", i+1, doc.Similarity))
		
		// Add metadata if available
		if filePath, ok := doc.Metadata["file_path"].(string); ok {
			contextBuilder.WriteString(fmt.Sprintf("File: %s\n", filePath))
		}
		
		if contentType, ok := doc.Metadata["content_type"].(string); ok {
			contextBuilder.WriteString(fmt.Sprintf("Type: %s\n", contentType))
		}
		
		// Add content preview if available
		if doc.Content != "" {
			preview := doc.Content
			if len(preview) > 200 {
				preview = preview[:200] + "..."
			}
			contextBuilder.WriteString(fmt.Sprintf("Content: %s\n", preview))
		}
		
		contextBuilder.WriteString("\n")
	}
	
	contextBuilder.WriteString("---\n\n")
	contextBuilder.WriteString("User Question: ")
	contextBuilder.WriteString(originalMessage)
	contextBuilder.WriteString("\n\nPlease answer the user's question using the provided context when relevant.")
	
	return contextBuilder.String()
}

// @decorator: generateDocumentID
// @description: Generate a unique document ID from file path
func (a *AIAssistant) generateDocumentID(filePath string) string {
	// Use MD5 hash of the file path as document ID
	hash := md5.Sum([]byte(filePath))
	return fmt.Sprintf("doc_%x", hash)
}

// @decorator: isTextFile
// @description: Check if a file is a text file that should be indexed
func (a *AIAssistant) isTextFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	textExtensions := []string{
		".go", ".py", ".js", ".ts", ".java", ".cpp", ".c", ".h",
		".txt", ".md", ".json", ".yaml", ".yml", ".xml", ".html",
		".css", ".sql", ".sh", ".bat", ".ps1", ".dockerfile",
		".gitignore", ".env", ".config", ".conf", ".ini",
	}
	
	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	
	// Also check files without extensions (like Dockerfile, Makefile)
	fileName := strings.ToLower(filepath.Base(filePath))
	specialFiles := []string{
		"dockerfile", "makefile", "readme", "license", "changelog",
		"contributing", "authors", "todo", "notes",
	}
	
	for _, special := range specialFiles {
		if fileName == special {
			return true
		}
	}
	
	return false
}

// @decorator: detectContentType
// @description: Detect content type based on file extension
func (a *AIAssistant) detectContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	contentTypes := map[string]string{
		".go":         "source_code",
		".py":         "source_code",
		".js":         "source_code",
		".ts":         "source_code",
		".java":       "source_code",
		".cpp":        "source_code",
		".c":          "source_code",
		".h":          "source_code",
		".md":         "documentation",
		".txt":        "text",
		".json":       "configuration",
		".yaml":       "configuration",
		".yml":        "configuration",
		".xml":        "configuration",
		".sql":        "database",
		".sh":         "script",
		".bat":        "script",
		".ps1":        "script",
		".dockerfile": "configuration",
		".env":        "configuration",
		".config":     "configuration",
		".conf":       "configuration",
		".ini":        "configuration",
	}
	
	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}
	
	return "text"
}

// @decorator: GetRAGStats
// @description: Get statistics about the RAG system
func (a *AIAssistant) GetRAGStats() *RAGStats {
	a.vectorStore.mutex.RLock()
	defer a.vectorStore.mutex.RUnlock()
	
	stats := &RAGStats{
		DocumentCount:   len(a.vectorStore.vectors),
		VectorDimension: a.vectorStore.dimensions,
		IndexSize:       len(a.vectorStore.index),
		CacheSize:       len(a.embeddings.cache),
		LastUpdated:     time.Now(),
	}
	
	// Calculate content type distribution
	stats.ContentTypes = make(map[string]int)
	for _, metadata := range a.vectorStore.metadata {
		if contentType, ok := metadata["content_type"].(string); ok {
			stats.ContentTypes[contentType]++
		}
	}
	
	return stats
}

// @decorator: RAGStats
// @description: Statistics about the RAG system
type RAGStats struct {
	DocumentCount   int            `json:"document_count"`
	VectorDimension int            `json:"vector_dimension"`
	IndexSize       int            `json:"index_size"`
	CacheSize       int            `json:"cache_size"`
	ContentTypes    map[string]int `json:"content_types"`
	LastUpdated     time.Time      `json:"last_updated"`
}

// @decorator: ClearRAG
// @description: Clear all data from the RAG system
func (a *AIAssistant) ClearRAG() error {
	a.vectorStore.mutex.Lock()
	defer a.vectorStore.mutex.Unlock()
	
	a.vectorStore.vectors = make(map[string][]float64)
	a.vectorStore.metadata = make(map[string]map[string]interface{})
	a.vectorStore.index = make(map[string][]string)
	
	a.embeddings.cacheMutex.Lock()
	a.embeddings.cache = make(map[string][]float64)
	a.embeddings.cacheMutex.Unlock()
	
	a.logger.Info("🧹 RAG system cleared")
	
	return nil
}
