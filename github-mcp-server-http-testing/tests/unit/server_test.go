package unit_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/github/github-mcp-server/internal/ghmcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for http.ListenAndServe
type MockHttp struct {
	mock.Mock
}

func (m *MockHttp) ListenAndServe(addr string, handler http.Handler) error {
	args := m.Called(addr, handler)
	return args.Error(0)
}

// Global variable to hold our mock for ListenAndServe
var mockListenAndServe func(addr string, handler http.Handler) error

// Overwrite http.ListenAndServe for testing
func init() {
	ghmcp.ListenAndServe = func(addr string, handler http.Handler) error {
		return mockListenAndServe(addr, handler)
	}
}

func TestRunHttpServerInitialization(t *testing.T) {
	// Mock the ListenAndServe function to just return nil, indicating success
	// without actually starting a server.
	mockListenAndServe = func(addr string, handler http.Handler) error {
		return nil
	}

	cfg := ghmcp.HttpServerConfig{
		Version:         "test-version",
		Host:            "github.com",
		Token:           "mock_token",
		EnabledToolsets: []string{"all"},
		Addr:            "127.0.0.1",
		Port:            8080,
	}

	// Run in a goroutine to not block the test, as it contains a select statement
	errC := make(chan error)
	go func() {
		errC <- ghmcp.RunHttpServer(cfg)
	}()

	// Give the server a moment to start and then send a signal to shut it down
	time.Sleep(100 * time.Millisecond)
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	select {
	case err := <-errC:
		assert.NoError(t, err, "RunHttpServer should not return an error on graceful shutdown")
	case <-time.After(5 * time.Second):
		t.Fatal("RunHttpServer did not exit in time")
	}
}

func TestRunHttpServerStartsHttpServer(t *testing.T) {
	// Use a channel to signal when ListenAndServe is called
	listenAndServeCalled := make(chan struct{})
	mockListenAndServe = func(addr string, handler http.Handler) error {
		close(listenAndServeCalled) // Signal that it was called
		// Keep the mock blocking until the test sends a signal to stop,
		// or return an error to simulate startup failure.
		<-time.After(1 * time.Second) // Simulate server running for a bit
		return http.ErrServerClosed   // Simulate graceful shutdown
	}

	cfg := ghmcp.HttpServerConfig{
		Version:         "test-version",
		Host:            "github.com",
		Token:           "mock_token",
		EnabledToolsets: []string{"all"},
		Addr:            "127.0.0.1",
		Port:            8081, // Use a different port
	}

	errC := make(chan error)
	go func() {
		errC <- ghmcp.RunHttpServer(cfg)
	}()

	select {
	case <-listenAndServeCalled:
		// ListenAndServe was called, now simulate graceful shutdown
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	case err := <-errC:
		t.Fatalf("RunHttpServer returned an error unexpectedly: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("ListenAndServe was not called in time")
	}

	// Wait for the server goroutine to exit after signal
	select {
	case err := <-errC:
		assert.NoError(t, err, "RunHttpServer should exit gracefully")
	case <-time.After(5 * time.Second):
		t.Fatal("RunHttpServer did not exit after signal")
	}
}

func TestRunHttpServerHandlesListenAndServeError(t *testing.T) {
	expectedErr := errors.New("mock listen error")
	mockListenAndServe = func(addr string, handler http.Handler) error {
		return expectedErr
	}

	cfg := ghmcp.HttpServerConfig{
		Version:         "test-version",
		Host:            "github.com",
		Token:           "mock_token",
		EnabledToolsets: []string{"all"},
		Addr:            "127.0.0.1",
		Port:            8082, // Use a different port
	}

	err := ghmcp.RunHttpServer(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
}

func TestRunHttpServerGracefulShutdown(t *testing.T) {
	// Use a real listener to ensure the server starts, then close it
	// to simulate external shutdown.
	listener, err := net.Listen("tcp", "127.0.0.1:0") // Listen on a random available port
	assert.NoError(t, err)
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	addr := listener.Addr().(*net.TCPAddr).IP.String()

	// Mock ListenAndServe to use our test listener
	mockListenAndServe = func(a string, handler http.Handler) error {
		server := &http.Server{Handler: handler}
		return server.Serve(listener) // Use the test listener
	}

	cfg := ghmcp.HttpServerConfig{
		Version:         "test-version",
		Host:            "github.com",
		Token:           "mock_token",
		EnabledToolsets: []string{"all"},
		Addr:            addr,
		Port:            port,
	}

	errC := make(chan error)
	go func() {
		errC <- ghmcp.RunHttpServer(cfg)
	}()

	// Give the server a moment to start
	time.Sleep(200 * time.Millisecond)

	// Send a SIGTERM signal to the process to trigger graceful shutdown
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	select {
	case err := <-errC:
		assert.NoError(t, err, "RunHttpServer should exit without error on graceful shutdown")
	case <-time.After(5 * time.Second):
		t.Fatal("RunHttpServer did not exit in time after SIGTERM")
	}
}
