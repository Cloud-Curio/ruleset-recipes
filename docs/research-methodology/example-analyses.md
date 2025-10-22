# Example Analyses

## Overview

This document provides practical examples demonstrating how to apply our research methodologies to real political data analysis tasks. Each example includes research question, methodology, code implementation, and interpretation.

## Example 1: Comparing Bipartisan Voting Between Two Politicians

### Research Question

How do Senator A and Senator B compare in terms of bipartisan voting behavior during the 118th Congress?

### Methodology

1. Collect voting records for both senators
2. Identify party-line votes (>80% of party votes same way)
3. Calculate percentage of times each senator voted with opposition
4. Perform statistical test to determine if difference is significant
5. Visualize results

### Implementation

```python
import pandas as pd
import numpy as np
from scipy import stats
import matplotlib.pyplot as plt

class BipartisanComparison:
    def __init__(self, api_client):
        self.api = api_client
    
    def analyze_bipartisan_voting(self, senator_id_a, senator_id_b, congress=118):
        """Compare bipartisan voting between two senators"""
        
        # 1. Fetch voting records
        votes_a = self.get_voting_records(senator_id_a, congress)
        votes_b = self.get_voting_records(senator_id_b, congress)
        
        # 2. Calculate bipartisan scores
        score_a = self.calculate_bipartisan_score(votes_a)
        score_b = self.calculate_bipartisan_score(votes_b)
        
        # 3. Statistical significance test
        significance = self.test_significance(votes_a, votes_b)
        
        # 4. Generate report
        report = self.generate_report(
            senator_id_a, senator_id_b,
            score_a, score_b, significance
        )
        
        # 5. Create visualization
        self.visualize_comparison(senator_id_a, senator_id_b, score_a, score_b)
        
        return report
    
    def calculate_bipartisan_score(self, votes):
        """Calculate percentage of party-line votes with opposition"""
        party = votes['party'].iloc[0]
        party_line_votes = []
        cross_party_count = 0
        
        for _, vote in votes.iterrows():
            # Get all votes for this roll call
            all_votes = self.get_roll_call_votes(vote['roll_call_id'])
            
            # Determine party majority position
            party_votes = all_votes[all_votes['party'] == party]
            majority_position = party_votes['position'].mode()[0]
            
            # Check if >80% voted with majority
            party_agreement = (
                party_votes['position'] == majority_position
            ).sum() / len(party_votes)
            
            if party_agreement >= 0.8:
                # This is a party-line vote
                party_line_votes.append(vote)
                
                # Did senator vote with opposition?
                opposition_votes = all_votes[all_votes['party'] != party]
                opposition_majority = opposition_votes['position'].mode()[0]
                
                if vote['position'] == opposition_majority:
                    cross_party_count += 1
        
        if len(party_line_votes) == 0:
            return 0
        
        return (cross_party_count / len(party_line_votes)) * 100
    
    def test_significance(self, votes_a, votes_b):
        """Test if difference is statistically significant"""
        # Create binary arrays (1 = voted with opposition, 0 = not)
        cross_party_a = self.get_cross_party_array(votes_a)
        cross_party_b = self.get_cross_party_array(votes_b)
        
        # Two-proportion z-test
        n1, n2 = len(cross_party_a), len(cross_party_b)
        p1, p2 = cross_party_a.mean(), cross_party_b.mean()
        
        # Pooled proportion
        p_pool = (cross_party_a.sum() + cross_party_b.sum()) / (n1 + n2)
        
        # Standard error
        se = np.sqrt(p_pool * (1 - p_pool) * (1/n1 + 1/n2))
        
        # Z-statistic
        z = (p1 - p2) / se
        
        # P-value (two-tailed)
        p_value = 2 * (1 - stats.norm.cdf(abs(z)))
        
        return {
            'z_statistic': z,
            'p_value': p_value,
            'significant': p_value < 0.05
        }
    
    def generate_report(self, id_a, id_b, score_a, score_b, significance):
        """Generate analysis report"""
        report = f"""
# Bipartisan Voting Comparison

## Senators Analyzed
- Senator A (ID: {id_a})
- Senator B (ID: {id_b})

## Bipartisan Scores
- Senator A: {score_a:.2f}%
- Senator B: {score_b:.2f}%
- Difference: {abs(score_a - score_b):.2f} percentage points

## Statistical Significance
- Z-statistic: {significance['z_statistic']:.3f}
- P-value: {significance['p_value']:.4f}
- Significant difference: {'Yes' if significance['significant'] else 'No'}

## Interpretation
"""
        if significance['significant']:
            if score_a > score_b:
                report += f"Senator A votes with the opposition party significantly more often than Senator B on party-line votes."
            else:
                report += f"Senator B votes with the opposition party significantly more often than Senator A on party-line votes."
        else:
            report += "There is no statistically significant difference in bipartisan voting between the two senators."
        
        return report
    
    def visualize_comparison(self, id_a, id_b, score_a, score_b):
        """Create bar chart comparison"""
        fig, ax = plt.subplots(figsize=(10, 6))
        
        senators = [f'Senator A\n({id_a})', f'Senator B\n({id_b})']
        scores = [score_a, score_b]
        
        bars = ax.bar(senators, scores, color=['#4A90E2', '#E24A4A'])
        ax.set_ylabel('Bipartisan Voting Rate (%)')
        ax.set_title('Bipartisan Voting Comparison')
        ax.set_ylim(0, max(scores) * 1.2)
        
        # Add value labels on bars
        for bar in bars:
            height = bar.get_height()
            ax.text(bar.get_x() + bar.get_width()/2., height,
                   f'{height:.1f}%',
                   ha='center', va='bottom')
        
        plt.tight_layout()
        plt.savefig('bipartisan_comparison.png', dpi=300)
        plt.show()
    
    # Helper methods (placeholders)
    def get_voting_records(self, senator_id, congress):
        # Fetch from database/API
        return pd.DataFrame()
    
    def get_roll_call_votes(self, roll_call_id):
        # Fetch all votes for a roll call
        return pd.DataFrame()
    
    def get_cross_party_array(self, votes):
        # Create binary array
        return np.array([])

# Usage example
api_client = CongressAPIClient(api_key='YOUR_KEY')
analyzer = BipartisanComparison(api_client)

report = analyzer.analyze_bipartisan_voting(
    senator_id_a='S000148',  # Schumer
    senator_id_b='M000355',  # McConnell
    congress=118
)

print(report)
```

### Expected Output

```
# Bipartisan Voting Comparison

## Senators Analyzed
- Senator A (ID: S000148)
- Senator B (ID: M000355)

## Bipartisan Scores
- Senator A: 8.3%
- Senator B: 6.1%
- Difference: 2.2 percentage points

## Statistical Significance
- Z-statistic: 1.823
- P-value: 0.0683
- Significant difference: No

## Interpretation
There is no statistically significant difference in bipartisan voting 
between the two senators.
```

## Example 2: Identifying Voting Coalitions Using Clustering

### Research Question

Can we identify distinct voting coalitions in the Senate based on voting patterns?

### Methodology

1. Create voting vectors for all senators
2. Calculate pairwise similarities
3. Apply hierarchical clustering
4. Identify and characterize clusters
5. Visualize coalition structure

### Implementation

```python
import numpy as np
from scipy.cluster.hierarchy import dendrogram, linkage
from scipy.spatial.distance import pdist, squareform
from sklearn.manifold import MDS
import matplotlib.pyplot as plt

class CoalitionAnalysis:
    def __init__(self, voting_data):
        """
        voting_data: DataFrame with senators as rows, votes as columns
                    Values: 1 (Yea), -1 (Nay), 0 (Absent/Present)
        """
        self.voting_data = voting_data
        self.senators = voting_data.index.tolist()
    
    def identify_coalitions(self, method='ward', n_clusters=4):
        """Identify voting coalitions"""
        
        # 1. Calculate similarity matrix
        similarity_matrix = self.calculate_similarity_matrix()
        
        # 2. Convert to distance matrix
        distance_matrix = 1 - similarity_matrix
        
        # 3. Hierarchical clustering
        linkage_matrix = linkage(
            squareform(distance_matrix), 
            method=method
        )
        
        # 4. Cut dendrogram to get clusters
        from scipy.cluster.hierarchy import fcluster
        clusters = fcluster(linkage_matrix, n_clusters, criterion='maxclust')
        
        # 5. Characterize clusters
        cluster_characteristics = self.characterize_clusters(clusters)
        
        # 6. Visualize
        self.visualize_dendrogram(linkage_matrix)
        self.visualize_2d_clusters(similarity_matrix, clusters)
        
        return {
            'clusters': clusters,
            'characteristics': cluster_characteristics,
            'linkage_matrix': linkage_matrix
        }
    
    def calculate_similarity_matrix(self):
        """Calculate pairwise cosine similarity"""
        n_senators = len(self.senators)
        similarity_matrix = np.zeros((n_senators, n_senators))
        
        for i in range(n_senators):
            for j in range(i, n_senators):
                sim = self.cosine_similarity(
                    self.voting_data.iloc[i].values,
                    self.voting_data.iloc[j].values
                )
                similarity_matrix[i, j] = sim
                similarity_matrix[j, i] = sim
        
        return similarity_matrix
    
    def cosine_similarity(self, vec1, vec2):
        """Calculate cosine similarity"""
        dot_product = np.dot(vec1, vec2)
        norm1 = np.linalg.norm(vec1)
        norm2 = np.linalg.norm(vec2)
        
        if norm1 == 0 or norm2 == 0:
            return 0
        
        return dot_product / (norm1 * norm2)
    
    def characterize_clusters(self, clusters):
        """Describe characteristics of each cluster"""
        characteristics = {}
        
        for cluster_id in np.unique(clusters):
            cluster_mask = clusters == cluster_id
            cluster_senators = self.voting_data[cluster_mask]
            
            # Get party composition
            # (Assuming we have party data)
            party_counts = self.get_party_composition(cluster_senators.index)
            
            # Average voting pattern
            avg_voting = cluster_senators.mean()
            
            characteristics[cluster_id] = {
                'size': cluster_mask.sum(),
                'members': cluster_senators.index.tolist(),
                'party_composition': party_counts,
                'avg_yea_rate': (avg_voting > 0).sum() / len(avg_voting),
            }
        
        return characteristics
    
    def visualize_dendrogram(self, linkage_matrix):
        """Plot hierarchical clustering dendrogram"""
        plt.figure(figsize=(12, 8))
        dendrogram(
            linkage_matrix,
            labels=self.senators,
            leaf_rotation=90,
            leaf_font_size=8
        )
        plt.title('Voting Coalition Dendrogram')
        plt.xlabel('Senator')
        plt.ylabel('Distance')
        plt.tight_layout()
        plt.savefig('coalition_dendrogram.png', dpi=300)
        plt.show()
    
    def visualize_2d_clusters(self, similarity_matrix, clusters):
        """Visualize clusters in 2D using MDS"""
        # Convert similarity to distance
        distance_matrix = 1 - similarity_matrix
        
        # Apply MDS
        mds = MDS(n_components=2, dissimilarity='precomputed', random_state=42)
        coords = mds.fit_transform(distance_matrix)
        
        # Plot
        plt.figure(figsize=(12, 10))
        
        colors = ['red', 'blue', 'green', 'orange', 'purple', 'brown']
        
        for cluster_id in np.unique(clusters):
            mask = clusters == cluster_id
            plt.scatter(
                coords[mask, 0], 
                coords[mask, 1],
                c=colors[cluster_id - 1],
                label=f'Coalition {cluster_id}',
                s=100,
                alpha=0.6
            )
            
            # Annotate senators
            for i, senator in enumerate(self.senators):
                if mask[i]:
                    plt.annotate(
                        senator, 
                        (coords[i, 0], coords[i, 1]),
                        fontsize=8
                    )
        
        plt.xlabel('Dimension 1')
        plt.ylabel('Dimension 2')
        plt.title('Voting Coalitions (2D Visualization)')
        plt.legend()
        plt.grid(True, alpha=0.3)
        plt.tight_layout()
        plt.savefig('coalition_2d.png', dpi=300)
        plt.show()
    
    def get_party_composition(self, senator_ids):
        """Get party breakdown (placeholder)"""
        # In real implementation, fetch from database
        return {'Democrat': 0, 'Republican': 0, 'Independent': 0}

# Usage example
# Create sample voting data
voting_data = pd.DataFrame(
    np.random.choice([1, -1, 0], size=(100, 500)),  # 100 senators, 500 votes
    index=[f'Senator_{i}' for i in range(100)]
)

analyzer = CoalitionAnalysis(voting_data)
results = analyzer.identify_coalitions(n_clusters=4)

# Print cluster characteristics
for cluster_id, chars in results['characteristics'].items():
    print(f"\nCoalition {cluster_id}:")
    print(f"  Size: {chars['size']} senators")
    print(f"  Members: {', '.join(chars['members'][:5])}...")
    print(f"  Party composition: {chars['party_composition']}")
```

## Example 3: Bill Classification and Topic Modeling

### Research Question

Can we automatically classify bills by policy area and extract key topics?

### Methodology

1. Collect bill texts
2. Preprocess text (cleaning, tokenization)
3. Extract TF-IDF features
4. Train classification model
5. Apply topic modeling (LDA)
6. Evaluate and visualize results

### Implementation

```python
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.decomposition import LatentDirichletAllocation
from sklearn.model_selection import train_test_split
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import classification_report
import re

class BillClassifier:
    def __init__(self):
        self.vectorizer = TfidfVectorizer(
            max_features=5000,
            ngram_range=(1, 2),
            min_df=2,
            max_df=0.8,
            stop_words='english'
        )
        self.classifier = None
        self.lda_model = None
    
    def preprocess_text(self, text):
        """Clean and preprocess bill text"""
        # Lowercase
        text = text.lower()
        
        # Remove section numbers and legal citations
        text = re.sub(r'section \d+', '', text)
        text = re.sub(r'\d+ u\.s\.c\.', '', text)
        
        # Remove special characters
        text = re.sub(r'[^a-z\s]', ' ', text)
        
        # Remove extra whitespace
        text = re.sub(r'\s+', ' ', text).strip()
        
        return text
    
    def train_classifier(self, bills_df):
        """
        Train bill classification model
        
        bills_df: DataFrame with 'text' and 'policy_area' columns
        """
        # Preprocess texts
        texts = bills_df['text'].apply(self.preprocess_text)
        labels = bills_df['policy_area']
        
        # Split data
        X_train, X_test, y_train, y_test = train_test_split(
            texts, labels, test_size=0.2, random_state=42, stratify=labels
        )
        
        # Vectorize
        X_train_tfidf = self.vectorizer.fit_transform(X_train)
        X_test_tfidf = self.vectorizer.transform(X_test)
        
        # Train classifier
        self.classifier = RandomForestClassifier(
            n_estimators=100,
            max_depth=20,
            random_state=42
        )
        self.classifier.fit(X_train_tfidf, y_train)
        
        # Evaluate
        y_pred = self.classifier.predict(X_test_tfidf)
        
        print("Classification Report:")
        print(classification_report(y_test, y_pred))
        
        # Feature importance
        feature_importance = self.get_feature_importance()
        
        return {
            'accuracy': self.classifier.score(X_test_tfidf, y_test),
            'classification_report': classification_report(y_test, y_pred),
            'feature_importance': feature_importance
        }
    
    def classify_bill(self, bill_text):
        """Classify a single bill"""
        if self.classifier is None:
            raise ValueError("Model not trained")
        
        # Preprocess
        processed_text = self.preprocess_text(bill_text)
        
        # Vectorize
        tfidf_vector = self.vectorizer.transform([processed_text])
        
        # Predict
        category = self.classifier.predict(tfidf_vector)[0]
        probabilities = self.classifier.predict_proba(tfidf_vector)[0]
        
        # Get top 3 predictions
        top_indices = probabilities.argsort()[-3:][::-1]
        top_predictions = [
            {
                'category': self.classifier.classes_[i],
                'probability': probabilities[i]
            }
            for i in top_indices
        ]
        
        return {
            'primary_category': category,
            'top_predictions': top_predictions
        }
    
    def extract_topics(self, bills_df, n_topics=10):
        """Extract topics using LDA"""
        # Preprocess
        texts = bills_df['text'].apply(self.preprocess_text)
        
        # Vectorize
        tfidf_matrix = self.vectorizer.fit_transform(texts)
        
        # LDA
        self.lda_model = LatentDirichletAllocation(
            n_components=n_topics,
            random_state=42,
            max_iter=20
        )
        self.lda_model.fit(tfidf_matrix)
        
        # Extract topic words
        topics = self.get_topic_words(n_words=10)
        
        return topics
    
    def get_topic_words(self, n_words=10):
        """Get top words for each topic"""
        if self.lda_model is None:
            raise ValueError("LDA model not trained")
        
        feature_names = self.vectorizer.get_feature_names_out()
        topics = []
        
        for topic_idx, topic in enumerate(self.lda_model.components_):
            top_indices = topic.argsort()[-n_words:][::-1]
            top_words = [feature_names[i] for i in top_indices]
            
            topics.append({
                'topic_id': topic_idx,
                'top_words': top_words,
                'word_scores': topic[top_indices].tolist()
            })
        
        return topics
    
    def get_feature_importance(self, n_features=20):
        """Get most important features for classification"""
        if self.classifier is None:
            raise ValueError("Classifier not trained")
        
        feature_names = self.vectorizer.get_feature_names_out()
        importances = self.classifier.feature_importances_
        
        # Get top features
        top_indices = importances.argsort()[-n_features:][::-1]
        
        return [
            {
                'feature': feature_names[i],
                'importance': importances[i]
            }
            for i in top_indices
        ]

# Usage example
# Load bill data
bills_df = pd.DataFrame({
    'text': ['Bill text 1...', 'Bill text 2...'],
    'policy_area': ['Healthcare', 'Education']
})

classifier = BillClassifier()

# Train
results = classifier.train_classifier(bills_df)
print(f"Model accuracy: {results['accuracy']:.2%}")

# Classify new bill
new_bill = "A bill to expand access to healthcare..."
classification = classifier.classify_bill(new_bill)
print(f"\nPredicted category: {classification['primary_category']}")
print("Top predictions:")
for pred in classification['top_predictions']:
    print(f"  {pred['category']}: {pred['probability']:.2%}")

# Extract topics
topics = classifier.extract_topics(bills_df, n_topics=5)
print("\nDiscovered Topics:")
for topic in topics:
    print(f"\nTopic {topic['topic_id']}:")
    print(f"  Top words: {', '.join(topic['top_words'][:5])}")
```

## Example 4: Influence Network Analysis

### Research Question

Who are the most influential politicians based on co-sponsorship networks?

### Methodology

1. Build co-sponsorship network
2. Calculate centrality metrics
3. Identify influential nodes
4. Visualize network
5. Analyze influence patterns

### Implementation

```python
import networkx as nx
import matplotlib.pyplot as plt
from community import community_louvain

class InfluenceNetworkAnalysis:
    def __init__(self, cosponsorship_data):
        """
        cosponsorship_data: List of bills with sponsor and cosponsors
        """
        self.data = cosponsorship_data
        self.network = None
    
    def build_network(self):
        """Build co-sponsorship network"""
        G = nx.Graph()
        
        for bill in self.data:
            sponsor = bill['sponsor']
            cosponsors = bill['cosponsors']
            
            # Add edges between sponsor and all cosponsors
            for cosponsor in cosponsors:
                if G.has_edge(sponsor, cosponsor):
                    # Increment weight
                    G[sponsor][cosponsor]['weight'] += 1
                else:
                    # Add new edge
                    G.add_edge(sponsor, cosponsor, weight=1)
        
        self.network = G
        return G
    
    def calculate_influence_metrics(self):
        """Calculate various influence metrics"""
        if self.network is None:
            self.build_network()
        
        metrics = {}
        
        # Degree centrality
        degree_cent = nx.degree_centrality(self.network)
        
        # Betweenness centrality
        between_cent = nx.betweenness_centrality(self.network, weight='weight')
        
        # Closeness centrality
        close_cent = nx.closeness_centrality(self.network, distance='weight')
        
        # PageRank
        pagerank = nx.pagerank(self.network, weight='weight')
        
        # Combine metrics
        for node in self.network.nodes():
            metrics[node] = {
                'degree_centrality': degree_cent[node],
                'betweenness_centrality': between_cent[node],
                'closeness_centrality': close_cent[node],
                'pagerank': pagerank[node],
                # Composite influence score
                'influence_score': (
                    degree_cent[node] * 0.3 +
                    between_cent[node] * 0.3 +
                    close_cent[node] * 0.2 +
                    pagerank[node] * 0.2
                )
            }
        
        return metrics
    
    def identify_top_influencers(self, n=10):
        """Identify most influential politicians"""
        metrics = self.calculate_influence_metrics()
        
        # Sort by influence score
        sorted_politicians = sorted(
            metrics.items(),
            key=lambda x: x[1]['influence_score'],
            reverse=True
        )
        
        return sorted_politicians[:n]
    
    def detect_communities(self):
        """Detect communities in the network"""
        if self.network is None:
            self.build_network()
        
        # Louvain community detection
        communities = community_louvain.best_partition(self.network)
        
        return communities
    
    def visualize_network(self, highlight_top_n=10):
        """Visualize co-sponsorship network"""
        if self.network is None:
            self.build_network()
        
        # Get top influencers
        top_influencers = self.identify_top_influencers(highlight_top_n)
        top_nodes = [node for node, _ in top_influencers]
        
        # Detect communities
        communities = self.detect_communities()
        
        # Create layout
        pos = nx.spring_layout(self.network, k=0.5, iterations=50)
        
        # Plot
        plt.figure(figsize=(16, 12))
        
        # Draw edges
        nx.draw_networkx_edges(
            self.network, pos,
            alpha=0.2,
            width=0.5
        )
        
        # Draw nodes colored by community
        for community_id in set(communities.values()):
            nodes_in_community = [
                n for n in self.network.nodes() 
                if communities[n] == community_id
            ]
            nx.draw_networkx_nodes(
                self.network, pos,
                nodelist=nodes_in_community,
                node_size=50,
                label=f'Community {community_id}'
            )
        
        # Highlight top influencers
        nx.draw_networkx_nodes(
            self.network, pos,
            nodelist=top_nodes,
            node_size=300,
            node_color='red',
            label='Top Influencers'
        )
        
        # Add labels for top influencers
        labels = {node: node for node in top_nodes}
        nx.draw_networkx_labels(
            self.network, pos,
            labels,
            font_size=10,
            font_weight='bold'
        )
        
        plt.title('Co-Sponsorship Network with Top Influencers')
        plt.legend()
        plt.axis('off')
        plt.tight_layout()
        plt.savefig('influence_network.png', dpi=300)
        plt.show()
    
    def generate_influence_report(self, n=20):
        """Generate comprehensive influence report"""
        top_influencers = self.identify_top_influencers(n)
        
        report = "# Influence Network Analysis Report\n\n"
        report += f"## Network Statistics\n"
        report += f"- Total politicians: {self.network.number_of_nodes()}\n"
        report += f"- Total co-sponsorships: {self.network.number_of_edges()}\n"
        report += f"- Network density: {nx.density(self.network):.3f}\n\n"
        
        report += f"## Top {n} Most Influential Politicians\n\n"
        
        for rank, (politician, metrics) in enumerate(top_influencers, 1):
            report += f"### {rank}. {politician}\n"
            report += f"- Influence Score: {metrics['influence_score']:.3f}\n"
            report += f"- Degree Centrality: {metrics['degree_centrality']:.3f}\n"
            report += f"- Betweenness Centrality: {metrics['betweenness_centrality']:.3f}\n"
            report += f"- PageRank: {metrics['pagerank']:.3f}\n\n"
        
        return report

# Usage example
# Sample data
cosponsorship_data = [
    {
        'bill_id': 'HR1',
        'sponsor': 'Sen_A',
        'cosponsors': ['Sen_B', 'Sen_C', 'Sen_D']
    },
    {
        'bill_id': 'HR2',
        'sponsor': 'Sen_B',
        'cosponsors': ['Sen_A', 'Sen_E']
    },
    # ... more bills
]

analyzer = InfluenceNetworkAnalysis(cosponsorship_data)
analyzer.build_network()

# Generate report
report = analyzer.generate_influence_report(n=10)
print(report)

# Visualize
analyzer.visualize_network(highlight_top_n=10)
```

## Best Practices for Analyses

### 1. Always Document Methodology

```python
def analyze_voting_trend(politician_id, start_date, end_date):
    """
    Analyze voting trend over time for a politician.
    
    Methodology:
    1. Fetch voting records in date range
    2. Calculate monthly bipartisan scores
    3. Fit linear regression to detect trend
    4. Report slope, R-squared, and p-value
    
    Args:
        politician_id: Bioguide ID
        start_date: Start of analysis period
        end_date: End of analysis period
    
    Returns:
        dict with trend statistics and visualization
    
    Limitations:
        - Excludes voice votes (no individual record)
        - Sensitive to short time periods
        - Does not account for changing political context
    """
    # Implementation
    pass
```

### 2. Include Confidence Intervals

```python
def report_bipartisan_score(score, total_votes):
    """Report score with confidence interval"""
    from statsmodels.stats.proportion import proportion_confint
    
    successes = int(score * total_votes / 100)
    ci_lower, ci_upper = proportion_confint(
        successes, total_votes, alpha=0.05, method='wilson'
    )
    
    return {
        'score': score,
        'ci_lower': ci_lower * 100,
        'ci_upper': ci_upper * 100,
        'n': total_votes
    }
```

### 3. Provide Context

```python
def contextualize_metric(value, historical_data):
    """Add historical context to metric"""
    percentile = stats.percentileofscore(historical_data, value)
    
    context = {
        'value': value,
        'percentile': percentile,
        'historical_mean': np.mean(historical_data),
        'historical_std': np.std(historical_data),
        'interpretation': ''
    }
    
    if percentile > 90:
        context['interpretation'] = 'Very high (top 10%)'
    elif percentile > 75:
        context['interpretation'] = 'High (top 25%)'
    elif percentile > 25:
        context['interpretation'] = 'Average'
    else:
        context['interpretation'] = 'Below average'
    
    return context
```

## Revision History

- v1.0 (Initial) - Comprehensive example analyses with code implementations
