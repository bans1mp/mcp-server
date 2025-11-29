package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type UpdateNotesInput struct {
	Changes string `json:"changes" jsonschema:"the changes to be added to the notes"`
}

type UpdateNotesOutput struct {
	Success bool `json:"success" jsonschema:"indicates if the notes were successfully updated"`
}

var UpdateNotesTool = &mcp.Tool{
	Name:        "update_notes",
	Description: "Update the notes with the latest changes",
}

func UpdateNotes(ctx context.Context, req *mcp.CallToolRequest, input *UpdateNotesInput) (*mcp.CallToolResult, *UpdateNotesOutput, error) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return nil, &UpdateNotesOutput{}, fmt.Errorf("failed to get current file path")
	}
	basePath := filepath.Dir(b)
	filePath := filepath.Join(basePath, "../assets/notes.txt")

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, &UpdateNotesOutput{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(input.Changes + "\n"); err != nil {
		return nil, &UpdateNotesOutput{}, fmt.Errorf("failed to write to file: %w", err)
	}

	return &mcp.CallToolResult{}, &UpdateNotesOutput{
		Success: true,
	}, nil
}