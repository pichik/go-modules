package tool

import (
	"flag"
)

type IUtil interface {
	SetupFlags() []UtilData
	SetupData()
}

type ITool interface {
	SetupFlags()
	SetupInput(urls []string)
}

type ToolData struct {
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

// ToolRegistry holds all the tools registered in the system
var ToolRegistry = make(map[string]ToolData)
var toolOrder []string

// RegisterTool registers a new tool dynamically
func RegisterTool(name, description string, examples map[string]string) ToolData {
	tool := ToolData{
		Name:        name,
		Description: description,
		Examples:    examples,
	}

	// Register the tool by adding it to the global registry
	ToolRegistry[name] = tool
	toolOrder = append(toolOrder, name)
	return tool
}

// GetTool returns a tool based on its name from the registry
func GetTool(toolName string) ToolData {
	tool, exists := ToolRegistry[toolName]
	if !exists {
		// If the tool doesn't exist, return a default error tool
		return ToolData{Name: "error"}
	}
	return tool
}

// GetTools returns a list of all tools in the registry
// func GetTools() []ToolData {
// 	var tools []ToolData
// 	for _, tool := range ToolRegistry {
// 		tools = append(tools, tool)
// 	}
// 	return tools
// }
