package main

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"maestro/internal/commands" // Import the CLI package
)

// Embed test fixtures
// go:embed fixtures/*
var fixturesFS embed.FS

// Embed schemas
// go:embed schemas/*
var schemasFS embed.FS

// TestFixtures handles test fixtures
type TestFixtures struct {
	t        *testing.T
	tempDir  string
	fixtures map[string]string
	schemas  map[string]string
}

// NewTestFixtures creates a new TestFixtures instance
func NewTestFixtures(t *testing.T) *TestFixtures {
	// Create a temporary directory for test fixtures
	tempDir, err := os.MkdirTemp("", "maestro-test-fixtures-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory for fixtures: %v", err)
	}

	// Register cleanup to remove the directory after the test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	fixtures := &TestFixtures{
		t:        t,
		tempDir:  tempDir,
		fixtures: make(map[string]string),
		schemas:  make(map[string]string),
	}

	// Extract fixtures to temporary directory
	fixtures.extractFixtures()
	fixtures.extractSchemas()

	return fixtures
}

// extractFixtures extracts fixtures from the embedded filesystem to the temporary directory
func (f *TestFixtures) extractFixtures() {
	// Walk through the embedded filesystem and extract files
	err := fs.WalkDir(fixturesFS, "fixtures", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Read file content
		content, err := fixturesFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Create relative path
		relPath := path[len("fixtures/"):]

		// Create directory structure
		dir := filepath.Dir(filepath.Join(f.tempDir, relPath))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Write file
		filePath := filepath.Join(f.tempDir, relPath)
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return err
		}

		// Store file path
		f.fixtures[relPath] = filePath

		return nil
	})

	if err != nil {
		f.t.Fatalf("Failed to extract fixtures: %v", err)
	}
}

// extractSchemas extracts schemas from the embedded filesystem to the temporary directory
func (f *TestFixtures) extractSchemas() {
	// Walk through the embedded filesystem and extract files
	err := fs.WalkDir(schemasFS, "schemas", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Read file content
		content, err := schemasFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Create relative path
		relPath := path[len("schemas/"):]

		// Create directory structure
		dir := filepath.Dir(filepath.Join(f.tempDir, "schemas", relPath))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Write file
		filePath := filepath.Join(f.tempDir, "schemas", relPath)
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return err
		}

		// Store file path
		f.schemas[relPath] = filePath

		return nil
	})

	if err != nil {
		f.t.Fatalf("Failed to extract schemas: %v", err)
	}
}

// GetFixture returns the path to a test fixture file
func (f *TestFixtures) GetFixture(fileName string) string {
	if path, ok := f.fixtures[fileName]; ok {
		return path
	}
	f.t.Fatalf("Fixture not found: %s", fileName)
	return ""
}

// GetSchema returns the path to a schema file
func (f *TestFixtures) GetSchema(fileName string) string {
	if path, ok := f.schemas[fileName]; ok {
		return path
	}
	f.t.Fatalf("Schema not found: %s", fileName)
	return ""
}

// CreateTempFile creates a temporary file with the given content and returns its path
func (f *TestFixtures) CreateTempFile(content string, suffix string) string {
	tempFile, err := os.CreateTemp(f.tempDir, "temp-*"+suffix)
	if err != nil {
		f.t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tempFile.WriteString(content); err != nil {
		tempFile.Close()
		f.t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		f.t.Fatalf("Failed to close temp file: %v", err)
	}

	return tempFile.Name()
}

// TestHelper provides utility methods for tests
type TestHelper struct {
	t        *testing.T
	fixtures *TestFixtures
}

// NewTestHelper creates a new TestHelper instance
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		t:        t,
		fixtures: NewTestFixtures(t),
	}
}

// GetFixture returns the path to a test fixture file
func (h *TestHelper) GetFixture(fileName string) string {
	return h.fixtures.GetFixture(fileName)
}

// GetSchema returns the path to a schema file
func (h *TestHelper) GetSchema(fileName string) string {
	return h.fixtures.GetSchema(fileName)
}

// CreateTempFile creates a temporary file with the given content and returns its path
func (h *TestHelper) CreateTempFile(content string, suffix string) string {
	return h.fixtures.CreateTempFile(content, suffix)
}

// SkipIfEnvNotSet skips the test if the specified environment variable is not set to the expected value
func (h *TestHelper) SkipIfEnvNotSet(envVar, expectedValue string) {
	if os.Getenv(envVar) != expectedValue {
		h.t.Skipf("Skipping test: %s is not set to %s", envVar, expectedValue)
	}
}

// CommandArgs represents the command-line arguments for a CLI command
type CommandArgs struct {
	// Common flags
	*commands.CommandOptions
	//DryRun   bool
	//Help     bool
	//Verbose  bool
	//Silent   bool
	//Version  bool

	// Command-specific flags
	K8s             bool
	Kubernetes      bool
	Docker          bool
	Streamlit       bool
	AutoPrompt      bool
	Prompt          bool
	URL             string
	SequenceDiagram bool
	FlowchartTD     bool
	FlowchartLR     bool

	// File arguments
	AgentsFile   string
	SchemaFile   string
	WorkflowFile string
	YAMLFile     string
	TextFile     string
	Env          string

	// Commands
	Deploy     bool
	Run        bool
	Create     bool
	Validate   bool
	Mermaid    bool
	MetaAgents bool
}

// NewCommandArgs creates a new CommandArgs with default values
func NewCommandArgs() *CommandArgs {
	return &CommandArgs{}
}

// TestDeployCommand tests the deploy command
func TestDeployCommand(t *testing.T) {
	helper := NewTestHelper(t)

	// Test deploy with dry run and k8s flag
	t.Run("DryRunK8s", func(t *testing.T) {
		// Skip if DEPLOY_KUBERNETES_TEST is not set to "1"
		if os.Getenv("DEPLOY_KUBERNETES_TEST") != "1" {
			t.Skip("Skipping Kubernetes deploy test")
		}

		// Create command args
		args := NewCommandArgs()
		args.DryRun = true
		args.Verbose = true
		args.URL = "127.0.0.1:5000"
		args.K8s = true
		args.Deploy = true
		args.AgentsFile = helper.GetFixture("yamls/agents/simple_agent.yaml")
		args.WorkflowFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewBaseCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "deploy" {
			t.Errorf("Expected command name 'deploy', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test deploy with dry run and kubernetes flag
	t.Run("DryRunKubernetes", func(t *testing.T) {
		// Skip if DEPLOY_KUBERNETES_TEST is not set to "1"
		if os.Getenv("DEPLOY_KUBERNETES_TEST") != "1" {
			t.Skip("Skipping Kubernetes deploy test")
		}

		// Create command args
		args := NewCommandArgs()
		args.DryRun = true
		args.Verbose = true
		args.URL = "127.0.0.1:5000"
		args.Kubernetes = true
		args.Deploy = true
		args.AgentsFile = helper.GetFixture("yamls/agents/simple_agent.yaml")
		args.WorkflowFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "deploy" {
			t.Errorf("Expected command name 'deploy', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test deploy with dry run and docker flag
	t.Run("DryRunDocker", func(t *testing.T) {
		// Skip if DEPLOY_DOCKER_TEST is not set to "1"
		if os.Getenv("DEPLOY_DOCKER_TEST") != "1" {
			t.Skip("Skipping Docker deploy test")
		}

		// Create command args
		args := NewCommandArgs()
		args.DryRun = true
		args.Verbose = true
		args.URL = "127.0.0.1:5000"
		args.Docker = true
		args.Deploy = true
		args.AgentsFile = helper.GetFixture("yamls/agents/simple_agent.yaml")
		args.WorkflowFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "deploy" {
			t.Errorf("Expected command name 'deploy', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test deploy with auto prompt
	t.Run("WithAutoPrompt", func(t *testing.T) {
		// Create a temporary workflow file with a prompt
		workflowContent := `
spec:
  template:
    prompt: "This is a test input"
`
		tempFile := helper.CreateTempFile(workflowContent, ".yaml")

		// Create command args
		args := NewCommandArgs()
		args.DryRun = true
		args.Verbose = true
		args.URL = "127.0.0.1:5000"
		args.Docker = true
		args.AutoPrompt = true
		args.Deploy = true
		args.AgentsFile = helper.GetFixture("yamls/agents/simple_agent.yaml")
		args.WorkflowFile = tempFile

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "deploy" {
			t.Errorf("Expected command name 'deploy', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})
}

// TestRunCommand tests the run command
func TestRunCommand(t *testing.T) {
	helper := NewTestHelper(t)

	// Test run with dry run
	t.Run("DryRun", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.DryRun = true
		args.Run = true
		args.AgentsFile = helper.GetFixture("yamls/agents/simple_agent.yaml")
		args.WorkflowFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "run" {
			t.Errorf("Expected command name 'run', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test run with dry run and prompt
	t.Run("DryRunPrompt", func(t *testing.T) {
		// Create a mock stdin reader that returns "test prompt"
		originalStdin := os.Stdin
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create pipe: %v", err)
		}
		os.Stdin = r

		// Write test prompt to the pipe
		go func() {
			defer w.Close()
			w.Write([]byte("test prompt\n"))
		}()

		// Restore stdin after the test
		defer func() {
			os.Stdin = originalStdin
		}()

		// Create command args
		args := NewCommandArgs()
		args.DryRun = true
		args.Prompt = true
		args.Run = true
		args.AgentsFile = helper.GetFixture("yamls/agents/simple_agent.yaml")
		args.WorkflowFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "run" {
			t.Errorf("Expected command name 'run', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test run with nil agents file
	t.Run("NilAgentsFile", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.DryRun = true
		args.Run = true
		args.AgentsFile = ""
		args.WorkflowFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "run" {
			t.Errorf("Expected command name 'run', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test run with "None" agents file
	t.Run("NoneAgentsFile", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.DryRun = true
		args.Run = true
		args.AgentsFile = "None"
		args.WorkflowFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "run" {
			t.Errorf("Expected command name 'run', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})
}

// TestCreateCommand tests the create command
func TestCreateCommand(t *testing.T) {
	helper := NewTestHelper(t)

	// Test create with dry run
	t.Run("DryRun", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.DryRun = true
		args.Create = true
		args.AgentsFile = helper.GetFixture("yamls/agents/simple_agent.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "create" {
			t.Errorf("Expected command name 'create', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})
}

// TestCreateAndRunCommand tests combinations of create and run commands
func TestCreateAndRunCommand(t *testing.T) {
	helper := NewTestHelper(t)

	// Test create with dry run and no workflow file
	t.Run("CreateDryRun", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.DryRun = true
		args.Create = true
		args.Run = false
		args.AgentsFile = helper.GetFixture("yamls/agents/simple_agent.yaml")
		args.WorkflowFile = ""

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "create" {
			t.Errorf("Expected command name 'create', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test run with dry run and "None" agents file
	t.Run("RunDryRun", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.DryRun = true
		args.Create = false
		args.Run = true
		args.AgentsFile = "None"
		args.WorkflowFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "run" {
			t.Errorf("Expected command name 'run', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})
}

// TestValidateCommand tests the validate command
func TestValidateCommand(t *testing.T) {
	helper := NewTestHelper(t)

	// Test validate agents file
	t.Run("ValidateAgentsFile", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.Validate = true
		args.YAMLFile = helper.GetFixture("yamls/agents/simple_agent.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "validate" {
			t.Errorf("Expected command name 'validate', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test validate agents file with schema
	t.Run("ValidateAgentsFileWithSchema", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.Validate = true
		args.YAMLFile = helper.GetFixture("yamls/agents/simple_agent.yaml")
		args.SchemaFile = helper.GetSchema("agent_schema.json")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "validate" {
			t.Errorf("Expected command name 'validate', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test validate workflow file
	t.Run("ValidateWorkflowFile", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.Validate = true
		args.YAMLFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "validate" {
			t.Errorf("Expected command name 'validate', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test validate workflow file with schema
	t.Run("ValidateWorkflowFileWithSchema", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.Validate = true
		args.YAMLFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")
		args.SchemaFile = helper.GetSchema("workflow_schema.json")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "validate" {
			t.Errorf("Expected command name 'validate', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})
}

// TestMermaidCommand tests the mermaid command
func TestMermaidCommand(t *testing.T) {
	helper := NewTestHelper(t)

	// Test mermaid with sequence diagram
	t.Run("SequenceDiagram", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.Verbose = true
		args.Mermaid = true
		args.SequenceDiagram = true
		args.FlowchartTD = false
		args.FlowchartLR = false
		args.WorkflowFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "mermaid" {
			t.Errorf("Expected command name 'mermaid', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test mermaid with flowchart TD
	t.Run("FlowchartTD", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.Verbose = true
		args.Mermaid = true
		args.SequenceDiagram = false
		args.FlowchartTD = true
		args.FlowchartLR = false
		args.WorkflowFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "mermaid" {
			t.Errorf("Expected command name 'mermaid', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	// Test mermaid with flowchart LR
	t.Run("FlowchartLR", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.Verbose = true
		args.Mermaid = true
		args.SequenceDiagram = false
		args.FlowchartTD = false
		args.FlowchartLR = true
		args.WorkflowFile = helper.GetFixture("yamls/workflows/simple_workflow.yaml")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "mermaid" {
			t.Errorf("Expected command name 'mermaid', got '%s'", cmd.Name())
		}

		// Execute command
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})
}

// TestMetaAgentsCommand tests the meta-agents command
func TestMetaAgentsCommand(t *testing.T) {
	helper := NewTestHelper(t)

	// Test meta-agents
	t.Run("MetaAgents", func(t *testing.T) {
		// Create command args
		args := NewCommandArgs()
		args.Verbose = true
		args.MetaAgents = true
		args.TextFile = helper.GetFixture("agents/meta_agents/simple_prompt.txt")

		// Create and execute command
		cmd, err := commands.NewCommand(args)
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}

		// Verify command name
		if cmd.Name() != "meta-agents" {
			t.Errorf("Expected command name 'meta-agents', got '%s'", cmd.Name())
		}

		// Create a channel to signal test completion
		done := make(chan bool)

		// Execute command in a goroutine
		go func() {
			exitCode := cmd.Execute()
			if exitCode != 0 {
				t.Errorf("Expected exit code 0, got %d", exitCode)
			}
			done <- true
		}()

		// Wait for command to complete or timeout
		select {
		case <-done:
			// Command completed successfully
		case <-time.After(5 * time.Second):
			// Command timed out, kill the process
			if process := cmd.GetProcess(); process != nil {
				process.Kill()
			}
			t.Fatal("Command timed out")
		}
	})
}
