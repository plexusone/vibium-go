package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// StartRecording tool

type StartRecordingInput struct {
	Name        string `json:"name,omitempty" jsonschema:"description=Name for the recorded script"`
	Description string `json:"description,omitempty" jsonschema:"description=Description of what the script tests"`
	BaseURL     string `json:"baseUrl,omitempty" jsonschema:"description=Base URL for relative URLs in the script"`
}

type StartRecordingOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleStartRecording(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StartRecordingInput,
) (*mcp.CallToolResult, StartRecordingOutput, error) {
	recorder := s.session.Recorder()

	if recorder.IsRecording() {
		return nil, StartRecordingOutput{}, fmt.Errorf("recording already in progress (use stop_recording first)")
	}

	recorder.Start(RecorderMetadata{
		Name:        input.Name,
		Description: input.Description,
		BaseURL:     input.BaseURL,
	})

	msg := "Recording started"
	if input.Name != "" {
		msg = fmt.Sprintf("Recording started: %s", input.Name)
	}

	return nil, StartRecordingOutput{Message: msg}, nil
}

// StopRecording tool

type StopRecordingInput struct{}

type StopRecordingOutput struct {
	Message   string `json:"message"`
	StepCount int    `json:"stepCount"`
}

func (s *Server) handleStopRecording(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StopRecordingInput,
) (*mcp.CallToolResult, StopRecordingOutput, error) {
	recorder := s.session.Recorder()

	if !recorder.IsRecording() {
		return nil, StopRecordingOutput{}, fmt.Errorf("no recording in progress")
	}

	recorder.Stop()
	count := recorder.StepCount()

	return nil, StopRecordingOutput{
		Message:   fmt.Sprintf("Recording stopped with %d steps", count),
		StepCount: count,
	}, nil
}

// ExportScript tool

type ExportScriptInput struct {
	Format string `json:"format,omitempty" jsonschema:"description=Output format: json or yaml (default: json),enum=json,enum=yaml"`
}

type ExportScriptOutput struct {
	Script    string `json:"script"`
	StepCount int    `json:"stepCount"`
	Format    string `json:"format"`
}

func (s *Server) handleExportScript(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ExportScriptInput,
) (*mcp.CallToolResult, ExportScriptOutput, error) {
	recorder := s.session.Recorder()
	count := recorder.StepCount()

	if count == 0 {
		return nil, ExportScriptOutput{}, fmt.Errorf("no steps recorded")
	}

	format := input.Format
	if format == "" {
		format = "json"
	}

	var scriptBytes []byte
	var err error

	switch format {
	case "json":
		scriptBytes, err = recorder.ExportJSON()
	case "yaml":
		// For now, just use JSON - could add YAML support later
		scriptBytes, err = recorder.ExportJSON()
		format = "json" // Report actual format used
	default:
		return nil, ExportScriptOutput{}, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return nil, ExportScriptOutput{}, fmt.Errorf("export failed: %w", err)
	}

	return nil, ExportScriptOutput{
		Script:    string(scriptBytes),
		StepCount: count,
		Format:    format,
	}, nil
}

// RecordingStatus tool

type RecordingStatusInput struct{}

type RecordingStatusOutput struct {
	Recording bool `json:"recording"`
	StepCount int  `json:"stepCount"`
}

func (s *Server) handleRecordingStatus(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input RecordingStatusInput,
) (*mcp.CallToolResult, RecordingStatusOutput, error) {
	recorder := s.session.Recorder()

	return nil, RecordingStatusOutput{
		Recording: recorder.IsRecording(),
		StepCount: recorder.StepCount(),
	}, nil
}

// ClearRecording tool

type ClearRecordingInput struct{}

type ClearRecordingOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClearRecording(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClearRecordingInput,
) (*mcp.CallToolResult, ClearRecordingOutput, error) {
	recorder := s.session.Recorder()
	recorder.Clear()

	return nil, ClearRecordingOutput{
		Message: "Recording cleared",
	}, nil
}
