package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define the suite
type HttpTransportTestSuite struct {
	suite.Suite
	serverCmd *exec.Cmd
	serverURL string
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestHttpTransportTestSuite(t *testing.T) {
	suite.Run(t, new(HttpTransportTestSuite))
}

// SetupSuite runs once before the entire test suite
func (suite *HttpTransportTestSuite) SetupSuite() {
	// Build the github-mcp-server-http executable
	buildCmd := exec.Command("go", "build", "-o", "github-mcp-server-http-test", "../../github-mcp-server-http/cmd/github-mcp-server")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	err := buildCmd.Run()
	suite.Require().NoError(err, "Failed to build github-mcp-server-http executable")

	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	suite.Require().NoError(err, "Failed to find available port")
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close() // Close the listener so the server can use the port

	suite.serverURL = fmt.Sprintf("http://127.0.0.1:%d", port)

	// Start the server in a goroutine
	suite.serverCmd = exec.Command("./github-mcp-server-http-test", "http", "--port", fmt.Sprintf("%d", port))
	suite.serverCmd.Stdout = os.Stdout
	suite.serverCmd.Stderr = os.Stderr
	err = suite.serverCmd.Start()
	suite.Require().NoError(err, "Failed to start github-mcp-server-http")

	// Wait for the server to be ready
	suite.T().Logf("Waiting for server to be ready at %s...", suite.serverURL)
	assert.Eventually(suite.T(), func() bool {
		resp, err := http.Get(suite.serverURL)
		if err == nil {
			resp.Body.Close()
			return true // Server is responding
		}
		return false
	}, 10*time.Second, 500*time.Millisecond, "Server did not become ready")
	suite.T().Log("Server is ready.")
}

// TearDownSuite runs once after the entire test suite
func (suite *HttpTransportTestSuite) TearDownSuite() {
	if suite.serverCmd != nil && suite.serverCmd.Process != nil {
		suite.T().Log("Stopping server...")
		err := suite.serverCmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			suite.T().Logf("Failed to send SIGTERM to server: %v", err)
		}
		_, err = suite.serverCmd.Process.Wait()
		if err != nil {
			suite.T().Logf("Error waiting for server to exit: %v", err)
		}
		suite.T().Log("Server stopped.")
	}
	// Clean up the executable
	os.Remove("github-mcp-server-http-test")
}

// Helper to send MCP requests
func (suite *HttpTransportTestSuite) sendMCPRequest(method string, params interface{}) (*mcp.RPCResponse, error) {
	reqBody := mcp.RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  method,
		Params:  params,
	}
	jsonBody, err := json.Marshal(reqBody)
	suite.Require().NoError(err, "Failed to marshal JSON request")

	resp, err := http.Post(suite.serverURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	suite.Require().NoError(err, "Failed to read response body")

	var rpcResp mcp.RPCResponse
	err = json.Unmarshal(bodyBytes, &rpcResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w, body: %s", err, string(bodyBytes))
	}
	return &rpcResp, nil
}

func (suite *HttpTransportTestSuite) TestListTools() {
	resp, err := suite.sendMCPRequest("mcp_listTools", nil)
	suite.Require().NoError(err)
	suite.Assert().NotNil(resp.Result, "Result should not be nil")
	suite.Assert().Nil(resp.Error, "Error should be nil")

	var listToolsResult mcp.ListToolsResult
	err = json.Unmarshal(resp.Result, &listToolsResult)
	suite.Require().NoError(err, "Failed to unmarshal ListToolsResult")

	suite.Assert().NotEmpty(listToolsResult.Tools, "Tools list should not be empty")
	suite.Assert().Contains(listToolsResult.Tools[0].Name, "github", "Expected GitHub tools to be present")
}

func (suite *HttpTransportTestSuite) TestListResources() {
	resp, err := suite.sendMCPRequest("mcp_listResources", nil)
	suite.Require().NoError(err)
	suite.Assert().NotNil(resp.Result, "Result should not be nil")
	suite.Assert().Nil(resp.Error, "Error should be nil")

	var listResourcesResult mcp.ListResourcesResult
	err = json.Unmarshal(resp.Result, &listResourcesResult)
	suite.Require().NoError(err, "Failed to unmarshal ListResourcesResult")

	suite.Assert().NotEmpty(listResourcesResult.Resources, "Resources list should not be empty")
	suite.Assert().Contains(listResourcesResult.Resources[0].URI, "repo://", "Expected repository resources to be present")
}

func (suite *HttpTransportTestSuite) TestCallToolInvalidMethod() {
	resp, err := suite.sendMCPRequest("non_existent_method", nil)
	suite.Require().NoError(err)
	suite.Assert().Nil(resp.Result, "Result should be nil for error")
	suite.Assert().NotNil(resp.Error, "Error should not be nil")
	suite.Assert().Equal(mcp.MethodNotFoundErrorCode, resp.Error.Code, "Expected MethodNotFound error code")
}

func (suite *HttpTransportTestSuite) TestCallToolMissingToken() {
	// This test assumes a GitHub token is required for most tools.
	// The server setup in SetupSuite doesn't provide a real token,
	// so most tool calls should fail due to missing token.
	resp, err := suite.sendMCPRequest("github.list_branches", map[string]interface{}{
		"owner": "octocat",
		"repo":  "Spoon-Knife",
	})
	suite.Require().NoError(err)
	suite.Assert().Nil(resp.Result, "Result should be nil for error")
	suite.Assert().NotNil(resp.Error, "Error should not be nil")
	suite.Assert().Equal(mcp.InternalErrorCode, resp.Error.Code, "Expected Internal error code for missing token")
	suite.Assert().Contains(resp.Error.Message, "GITHUB_PERSONAL_ACCESS_TOKEN not set", "Expected missing token error message")
}

func (suite *HttpTransportTestSuite) TestConcurrentRequests() {
	numRequests := 10
	results := make(chan *mcp.RPCResponse, numRequests)
	errs := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			resp, err := suite.sendMCPRequest("mcp_listTools", nil)
			if err != nil {
				errs <- err
				return
			}
			results <- resp
		}()
	}

	for i := 0; i < numRequests; i++ {
		select {
		case resp := <-results:
			suite.Assert().NotNil(resp.Result, "Result should not be nil for concurrent request")
			suite.Assert().Nil(resp.Error, "Error should be nil for concurrent request")
		case err := <-errs:
			suite.Fail("Error during concurrent request", err.Error())
		case <-time.After(5 * time.Second):
			suite.Fail("Timeout waiting for concurrent request")
		}
	}
}
