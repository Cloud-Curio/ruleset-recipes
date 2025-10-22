# Data Collection Methods

## Overview

This document outlines standardized procedures for collecting political data from various government APIs and sources. Proper data collection is fundamental to reliable analysis and ensures reproducibility, completeness, and quality.

## Primary Data Sources

### 1. Congress.gov API

**Base URL**: `https://api.congress.gov/v3/`

**Authentication**: API key required
- Request key at: https://api.congress.gov/sign-up/
- Include in header: `X-Api-Key: YOUR_API_KEY`

#### Key Endpoints

**Bills and Resolutions**:
```
GET /bill/{congress}/{billType}/{billNumber}
GET /bill/{congress}/{billType}/{billNumber}/text
GET /bill/{congress}/{billType}/{billNumber}/actions
GET /bill/{congress}/{billType}/{billNumber}/amendments
GET /bill/{congress}/{billType}/{billNumber}/cosponsors
GET /bill/{congress}/{billType}/{billNumber}/relatedbills
GET /bill/{congress}/{billType}/{billNumber}/subjects
GET /bill/{congress}/{billType}/{billNumber}/summaries
GET /bill/{congress}/{billType}/{billNumber}/titles
```

**Members**:
```
GET /member
GET /member/{bioguideId}
GET /member/{bioguideId}/sponsored-legislation
GET /member/{bioguideId}/cosponsored-legislation
```

**Votes**:
```
GET /nomination/{congress}/{nominationNumber}
GET /nomination/{congress}/{nominationNumber}/actions
```

**Committees**:
```
GET /committee/{chamber}
GET /committee/{chamber}/{committeeCode}
GET /committee/{chamber}/{committeeCode}/bills
```

#### Rate Limits

- **Default**: 5,000 requests per hour
- **Recommended**: 1 request per second to avoid throttling
- **Headers**: Check `X-RateLimit-Remaining` and `X-RateLimit-Reset`

#### Data Collection Strategy

```python
import time
import requests
from typing import Dict, List, Optional

class CongressAPIClient:
    BASE_URL = "https://api.congress.gov/v3"
    
    def __init__(self, api_key: str):
        self.api_key = api_key
        self.session = requests.Session()
        self.session.headers.update({'X-Api-Key': api_key})
        self.last_request_time = 0
        self.min_interval = 1.0  # seconds between requests
    
    def _rate_limit(self):
        """Ensure minimum interval between requests"""
        elapsed = time.time() - self.last_request_time
        if elapsed < self.min_interval:
            time.sleep(self.min_interval - elapsed)
        self.last_request_time = time.time()
    
    def _make_request(self, endpoint: str, params: Dict = None) -> Dict:
        """Make API request with rate limiting and error handling"""
        self._rate_limit()
        
        url = f"{self.BASE_URL}/{endpoint}"
        
        try:
            response = self.session.get(url, params=params, timeout=30)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            if e.response.status_code == 429:
                # Rate limited - wait and retry
                retry_after = int(e.response.headers.get('Retry-After', 60))
                time.sleep(retry_after)
                return self._make_request(endpoint, params)
            else:
                raise
        except requests.exceptions.RequestException as e:
            # Log error and decide on retry strategy
            print(f"Request failed: {e}")
            raise
    
    def get_bill(self, congress: int, bill_type: str, 
                 bill_number: int) -> Dict:
        """Retrieve bill details"""
        endpoint = f"bill/{congress}/{bill_type}/{bill_number}"
        return self._make_request(endpoint)
    
    def get_bill_text(self, congress: int, bill_type: str, 
                      bill_number: int) -> Dict:
        """Retrieve bill full text"""
        endpoint = f"bill/{congress}/{bill_type}/{bill_number}/text"
        return self._make_request(endpoint)
    
    def get_member(self, bioguide_id: str) -> Dict:
        """Retrieve member information"""
        endpoint = f"member/{bioguide_id}"
        return self._make_request(endpoint)
    
    def get_bills_by_congress(self, congress: int, 
                              limit: int = 250) -> List[Dict]:
        """Retrieve all bills for a congress (paginated)"""
        bills = []
        offset = 0
        
        while True:
            params = {'limit': limit, 'offset': offset}
            endpoint = f"bill/{congress}"
            data = self._make_request(endpoint, params)
            
            batch = data.get('bills', [])
            bills.extend(batch)
            
            # Check if more results available
            if len(batch) < limit:
                break
            
            offset += limit
        
        return bills
```

### 2. GovInfo.gov API

**Base URL**: `https://api.govinfo.gov/`

**Authentication**: API key required
- Request key at: https://api.govinfo.gov/docs/
- Include in query param: `api_key=YOUR_API_KEY`

#### Key Collections

**Congressional Bills**:
```
GET /collections/BILLS/{congress}/json
GET /packages/BILLS-{congress}{billtype}{number}/summary
GET /packages/BILLS-{congress}{billtype}{number}/pdf
```

**Congressional Record**:
```
GET /collections/CREC/{publishDate}/json
```

**Federal Register**:
```
GET /collections/FR/{publishDate}/json
```

#### Bulk Data Downloads

For large-scale data collection:
```python
def download_bulk_data(collection: str, congress: int, 
                       output_dir: str):
    """Download bulk data for entire congress"""
    base_url = f"https://api.govinfo.gov/collections/{collection}"
    
    # Get list of packages
    params = {'congress': congress, 'api_key': API_KEY}
    response = requests.get(base_url, params=params)
    packages = response.json().get('packages', [])
    
    # Download each package
    for package in packages:
        package_id = package['packageId']
        
        # Download metadata
        metadata_url = f"{base_url}/{package_id}/summary"
        metadata = requests.get(metadata_url, 
                               params={'api_key': API_KEY}).json()
        
        # Download PDF
        pdf_url = f"{base_url}/{package_id}/pdf"
        pdf_response = requests.get(pdf_url, 
                                   params={'api_key': API_KEY})
        
        # Save files
        save_path = f"{output_dir}/{package_id}"
        with open(f"{save_path}.json", 'w') as f:
            json.dump(metadata, f)
        with open(f"{save_path}.pdf", 'wb') as f:
            f.write(pdf_response.content)
        
        # Rate limiting
        time.sleep(1)
```

### 3. OpenStates GraphQL API

**Base URL**: `https://v3.openstates.org/graphql`

**Authentication**: API key in header
- Get key at: https://openstates.org/accounts/profile/
- Header: `X-API-KEY: YOUR_API_KEY`

#### Example Queries

**Get State Bills**:
```graphql
query GetBills($jurisdiction: String!, $session: String!) {
  bills(jurisdiction: $jurisdiction, session: $session, first: 100) {
    edges {
      node {
        id
        identifier
        title
        classification
        subject
        abstracts {
          abstract
        }
        sponsorships {
          name
          classification
        }
        actions {
          description
          date
          classification
        }
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

**Get Legislator**:
```graphql
query GetLegislator($id: String!) {
  person(id: $id) {
    id
    name
    currentMemberships {
      organization {
        name
        classification
      }
      post {
        label
      }
    }
    contactDetails {
      type
      value
    }
  }
}
```

#### Python Client

```python
from gql import gql, Client
from gql.transport.requests import RequestsHTTPTransport

class OpenStatesClient:
    def __init__(self, api_key: str):
        transport = RequestsHTTPTransport(
            url='https://v3.openstates.org/graphql',
            headers={'X-API-KEY': api_key}
        )
        self.client = Client(transport=transport, 
                            fetch_schema_from_transport=True)
    
    def get_bills(self, jurisdiction: str, session: str):
        query = gql("""
            query GetBills($jurisdiction: String!, $session: String!) {
                bills(jurisdiction: $jurisdiction, 
                      session: $session, first: 100) {
                    edges {
                        node {
                            id
                            identifier
                            title
                        }
                    }
                }
            }
        """)
        
        params = {
            'jurisdiction': jurisdiction,
            'session': session
        }
        
        result = self.client.execute(query, variable_values=params)
        return result['bills']['edges']
```

## Data Quality Procedures

### 1. Validation Rules

**Required Fields Check**:
```python
def validate_bill_data(bill: Dict) -> bool:
    """Validate bill has required fields"""
    required_fields = [
        'congress',
        'type',
        'number',
        'title',
        'introducedDate'
    ]
    
    for field in required_fields:
        if field not in bill or bill[field] is None:
            return False
    
    return True
```

**Data Type Validation**:
```python
def validate_data_types(bill: Dict) -> List[str]:
    """Check data types match expected schema"""
    errors = []
    
    if not isinstance(bill.get('congress'), int):
        errors.append('Congress must be integer')
    
    if not isinstance(bill.get('cosponsors'), list):
        errors.append('Cosponsors must be list')
    
    # Add more validations as needed
    
    return errors
```

### 2. Completeness Checks

**Missing Data Detection**:
```python
def check_completeness(bills: List[Dict]) -> Dict:
    """Analyze data completeness across collection"""
    stats = {
        'total_bills': len(bills),
        'missing_text': 0,
        'missing_sponsors': 0,
        'missing_subjects': 0
    }
    
    for bill in bills:
        if not bill.get('text'):
            stats['missing_text'] += 1
        if not bill.get('sponsors'):
            stats['missing_sponsors'] += 1
        if not bill.get('subjects'):
            stats['missing_subjects'] += 1
    
    # Calculate percentages
    for key in ['missing_text', 'missing_sponsors', 'missing_subjects']:
        pct = (stats[key] / stats['total_bills']) * 100
        stats[f'{key}_pct'] = round(pct, 2)
    
    return stats
```

### 3. Duplicate Detection

**Identify Duplicates**:
```python
def find_duplicates(bills: List[Dict]) -> List[tuple]:
    """Find duplicate bills in collection"""
    seen = {}
    duplicates = []
    
    for bill in bills:
        # Create unique identifier
        bill_id = f"{bill['congress']}-{bill['type']}-{bill['number']}"
        
        if bill_id in seen:
            duplicates.append((bill_id, seen[bill_id], bill))
        else:
            seen[bill_id] = bill
    
    return duplicates
```

### 4. Consistency Checks

**Cross-reference Validation**:
```python
def validate_relationships(bill: Dict, members: Dict) -> List[str]:
    """Verify relationships are consistent"""
    errors = []
    
    # Check sponsor exists in members database
    if bill.get('sponsor'):
        sponsor_id = bill['sponsor']['bioguideId']
        if sponsor_id not in members:
            errors.append(f'Unknown sponsor: {sponsor_id}')
    
    # Check cosponsors exist
    for cosponsor in bill.get('cosponsors', []):
        cosponsor_id = cosponsor['bioguideId']
        if cosponsor_id not in members:
            errors.append(f'Unknown cosponsor: {cosponsor_id}')
    
    return errors
```

## Data Storage

### Database Schema

```sql
-- Raw data storage
CREATE TABLE raw_bills (
    id SERIAL PRIMARY KEY,
    api_source VARCHAR(50),
    collection_date TIMESTAMP,
    congress INTEGER,
    bill_type VARCHAR(10),
    bill_number INTEGER,
    raw_json JSONB,
    UNIQUE(api_source, congress, bill_type, bill_number)
);

-- Processed bills
CREATE TABLE bills (
    id SERIAL PRIMARY KEY,
    raw_bill_id INTEGER REFERENCES raw_bills(id),
    congress INTEGER,
    bill_type VARCHAR(10),
    bill_number INTEGER,
    title TEXT,
    summary TEXT,
    introduced_date DATE,
    status VARCHAR(50),
    sponsor_bioguide_id VARCHAR(10),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    UNIQUE(congress, bill_type, bill_number)
);

-- Data quality tracking
CREATE TABLE data_quality_log (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(100),
    record_id INTEGER,
    check_type VARCHAR(50),
    status VARCHAR(20),
    error_message TEXT,
    checked_at TIMESTAMP
);
```

### File Storage Structure

```
data/
├── raw/                          # Original API responses
│   ├── congress/
│   │   ├── 118/
│   │   │   ├── bills/
│   │   │   │   ├── hr/
│   │   │   │   │   └── hr1.json
│   │   │   │   └── s/
│   │   │   ├── members/
│   │   │   └── votes/
│   │   └── 117/
│   ├── govinfo/
│   └── openstates/
├── processed/                    # Cleaned and transformed
│   ├── bills/
│   ├── members/
│   └── votes/
└── metadata/                     # Collection metadata
    ├── collection_runs.json
    └── quality_reports/
```

## Data Collection Workflows

### Initial Data Load

```python
async def initial_data_load(congress: int):
    """Complete data collection for a congress"""
    
    # 1. Collect bills
    print(f"Collecting bills for Congress {congress}...")
    bills = await collect_bills(congress)
    print(f"Collected {len(bills)} bills")
    
    # 2. For each bill, get full details
    for bill in bills:
        bill_details = await get_full_bill_data(
            congress, 
            bill['type'], 
            bill['number']
        )
        await save_bill(bill_details)
    
    # 3. Collect members
    print("Collecting members...")
    members = await collect_members(congress)
    for member in members:
        member_details = await get_full_member_data(member['bioguideId'])
        await save_member(member_details)
    
    # 4. Collect votes
    print("Collecting votes...")
    votes = await collect_votes(congress)
    for vote in votes:
        await save_vote(vote)
    
    # 5. Generate quality report
    report = await generate_quality_report(congress)
    print(report)
```

### Incremental Updates

```python
async def incremental_update():
    """Daily update to collect new/changed data"""
    
    # 1. Get current congress
    current_congress = get_current_congress()
    
    # 2. Fetch bills updated since last run
    last_update = get_last_update_time()
    updated_bills = await get_updated_bills(
        congress=current_congress,
        since=last_update
    )
    
    # 3. Update bill records
    for bill in updated_bills:
        await update_bill(bill)
    
    # 4. Check for new votes
    new_votes = await get_new_votes(since=last_update)
    for vote in new_votes:
        await save_vote(vote)
    
    # 5. Update last run timestamp
    update_last_run_time()
```

### Error Handling and Retry Logic

```python
import asyncio
from typing import Callable, Any

async def retry_with_backoff(
    func: Callable,
    max_retries: int = 3,
    base_delay: float = 1.0,
    *args,
    **kwargs
) -> Any:
    """Retry function with exponential backoff"""
    
    for attempt in range(max_retries):
        try:
            return await func(*args, **kwargs)
        except Exception as e:
            if attempt == max_retries - 1:
                # Last attempt failed, raise error
                raise
            
            # Calculate delay with exponential backoff
            delay = base_delay * (2 ** attempt)
            print(f"Attempt {attempt + 1} failed: {e}")
            print(f"Retrying in {delay} seconds...")
            await asyncio.sleep(delay)
```

## Monitoring and Logging

### Collection Metrics

```python
class CollectionMetrics:
    def __init__(self):
        self.start_time = time.time()
        self.records_processed = 0
        self.errors = 0
        self.api_calls = 0
    
    def record_success(self):
        self.records_processed += 1
    
    def record_error(self, error: Exception):
        self.errors += 1
        # Log error details
        logging.error(f"Collection error: {error}")
    
    def record_api_call(self):
        self.api_calls += 1
    
    def get_summary(self) -> Dict:
        elapsed = time.time() - self.start_time
        return {
            'duration_seconds': round(elapsed, 2),
            'records_processed': self.records_processed,
            'errors': self.errors,
            'api_calls': self.api_calls,
            'records_per_second': round(
                self.records_processed / elapsed, 2
            )
        }
```

### Logging Configuration

```python
import logging

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('data_collection.log'),
        logging.StreamHandler()
    ]
)

logger = logging.getLogger('data_collection')

# Log collection events
logger.info(f"Starting collection for Congress {congress}")
logger.info(f"Collected {len(bills)} bills")
logger.error(f"Failed to fetch bill: {bill_id}")
logger.warning(f"Missing required field: {field_name}")
```

## Best Practices

### 1. API Usage

- **Respect rate limits**: Implement proper throttling
- **Cache responses**: Avoid redundant API calls
- **Handle errors gracefully**: Implement retry logic
- **Monitor API changes**: Subscribe to API update notifications

### 2. Data Storage

- **Store raw data**: Keep original API responses
- **Version data**: Track when data was collected
- **Document transformations**: Record all data processing steps
- **Backup regularly**: Implement backup strategy

### 3. Quality Assurance

- **Validate continuously**: Check data quality at collection time
- **Monitor completeness**: Track missing data patterns
- **Generate reports**: Regular quality assessment reports
- **Alert on issues**: Notify when quality degrades

### 4. Performance

- **Use async/await**: Parallel data collection
- **Batch requests**: Group related API calls
- **Implement caching**: Redis for frequently accessed data
- **Optimize queries**: Index database properly

## Troubleshooting

### Common Issues

**Rate Limiting**:
```python
# Solution: Implement backoff strategy
if response.status_code == 429:
    retry_after = int(response.headers.get('Retry-After', 60))
    time.sleep(retry_after)
    # Retry request
```

**Incomplete Data**:
```python
# Solution: Retry failed collections
failed_bills = get_failed_collections()
for bill in failed_bills:
    try:
        recollect_bill(bill)
    except Exception as e:
        log_permanent_failure(bill, e)
```

**API Changes**:
```python
# Solution: Version your API client
class CongressAPIClientV3:
    API_VERSION = "v3"
    # Implement versioned client
```

## References

### API Documentation
- Congress.gov API: https://api.congress.gov/
- GovInfo.gov API: https://www.govinfo.gov/help/api
- OpenStates API: https://docs.openstates.org/api-v3/

### Tools and Libraries
- `requests`: HTTP library for Python
- `aiohttp`: Async HTTP client
- `gql`: GraphQL client
- `pandas`: Data manipulation
- `sqlalchemy`: Database ORM

## Revision History

- v1.0 (Initial) - Comprehensive data collection methodology
