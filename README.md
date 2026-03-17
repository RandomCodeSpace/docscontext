# DocsGraphContext

A pure Go ([CGO_ENABLED=0](https://pkg.go.dev/cmd/cgo)) GraphRAG tool inspired by [Microsoft GraphRAG](https://github.com/microsoft/graphrag).
Ingests unstructured documents, builds a hierarchical knowledge graph with community detection, and exposes an **MCP server + embedded Web UI** on a single port.

## Features

- **GraphRAG pipeline** — 5-phase: load → chunk → embed → graph extraction → community detection
- **Knowledge graph** — entity/relationship/claim extraction via LLM (JSON mode)
- **Louvain community detection** — pure Go, hierarchical, no external dependencies
- **Three LLM providers** — Azure OpenAI, Ollama (local), HuggingFace TGI
- **12 MCP tools** — local search, global search, graph walk, community reports, and more
- **Embedded Web UI** — vis-network graph explorer, semantic search, document browser
- **Single binary** — zero CGO, cross-compiles to Linux / macOS / Windows

## Install

```bash
go install github.com/RandomCodeSpace/docsgraphcontext@latest
```

Or build from source:

```bash
git clone https://github.com/RandomCodeSpace/docsgraphcontext.git
cd docsgraphcontext
CGO_ENABLED=0 go build -o docsgraph .
```

## Quick Start

```bash
# 1. Create config
mkdir -p ~/.docsgraph
cp config.example.yaml ~/.docsgraph/config.yaml
# Edit ~/.docsgraph/config.yaml — set your LLM provider

# 2. Index documents (Phases 1-2: load, chunk, embed, extract entities)
docsgraph index ./your-docs/ --workers 4

# 3. Build knowledge graph (Phases 3-4: community detection + summaries)
docsgraph index --finalize

# 4. Check stats
docsgraph stats

# 5. Start server
docsgraph serve --port 8080
```

Open **http://localhost:8080** for the Web UI.

## Configuration

Copy `config.example.yaml` to `~/.docsgraph/config.yaml` and edit:

```yaml
data_dir: ~/.docsgraph/data

llm:
  provider: ollama          # azure | ollama | huggingface

  ollama:
    base_url: http://localhost:11434
    chat_model: llama3.2
    embed_model: nomic-embed-text

  azure:
    endpoint: https://myresource.openai.azure.com
    api_key: ${AZURE_OPENAI_API_KEY}
    api_version: "2024-02-01"
    chat_model: gpt-4o
    embed_model: text-embedding-3-small

indexing:
  chunk_size: 512
  chunk_overlap: 50
  workers: 4

server:
  host: 127.0.0.1
  port: 8080
```

Environment variable overrides use the `DOCSGRAPH_` prefix:

```bash
DOCSGRAPH_LLM_PROVIDER=azure
DOCSGRAPH_LLM_AZURE_API_KEY=sk-...
DOCSGRAPH_SERVER_PORT=9090
```

## CLI

```bash
# Index a file or directory
docsgraph index ./docs/ [--force] [--workers 4] [--verbose]

# Run community detection + LLM summaries
docsgraph index --finalize

# Show statistics
docsgraph stats
docsgraph stats --json

# Start MCP + Web UI server
docsgraph serve [--port 8080] [--host 127.0.0.1]
```

## MCP Tools

Connect any MCP client to `http://localhost:8080/mcp/sse`.

| Tool | Description |
|---|---|
| `search_documents` | Vector similarity search over chunks |
| `local_search` | Vector + graph walk (GraphRAG local) |
| `global_search` | Community summary aggregation with LLM synthesis |
| `query_entity` | Entity details + relationships by name |
| `find_relationships` | Relationship lookup by source / target / predicate |
| `get_graph_neighborhood` | Subgraph JSON for visualization |
| `get_document_structure` | LLM-generated structured summary |
| `list_entities` | Browse entities with type filter |
| `list_documents` | Browse indexed documents |
| `get_community_report` | Community summary + member entities |
| `get_chunk` | Retrieve chunk by ID |
| `stats` | Full index statistics |

## REST API

```
GET  /api/stats
GET  /api/documents
GET  /api/documents/{id}
POST /api/search          {"query":"...","mode":"local|global","top_k":5}
GET  /api/graph/neighborhood?entity=<name>&depth=2
GET  /api/entities
GET  /api/communities
GET  /api/communities/{id}
POST /api/upload
```

## Architecture

```
Document In
    │
    ▼ Phase 1 — Text Units
  Loader (PDF/DOCX/TXT/MD) → Chunker → Embedder → SQLite

    ▼ Phase 2 — Graph Extraction  [parallel per document]
  LLM → Entities + Relationships + Claims → SQLite

    ▼ Phase 3 — Community Detection  [post-index finalization]
  Louvain algorithm → hierarchical community assignments

    ▼ Phase 4 — Community Summaries  [parallel]
  LLM → CommunityReport → embed summary → SQLite

    ▼ Phase 5 — Structured Doc
  LLM → JSON summary → SQLite
```

All data lives in a single SQLite file at `$DATA_DIR/docsgraph.db`.

## Supported File Types

| Type | Extension |
|---|---|
| PDF | `.pdf` |
| Word | `.docx` |
| Markdown | `.md`, `.markdown` |
| Plain text | `.txt`, `.text` |

## License

MIT
