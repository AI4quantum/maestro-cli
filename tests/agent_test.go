package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestAgentCreateCommand tests the agent create command
func TestAgentCreateCommand(t *testing.T) {
	// Create a valid YAML file for testing
	validYAML := `---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent
spec:
  framework: openai
  description: "Test agent for unit tests"
  model: gpt-4
  tools:
    - name: test-tool
      description: "A test tool"
`

	tempFile := createTempFile(t, "valid-agent-*.yaml", validYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "agent", "create", tempFile, "--dry-run")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// This test is expected to fail if no MCP server is running
	if err != nil {
		// Check if the error is due to MCP server not being available
		if strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running (expected): %s", outputStr)
			return
		}
		t.Fatalf("Agent create command failed with unexpected error: %v, output: %s", err, string(output))
	}

	if !strings.Contains(outputStr, "Creating agents from YAML configuration") {
		t.Errorf("Should show agent creation message, got: %s", outputStr)
	}
}

// TestAgentCreateWithInvalidYAML tests the agent create command with invalid YAML
func TestAgentCreateWithInvalidYAML(t *testing.T) {
	// Create an invalid YAML file
	invalidYAML := `---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent
spec:
  framework: "openai
  description: "Test agent with invalid YAML"
  model: gpt-4
`

	tempFile := createTempFile(t, "invalid-agent-*.yaml", invalidYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "agent", "create", tempFile)
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// Should fail with invalid YAML
	if err == nil {
		t.Error("Agent create command should fail with invalid YAML")
	}

	if !strings.Contains(outputStr, "no valid YAML documents found") {
		t.Errorf("Error message should mention YAML parsing error, got: %s", outputStr)
	}
}

// TestAgentServeCommand tests the agent serve command
func TestAgentServeCommand(t *testing.T) {
	// Create a valid YAML file for testing
	validYAML := `---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent
spec:
  framework: openai
  description: "Test agent for unit tests"
  model: gpt-4
  tools:
    - name: test-tool
      description: "A test tool"
`

	tempFile := createTempFile(t, "valid-agent-*.yaml", validYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "agent", "serve", tempFile, "--dry-run")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// This test is expected to fail if no MCP server is running
	if err != nil {
		// Check if the error is due to MCP server not being available
		if strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running (expected): %s", outputStr)
			return
		}
		t.Fatalf("Agent serve command failed with unexpected error: %v, output: %s", err, string(output))
	}

	if !strings.Contains(outputStr, "Agent server started successfully") {
		t.Errorf("Should show agent serving message, got: %s", outputStr)
	}
}

// TestAgentServeWithCustomPort tests the agent serve command with custom port
func TestAgentServeWithCustomPort(t *testing.T) {
	// Create a valid YAML file for testing
	validYAML := `---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent
spec:
  framework: openai
  description: "Test agent for unit tests"
  model: gpt-4
  tools:
    - name: test-tool
      description: "A test tool"
`

	tempFile := createTempFile(t, "valid-agent-*.yaml", validYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "agent", "serve", tempFile, "--port=8080", "--dry-run")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// This test is expected to fail if no MCP server is running
	if err != nil {
		// Check if the error is due to MCP server not being available
		if strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running (expected): %s", outputStr)
			return
		}
		t.Fatalf("Agent serve command failed with unexpected error: %v, output: %s", err, string(output))
	}

	if !strings.Contains(outputStr, "Agent server started successfully") {
		t.Errorf("Should show agent serving message, got: %s", outputStr)
	}
}

// TestAgentHelpCommand tests the agent help command
func TestAgentHelpCommand(t *testing.T) {
	cmd := exec.Command("../maestro", "agent", "--help")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Failed to run agent help command: %v", err)
	}

	helpOutput := string(output)

	// Check for expected help content
	expectedContent := []string{
		"agent",
		"create",
		"serve",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(helpOutput, expected) {
			t.Errorf("Help output should contain '%s'", expected)
		}
	}
}

// Made with Bob
