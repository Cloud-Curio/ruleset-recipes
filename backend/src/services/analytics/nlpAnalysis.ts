/**
 * NLP Analysis Service
 * 
 * Provides comprehensive natural language processing capabilities for
 * analyzing political text content including sentiment analysis, topic
 * extraction, entity recognition, and toxicity detection.
 */

import natural from 'natural';
import Sentiment from 'sentiment';
import compromise from 'compromise';
import { Knex } from 'knex';
import db from '../../database/connection';

const TfIdf = natural.TfIdf;
const tokenizer = new natural.WordTokenizer();
const sentiment = new Sentiment();

interface NLPAnalysisResult {
  sentimentScore: number;
  sentimentLabel: 'positive' | 'negative' | 'neutral';
  sentimentMagnitude: number;
  topics: string[];
  entities: {
    people: string[];
    organizations: string[];
    places: string[];
  };
  keyPhrases: string[];
  toxicityScore: number;
  profanityScore: number;
  identityAttackScore: number;
  insultScore: number;
  threatScore: number;
  containsHateSpeech: boolean;
  hateSpeechTarget?: string;
  policyPositions: PolicyPosition[];
  complexityMetrics: ComplexityMetrics;
}

interface PolicyPosition {
  topic: string;
  stance: 'support' | 'oppose' | 'neutral';
  confidence: number;
  evidence: string[];
}

interface ComplexityMetrics {
  readingLevel: string;
  avgWordLength: number;
  avgSentenceLength: number;
  syllablesPerWord: number;
  fleschReadingEase: number;
}

export class NLPAnalysisService {
  private db: Knex;
  private tfidf: any;

  // Common political topics for classification
  private policyTopics = [
    'healthcare', 'immigration', 'economy', 'environment', 'education',
    'defense', 'foreign policy', 'taxes', 'gun control', 'abortion',
    'criminal justice', 'infrastructure', 'social security', 'energy',
    'civil rights', 'veterans', 'agriculture', 'technology', 'housing'
  ];

  // Hate speech indicators and targets
  private hateSpeechIndicators = {
    general: ['hate', 'disgust', 'destroy', 'eliminate', 'purge'],
    racial: ['racial slur', 'ethnic slur'], // Simplified - would use comprehensive list
    religious: ['religious slur'],
    gender: ['gender slur'],
    orientation: ['orientation slur']
  };

  constructor() {
    this.db = db;
    this.tfidf = new TfIdf();
  }

  /**
   * Analyze text content comprehensively
   */
  async analyzeText(
    text: string,
    contentType: string,
    contentId: string,
    politicianId: string
  ): Promise<NLPAnalysisResult> {
    // Sentiment analysis
    const sentimentResult = this.analyzeSentiment(text);

    // Topic extraction
    const topics = await this.extractTopics(text);

    // Entity recognition
    const entities = this.extractEntities(text);

    // Key phrases
    const keyPhrases = this.extractKeyPhrases(text);

    // Toxicity detection
    const toxicityScores = await this.detectToxicity(text);

    // Hate speech detection
    const hateSpeech = this.detectHateSpeech(text);

    // Policy position extraction
    const policyPositions = this.extractPolicyPositions(text);

    // Complexity metrics
    const complexityMetrics = this.calculateComplexity(text);

    const result: NLPAnalysisResult = {
      sentimentScore: sentimentResult.score,
      sentimentLabel: sentimentResult.label,
      sentimentMagnitude: sentimentResult.magnitude,
      topics,
      entities,
      keyPhrases,
      toxicityScore: toxicityScores.overall,
      profanityScore: toxicityScores.profanity,
      identityAttackScore: toxicityScores.identityAttack,
      insultScore: toxicityScores.insult,
      threatScore: toxicityScores.threat,
      containsHateSpeech: hateSpeech.detected,
      hateSpeechTarget: hateSpeech.target,
      policyPositions,
      complexityMetrics
    };

    // Store results in database
    await this.storeAnalysis(result, contentType, contentId, politicianId, text);

    return result;
  }

  /**
   * Sentiment analysis using multiple approaches
   */
  private analyzeSentiment(text: string): {
    score: number;
    label: 'positive' | 'negative' | 'neutral';
    magnitude: number;
  } {
    const result = sentiment.analyze(text);
    
    // Normalize score to -1 to 1 range
    const normalizedScore = Math.max(-1, Math.min(1, result.score / 10));
    
    // Calculate magnitude (strength of sentiment)
    const magnitude = Math.abs(normalizedScore);
    
    // Determine label
    let label: 'positive' | 'negative' | 'neutral' = 'neutral';
    if (normalizedScore > 0.2) label = 'positive';
    else if (normalizedScore < -0.2) label = 'negative';
    
    return {
      score: normalizedScore,
      label,
      magnitude
    };
  }

  /**
   * Extract topics using TF-IDF and keyword analysis
   */
  private async extractTopics(text: string): Promise<string[]> {
    const doc = compromise(text);
    const topics: string[] = [];

    // Check for policy topics
    for (const topic of this.policyTopics) {
      if (text.toLowerCase().includes(topic.replace('_', ' '))) {
        topics.push(topic);
      }
    }

    // Extract topics from noun phrases
    const nounPhrases = doc.nouns().out('array');
    const importantPhrases = nounPhrases.slice(0, 10); // Top 10
    
    return [...new Set([...topics, ...importantPhrases])];
  }

  /**
   * Extract named entities (people, organizations, places)
   */
  private extractEntities(text: string): {
    people: string[];
    organizations: string[];
    places: string[];
  } {
    const doc = compromise(text);

    return {
      people: doc.people().out('array'),
      organizations: doc.organizations().out('array'),
      places: doc.places().out('array')
    };
  }

  /**
   * Extract key phrases
   */
  private extractKeyPhrases(text: string): string[] {
    const doc = compromise(text);
    
    // Get important noun phrases and verb phrases
    const nounPhrases = doc.nouns().out('array');
    const verbPhrases = doc.verbs().out('array');
    
    // Combine and return top phrases
    const allPhrases = [...nounPhrases, ...verbPhrases];
    return allPhrases.slice(0, 20);
  }

  /**
   * Detect toxicity in text
   * Note: In production, use Perspective API or similar service
   */
  private async detectToxicity(text: string): Promise<{
    overall: number;
    profanity: number;
    identityAttack: number;
    insult: number;
    threat: number;
  }> {
    // Simplified toxicity detection
    // In production, integrate with Perspective API or similar
    
    const lowerText = text.toLowerCase();
    
    // Simple profanity detection
    const profanityWords = ['damn', 'hell']; // Simplified list
    const profanityCount = profanityWords.filter(word => 
      lowerText.includes(word)
    ).length;
    
    // Threat detection keywords
    const threatWords = ['destroy', 'attack', 'eliminate', 'kill'];
    const threatCount = threatWords.filter(word => 
      lowerText.includes(word)
    ).length;
    
    // Insult detection keywords
    const insultWords = ['stupid', 'idiot', 'fool', 'incompetent'];
    const insultCount = insultWords.filter(word => 
      lowerText.includes(word)
    ).length;
    
    return {
      overall: Math.min(1, (profanityCount + threatCount + insultCount) / 10),
      profanity: Math.min(1, profanityCount / 5),
      identityAttack: 0, // Would require more sophisticated analysis
      insult: Math.min(1, insultCount / 5),
      threat: Math.min(1, threatCount / 5)
    };
  }

  /**
   * Detect hate speech and identify targets
   */
  private detectHateSpeech(text: string): {
    detected: boolean;
    target?: string;
  } {
    const lowerText = text.toLowerCase();
    
    // Check for hate speech indicators
    for (const [category, keywords] of Object.entries(this.hateSpeechIndicators)) {
      for (const keyword of keywords) {
        if (lowerText.includes(keyword.toLowerCase())) {
          return {
            detected: true,
            target: category
          };
        }
      }
    }
    
    return { detected: false };
  }

  /**
   * Extract policy positions from text
   */
  private extractPolicyPositions(text: string): PolicyPosition[] {
    const positions: PolicyPosition[] = [];
    const doc = compromise(text);
    const lowerText = text.toLowerCase();

    // Support indicators
    const supportWords = ['support', 'advocate', 'favor', 'endorse', 'back', 'promote'];
    const opposeWords = ['oppose', 'against', 'reject', 'condemn', 'fight'];

    for (const topic of this.policyTopics) {
      const topicPhrase = topic.replace('_', ' ');
      
      if (lowerText.includes(topicPhrase)) {
        let stance: 'support' | 'oppose' | 'neutral' = 'neutral';
        let confidence = 0.5;

        // Find sentences containing the topic
        const sentences = text.split(/[.!?]+/);
        const topicSentences = sentences.filter(s => 
          s.toLowerCase().includes(topicPhrase)
        );

        for (const sentence of topicSentences) {
          const sentLower = sentence.toLowerCase();
          
          if (supportWords.some(word => sentLower.includes(word))) {
            stance = 'support';
            confidence = 0.7;
          } else if (opposeWords.some(word => sentLower.includes(word))) {
            stance = 'oppose';
            confidence = 0.7;
          }
        }

        if (stance !== 'neutral') {
          positions.push({
            topic,
            stance,
            confidence,
            evidence: topicSentences
          });
        }
      }
    }

    return positions;
  }

  /**
   * Calculate text complexity metrics
   */
  private calculateComplexity(text: string): ComplexityMetrics {
    const sentences = text.split(/[.!?]+/).filter(s => s.trim().length > 0);
    const words = tokenizer.tokenize(text.toLowerCase()) || [];
    
    const avgSentenceLength = words.length / sentences.length;
    const avgWordLength = words.reduce((sum, w) => sum + w.length, 0) / words.length;
    
    // Simplified syllable counting
    const syllablesPerWord = words.reduce((sum, word) => {
      return sum + this.countSyllables(word);
    }, 0) / words.length;
    
    // Flesch Reading Ease Score
    const fleschReadingEase = 206.835 - 1.015 * avgSentenceLength - 84.6 * syllablesPerWord;
    
    // Determine reading level
    let readingLevel = 'college';
    if (fleschReadingEase > 90) readingLevel = 'elementary';
    else if (fleschReadingEase > 80) readingLevel = 'middle_school';
    else if (fleschReadingEase > 70) readingLevel = 'high_school';
    else if (fleschReadingEase > 60) readingLevel = 'some_college';
    
    return {
      readingLevel,
      avgWordLength,
      avgSentenceLength,
      syllablesPerWord,
      fleschReadingEase
    };
  }

  /**
   * Simple syllable counter
   */
  private countSyllables(word: string): number {
    word = word.toLowerCase();
    if (word.length <= 3) return 1;
    
    const vowels = 'aeiouy';
    let count = 0;
    let previousWasVowel = false;
    
    for (let i = 0; i < word.length; i++) {
      const isVowel = vowels.includes(word[i]);
      if (isVowel && !previousWasVowel) {
        count++;
      }
      previousWasVowel = isVowel;
    }
    
    // Adjust for silent 'e'
    if (word.endsWith('e')) count--;
    
    return Math.max(1, count);
  }

  /**
   * Store analysis results in database
   */
  private async storeAnalysis(
    result: NLPAnalysisResult,
    contentType: string,
    contentId: string,
    politicianId: string,
    analyzedText: string
  ): Promise<void> {
    await this.db('nlp_analysis').insert({
      politician_id: politicianId,
      content_type: contentType,
      content_id: contentId,
      analyzed_text: analyzedText.substring(0, 10000), // Limit storage
      sentiment_score: result.sentimentScore,
      sentiment_label: result.sentimentLabel,
      sentiment_magnitude: result.sentimentMagnitude,
      topics: JSON.stringify(result.topics),
      entities: JSON.stringify(result.entities),
      key_phrases: JSON.stringify(result.keyPhrases),
      toxicity_score: result.toxicityScore,
      profanity_score: result.profanityScore,
      identity_attack_score: result.identityAttackScore,
      insult_score: result.insultScore,
      threat_score: result.threatScore,
      contains_hate_speech: result.containsHateSpeech,
      hate_speech_target: result.hateSpeechTarget,
      policy_positions: JSON.stringify(result.policyPositions),
      complexity_metrics: JSON.stringify(result.complexityMetrics),
      language: 'en',
      analyzed_at: new Date()
    });
  }

  /**
   * Batch analyze multiple texts
   */
  async batchAnalyze(
    texts: Array<{
      text: string;
      contentType: string;
      contentId: string;
      politicianId: string;
    }>
  ): Promise<NLPAnalysisResult[]> {
    const results: NLPAnalysisResult[] = [];
    
    for (const item of texts) {
      try {
        const result = await this.analyzeText(
          item.text,
          item.contentType,
          item.contentId,
          item.politicianId
        );
        results.push(result);
      } catch (error) {
        console.error(`Error analyzing text for ${item.contentId}:`, error);
      }
    }
    
    return results;
  }

  /**
   * Compare sentiment across time periods
   */
  async analyzeSentimentTrend(
    politicianId: string,
    startDate: Date,
    endDate: Date
  ): Promise<any> {
    const analyses = await this.db('nlp_analysis')
      .where('politician_id', politicianId)
      .whereBetween('analyzed_at', [startDate, endDate])
      .orderBy('analyzed_at', 'asc');

    // Group by month
    const monthlyTrends = new Map();
    
    for (const analysis of analyses) {
      const month = new Date(analysis.analyzed_at).toISOString().slice(0, 7);
      
      if (!monthlyTrends.has(month)) {
        monthlyTrends.set(month, {
          count: 0,
          totalSentiment: 0,
          totalToxicity: 0,
          hateSpeechCount: 0
        });
      }
      
      const trend = monthlyTrends.get(month);
      trend.count++;
      trend.totalSentiment += analysis.sentiment_score;
      trend.totalToxicity += analysis.toxicity_score;
      if (analysis.contains_hate_speech) trend.hateSpeechCount++;
    }
    
    // Calculate averages
    const trends = Array.from(monthlyTrends.entries()).map(([month, data]) => ({
      month,
      avgSentiment: data.totalSentiment / data.count,
      avgToxicity: data.totalToxicity / data.count,
      hateSpeechCount: data.hateSpeechCount,
      postCount: data.count
    }));
    
    return trends;
  }
}

// Export singleton instance
export default new NLPAnalysisService();
