package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestToolCreateCommand tests the tool create command
func TestToolCreateCommand(t *testing.T) {
	// Create a valid YAML file for testing
	validYAML := `---
apiVersion: maestro/v1alpha1
kind: MCPTool
metadata:
  name: test-tool
spec:
  description: "Test tool for unit tests"
  parameters:
    - name: param1
      description: "A test parameter"
      required: true
      type: string
  returns:
    description: "Test return value"
    type: string
`

	tempFile := createTempFile(t, "valid-tool-*.yaml", validYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "tool", "create", tempFile, "--dry-run")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// This test is expected to fail if no MCP server is running
	if err != nil {
		// Check if the error is due to MCP server not being available
		if strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running (expected): %s", outputStr)
			return
		}
		t.Fatalf("Tool create command failed with unexpected error: %v, output: %s", err, string(output))
	}

	if !strings.Contains(outputStr, "Creating MCP tools from YAML configuration") {
		t.Errorf("Should show MCP tools creation message, got: %s", outputStr)
	}
}

// TestToolCreateWithNonExistentFile tests with non-existent file
func TestToolCreateWithNonExistentFile(t *testing.T) {
	cmd := exec.Command("../maestro", "tool", "create", "nonexistent.yaml")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// Should fail with non-existent file
	if err == nil {
		t.Error("Tool create command should fail with non-existent file")
	}

	if !strings.Contains(outputStr, "no such file or directory") {
		t.Errorf("Error message should mention file not found, got: %s", outputStr)
	}
}

// TestToolCreateWithInvalidYAML tests with invalid YAML
func TestToolCreateWithInvalidYAML(t *testing.T) {
	// Create an invalid YAML file
	invalidYAML := `---
apiVersion: maestro/v1alpha1
kind: MCPTool
metadata:
  name: test-tool
spec:
  description: "Test tool with invalid YAML
  parameters:
    - name: param1
      description: "A test parameter"
      required: true
      type: string
`

	tempFile := createTempFile(t, "invalid-tool-*.yaml", invalidYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "tool", "create", tempFile)
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// Should fail with invalid YAML
	if err == nil {
		t.Error("Tool create command should fail with invalid YAML")
	}

	if !strings.Contains(outputStr, "no valid YAML documents found") {
		t.Errorf("Error message should mention YAML parsing error, got: %s", outputStr)
	}
}

// TestToolHelpCommand tests the tool help command
func TestToolHelpCommand(t *testing.T) {
	cmd := exec.Command("../maestro", "tool", "--help")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Failed to run tool help command: %v", err)
	}

	helpOutput := string(output)

	// Check for expected help content
	expectedContent := []string{
		"tool",
		"create",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(helpOutput, expected) {
			t.Errorf("Help output should contain '%s'", expected)
		}
	}
}

// TestToolCreateWithMultipleTools tests creating multiple tools from a single file
func TestToolCreateWithMultipleTools(t *testing.T) {
	// Create a valid YAML file with multiple tools
	validYAML := `---
apiVersion: maestro/v1alpha1
kind: MCPTool
metadata:
  name: test-tool-1
spec:
  description: "Test tool 1"
  parameters:
    - name: param1
      description: "A test parameter"
      required: true
      type: string
  returns:
    description: "Test return value"
    type: string
---
apiVersion: maestro/v1alpha1
kind: MCPTool
metadata:
  name: test-tool-2
spec:
  description: "Test tool 2"
  parameters:
    - name: param1
      description: "A test parameter"
      required: true
      type: string
  returns:
    description: "Test return value"
    type: string
`

	tempFile := createTempFile(t, "valid-tools-*.yaml", validYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "tool", "create", tempFile, "--dry-run")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// This test is expected to fail if no MCP server is running
	if err != nil {
		// Check if the error is due to MCP server not being available
		if strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running (expected): %s", outputStr)
			return
		}
		t.Fatalf("Tool create command failed with unexpected error: %v, output: %s", err, string(output))
	}

	if !strings.Contains(outputStr, "Creating MCP tools from YAML configuration") {
		t.Errorf("Should show MCP tools creation message, got: %s", outputStr)
	}
}

// Made with Bob
