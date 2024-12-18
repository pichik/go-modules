package tool

import (
	"flag"
)

var toolData *ToolData

type FlagData struct {
	Name        string
	Description string
	Required    bool
	Data        string
	VarStr      *string
	VarInt      *int
	VarBool     *bool
	VarAStr     *ArrayStringFlag
	Def         any
}

type Section struct {
	Name      string
	FlagDatas []FlagData
	Examples  map[string]string
}

type ArrayStringFlag []string

func (i *ArrayStringFlag) String() string {
	return ""
}

func (i *ArrayStringFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func CreateFlagSet(t ToolData) *ToolData {
	s := t.Description
	flagSet := flag.NewFlagSet(s, flag.ExitOnError)
	// toolData = &ToolData{FlagSet: flagSet}
	toolData = &t
	t.FlagSet = flagSet
	return toolData
}

func ParseFlags(args []string) {
	toolData.FlagSet.Parse(args)
}

func SetupFlags() {
	for _, fl := range toolData.FlagDatas {
		switch fl.Def.(type) {
		case string:
			toolData.FlagSet.StringVar(fl.VarStr, fl.Name, fl.Def.(string), fl.Description)
		case int:
			toolData.FlagSet.IntVar(fl.VarInt, fl.Name, fl.Def.(int), fl.Description)
		case bool:
			toolData.FlagSet.BoolVar(fl.VarBool, fl.Name, fl.Def.(bool), fl.Description)
		case ArrayStringFlag:
			toolData.FlagSet.Var(fl.VarAStr, fl.Name, fl.Description)
		}
	}
	for _, util := range toolData.Utils {
		for _, fl := range util.FlagDatas {
			switch fl.Def.(type) {
			case string:
				toolData.FlagSet.StringVar(fl.VarStr, fl.Name, fl.Def.(string), fl.Description)
			case int:
				toolData.FlagSet.IntVar(fl.VarInt, fl.Name, fl.Def.(int), fl.Description)
			case bool:
				toolData.FlagSet.BoolVar(fl.VarBool, fl.Name, fl.Def.(bool), fl.Description)
			case ArrayStringFlag:
				toolData.FlagSet.Var(fl.VarAStr, fl.Name, fl.Description)
			}
		}
	}
}

func UpdateFlagUsageHelp() {
	printToolName(toolData.AName)
	toolData.FlagSet.Usage = func() {
		PrintToolHelp(*toolData)
	}
}
