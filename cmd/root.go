package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/manifoldco/promptui"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config File Name
var configFile string

var rootCmd = &cobra.Command{
	Use:   "ll",
	Short: "CLI base translator",
	Long:  "This CLI is translator. You should translate a word and You can see type of word",
	Run: func(cmd *cobra.Command, args []string) {
		content := search()

		if content != "" {
			phrases, _ := FetchFromTureng(content)
			if len(phrases) > 0 {
				printPhrases(phrases[:10])
			}

			if viper.GetString("search") == "Always" {
				_ = cmd.Execute()
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.ll.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Create config file
		err = createConfigFile(home + "/.ll")

		if err != nil {
			fmt.Println(err)
		}

		viper.SetConfigName(".ll")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("couldn't find file")
	}

	// Default mode
	if viper.GetString("mode") == "" {
		viper.Set("mode", "Turkish - English")

		_ = viper.WriteConfig()
	}
}

// File exist control
func createConfigFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}

		err = f.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func printPhrases(phrases Phrases) {
	var rows []table.Row

	for _, item := range phrases[1:] {
		rows = append(rows, table.Row{item.Source, item.Target, item.Category, item.Type})
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{phrases[0].Source, phrases[0].Target, phrases[0].Category, phrases[0].Type})
	t.AppendRows(rows)
	t.AppendSeparator()
	t.Render()
}

func search() string {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("content must have more than 3 characters")
		}

		return nil
	}

	var search string

	prompt := promptui.Prompt{
		Label:    "Search ",
		Validate: validate,
		Default:  search,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}
