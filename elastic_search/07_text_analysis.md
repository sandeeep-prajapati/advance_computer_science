# Chapter 7: Text Analysis & Analyzers

## What is Text Analysis?

When you index a `text` field, Elasticsearch doesn't store the raw string — it runs the string through an **analyzer** that transforms it into **tokens** (terms) stored in the inverted index.

The same analyzer is applied at **index time** and **search time** to ensure consistent matching.

```
Input:    "Quick Brown FOX jumped over the lazy Dog!"
          ↓ Analyzer
Tokens:   ["quick", "brown", "fox", "jump", "lazi", "dog"]
          ↓ Stored in inverted index
```

---

## Anatomy of an Analyzer

An analyzer is a pipeline with three components:

```
Input Text
    ↓
Character Filters   (strip HTML, replace chars)
    ↓
Tokenizer           (split into tokens — exactly ONE tokenizer)
    ↓
Token Filters       (lowercase, stop words, stemming — zero or more)
    ↓
Token Stream        (stored in inverted index)
```

---

## Built-In Analyzers

### `standard` (Default)

Splits on whitespace and punctuation, lowercases, removes common stop words.

```bash
POST /_analyze
{
  "analyzer": "standard",
  "text": "Quick Brown FOX jumped over the lazy Dog!"
}
```

Tokens: `["quick", "brown", "fox", "jumped", "over", "the", "lazy", "dog"]`

### `simple`

Splits on any non-letter character, lowercases. No stop words.

Tokens: `["quick", "brown", "fox", "jumped", "over", "the", "lazy", "dog"]`

### `whitespace`

Splits on whitespace only. Preserves case and punctuation.

Tokens: `["Quick", "Brown", "FOX!", "jumped", "over", "the", "lazy", "Dog!"]`

### `english` (and other language analyzers)

Applies stemming and language-specific stop words:

```bash
POST /_analyze
{
  "analyzer": "english",
  "text": "Running quickly and jumping over fences"
}
```

Tokens: `["run", "quick", "jump", "over", "fenc"]` — stemmed to root forms.

### `keyword`

No tokenization — entire string as one token. Used with `keyword` type fields.

---

## Testing Analyzers

Use the `_analyze` API to test how text is tokenized:

```bash
POST /_analyze
{
  "analyzer": "standard",
  "text": "My email is user@example.com"
}

# Test with a specific field's analyzer
POST /products/_analyze
{
  "field": "description",
  "text": "Wireless Bluetooth Headphones"
}

# Test a custom pipeline
POST /_analyze
{
  "tokenizer": "standard",
  "filter": ["lowercase", "stop", "snowball"],
  "text": "Running quickly through the forest"
}
```

---

## Custom Analyzers

Define custom analyzers in index settings:

```bash
PUT /my_index
{
  "settings": {
    "analysis": {
      "analyzer": {
        "my_custom_analyzer": {
          "type":            "custom",
          "char_filter":     ["html_strip"],
          "tokenizer":       "standard",
          "filter":          ["lowercase", "stop", "my_stemmer"]
        },
        "autocomplete_analyzer": {
          "type":      "custom",
          "tokenizer": "standard",
          "filter":    ["lowercase", "autocomplete_filter"]
        },
        "autocomplete_search": {
          "type":      "custom",
          "tokenizer": "standard",
          "filter":    ["lowercase"]
        }
      },
      "filter": {
        "my_stemmer": {
          "type":     "stemmer",
          "language": "english"
        },
        "autocomplete_filter": {
          "type":     "edge_ngram",
          "min_gram": 1,
          "max_gram": 20
        }
      }
    }
  }
}
```

---

## Common Tokenizers

| Tokenizer | Description | Example Input → Tokens |
|---|---|---|
| `standard` | Word boundaries, punctuation | `"hello world"` → `[hello, world]` |
| `whitespace` | Split on whitespace only | `"foo bar"` → `[foo, bar]` |
| `keyword` | Entire string as one token | `"foo bar"` → `[foo bar]` |
| `letter` | Split on non-letter chars | `"foo123bar"` → `[foo, bar]` |
| `ngram` | Generate n-grams | `"cat"` (n=2) → `[ca, at]` |
| `edge_ngram` | N-grams from start | `"cat"` → `[c, ca, cat]` |
| `pattern` | Split on regex | Split on comma: `"a,b,c"` → `[a, b, c]` |
| `path_hierarchy` | File path segments | `"/a/b/c"` → `[/a, /a/b, /a/b/c]` |

---

## Common Token Filters

| Filter | Description |
|---|---|
| `lowercase` | Converts tokens to lowercase |
| `stop` | Removes common stop words (the, is, are) |
| `stemmer` | Reduces words to root form (running→run) |
| `snowball` | Snowball stemmer algorithm |
| `synonym` | Replace tokens with synonyms |
| `asciifolding` | Convert accented chars (café→cafe) |
| `unique` | Remove duplicate tokens |
| `length` | Remove tokens outside min/max length |
| `edge_ngram` | Generate edge n-grams (for autocomplete) |
| `ngram` | Generate n-grams |
| `shingle` | Generate word n-grams (bigrams, trigrams) |

---

## Character Filters

| Filter | Description |
|---|---|
| `html_strip` | Remove HTML tags |
| `mapping` | Replace specific characters/strings |
| `pattern_replace` | Replace using regex |

```bash
# Remove HTML and replace & with and
{
  "char_filter": ["html_strip"],
  "text": "<p>Soup &amp; Salad</p>"
}
# → "Soup & Salad"
```

---

## Autocomplete with Edge N-Gram

A classic use case: "as-you-type" search.

```bash
PUT /products
{
  "settings": {
    "analysis": {
      "filter": {
        "autocomplete_filter": {
          "type":     "edge_ngram",
          "min_gram": 2,
          "max_gram": 20
        }
      },
      "analyzer": {
        "autocomplete": {
          "type":      "custom",
          "tokenizer": "standard",
          "filter":    ["lowercase", "autocomplete_filter"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "name": {
        "type": "text",
        "analyzer": "autocomplete",
        "search_analyzer": "standard"
      }
    }
  }
}
```

Index time: `"Laptop"` → `["la", "lap", "lapt", "lapto", "laptop"]`
Search time: `"lapt"` → `["lapt"]` — matched exactly.

---

## Synonym Filters

```bash
PUT /products
{
  "settings": {
    "analysis": {
      "filter": {
        "tech_synonyms": {
          "type": "synonym",
          "synonyms": [
            "phone, mobile, smartphone",
            "laptop => notebook, portable computer",
            "tv, television, smart tv"
          ]
        }
      },
      "analyzer": {
        "synonym_analyzer": {
          "tokenizer": "standard",
          "filter":    ["lowercase", "tech_synonyms"]
        }
      }
    }
  }
}
```

Now searching for "phone" also matches documents containing "mobile" and "smartphone".

---

## Field Mapping with Multiple Analyzers

You can use different analyzers for the same field:

```bash
"properties": {
  "title": {
    "type":            "text",
    "analyzer":        "english",         ← index-time analyzer
    "search_analyzer": "standard"         ← query-time analyzer
  }
}
```

---

## Multi-Fields: text + keyword

A very common pattern — store a field as both analyzed `text` and exact `keyword`:

```bash
"properties": {
  "name": {
    "type": "text",
    "analyzer": "standard",
    "fields": {
      "keyword": {
        "type": "keyword",
        "ignore_above": 256
      },
      "english": {
        "type": "text",
        "analyzer": "english"
      }
    }
  }
}
```

Usage:
- `name` → full-text search (analyzed)
- `name.keyword` → exact match, sorting, aggregations
- `name.english` → stemmed English search

---

## Update Analyzers on Existing Index

You **cannot change analyzers without reindexing**, but you can add new analyzers:

```bash
# 1. Close the index
POST /products/_close

# 2. Update settings
PUT /products/_settings
{
  "analysis": {
    "analyzer": {
      "new_analyzer": { ... }
    }
  }
}

# 3. Reopen
POST /products/_open

# 4. Reindex existing docs to apply new analyzer
POST /products/_update_by_query?conflicts=proceed
```
