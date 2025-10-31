package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMermaidCommandWithValidWorkflow tests the mermaid command with a valid workflow file
func TestMermaidCommandWithValidWorkflow(t *testing.T) {
	// Create a valid workflow YAML file for testing
	validWorkflowYAML := `---
apiVersion: maestro/v1alpha1
kind: Workflow
metadata:
  name: test-workflow
spec:
  template:
    prompt: "Test prompt"
    steps:
      - name: test-step
        agent: test-agent
        input: "{{ .prompt }}"
`

	tempFile := createTempFile(t, "valid-workflow-*.yaml", validWorkflowYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "mermaid", tempFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Mermaid command failed with error: %v, output: %s", err, string(output))
	}

	outputStr := string(output)

	// Check for expected output
	if !strings.Contains(outputStr, "sequenceDiagram") {
		t.Errorf("Expected sequenceDiagram in output, got: %s", outputStr)
	}
}

// TestMermaidCommandWithSequenceDiagram tests the mermaid command with sequenceDiagram flag
func TestMermaidCommandWithSequenceDiagram(t *testing.T) {
	// Create a valid workflow YAML file for testing
	validWorkflowYAML := `---
apiVersion: maestro/v1alpha1
kind: Workflow
metadata:
  name: test-workflow
spec:
  template:
    prompt: "Test prompt"
    steps:
      - name: test-step
        agent: test-agent
        input: "{{ .prompt }}"
`

	tempFile := createTempFile(t, "valid-workflow-*.yaml", validWorkflowYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "mermaid", tempFile, "--sequenceDiagram")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Mermaid command with sequenceDiagram flag failed with error: %v, output: %s", err, string(output))
	}

	outputStr := string(output)

	// Check for expected output
	if !strings.Contains(outputStr, "sequenceDiagram") {
		t.Errorf("Expected sequenceDiagram in output, got: %s", outputStr)
	}
}

// TestMermaidCommandWithFlowchartTD tests the mermaid command with flowchart-td flag
func TestMermaidCommandWithFlowchartTD(t *testing.T) {
	// Create a valid workflow YAML file for testing
	validWorkflowYAML := `---
apiVersion: maestro/v1alpha1
kind: Workflow
metadata:
  name: test-workflow
spec:
  template:
    prompt: "Test prompt"
    steps:
      - name: test-step
        agent: test-agent
        input: "{{ .prompt }}"
`

	tempFile := createTempFile(t, "valid-workflow-*.yaml", validWorkflowYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "mermaid", tempFile, "--flowchart-td")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Mermaid command with flowchart-td flag failed with error: %v, output: %s", err, string(output))
	}

	outputStr := string(output)

	// Check for expected output
	if !strings.Contains(outputStr, "flowchart TD") {
		t.Errorf("Expected flowchart TD in output, got: %s", outputStr)
	}
}

// TestMermaidCommandWithFlowchartLR tests the mermaid command with flowchart-lr flag
func TestMermaidCommandWithFlowchartLR(t *testing.T) {
	// Create a valid workflow YAML file for testing
	validWorkflowYAML := `---
apiVersion: maestro/v1alpha1
kind: Workflow
metadata:
  name: test-workflow
spec:
  template:
    prompt: "Test prompt"
    steps:
      - name: test-step
        agent: test-agent
        input: "{{ .prompt }}"
`

	tempFile := createTempFile(t, "valid-workflow-*.yaml", validWorkflowYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "mermaid", tempFile, "--flowchart-lr")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Mermaid command with flowchart-lr flag failed with error: %v, output: %s", err, string(output))
	}

	outputStr := string(output)

	// Check for expected output
	if !strings.Contains(outputStr, "flowchart LR") {
		t.Errorf("Expected flowchart LR in output, got: %s", outputStr)
	}
}

// TestMermaidCommandWithNonExistentFile tests the mermaid command with a non-existent file
func TestMermaidCommandWithNonExistentFile(t *testing.T) {
	cmd := exec.Command("../maestro", "mermaid", "nonexistent.yaml")
	output, err := cmd.CombinedOutput()

	// Should fail with non-existent file
	if err == nil {
		t.Error("Mermaid command should fail with non-existent file")
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "no such file or directory") {
		t.Errorf("Error message should mention file not found, got: %s", outputStr)
	}
}

// TestMermaidCommandWithInvalidYAML tests the mermaid command with invalid YAML
func TestMermaidCommandWithInvalidYAML(t *testing.T) {
	// Create an invalid YAML file
	invalidYAML := `---
apiVersion: maestro/v1alpha1
kind: Workflow
metadata:
  name: test-workflow
spec:
  template:
    prompt: "Test prompt
    steps:
      - name: test-step
        agent: test-agent
        input: "{{ .prompt }}"
`

	tempFile := createTempFile(t, "invalid-workflow-*.yaml", invalidYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "mermaid", tempFile)
	output, err := cmd.CombinedOutput()

	// Should fail with invalid YAML
	if err == nil {
		t.Error("Mermaid command should fail with invalid YAML")
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "Unable to parse workflow file") {
		t.Errorf("Error message should mention parsing error, got: %s", outputStr)
	}
}

// TestMermaidHelp tests the mermaid help command
func TestMermaidHelp(t *testing.T) {
	cmd := exec.Command("../maestro", "mermaid", "--help")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Failed to run mermaid help command: %v", err)
	}

	helpOutput := string(output)

	// Check for expected help content
	expectedContent := []string{
		"mermaid",
		"Generate mermaid diagrams",
		"--sequenceDiagram",
		"--flowchart-td",
		"--flowchart-lr",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(helpOutput, expected) {
			t.Errorf("Help output should contain '%s'", expected)
		}
	}
}

// Using the common createTempFile function from test_utils.go

// Made with Bob
