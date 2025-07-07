# Phase 1: Local Development and Testing

This phase focuses on ensuring the new HTTP server functionality is working correctly in a local development environment before containerization.

## 1. Unit Testing

### Goal
Verify the correctness of individual functions and components in isolation.

### Steps

1.  **`cmd/github-mcp-server/main.go` - HTTP Command Tests:**
    *   **Test Case 1: Command Registration:** Verify that the `http` command is correctly added as a subcommand to the `rootCmd`.
    *   **Test Case 2: Flag Parsing:** Ensure that the `--addr` and `--port` flags are correctly parsed when provided on the command line.
    *   **Test Case 3: RunE Function Call:** Confirm that the `RunE` function of the `http` command is invoked with the expected configuration (e.g., correct token, toolsets, and HTTP server parameters).

2.  **`internal/ghmcp/server.go` - RunHttpServer Tests:**
    *   **Test Case 1: Server Initialization:** Verify that `RunHttpServer` successfully initializes an `MCPServer` instance with the provided configuration.
    *   **Test Case 2: HTTP Server Start:** Ensure that the `http.ListenAndServe` function is called with the correct address and port, indicating the HTTP server attempts to start. (Note: This might require mocking `http.ListenAndServe` for true unit testing, or using a short timeout).
    *   **Test Case 3: Graceful Shutdown:** Test that the server responds to OS signals (e.g., `SIGTERM`, `SIGINT`) by attempting a graceful shutdown.

### Expected Results

*	All unit tests pass without errors.
*	Test coverage for the new `http` command and `RunHttpServer` function is high.
