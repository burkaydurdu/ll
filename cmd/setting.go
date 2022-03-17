package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/spf13/viper"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// settingCmd represents the setting command
var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "A brief description of your command",
	Long:  `Settings`,
	Run: func(cmd *cobra.Command, args []string) {
		if contains(args, "list") {
			printSettingList()
		} else {
			setupSettings()
		}
	},
}

func init() {
	settingCmd.Flags().BoolP("list", "l", false, "list of setting")

	rootCmd.AddCommand(settingCmd)
}

func printSettingList() {
	var rows []table.Row

	for _, key := range viper.AllKeys() {
		rows = append(rows, table.Row{key, viper.Get(key)})
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "Value"})
	t.AppendRows(rows)
	t.AppendSeparator()
	t.AppendFooter(table.Row{"", ""})
	t.Render()
}

func promptGetSelect(command string, items []string) string {
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: command,
			Items: items,
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func setupSettings() {
	mode := promptGetSelect(
		"Please, select a translate mode",
		[]string{"English - Turkish", "Turkish - English"},
	)

	viper.Set("mode", mode)

	_ = viper.WriteConfig()
}
