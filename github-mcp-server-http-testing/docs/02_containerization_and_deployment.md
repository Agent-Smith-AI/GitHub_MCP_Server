# Phase 2: Containerization and Deployment Testing

This phase ensures that the server runs correctly within a Docker container.

## 1. Build Docker Image

### Goal
Successfully build a Docker image from the updated `Dockerfile`.

### Steps
1.  Navigate to the `github-mcp-server-http/` directory.
2.  Execute the Docker build command:
    ```bash
    docker build -t github-mcp-server-http:latest .
    ```

### Expected Results
*	The Docker image builds successfully without errors.
*	A new image named `github-mcp-server-http` with the tag `latest` is present in your local Docker images.

## 2. Run Docker Container

### Goal
Start the container and ensure the HTTP server is accessible.

### Steps
1.  Execute the Docker run command, mapping the container's exposed port (8080) to a port on the host machine (e.g., 8080):
    ```bash
    docker run -p 8080:8080 --name github-mcp-http-server github-mcp-server-http:latest http
    ```
2.  Monitor the container logs for any errors or indications of successful server startup.

### Expected Results
*	The Docker container starts successfully.
*	The container logs show "GitHub MCP Server running on http://0.0.0.0:8080" (or similar, depending on the `Addr` setting).
*	No critical errors are reported in the container logs.

## 3. Container Sanity Checks

### Goal
Perform basic checks to ensure the containerized server is responsive.

### Steps
1.  From your host machine's terminal, use `curl` to send a simple GET request to the mapped port:
    ```bash
    curl http://localhost:8080/
    ```

### Expected Results
*	The `curl` command returns an HTTP response (e.g., a 404 Not Found or a message indicating a malformed request, as it's an MCP endpoint, not a web page). This confirms network connectivity to the container and that the server is listening.
*	The curl command should not hang or return connection refused errors.
