# Chapter 10: Full-Text Search

## What is Full-Text Search?

Full-text search (FTS) finds documents that match natural language queries, with ranking by relevance. It's smarter than `LIKE` — it understands language:
- **Stemming:** "running" matches "run", "runs", "runner"
- **Stop words:** ignores "the", "a", "is"
- **Ranking:** most relevant results first
- **Phrase search:** exact phrase matching
- **Dictionaries:** language-specific rules

---

## Core Types

```sql
-- tsvector: preprocessed document representation
SELECT to_tsvector('english', 'The quick brown fox jumps over the lazy dog');
-- 'brown':3 'dog':9 'fox':4 'jump':5 'lazi':8 'quick':2
-- (stop words "the", "over" removed; "jumps"→"jump", "lazy"→"lazi")

-- tsquery: search query
SELECT to_tsquery('english', 'quick & fox');         -- AND
SELECT to_tsquery('english', 'quick | fox');         -- OR
SELECT to_tsquery('english', 'quick & !cat');        -- NOT
SELECT plainto_tsquery('english', 'quick brown fox'); -- converts plain text
SELECT phraseto_tsquery('english', 'quick brown fox'); -- phrase search
SELECT websearch_to_tsquery('english', '"quick brown" -cat'); -- web-style
```

---

## Basic Full-Text Search

```sql
-- Simple search without index
SELECT title, body
FROM posts
WHERE to_tsvector('english', title || ' ' || body) @@ to_tsquery('english', 'postgresql');

-- Using plainto_tsquery (easier for user input)
SELECT title
FROM posts
WHERE to_tsvector('english', title || ' ' || body) @@ plainto_tsquery('english', 'learn postgresql tutorial');

-- Ranking results
SELECT
    title,
    ts_rank(to_tsvector('english', title || ' ' || body), query) AS rank
FROM posts,
     to_tsquery('english', 'postgresql') query
WHERE to_tsvector('english', title || ' ' || body) @@ query
ORDER BY rank DESC
LIMIT 10;
```

---

## Adding a tsvector Column (Recommended)

For performance, store the computed tsvector in a column:

```sql
-- Add tsvector column
ALTER TABLE posts ADD COLUMN search_vector TSVECTOR;

-- Populate it
UPDATE posts SET search_vector =
    setweight(to_tsvector('english', COALESCE(title, '')), 'A') ||
    setweight(to_tsvector('english', COALESCE(body, '')), 'B') ||
    setweight(to_tsvector('english', COALESCE(tags::TEXT, '')), 'C');
-- setweight assigns relevance: A > B > C > D
-- Matches in title rank higher than body

-- GIN index on the column
CREATE INDEX idx_posts_search ON posts USING GIN(search_vector);

-- Trigger to keep it updated automatically
CREATE OR REPLACE FUNCTION posts_search_trigger()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.body, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER posts_search_update
BEFORE INSERT OR UPDATE ON posts
FOR EACH ROW EXECUTE FUNCTION posts_search_trigger();
```

---

## Searching with the Index

```sql
-- Fast GIN index search
SELECT id, title,
    ts_rank(search_vector, query) AS rank
FROM posts,
     plainto_tsquery('english', 'learn postgresql') query
WHERE search_vector @@ query
ORDER BY rank DESC
LIMIT 20;

-- Web-style search (supports quotes and minus)
SELECT id, title, ts_rank(search_vector, query) AS rank
FROM posts,
     websearch_to_tsquery('english', '"full text" -mysql') query
WHERE search_vector @@ query
ORDER BY rank DESC;
```

---

## Ranking Functions

```sql
-- ts_rank: ranks based on term frequency
SELECT ts_rank(search_vector, query) AS rank

-- ts_rank_cd: ranks using cover density (considers proximity of terms)
SELECT ts_rank_cd(search_vector, query) AS rank

-- Normalize rank by document length (prevents long docs from always winning)
SELECT ts_rank(search_vector, query, 1) AS rank  -- 1 = divide by 1 + log(length)
SELECT ts_rank(search_vector, query, 32) AS rank -- 32 = divide by ndoc (# documents)
```

---

## Highlighting — ts_headline

```sql
SELECT
    id,
    title,
    ts_headline(
        'english',
        body,
        plainto_tsquery('english', 'postgresql index'),
        'MaxFragments=3, MaxWords=15, MinWords=5, StartSel=<b>, StopSel=</b>'
    ) AS excerpt
FROM posts
WHERE search_vector @@ plainto_tsquery('english', 'postgresql index')
ORDER BY ts_rank(search_vector, plainto_tsquery('english', 'postgresql index')) DESC;

-- Result: "...using a <b>PostgreSQL</b> <b>index</b> for full-text..."
```

---

## Phrase Search

```sql
-- Exact phrase search
SELECT * FROM posts
WHERE search_vector @@ phraseto_tsquery('english', 'getting started with postgresql');

-- Proximity search — within N positions of each other
SELECT * FROM posts
WHERE search_vector @@ to_tsquery('english', 'getting <-> started');  -- adjacent
WHERE search_vector @@ to_tsquery('english', 'getting <3> postgresql'); -- within 3 words
```

---

## Dictionaries and Languages

```sql
-- List available dictionaries
SELECT dictname FROM pg_ts_dict;

-- Available languages
SELECT cfgname FROM pg_ts_config;

-- Use different languages
SELECT to_tsvector('spanish', 'Los perros corren en el parque');
SELECT to_tsvector('german', 'Die Hunde laufen im Park');
SELECT to_tsvector('french', 'Les chiens courent dans le parc');

-- Simple dictionary (no stemming, just lowercase)
SELECT to_tsvector('simple', 'The Quick Brown Fox');
-- 'brown':3 'fox':4 'quick':2 'the':1  (no stemming, no stop words)
```

---

## Multi-language Search

```sql
CREATE TABLE articles (
    id BIGSERIAL PRIMARY KEY,
    title TEXT,
    body TEXT,
    language TEXT DEFAULT 'english',
    search_vector TSVECTOR
);

-- Trigger with dynamic language
CREATE OR REPLACE FUNCTION articles_search_trigger()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector(NEW.language::regconfig, COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector(NEW.language::regconfig, COALESCE(NEW.body, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Search in specific language
SELECT * FROM articles
WHERE search_vector @@ plainto_tsquery('spanish', 'perros')
AND language = 'spanish';
```

---

## Full-Text Search Performance Tips

```sql
-- Always use GIN index for full-text search (not B-tree, not GiST unless using phrase)
CREATE INDEX idx_fts ON posts USING GIN(search_vector);

-- Use GiST if you need phrase search (slower build, supports <->)
CREATE INDEX idx_fts_gist ON posts USING GIST(search_vector);

-- Check if index is being used
EXPLAIN ANALYZE
SELECT * FROM posts
WHERE search_vector @@ plainto_tsquery('english', 'postgresql')
LIMIT 10;
-- Should show: Bitmap Index Scan on idx_fts

-- Use LIMIT to avoid ranking millions of rows
SELECT id, title, ts_rank(search_vector, q) AS rank
FROM posts, plainto_tsquery('english', 'postgresql') q
WHERE search_vector @@ q
ORDER BY rank DESC
LIMIT 20;   -- rank only what you need
```

---

## When to Use PostgreSQL FTS vs Elasticsearch

| Feature | PostgreSQL FTS | Elasticsearch |
|---|---|---|
| Setup | None (built-in) | Separate cluster |
| Sync | Automatic (trigger) | ETL/CDC required |
| ACID | Yes | No |
| Performance | Good for millions | Better for billions |
| Aggregations | Basic | Advanced |
| Fuzzy search | Limited | Excellent |
| Relevance tuning | Limited | Extensive |
| Cost | Free (already running) | Additional infra |

**Use PostgreSQL FTS when:** you have < 10M documents, already use PostgreSQL, need transactional consistency.

**Use Elasticsearch when:** you have billions of documents, need advanced relevance tuning, or need complex aggregation analytics.
