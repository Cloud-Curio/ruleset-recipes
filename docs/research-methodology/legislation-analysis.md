# Legislation Analysis Methodology

## Overview

This document outlines comprehensive methodologies for analyzing bills, resolutions, and other legislative texts. Our approach combines Natural Language Processing (NLP), statistical analysis, and domain expertise to extract meaningful insights from legislation.

## Research Objectives

When analyzing legislation, we aim to:

1. **Classify and Categorize**: Automatically sort legislation by policy area and topic
2. **Summarize**: Extract key provisions and generate human-readable summaries
3. **Predict Impact**: Assess potential effects on various stakeholders
4. **Track Evolution**: Monitor bills through the legislative process
5. **Identify Relationships**: Find similar bills and legislative patterns
6. **Measure Complexity**: Assess readability and complexity of legislative text

## Data Collection

### Primary Sources

#### 1. Congress.gov API
```
Endpoints:
- /bill/{congress}/{billType}/{billNumber}
- /bill/{congress}/{billType}/{billNumber}/text
- /bill/{congress}/{billType}/{billNumber}/actions
- /bill/{congress}/{billType}/{billNumber}/amendments
```

**Data Elements**:
- Bill text (XML and plain text)
- Bill title and summary
- Sponsors and co-sponsors
- Actions and status updates
- Committee referrals
- Related bills
- Subject tags

#### 2. GovInfo.gov API
- **Full text documents**: PDF, XML, HTML formats
- **Congressional Record**: Floor debates and statements
- **Committee Reports**: Analysis and recommendations
- **Bill versions**: Introduced, engrossed, enrolled, etc.

#### 3. OpenStates API (for state legislation)
- State bill texts
- State legislative calendars
- State committee actions
- State voting records

### Data Preprocessing

#### Text Extraction
1. **Parse XML/HTML**: Extract structured bill text
2. **Section Identification**: Identify sections, subsections, paragraphs
3. **Remove Boilerplate**: Strip standard legislative formatting
4. **Clean Text**: Normalize whitespace, fix encoding issues

#### Version Management
Track bill versions:
- **Introduced (IH/IS)**: Original version
- **Reported (RH/RS)**: Committee-amended version
- **Engrossed (EH/ES)**: House/Senate-passed version
- **Enrolled (ENR)**: Final version sent to President
- **Public Law**: Enacted version

## Analytical Framework

### 1. Text Classification

#### Policy Area Classification

**Primary Categories** (based on Congressional Research Service taxonomy):
1. Agriculture and Food
2. Armed Forces and National Security
3. Civil Rights and Liberties
4. Commerce and Trade
5. Crime and Law Enforcement
6. Economics and Public Finance
7. Education
8. Energy
9. Environmental Protection
10. Families and Family Issues
11. Finance and Financial Sector
12. Foreign Trade and International Finance
13. Government Operations and Politics
14. Health
15. Housing and Community Development
16. Immigration
17. International Affairs
18. Labor and Employment
19. Law
20. Native Americans
21. Public Lands and Natural Resources
22. Science, Technology, Communications
23. Social Welfare
24. Sports and Recreation
25. Taxation
26. Transportation and Public Works
27. Water Resources Development

#### Machine Learning Classification

**Approach 1: Supervised Learning**

```python
# Pseudocode for bill classification
def classify_bill(bill_text):
    # 1. Preprocess text
    cleaned_text = preprocess(bill_text)
    
    # 2. Extract features
    features = extract_features(cleaned_text)
    # - TF-IDF vectors
    # - N-grams (unigrams, bigrams, trigrams)
    # - Named entities
    # - Key phrases
    
    # 3. Apply trained classifier
    category = classifier.predict(features)
    confidence = classifier.predict_proba(features)
    
    return category, confidence
```

**Models to Consider**:
- **Logistic Regression**: Baseline, interpretable
- **Random Forest**: Good performance, feature importance
- **Support Vector Machines**: High accuracy for text
- **Neural Networks**: State-of-art (BERT, RoBERTa)
- **Ensemble Methods**: Combine multiple models

**Training Data**:
- Use existing CRS subject tags as labels
- Minimum 1000 examples per category
- Balance dataset or use weighted sampling
- Split: 70% train, 15% validation, 15% test

**Evaluation Metrics**:
- Accuracy: Overall correctness
- Precision: Correctness of positive predictions
- Recall: Coverage of actual positives
- F1-Score: Harmonic mean of precision and recall
- Confusion Matrix: Detailed error analysis

#### Topic Modeling

**Approach 2: Unsupervised Learning (LDA)**

Latent Dirichlet Allocation for discovering topics:

```python
# Pseudocode for topic modeling
def discover_topics(bill_corpus, num_topics=50):
    # 1. Preprocess corpus
    processed_docs = [preprocess(doc) for doc in bill_corpus]
    
    # 2. Create document-term matrix
    vectorizer = CountVectorizer(max_features=5000)
    dtm = vectorizer.fit_transform(processed_docs)
    
    # 3. Train LDA model
    lda_model = LDA(n_components=num_topics)
    lda_model.fit(dtm)
    
    # 4. Extract topics
    topics = lda_model.components_
    topic_words = get_top_words(topics, vectorizer)
    
    return lda_model, topic_words
```

**Parameter Tuning**:
- Number of topics: Use coherence score optimization
- Alpha/Beta parameters: Control topic/word distribution
- Iterations: Ensure convergence

### 2. Text Summarization

#### Extractive Summarization

Select most important sentences from bill text:

**Methods**:

1. **TextRank Algorithm**
```python
def textrank_summarize(text, num_sentences=5):
    # 1. Split into sentences
    sentences = sent_tokenize(text)
    
    # 2. Calculate sentence similarity matrix
    similarity_matrix = calculate_similarity(sentences)
    
    # 3. Apply PageRank
    scores = pagerank(similarity_matrix)
    
    # 4. Select top-ranked sentences
    top_sentences = get_top_n(sentences, scores, num_sentences)
    
    return " ".join(top_sentences)
```

2. **TF-IDF-based Selection**
- Calculate TF-IDF for each sentence
- Score sentences by average TF-IDF of constituent words
- Select highest-scoring sentences

3. **Position-based Selection**
- First sentence of each section (often contains key info)
- Title and purpose sections
- Effective date and applicability sections

#### Abstractive Summarization

Generate new summary text:

**Approaches**:

1. **Transformer Models** (T5, BART, Pegasus)
```python
def abstractive_summarize(text, max_length=150):
    # Use pre-trained summarization model
    summarizer = pipeline("summarization", model="facebook/bart-large-cnn")
    
    # Handle long documents with chunking
    chunks = chunk_text(text, max_chunk_size=1024)
    
    # Summarize each chunk
    chunk_summaries = [summarizer(chunk, max_length=max_length) 
                       for chunk in chunks]
    
    # Combine chunk summaries
    final_summary = combine_summaries(chunk_summaries)
    
    return final_summary
```

2. **Key Provisions Extraction**
Focus on specific legislative elements:
- Appropriations amounts
- Effective dates
- Affected agencies/departments
- New requirements or prohibitions
- Penalties and enforcement mechanisms

#### Summary Quality Metrics

**Automatic Evaluation**:
- **ROUGE Scores**: Compare to reference summaries
  - ROUGE-1: Unigram overlap
  - ROUGE-2: Bigram overlap
  - ROUGE-L: Longest common subsequence
- **BLEU Score**: Precision-based metric
- **BERTScore**: Semantic similarity using embeddings

**Human Evaluation**:
- Accuracy: Does summary reflect bill content?
- Coverage: Are key provisions included?
- Coherence: Is summary readable and logical?
- Conciseness: Is summary appropriately brief?

### 3. Similarity Analysis

#### Bill-to-Bill Similarity

**Cosine Similarity on TF-IDF Vectors**:

```python
def calculate_bill_similarity(bill1_text, bill2_text):
    # 1. Create TF-IDF vectors
    vectorizer = TfidfVectorizer()
    vectors = vectorizer.fit_transform([bill1_text, bill2_text])
    
    # 2. Calculate cosine similarity
    similarity = cosine_similarity(vectors[0:1], vectors[1:2])[0][0]
    
    return similarity
```

**Interpretation**:
- >0.9: Highly similar (likely same bill, different versions)
- 0.7-0.9: Very similar (companion bills, related legislation)
- 0.5-0.7: Moderately similar (same policy area)
- 0.3-0.5: Somewhat similar (overlapping topics)
- <0.3: Not similar

#### Semantic Similarity with Embeddings

**BERT-based Similarity**:

```python
def semantic_similarity(bill1_text, bill2_text):
    # 1. Generate embeddings
    model = SentenceTransformer('all-mpnet-base-v2')
    embeddings = model.encode([bill1_text, bill2_text])
    
    # 2. Calculate cosine similarity
    similarity = cosine_similarity(
        embeddings[0].reshape(1, -1),
        embeddings[1].reshape(1, -1)
    )[0][0]
    
    return similarity
```

**Advantages**:
- Captures semantic meaning, not just keywords
- Better handling of paraphrases
- More robust to wording variations

#### Historical Pattern Matching

Find similar bills from past sessions:
1. Build index of historical bills
2. Query index with new bill
3. Return top-k most similar bills
4. Analyze outcomes of similar historical bills

### 4. Complexity Analysis

#### Readability Metrics

**Flesch Reading Ease**:
```
Score = 206.835 - 1.015(total words / total sentences) 
        - 84.6(total syllables / total words)
```

Interpretation:
- 90-100: Very easy (5th grade)
- 60-70: Standard (8th-9th grade)
- 30-50: Difficult (college level)
- 0-30: Very difficult (graduate level)

**Flesch-Kincaid Grade Level**:
```
Grade = 0.39(total words / total sentences) 
        + 11.8(total syllables / total words) - 15.59
```

**Gunning Fog Index**:
```
Grade = 0.4 × [(words/sentences) + 100 × (complex words/words)]
```

#### Legislative Complexity Metrics

**Custom Metrics for Bills**:

1. **Length Metrics**:
   - Total word count
   - Total page count
   - Number of sections
   - Average section length

2. **Structural Complexity**:
   - Depth of section nesting
   - Cross-references to other laws
   - Number of definitions
   - Amendments to existing law

3. **Legal Complexity**:
   - References to U.S. Code sections
   - Number of exceptions and conditions
   - "Notwithstanding" clauses
   - Contingent provisions

```python
def calculate_legislative_complexity(bill_text):
    metrics = {
        'word_count': count_words(bill_text),
        'section_count': count_sections(bill_text),
        'cross_references': count_cross_references(bill_text),
        'amendments': count_amendments(bill_text),
        'definitions': count_definitions(bill_text),
        'exceptions': count_exceptions(bill_text),
        'readability_score': flesch_reading_ease(bill_text),
        'grade_level': flesch_kincaid_grade(bill_text)
    }
    
    # Calculate composite complexity score
    complexity_score = calculate_composite_score(metrics)
    
    return metrics, complexity_score
```

### 5. Impact Assessment

#### Stakeholder Analysis

Identify affected parties:

1. **Named Entity Recognition (NER)**:
   - Organizations mentioned
   - Government agencies
   - Geographic locations
   - Specific populations

2. **Regulatory Impact**:
   - New regulations created
   - Existing regulations modified
   - Compliance requirements

3. **Financial Impact**:
   - Appropriations amounts
   - Tax provisions
   - Fee structures
   - Cost estimates (CBO scores)

#### Sentiment and Tone Analysis

**Sentiment Classification**:
```python
def analyze_bill_sentiment(bill_text):
    # 1. Split into provisions
    provisions = extract_provisions(bill_text)
    
    # 2. Classify sentiment for each
    sentiments = []
    for provision in provisions:
        sentiment = sentiment_classifier(provision)
        # Returns: positive, negative, neutral
        sentiments.append(sentiment)
    
    # 3. Aggregate results
    overall_sentiment = aggregate_sentiments(sentiments)
    
    return overall_sentiment, sentiments
```

**Tone Indicators**:
- Mandatory vs. permissive language
- Prohibitions and requirements
- Incentives and penalties
- Emergency/urgency language

### 6. Legislative Lifecycle Tracking

#### Status Categorization

**Major Milestones**:
1. Introduced
2. Referred to Committee
3. Committee Hearings Held
4. Reported by Committee
5. Floor Consideration
6. Passed House/Senate
7. Conference Committee (if needed)
8. Passed Both Chambers
9. Sent to President
10. Signed/Vetoed
11. Law/Failed

#### Success Prediction

**Features for Prediction Models**:

```python
features = {
    'sponsor_seniority': years_in_congress,
    'sponsor_party': majority_or_minority,
    'sponsor_committee_position': chair_or_member,
    'cosponsors_count': total_cosponsors,
    'bipartisan_support': opposite_party_cosponsors_pct,
    'bill_length': word_count,
    'bill_complexity': complexity_score,
    'policy_area': primary_category,
    'session_period': days_remaining_in_session,
    'similar_bills_success': historical_success_rate
}
```

**Model Training**:
- Use historical bill data (10+ years)
- Binary classification: passed/not passed
- Evaluate with precision, recall, F1-score
- Calibrate probabilities for reliable predictions

## Implementation Framework

### Data Pipeline

```
1. Data Ingestion
   ├── API calls to Congress.gov, GovInfo.gov
   ├── Rate limiting and error handling
   └── Store raw data

2. Preprocessing
   ├── Text extraction and cleaning
   ├── Version control
   └── Metadata extraction

3. Analysis
   ├── Classification
   ├── Summarization
   ├── Similarity calculation
   └── Complexity analysis

4. Storage
   ├── Processed data
   ├── Analysis results
   └── Cached computations

5. API/Interface
   ├── Query processed data
   ├── Real-time analysis
   └── Batch processing
```

### Database Schema

```sql
-- Bills table
CREATE TABLE bills (
    id SERIAL PRIMARY KEY,
    congress INTEGER,
    bill_type VARCHAR(10),
    bill_number INTEGER,
    title TEXT,
    summary TEXT,
    full_text TEXT,
    introduced_date DATE,
    status VARCHAR(50),
    policy_area VARCHAR(100),
    complexity_score FLOAT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Bill classifications
CREATE TABLE bill_classifications (
    id SERIAL PRIMARY KEY,
    bill_id INTEGER REFERENCES bills(id),
    category VARCHAR(100),
    confidence FLOAT,
    model_version VARCHAR(50)
);

-- Bill summaries
CREATE TABLE bill_summaries (
    id SERIAL PRIMARY KEY,
    bill_id INTEGER REFERENCES bills(id),
    summary_type VARCHAR(50), -- extractive, abstractive, official
    summary_text TEXT,
    generation_method VARCHAR(100)
);

-- Bill similarity
CREATE TABLE bill_similarities (
    id SERIAL PRIMARY KEY,
    bill_id_1 INTEGER REFERENCES bills(id),
    bill_id_2 INTEGER REFERENCES bills(id),
    similarity_score FLOAT,
    similarity_method VARCHAR(50)
);
```

### Processing Queue

Use background job processing for expensive operations:

```python
# Example with Bull queue
@job('analyze-bill')
async def analyze_bill_job(bill_id):
    # 1. Fetch bill
    bill = await db.get_bill(bill_id)
    
    # 2. Run analyses
    classification = classify_bill(bill.text)
    summary = summarize_bill(bill.text)
    complexity = analyze_complexity(bill.text)
    
    # 3. Store results
    await db.save_analysis(bill_id, {
        'classification': classification,
        'summary': summary,
        'complexity': complexity
    })
    
    # 4. Calculate similarities (separate job)
    await queue.add('calculate-similarities', {'bill_id': bill_id})
```

## Quality Assurance

### Validation Methods

1. **Manual Review Sample**:
   - Review 5% of automated classifications
   - Verify summaries for accuracy
   - Check for obvious errors

2. **Cross-validation**:
   - Compare results across different methods
   - Verify consistency
   - Investigate discrepancies

3. **Expert Review**:
   - Subject matter expert validation
   - Policy area specialists
   - Legal review when needed

### Error Monitoring

Track and alert on:
- Classification confidence below threshold
- Unusually high/low complexity scores
- Failed API calls or data retrieval
- Processing errors and exceptions
- Drift in model performance

## Limitations and Considerations

### Known Limitations

1. **Text Quality**: OCR errors in older bills
2. **Context**: Amendments require context of original law
3. **Complexity**: Some provisions are inherently ambiguous
4. **Timeliness**: Analysis may lag behind rapid legislative changes
5. **Scope**: Focus on federal legislation; state coverage varies

### Ethical Considerations

1. **Neutrality**: Avoid partisan bias in classification and summary
2. **Transparency**: Document methodology and limitations
3. **Accuracy**: Prioritize correctness over speed
4. **Context**: Provide sufficient context for interpretation
5. **Updates**: Keep analysis current with bill status changes

## References

### Academic Literature
- Kornilova, A., et al. (2018). "BillSum: A Corpus for Automatic Summarization of US Legislation". *EMNLP*.
- Nay, J. (2016). "Gov2Vec: Learning Distributed Representations of Institutions and Their Legal Text". *ACL*.
- Ash, E., & Chen, D. L. (2019). "Measuring the Complexity of the Law". *Journal of Legal Analysis*.

### Technical Resources
- Hugging Face Transformers library
- spaCy NLP framework
- NLTK documentation
- scikit-learn documentation
- BERT, GPT models for text analysis

### Government Resources
- Congress.gov API documentation
- Congressional Research Service reports
- Government Publishing Office style guide
- Legislative drafting manuals

## Revision History

- v1.0 (Initial) - Comprehensive legislation analysis methodology
