# Project: GitHub MCP Server - HTTP Version

This document outlines the tasks and instructions for creating an HTTP version of the GitHub MCP server.

## Task Overview

The primary goal of this project is to create a new version of the `github-mcp-server` that communicates over HTTP instead of stdio. This will allow the server to be run as a remote service, accessible over a network.

## TODO

- [ ] **Phase 1: Analysis and Dependency Verification**
    - [X] Analyze the codebase for existing HTTP connectors.
    - [X] Identify all dependencies in `go.mod`.
    - [ ] Use Context7 to verify dependency compatibility.

### Dependency Analysis Log

- `github.com/google/go-github/v72`: Resolved to `/google/go-github`.
- `github.com/josephburnett/jd`: No direct match found.
- `github.com/mark3labs/mcp-go`: Resolved to `/mark3labs/mcp-go`.
- `github.com/migueleliasweb/go-github-mock`: No direct match found.
- `github.com/sirupsen/logrus`: Resolved to `/sirupsen/logrus`.
- `github.com/spf13/cobra`: Resolved to `/spf13/cobra`.
- `github.com/spf13/viper`: Resolved to `/spf13/viper`.
- `github.com/stretchr/testify`: Resolved to `/stretchr/testify`.

### Dependency Documentation

- **`/google/go-github`**: The documentation for this library is extensive and shows how to create clients, authenticate, and interact with the GitHub API. The snippets on handling HTTP requests and webhooks are particularly relevant to our task. It is compatible with our goal of creating a streamable HTTP server.
- **`/mark3labs/mcp-go`**: This library is the core of our MCP server. The documentation provides clear examples of how to create both `stdio` and `streamable-http` servers. The snippets on creating a `StreamableHTTPClientPool` and adding middleware will be very useful. It is fully compatible with our task.
- **`/sirupsen/logrus`**: This is a structured logging library for Go. The documentation shows how to create loggers, set formatters, and use hooks. This will be useful for logging in our HTTP server. It is compatible with our task.
- **`/spf13/cobra`**: This library is used for creating CLI applications in Go. The documentation shows how to create commands, subcommands, and flags. We will use this to create the new `http` command. It is compatible with our task.
- **`/spf13/viper`**: This library is used for configuration management. The documentation shows how to read configuration files, environment variables, and flags. We will use this to configure the HTTP server. It is compatible with our task.
- **`/stretchr/testify`**: This library is a testing toolkit for Go. The documentation shows how to use assertions and mocks. This will be useful for writing tests for our HTTP server. It is compatible with our task.

- [ ] **Phase 2: Implementation**
    - [ ] Create a new `http` command in `cmd/github-mcp-server/main.go`.
    - [ ] Implement `RunHttpServer` in `internal/ghmcp/server.go`.
    - [ ] Configure a unique server name and port.

- [ ] **Phase 3: Docker and Deployment**
    - [ ] Update `Dockerfile` to support both `http` and `stdio` commands.
    - [ ] Expose the HTTP port in the `Dockerfile`.

- [ ] **Phase 4: Git Workflow**
    - [ ] Create and work on a `feature/http-server` branch.

## Instructions for AI

1.  **Work on the `github-mcp-server-http` directory:** All modifications should be made in this directory, not the original `github-mcp-server` directory.
2.  **Use a feature branch:** All Git commits should be made to a branch named `feature/http-server`.
3.  **Dependency Analysis:** Before making any code changes, use the Context7 MCP server to analyze the dependencies in `go.mod` to ensure compatibility.
4.  **Create a secondary server:** The new HTTP server should be a secondary option, meaning the existing `stdio` server should still be functional.
5.  **Minimize API calls:** Be mindful of API call limits when interacting with external services. Use bulk operations and caching where possible.