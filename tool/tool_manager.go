package tool

import (
	"flag"
)

type ToolData struct {
	Tool      Tool
	FlagSet   *flag.FlagSet
	FlagDatas []FlagData
	// Sections    []Section
	Utils       []UtilData
	AName       string
	Name        string
	Description string
	Examples    map[string]string
}

type UtilData struct {
	Name      string
	FlagDatas []FlagData
	Examples  map[string]string
}

// Tool represents a type of tool with a name and a description
type Tool struct {
	Name        string
	Description string
	Examples    map[string]string
}

// ToolRegistry holds all the tools registered in the system
var ToolRegistry = make(map[string]Tool)

// RegisterTool registers a new tool dynamically
func RegisterTool(name, description string, examples map[string]string) Tool {
	tool := Tool{
		Name:        name,
		Description: description,
		Examples:    examples,
	}

	// Register the tool by adding it to the global registry
	ToolRegistry[name] = tool
	return tool
}

// GetTool returns a tool based on its name from the registry
func GetTool(toolName string) Tool {
	tool, exists := ToolRegistry[toolName]
	if !exists {
		// If the tool doesn't exist, return a default error tool
		return Tool{Name: "error", Description: "Tool not found"}
	}
	return tool
}

// GetTools returns a list of all tools in the registry
func GetTools() []Tool {
	var tools []Tool
	for _, tool := range ToolRegistry {
		tools = append(tools, tool)
	}
	return tools
}
