# Analytical Methods and Techniques

## Overview

This document details the statistical methods, machine learning algorithms, and computational techniques used to analyze political data. All methods are chosen for their scientific validity, interpretability, and relevance to political science research.

## Statistical Foundations

### 1. Descriptive Statistics

#### Central Tendency Measures

**Mean (Average)**:
```python
def calculate_mean(values: List[float]) -> float:
    """Calculate arithmetic mean"""
    return sum(values) / len(values)
```

**Median**:
```python
def calculate_median(values: List[float]) -> float:
    """Calculate median (middle value)"""
    sorted_values = sorted(values)
    n = len(sorted_values)
    if n % 2 == 0:
        return (sorted_values[n//2 - 1] + sorted_values[n//2]) / 2
    else:
        return sorted_values[n//2]
```

**Mode**:
```python
from collections import Counter

def calculate_mode(values: List) -> Any:
    """Calculate mode (most frequent value)"""
    counter = Counter(values)
    return counter.most_common(1)[0][0]
```

#### Dispersion Measures

**Standard Deviation**:
```python
import math

def calculate_std_dev(values: List[float]) -> float:
    """Calculate standard deviation"""
    mean = calculate_mean(values)
    variance = sum((x - mean) ** 2 for x in values) / len(values)
    return math.sqrt(variance)
```

**Interquartile Range (IQR)**:
```python
import numpy as np

def calculate_iqr(values: List[float]) -> float:
    """Calculate IQR (Q3 - Q1)"""
    q1 = np.percentile(values, 25)
    q3 = np.percentile(values, 75)
    return q3 - q1
```

### 2. Inferential Statistics

#### Hypothesis Testing

**T-Test for Comparing Two Groups**:
```python
from scipy import stats

def compare_voting_patterns(group1_votes: List[int], 
                           group2_votes: List[int]) -> Dict:
    """
    Compare voting patterns between two groups
    
    Returns:
        dict with t-statistic, p-value, and interpretation
    """
    t_stat, p_value = stats.ttest_ind(group1_votes, group2_votes)
    
    result = {
        't_statistic': t_stat,
        'p_value': p_value,
        'significant': p_value < 0.05,
        'interpretation': (
            'Significant difference' if p_value < 0.05 
            else 'No significant difference'
        )
    }
    
    return result
```

**Chi-Square Test for Independence**:
```python
def test_party_vote_independence(contingency_table: np.ndarray) -> Dict:
    """
    Test if voting is independent of party affiliation
    
    Args:
        contingency_table: 2D array of [Party x Vote]
    
    Returns:
        dict with chi-square statistic, p-value, and interpretation
    """
    chi2, p_value, dof, expected = stats.chi2_contingency(contingency_table)
    
    return {
        'chi_square': chi2,
        'p_value': p_value,
        'degrees_of_freedom': dof,
        'significant': p_value < 0.05
    }
```

#### Confidence Intervals

**Calculate Confidence Interval for Proportion**:
```python
def proportion_confidence_interval(successes: int, 
                                  total: int, 
                                  confidence: float = 0.95) -> tuple:
    """
    Calculate confidence interval for a proportion
    
    Example: Confidence interval for bipartisan voting rate
    """
    from statsmodels.stats.proportion import proportion_confint
    
    lower, upper = proportion_confint(
        successes, 
        total, 
        alpha=1-confidence, 
        method='wilson'
    )
    
    return (lower, upper)
```

**Example Usage**:
```python
# Politician voted with opposition 15 times out of 100 votes
bipartisan_votes = 15
total_votes = 100

ci_lower, ci_upper = proportion_confidence_interval(
    bipartisan_votes, 
    total_votes
)

print(f"Bipartisan rate: {bipartisan_votes/total_votes:.1%}")
print(f"95% CI: [{ci_lower:.1%}, {ci_upper:.1%}]")
```

### 3. Correlation Analysis

#### Pearson Correlation

```python
def calculate_pearson_correlation(x: List[float], 
                                 y: List[float]) -> Dict:
    """
    Calculate Pearson correlation coefficient
    
    Use for: Relationship between continuous variables
    Example: Correlation between voting similarity and co-sponsorship
    """
    correlation, p_value = stats.pearsonr(x, y)
    
    return {
        'correlation': correlation,
        'p_value': p_value,
        'strength': interpret_correlation_strength(correlation)
    }

def interpret_correlation_strength(r: float) -> str:
    """Interpret correlation coefficient strength"""
    abs_r = abs(r)
    if abs_r >= 0.7:
        return 'Strong'
    elif abs_r >= 0.4:
        return 'Moderate'
    elif abs_r >= 0.2:
        return 'Weak'
    else:
        return 'Very weak'
```

#### Spearman Rank Correlation

```python
def calculate_spearman_correlation(x: List[float], 
                                  y: List[float]) -> Dict:
    """
    Calculate Spearman rank correlation
    
    Use for: Ordinal data or when relationship may not be linear
    Example: Relationship between seniority rank and influence score
    """
    correlation, p_value = stats.spearmanr(x, y)
    
    return {
        'correlation': correlation,
        'p_value': p_value
    }
```

## Similarity and Distance Metrics

### 1. Cosine Similarity

**Formula**: `similarity = cos(θ) = (A · B) / (||A|| × ||B||)`

```python
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity

def calculate_voting_similarity(votes_a: np.ndarray, 
                               votes_b: np.ndarray) -> float:
    """
    Calculate cosine similarity between two voting vectors
    
    Args:
        votes_a, votes_b: Arrays where 1=Yea, -1=Nay, 0=Absent/Present
    
    Returns:
        Similarity score between -1 and 1
    """
    # Reshape for sklearn
    votes_a = votes_a.reshape(1, -1)
    votes_b = votes_b.reshape(1, -1)
    
    similarity = cosine_similarity(votes_a, votes_b)[0][0]
    return similarity
```

**Interpretation**:
- 1.0: Identical voting patterns
- 0.5-0.99: High similarity
- 0.0-0.5: Moderate similarity
- -0.5-0.0: Low similarity (opposite voting)
- -1.0: Perfectly opposite voting

### 2. Jaccard Similarity

**Formula**: `J(A,B) = |A ∩ B| / |A ∪ B|`

```python
def jaccard_similarity(set_a: set, set_b: set) -> float:
    """
    Calculate Jaccard similarity between two sets
    
    Use for: Bill co-sponsorship overlap, committee membership overlap
    """
    intersection = len(set_a & set_b)
    union = len(set_a | set_b)
    
    if union == 0:
        return 0.0
    
    return intersection / union
```

**Example**:
```python
# Bills co-sponsored by politician A and B
bills_a = {'HR1', 'HR5', 'HR10', 'S2'}
bills_b = {'HR1', 'HR10', 'S2', 'S5'}

similarity = jaccard_similarity(bills_a, bills_b)
# Result: 0.6 (3 common bills out of 5 total unique bills)
```

### 3. Euclidean Distance

```python
def euclidean_distance(vector_a: np.ndarray, 
                       vector_b: np.ndarray) -> float:
    """
    Calculate Euclidean distance between two vectors
    
    Use for: Overall difference in voting patterns, policy positions
    """
    return np.linalg.norm(vector_a - vector_b)
```

### 4. Manhattan Distance

```python
def manhattan_distance(vector_a: np.ndarray, 
                       vector_b: np.ndarray) -> float:
    """
    Calculate Manhattan (L1) distance
    
    Use for: Count of differing votes
    """
    return np.sum(np.abs(vector_a - vector_b))
```

## Machine Learning Methods

### 1. Classification

#### Logistic Regression

```python
from sklearn.linear_model import LogisticRegression
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report, confusion_matrix

def train_bill_classifier(X: np.ndarray, y: np.ndarray) -> Dict:
    """
    Train logistic regression classifier for bill categorization
    
    Args:
        X: Feature matrix (TF-IDF vectors, bill characteristics)
        y: Labels (policy categories)
    
    Returns:
        Trained model and evaluation metrics
    """
    # Split data
    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, random_state=42, stratify=y
    )
    
    # Train model
    model = LogisticRegression(max_iter=1000, random_state=42)
    model.fit(X_train, y_train)
    
    # Evaluate
    y_pred = model.predict(X_test)
    
    return {
        'model': model,
        'accuracy': model.score(X_test, y_test),
        'classification_report': classification_report(y_test, y_pred),
        'confusion_matrix': confusion_matrix(y_test, y_pred)
    }
```

#### Random Forest

```python
from sklearn.ensemble import RandomForestClassifier

def train_random_forest(X: np.ndarray, y: np.ndarray, 
                       n_estimators: int = 100) -> Dict:
    """
    Train Random Forest classifier
    
    Advantages:
    - Handles non-linear relationships
    - Feature importance scores
    - Robust to outliers
    """
    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, random_state=42
    )
    
    model = RandomForestClassifier(
        n_estimators=n_estimators,
        random_state=42,
        max_depth=10
    )
    model.fit(X_train, y_train)
    
    # Get feature importance
    feature_importance = model.feature_importances_
    
    return {
        'model': model,
        'accuracy': model.score(X_test, y_test),
        'feature_importance': feature_importance
    }
```

### 2. Clustering

#### K-Means Clustering

```python
from sklearn.cluster import KMeans
from sklearn.metrics import silhouette_score

def cluster_politicians(voting_vectors: np.ndarray, 
                       n_clusters: int = 5) -> Dict:
    """
    Cluster politicians by voting patterns
    
    Use for: Identify voting coalitions, detect policy groups
    """
    # Apply K-Means
    kmeans = KMeans(n_clusters=n_clusters, random_state=42)
    cluster_labels = kmeans.fit_predict(voting_vectors)
    
    # Calculate quality metric
    silhouette = silhouette_score(voting_vectors, cluster_labels)
    
    return {
        'labels': cluster_labels,
        'centroids': kmeans.cluster_centers_,
        'silhouette_score': silhouette,
        'inertia': kmeans.inertia_
    }
```

**Choosing Number of Clusters**:
```python
def find_optimal_clusters(data: np.ndarray, 
                         max_k: int = 10) -> int:
    """
    Use elbow method to find optimal number of clusters
    """
    inertias = []
    silhouette_scores = []
    
    for k in range(2, max_k + 1):
        kmeans = KMeans(n_clusters=k, random_state=42)
        labels = kmeans.fit_predict(data)
        
        inertias.append(kmeans.inertia_)
        silhouette_scores.append(silhouette_score(data, labels))
    
    # Find elbow point (largest drop in inertia)
    # Or highest silhouette score
    optimal_k = np.argmax(silhouette_scores) + 2
    
    return optimal_k
```

#### Hierarchical Clustering

```python
from scipy.cluster.hierarchy import dendrogram, linkage
from scipy.spatial.distance import pdist

def hierarchical_clustering(data: np.ndarray, 
                           method: str = 'ward') -> Dict:
    """
    Perform hierarchical clustering
    
    Args:
        data: Feature matrix
        method: Linkage method ('ward', 'average', 'complete')
    
    Use for: Visualize politician relationships, 
             identify hierarchical coalition structure
    """
    # Calculate linkage
    linkage_matrix = linkage(data, method=method)
    
    return {
        'linkage_matrix': linkage_matrix,
        'method': method
    }
```

### 3. Dimensionality Reduction

#### Principal Component Analysis (PCA)

```python
from sklearn.decomposition import PCA

def reduce_dimensions_pca(data: np.ndarray, 
                         n_components: int = 2) -> Dict:
    """
    Reduce dimensions using PCA
    
    Use for: Visualization, feature reduction, identifying main axes of variation
    """
    pca = PCA(n_components=n_components)
    transformed_data = pca.fit_transform(data)
    
    return {
        'transformed_data': transformed_data,
        'explained_variance': pca.explained_variance_ratio_,
        'components': pca.components_,
        'model': pca
    }
```

**Visualizing Political Spectrum**:
```python
def visualize_political_spectrum(voting_data: np.ndarray, 
                                labels: List[str]):
    """
    Create 2D visualization of political positions
    """
    import matplotlib.pyplot as plt
    
    # Reduce to 2 dimensions
    pca_result = reduce_dimensions_pca(voting_data, n_components=2)
    coords = pca_result['transformed_data']
    
    # Plot
    plt.figure(figsize=(10, 8))
    plt.scatter(coords[:, 0], coords[:, 1])
    
    for i, label in enumerate(labels):
        plt.annotate(label, (coords[i, 0], coords[i, 1]))
    
    plt.xlabel(f"PC1 ({pca_result['explained_variance'][0]:.1%} variance)")
    plt.ylabel(f"PC2 ({pca_result['explained_variance'][1]:.1%} variance)")
    plt.title("Political Spectrum (PCA)")
    plt.show()
```

## Natural Language Processing

### 1. Text Preprocessing

```python
import re
from nltk.corpus import stopwords
from nltk.tokenize import word_tokenize
from nltk.stem import WordNetLemmatizer

def preprocess_bill_text(text: str) -> str:
    """
    Preprocess legislative text for NLP analysis
    """
    # Convert to lowercase
    text = text.lower()
    
    # Remove special characters and digits
    text = re.sub(r'[^a-z\s]', ' ', text)
    
    # Tokenize
    tokens = word_tokenize(text)
    
    # Remove stopwords
    stop_words = set(stopwords.words('english'))
    tokens = [t for t in tokens if t not in stop_words]
    
    # Lemmatize
    lemmatizer = WordNetLemmatizer()
    tokens = [lemmatizer.lemmatize(t) for t in tokens]
    
    return ' '.join(tokens)
```

### 2. Feature Extraction

#### TF-IDF Vectorization

```python
from sklearn.feature_extraction.text import TfidfVectorizer

def create_tfidf_features(documents: List[str], 
                         max_features: int = 1000) -> tuple:
    """
    Create TF-IDF feature matrix from documents
    
    Use for: Bill similarity, classification, topic modeling
    """
    vectorizer = TfidfVectorizer(
        max_features=max_features,
        ngram_range=(1, 2),  # Unigrams and bigrams
        min_df=2,  # Minimum document frequency
        max_df=0.8  # Maximum document frequency
    )
    
    tfidf_matrix = vectorizer.fit_transform(documents)
    feature_names = vectorizer.get_feature_names_out()
    
    return tfidf_matrix, feature_names, vectorizer
```

#### Word Embeddings

```python
from gensim.models import Word2Vec

def train_word_embeddings(tokenized_documents: List[List[str]], 
                         vector_size: int = 100) -> Word2Vec:
    """
    Train custom word embeddings on legislative corpus
    
    Use for: Semantic similarity, document representation
    """
    model = Word2Vec(
        sentences=tokenized_documents,
        vector_size=vector_size,
        window=5,
        min_count=2,
        workers=4
    )
    
    return model
```

### 3. Named Entity Recognition

```python
import spacy

def extract_entities(text: str) -> Dict[str, List[str]]:
    """
    Extract named entities from bill text
    
    Entities:
    - ORG: Organizations, agencies
    - GPE: Geographic locations
    - MONEY: Monetary values
    - LAW: Laws and regulations
    - DATE: Dates and time periods
    """
    nlp = spacy.load('en_core_web_sm')
    doc = nlp(text)
    
    entities = {}
    for ent in doc.ents:
        if ent.label_ not in entities:
            entities[ent.label_] = []
        entities[ent.label_].append(ent.text)
    
    return entities
```

## Network Analysis

### 1. Co-sponsorship Networks

```python
import networkx as nx

def build_cosponsorship_network(bills: List[Dict]) -> nx.Graph:
    """
    Build network of politicians connected by bill co-sponsorship
    
    Nodes: Politicians
    Edges: Co-sponsored at least one bill together
    Edge weights: Number of bills co-sponsored
    """
    G = nx.Graph()
    
    for bill in bills:
        sponsors = [bill['sponsor']] + bill['cosponsors']
        
        # Add edges between all pairs of sponsors
        for i, sponsor1 in enumerate(sponsors):
            for sponsor2 in sponsors[i+1:]:
                if G.has_edge(sponsor1, sponsor2):
                    G[sponsor1][sponsor2]['weight'] += 1
                else:
                    G.add_edge(sponsor1, sponsor2, weight=1)
    
    return G
```

### 2. Network Metrics

#### Centrality Measures

```python
def calculate_network_centrality(G: nx.Graph, 
                                node: str) -> Dict[str, float]:
    """
    Calculate various centrality metrics for a politician
    
    Metrics:
    - Degree: Number of connections
    - Betweenness: Bridging role
    - Closeness: Average distance to others
    - PageRank: Importance based on connections
    """
    return {
        'degree': G.degree(node),
        'betweenness': nx.betweenness_centrality(G)[node],
        'closeness': nx.closeness_centrality(G)[node],
        'pagerank': nx.pagerank(G)[node]
    }
```

#### Community Detection

```python
from networkx.algorithms import community

def detect_communities(G: nx.Graph) -> List[set]:
    """
    Detect communities in co-sponsorship network
    
    Use for: Identify political coalitions, working groups
    """
    communities = community.greedy_modularity_communities(G)
    return list(communities)
```

## Time Series Analysis

### 1. Trend Analysis

```python
import pandas as pd
from scipy import stats

def analyze_voting_trend(dates: List[str], 
                        bipartisan_scores: List[float]) -> Dict:
    """
    Analyze trend in bipartisan voting over time
    """
    # Create time series
    df = pd.DataFrame({
        'date': pd.to_datetime(dates),
        'score': bipartisan_scores
    })
    df = df.set_index('date').sort_index()
    
    # Calculate trend (linear regression)
    x = np.arange(len(df))
    slope, intercept, r_value, p_value, std_err = stats.linregress(
        x, df['score']
    )
    
    return {
        'slope': slope,
        'r_squared': r_value**2,
        'p_value': p_value,
        'trend': 'increasing' if slope > 0 else 'decreasing',
        'significant': p_value < 0.05
    }
```

### 2. Moving Averages

```python
def calculate_moving_average(values: pd.Series, 
                            window: int = 30) -> pd.Series:
    """
    Calculate moving average to smooth time series
    
    Args:
        values: Time series data
        window: Number of periods for averaging
    
    Use for: Smoothing voting patterns, identifying overall trends
    """
    return values.rolling(window=window, min_periods=1).mean()
```

## Model Evaluation

### Classification Metrics

```python
from sklearn.metrics import (
    accuracy_score, precision_score, recall_score, 
    f1_score, roc_auc_score
)

def evaluate_classifier(y_true: np.ndarray, 
                       y_pred: np.ndarray,
                       y_proba: np.ndarray = None) -> Dict:
    """
    Comprehensive classifier evaluation
    """
    metrics = {
        'accuracy': accuracy_score(y_true, y_pred),
        'precision': precision_score(y_true, y_pred, average='weighted'),
        'recall': recall_score(y_true, y_pred, average='weighted'),
        'f1_score': f1_score(y_true, y_pred, average='weighted')
    }
    
    if y_proba is not None:
        metrics['roc_auc'] = roc_auc_score(y_true, y_proba, 
                                          multi_class='ovr')
    
    return metrics
```

### Cross-Validation

```python
from sklearn.model_selection import cross_val_score

def cross_validate_model(model, X: np.ndarray, 
                        y: np.ndarray, cv: int = 5) -> Dict:
    """
    Perform k-fold cross-validation
    """
    scores = cross_val_score(model, X, y, cv=cv, 
                            scoring='accuracy')
    
    return {
        'scores': scores,
        'mean': scores.mean(),
        'std': scores.std(),
        'confidence_interval': (
            scores.mean() - 1.96 * scores.std(),
            scores.mean() + 1.96 * scores.std()
        )
    }
```

## Interpretability

### Feature Importance

```python
def get_feature_importance(model, feature_names: List[str], 
                          top_n: int = 10) -> pd.DataFrame:
    """
    Extract and display feature importance
    
    Works with: Random Forest, Gradient Boosting, etc.
    """
    importance = model.feature_importances_
    
    feature_importance = pd.DataFrame({
        'feature': feature_names,
        'importance': importance
    }).sort_values('importance', ascending=False)
    
    return feature_importance.head(top_n)
```

### SHAP Values

```python
import shap

def explain_prediction(model, X: np.ndarray, 
                      instance_idx: int) -> Dict:
    """
    Explain individual prediction using SHAP
    
    Use for: Understanding why a bill was classified a certain way
    """
    explainer = shap.TreeExplainer(model)
    shap_values = explainer.shap_values(X)
    
    return {
        'shap_values': shap_values[instance_idx],
        'base_value': explainer.expected_value,
        'feature_values': X[instance_idx]
    }
```

## Best Practices

### 1. Validation

- Always split data into train/test sets
- Use cross-validation for robust evaluation
- Test on out-of-sample data
- Report confidence intervals

### 2. Interpretation

- Prioritize interpretable models when possible
- Use SHAP/LIME for complex model explanation
- Validate results with domain experts
- Present uncertainty clearly

### 3. Reproducibility

- Set random seeds
- Version control code and data
- Document hyperparameters
- Save trained models

### 4. Performance

- Optimize algorithms for large datasets
- Use vectorized operations
- Cache intermediate results
- Parallelize when possible

## References

### Books
- "The Elements of Statistical Learning" - Hastie, Tibshirani, Friedman
- "Python for Data Analysis" - Wes McKinney
- "Speech and Language Processing" - Jurafsky & Martin

### Libraries
- scikit-learn: Machine learning
- scipy: Scientific computing
- statsmodels: Statistical models
- spaCy: NLP
- networkx: Network analysis
- gensim: Topic modeling

## Revision History

- v1.0 (Initial) - Comprehensive analytical methods documentation
