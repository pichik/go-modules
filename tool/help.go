package tool

import (
	"fmt"
	"os"

	"github.com/pichik/go-modules/misc"
)

func PrintLogo() {
	fmt.Printf("%s%s%s\n", misc.Purple, logo(), misc.White)
}

func printToolName(name string) {
	if name != "" {
		fmt.Printf("%s%s%s\n", misc.Red, name, misc.White)
	}
}

func PrintDefaultHelp() {
	var tools string

	PrintLogo()

	for _, t := range toolOrder {
		if t == "error" {
			continue
		}
		tools = tools + fmt.Sprintf("\n\t%s\t\t- %s", ToolRegistry[t].Name, ToolRegistry[t].Description)
	}

	fmt.Printf(`
   %sSelect tool:%s %s

   %sExamples:%s
	 Read from file:
	   %scat urls.txt | tool [tool] [-h]%s
	 Make request from string:
	   %secho 'https://google.com' | tool [tool] [-h]%s
  `, misc.Yellow, misc.White, tools, misc.Yellow, misc.White, misc.Blue, misc.White, misc.Blue, misc.White)
	os.Exit(0)
}

func PrintToolHelp(toolData ToolData) {
	fmt.Printf("%s\n", toolData.Description)

	//INCLUDED UTILS INFO
	for _, util := range toolData.Utils {
		fmt.Printf("%s%s\n%s", misc.Green, util.Name, misc.White)

		//Print Flags
		if len(util.FlagDatas) > 0 {
			if util.Name != "" {
				fmt.Printf(" %sFlags:\n%s", misc.Yellow, misc.White)
			}
			for _, flag := range util.FlagDatas {
				if flag.Def != "" {
					flag.Def = fmt.Sprintf("%s\n\t\t(Default: %v)%s", misc.Gray, flag.Def, misc.White)
				}
				fmt.Printf("\t-%s%s\t%s%s\n", flag.Name, required(flag.Required), flag.Description, flag.Def)
			}
		}

		printExamples(util.Examples)
	}

	//MAIN TOOL INFO
	if len(toolData.FlagDatas) > 0 {
		fmt.Printf(" %sFlags:\n%s", misc.Yellow, misc.White)

		for _, flag := range toolData.FlagDatas {
			if flag.Def != "" {
				flag.Def = fmt.Sprintf("%s\n\t\t(Default: %v)%s", misc.Gray, flag.Def, misc.White)
			}
			fmt.Printf("\t-%s%s\t%s%s\n", flag.Name, required(flag.Required), flag.Description, flag.Def)
		}
	}
	printExamples(toolData.Examples)

	fmt.Printf("\n !!! Flags marked with %s, are required !!!\n", required(true))
	os.Exit(0)
}

func printExamples(examples map[string]string) {

	if len(examples) > 0 {
		fmt.Printf(" %sExamples:\n%s", misc.Yellow, misc.White)

		for description, example := range examples {
			fmt.Printf("\t%s:\n\t %s%s%s\n", description, misc.Blue, example, misc.White)
		}
	}
}

func logo() string {
	return fmt.Sprintf(`
			 /$$$$$$$$ /$$                       /$$$$$$$$                  /$$
			|__  $$__/| $$                      |__  $$__/                 | $$
			   | $$   | $$$$$$$   /$$$$$$          | $$  /$$$$$$   /$$$$$$ | $$
			   | $$   | $$__  $$ /$$__  $$         | $$ /$$__  $$ /$$__  $$| $$
			   | $$   | $$  \ $$| $$$$$$$$         | $$| $$  \ $$| $$  \ $$| $$
			   | $$   | $$  | $$| $$_____/         | $$| $$  | $$| $$  | $$| $$
			   | $$   | $$  | $$|  $$$$$$$         | $$|  $$$$$$/|  $$$$$$/| $$
			   |__/   |__/  |__/ \_______/         |__/ \______/  \______/ |__/
			`)
}

func required(is bool) string {
	if is {
		return fmt.Sprintf("%s%s%s", misc.Red, "!R", misc.White)
	}
	return ""
}
