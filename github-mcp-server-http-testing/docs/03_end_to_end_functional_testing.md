# Phase 3: End-to-End Functional Testing

This phase focuses on testing the full functionality of the server as an independent service.

## 1. Manual Testing with `curl`

### Goal
Manually send valid MCP requests to test specific tool calls.

### Steps
1.  Ensure the `github-mcp-http-server` Docker container is running (as per Phase 2).
2.  Construct a valid JSON-RPC 2.0 request for an MCP tool call. For example, to list tools:
    ```json
    {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "mcp_listTools",
        "params": {}
    }
    ```
3.  Use `curl` to POST this request to the server's HTTP endpoint:
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{ "jsonrpc": "2.0", "id": 1, "method": "mcp_listTools", "params": {} }' http://localhost:8080/
    ```
4.  Repeat for other basic MCP methods like `mcp_listResources`.
5.  Construct requests for specific GitHub tools (e.g., `github.list_branches`, `github.get_issue`) with appropriate parameters, ensuring you have a valid GitHub Personal Access Token configured for the server.

### Expected Results
*	The `curl` command returns a valid JSON-RPC 2.0 response with the `result` field containing the expected data for the tool call.
*	Error responses (`error` field) are returned for invalid requests or missing parameters.

## 2. Automated Testing with an MCP Client

### Goal
Systematically test all available tools and resources.

### Steps
1.  **Prerequisites:**
    *	Ensure the `github-mcp-http-server` Docker container is running.
    *	You will need an MCP client library or a custom script capable of making MCP calls over HTTP. For Go, you can use `github.com/mark3labs/mcp-go/client`.
2.  **Test Setup:**
    *	Initialize an HTTP MCP client pointing to `http://localhost:8080/`.
3.  **Test Cases:**
    *	**Tool Discovery:** Call `mcp_listTools` and verify that all expected GitHub tools are listed with correct descriptions and input schemas.
    *	**Resource Discovery:** Call `mcp_listResources` and verify that all expected GitHub resources are listed.
    *	**Tool Call Validation (Positive):** For a selection of key tools (e.g., `github.list_branches`, `github.get_issue`), call them with valid parameters and assert that the returned data is correct.
    *	**Tool Call Validation (Negative):** For selected tools, call them with missing or invalid parameters and assert that appropriate MCP errors are returned.
    *	**Resource Access Validation (Positive):** For a selection of key resources (e.g., `repo://owner/repo/contents/path`), read them with valid URIs and assert that the content is correct.
    *	**Resource Access Validation (Negative):** Attempt to read non-existent resources or resources with invalid URIs and assert that appropriate MCP errors are returned.

### Expected Results
*	All automated tests pass, indicating that the server correctly implements the MCP specification for HTTP transport and exposes the GitHub tools as expected.
*	Comprehensive logging during automated tests helps identify any issues.
