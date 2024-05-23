package cmd

import (
	"errors"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "conn-gen",
	Short: "Connector generator",
	Long:  "Generates the base template to start building connector with",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("unhandled error: %v", err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("provider", "n", "", "Provider name")
	rootCmd.PersistentFlags().StringP("package", "p", "", "Package name")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Toggle manual test generation")

	if err := errors.Join(
		rootCmd.MarkPersistentFlagRequired("provider"),
		rootCmd.MarkPersistentFlagRequired("package"),
		rootCmd.MarkPersistentFlagRequired("output"),
		viper.BindPFlag("provider", rootCmd.PersistentFlags().Lookup("provider")),
		viper.BindPFlag("package", rootCmd.PersistentFlags().Lookup("package")),
		viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output")),
	); err != nil {
		log.Fatal(err)
	}
}

type Recipe struct {
	Provider string
	Package  string
	Output   string
}

func GetRecipe() *Recipe {
	return &Recipe{
		Provider: viper.GetString("provider"),
		Package:  viper.GetString("package"),
		Output:   viper.GetString("output"),
	}
}
