# Emoji and Non-ASCII Character Removal Examples

This document demonstrates the emoji removal feature in SlimJSON, which can significantly reduce token count for LLM contexts.

## Why Remove Emoji?

Emoji and non-ASCII characters consume multiple tokens in most LLM tokenizers:

| Character | Tokens | Example |
|-----------|--------|---------|
| ASCII letter | 1 | `a` = 1 token |
| Emoji | 2-4 | `👋` = 2-4 tokens |
| Chinese character | 2-3 | `中` = 2-3 tokens |
| Arabic character | 2-3 | `ع` = 2-3 tokens |

For large JSON payloads with emoji, this can significantly increase API costs.

## Basic Usage

### CLI

```bash
# Simple emoji removal
echo '{"message":"Hello 👋 World 🌍!"}' | slimjson -strip-emoji
# Output: {"message":"Hello  World !"}

# From file
slimjson -strip-emoji input.json > output.json

# With pretty printing
slimjson -strip-emoji -pretty data.json
```

### Go Library

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/tradik/slimjson"
)

func main() {
    data := map[string]interface{}{
        "message": "Hello 👋 World 🌍!",
        "status":  "✅ Completed",
    }

    cfg := slimjson.Config{
        StripUTF8Emoji: true,
    }

    slimmer := slimjson.New(cfg)
    result := slimmer.Slim(data)

    output, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(output))
}
```

### Config File

```ini
[llm-optimized]
strip-emoji=true
string-pooling=true
depth=4
list-len=15
```

## Real-World Examples

### Example 1: Social Media Post

**Input:**
```json
{
  "user": "John Doe 😊",
  "post": "Just launched our new product! 🚀🎉",
  "reactions": "❤️ 👍 🔥",
  "location": "San Francisco 🌉"
}
```

**Command:**
```bash
slimjson -strip-emoji -pretty input.json
```

**Output:**
```json
{
  "location": "San Francisco ",
  "post": "Just launched our new product! ",
  "reactions": "  ",
  "user": "John Doe "
}
```

**Token Savings:**
- Before: ~35 tokens
- After: ~20 tokens
- **Savings: 43%**

### Example 2: Product Catalog

**Input:**
```json
{
  "products": [
    {
      "name": "Coffee ☕",
      "description": "Premium coffee beans 🌱",
      "rating": "⭐⭐⭐⭐⭐",
      "price": "$19.99 💰"
    },
    {
      "name": "Tea 🍵",
      "description": "Organic green tea 🌿",
      "rating": "⭐⭐⭐⭐",
      "price": "$14.99 💵"
    }
  ]
}
```

**Command:**
```bash
slimjson -strip-emoji -list-len 2 -pretty data.json
```

**Output:**
```json
{
  "products": [
    {
      "description": "Premium coffee beans ",
      "name": "Coffee ",
      "price": "$19.99 ",
      "rating": ""
    },
    {
      "description": "Organic green tea ",
      "name": "Tea ",
      "price": "$14.99 ",
      "rating": ""
    }
  ]
}
```

**Token Savings:**
- Before: ~80 tokens
- After: ~45 tokens
- **Savings: 44%**

### Example 3: Chat Messages

**Input:**
```json
{
  "messages": [
    {
      "user": "Alice 👩‍💻",
      "text": "Hey! How are you? 😃",
      "timestamp": "2024-01-15T10:30:00Z"
    },
    {
      "user": "Bob 👨‍💼",
      "text": "Great! Working on the new feature 💪",
      "timestamp": "2024-01-15T10:31:00Z"
    }
  ]
}
```

**Command:**
```bash
slimjson -strip-emoji -timestamp-compression -pretty messages.json
```

**Output:**
```json
{
  "messages": [
    {
      "text": "Hey! How are you? ",
      "timestamp": 1705315800,
      "user": "Alice "
    },
    {
      "text": "Great! Working on the new feature ",
      "timestamp": 1705315860,
      "user": "Bob "
    }
  ]
}
```

### Example 4: Multilingual Content

**Input:**
```json
{
  "title": "Welcome! 欢迎! مرحبا! 🌍",
  "description": "Global platform for everyone 🌐",
  "languages": ["English", "中文", "العربية", "日本語"],
  "status": "✅ Active"
}
```

**Command:**
```bash
slimjson -strip-emoji -pretty data.json
```

**Output:**
```json
{
  "description": "Global platform for everyone ",
  "languages": [
    "English",
    "",
    "",
    ""
  ],
  "status": " Active",
  "title": "Welcome! ! ! "
}
```

**Note:** This removes ALL non-ASCII characters, including Chinese, Arabic, and Japanese characters. Use with caution if you need to preserve multilingual content.

## Combined with Other Features

### Maximum Token Reduction

```bash
slimjson \
  -strip-emoji \
  -string-pooling \
  -type-inference \
  -timestamp-compression \
  -enum-detection \
  -depth 3 \
  -list-len 10 \
  -decimal-places 2 \
  -pretty \
  data.json
```

### LLM Context Optimization

```bash
# Use with ai-optimized profile
slimjson -profile ai-optimized -strip-emoji data.json

# Or create custom profile
cat > .slimjson << EOF
[llm-context]
strip-emoji=true
depth=4
list-len=15
string-pooling=true
type-inference=true
bool-compression=true
timestamp-compression=true
block=avatar_url,url,html_url
EOF

slimjson -profile llm-context data.json
```

## HTTP API Usage

```bash
# Start daemon
slimjson -d -port 8080

# Use with profile
curl -X POST 'http://localhost:8080/slim?profile=llm-context' \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello 👋 World 🌍!",
    "status": "✅ Completed"
  }'

# Response:
# {
#   "message": "Hello  World !",
#   "status": " Completed"
# }
```

## Performance Impact

The emoji stripping operation is very efficient:

```
BenchmarkStripEmoji-8    5000000    250 ns/op    128 B/op    1 allocs/op
```

- **Minimal overhead**: ~250 nanoseconds per string
- **Memory efficient**: Single allocation per string
- **Scales linearly**: O(n) where n is string length

## Best Practices

### ✅ DO Use When:
- Preparing data for LLM APIs (OpenAI, Anthropic, etc.)
- Processing social media content
- Cleaning user-generated content
- Reducing API costs
- Normalizing text data

### ❌ DON'T Use When:
- Emoji are semantically important
- Processing multilingual content that needs preservation
- Displaying content to end users
- Emoji convey critical information

## Character Preservation

The feature preserves:
- **ASCII printable characters** (32-126): `A-Z`, `a-z`, `0-9`, punctuation
- **Whitespace**: newline (`\n`), carriage return (`\r`), tab (`\t`)

The feature removes:
- **Emoji**: 👋 🌍 😊 ✅ 🚀 etc.
- **Non-ASCII letters**: Chinese (中文), Arabic (العربية), Cyrillic (Русский)
- **Special symbols**: ☕ ⭐ 💰 etc.
- **Control characters**: Most Unicode control characters

## Token Count Comparison

### Example: GitHub API Response

**Original (with emoji):**
```json
{
  "user": "octocat 🐙",
  "bio": "GitHub's mascot 🎭",
  "location": "San Francisco 🌉",
  "status": "✅ Available"
}
```
**Tokens:** ~28 tokens

**After stripping:**
```json
{
  "bio": "GitHub's mascot ",
  "location": "San Francisco ",
  "status": " Available",
  "user": "octocat "
}
```
**Tokens:** ~18 tokens

**Savings: 36% fewer tokens**

## Integration Examples

### Python

```python
import requests
import json

data = {
    "message": "Hello 👋 World 🌍!",
    "status": "✅ Completed"
}

response = requests.post(
    'http://localhost:8080/slim?profile=llm-context',
    json=data
)

cleaned = response.json()
print(json.dumps(cleaned, indent=2))
```

### JavaScript

```javascript
const data = {
  message: "Hello 👋 World 🌍!",
  status: "✅ Completed"
};

const response = await fetch('http://localhost:8080/slim?profile=llm-context', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(data)
});

const cleaned = await response.json();
console.log(cleaned);
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

func main() {
    data := map[string]interface{}{
        "message": "Hello 👋 World 🌍!",
        "status":  "✅ Completed",
    }

    jsonData, _ := json.Marshal(data)
    
    resp, _ := http.Post(
        "http://localhost:8080/slim?profile=llm-context",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    defer resp.Body.Close()

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    
    fmt.Printf("%+v\n", result)
}
```

## Troubleshooting

### Issue: Too much whitespace after removal

**Problem:**
```json
{"text": "Hello     World  "}
```

**Solution:** Combine with string trimming in post-processing or use additional text normalization.

### Issue: Need to preserve some non-ASCII characters

**Problem:** Accented characters (é, ñ, ü) are also removed.

**Solution:** Currently, the feature removes ALL non-ASCII. If you need to preserve accented Latin characters, you can modify the `stripEmoji` function in `slimjson.go` to include range 128-255:

```go
// Uncomment in stripEmoji function:
else if r >= 128 && r <= 255 {
    result.WriteRune(r)
}
```

## See Also

- [Main README](../README.md)
- [Library Examples](../LIBRARY_EXAMPLES.md)
- [API Documentation](../api/README.md)
- [Configuration File Guide](../.slimjson.example)
