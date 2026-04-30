# envdiff

> Compare `.env` files across environments and flag missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff && go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <base-file> <compare-file> [compare-file...]
```

### Example

```bash
envdiff .env.example .env.production .env.staging
```

**Output:**

```
[MISSING]  .env.production  → DATABASE_URL
[MISSING]  .env.staging     → REDIS_URL, SECRET_KEY
[MISMATCH] .env.staging     → APP_ENV (expected: "production", got: "staging")

2 file(s) checked. 3 issue(s) found.
```

### Flags

| Flag | Description |
|------|-------------|
| `--strict` | Exit with non-zero code if any issues are found |
| `--ignore-values` | Check keys only, skip value comparison |
| `--json` | Output results as JSON |

---

## Why envdiff?

Keeping `.env` files in sync across environments is error-prone. `envdiff` makes it easy to catch missing variables before they cause runtime failures — useful in CI pipelines and deployment workflows.

---

## License

MIT © [yourusername](https://github.com/yourusername)