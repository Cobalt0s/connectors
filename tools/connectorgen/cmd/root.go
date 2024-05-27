package cmd

import (
	"errors"
	"fmt"
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

	if err := errors.Join(
		rootCmd.MarkPersistentFlagRequired("provider"),
		rootCmd.MarkPersistentFlagRequired("package"),
		viper.BindPFlag("provider", rootCmd.PersistentFlags().Lookup("provider")),
		viper.BindPFlag("package", rootCmd.PersistentFlags().Lookup("package")),
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
	result := &Recipe{
		Provider: viper.GetString("provider"),
		Package:  viper.GetString("package"),
	}
	result.Output = fmt.Sprintf("%v-output-gen", result.Package)

	return result
}
