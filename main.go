package main

import (
	"context"
	"log"

	"github.com/bans1mp/mcp-server/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := initMcpServer()

	initMcpTools(server)

	err := server.Run(context.Background(), &mcp.StdioTransport{}); 
	if err != nil {
		log.Fatal(err)
	}
}

func initMcpServer() *mcp.Server {
	return mcp.NewServer(&mcp.Implementation{Name: "mcp-server", Version: "v1.0.0"}, nil)
}

func initMcpTools(server *mcp.Server) {
	mcp.AddTool(server, tools.UpdateNotesTool, tools.UpdateNotes)
}