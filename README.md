# 🔗 404 Hunter

**404 Hunter** is a versatile tool for detecting dead or broken links in websites or projects. It can be used as a **CLI tool** for developers or as a **web application** for interactive usage. Ensure your links are always alive and your users never hit a dead end!

---

## Features

- ✅ Scan single URLs or multiple links at once
- ✅ Detect **broken links** (HTTP 4xx/5xx errors)
- ✅ CLI support for automation and scripting
- ✅ Web interface for easy scanning
- ✅ Export results in **JSON** or **CSV**
- ✅ Works with large websites with **multi-threaded scanning**

---

## Installation 🛠️

### Using Go

```bash
# Clone the repo
git clone https://github.com/homocodian/404-hunter.git
cd 404-hunter

# Build the CLI tool
go build -o 404-hunter main.go
```

---

## Usage 💻

### CLI Mode

```bash
# Scan a single URL
./404-hunter https://example.com

# Export results
./404-hunter https://example.com -o results.json
```

**Options:**

| Flag | Description                    |
| ---- | ------------------------------ |
| `-o` | Output file                    |
| `-w` | Number of workers (default: 5) |

---

### Web Mode 🌐

1. Start the web server:

```bash
./404-hunter web
```

2. Open your browser at `http://localhost:5000`
3. Enter the URLs you want to scan
4. Download results in **JSON**

---

## Example Output (Broken links)

```json
["https://example.com/page1", "https://example.com/broken"]
```
