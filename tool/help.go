package tool

import (
	"fmt"
	"os"

	"github.com/pichik/go-modules/output"
)

func PrintLogo() {
	fmt.Printf("%s%s%s\n", output.Purple, logo(), output.White)
}

func PrintToolName(name string) {
	fmt.Printf("%s%s%s\n", output.Red, name, output.White)
}

func PrintDefaultHelp() {
	var tools string

	for _, t := range GetTools() {
		if t.Name == "error" {
			continue
		}
		tools = tools + fmt.Sprintf("\n\t%s\t\t- %s", t.Name, t.Description)
	}

	fmt.Printf(`
   %sSelect tool:%s %s

   %sExamples:%s
	 Read from file:
	   %scat urls.txt | thetool [tool] [-h]%s
	 Make request from string:
	   %secho 'google.com' | thetool [tool] [-h]%s
  `, output.Yellow, output.White, tools, output.Yellow, output.White, output.Blue, output.White, output.Blue, output.White)
	os.Exit(0)
}

func PrintToolHelp(toolData ToolData) {
	fmt.Printf("%s\n", toolData.Description)

	//INCLUDED UTILS INFO
	for _, util := range toolData.Utils {
		fmt.Printf("%s%s\n%s", output.Green, util.Name, output.White)

		//Print Flags
		if len(util.FlagDatas) > 0 {
			if util.Name != "" {
				fmt.Printf(" %sFlags:\n%s", output.Yellow, output.White)
			}
			for _, flag := range util.FlagDatas {
				if flag.Def != "" {
					flag.Def = fmt.Sprintf("%s\n\t\t(Default: %v)%s", output.Gray, flag.Def, output.White)
				}
				fmt.Printf("\t-%s%s\t%s%s\n", flag.Name, required(flag.Required), flag.Description, flag.Def)
			}
		}

		printExamples(util.Examples)
	}

	//MAIN TOOL INFO
	if len(toolData.FlagDatas) > 0 {
		fmt.Printf(" %sFlags:\n%s", output.Yellow, output.White)

		for _, flag := range toolData.FlagDatas {
			if flag.Def != "" {
				flag.Def = fmt.Sprintf("%s\n\t\t(Default: %v)%s", output.Gray, flag.Def, output.White)
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
		fmt.Printf(" %sExamples:\n%s", output.Yellow, output.White)

		for description, example := range examples {
			fmt.Printf("\t%s:\n\t %s%s%s\n", description, output.Blue, example, output.White)
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
		return fmt.Sprintf("%s%s%s", output.Red, "!R", output.White)
	}
	return ""
}
