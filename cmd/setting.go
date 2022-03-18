package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/spf13/viper"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "Set the setting",
	Run: func(cmd *cobra.Command, args []string) {
		setupSettings()
	},
}

var settingListCmd = &cobra.Command{
	Use:   "list",
	Short: "list of setting",
	Run: func(cmd *cobra.Command, args []string) {
		printSettingList()
	},
}

func init() {
	settingCmd.AddCommand(settingListCmd)
	rootCmd.AddCommand(settingCmd)
}

func printSettingList() {
	var rows []table.Row

	for _, key := range viper.AllKeys() {
		rows = append(rows, table.Row{key, viper.Get(key)})
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("List of Settings")
	t.AppendHeader(table.Row{"Key", "Value"})
	t.AppendRows(rows)
	t.AppendSeparator()
	t.SetCaption("Total: %d", len(rows))
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

	search := promptGetSelect(
		"Please, select a search mode",
		[]string{"Single", "Always"},
	)

	viper.Set("search", search)

	_ = viper.WriteConfig()
}
