package cmd

import (
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var baseCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "base",
	Short: "Create connector",
	Long:  "Provides a template of connector struct with sample constructor",
	Run: func(cmd *cobra.Command, args []string) {
		recipe := GetRecipe()
		applyTemplatesFromDirectory("base", recipe,
			filepath.Join(recipe.Output, recipe.Package),
		)
	},
}

var readCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "read objectName <objectName, ex: contact, user>",
	Short: "Create read method",
	Long: "Provides a template to start implementing a read method. " +
		"Manual test will have read template for objectName",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		recipe := GetRecipe()
		applyTemplatesFromDirectory("read", recipe,
			filepath.Join(recipe.Output, recipe.Package),
		)
		createManualTest(recipe, strcase.ToCamel(args[0]), "read")
	},
}

var writeCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "write objectName <objectName, ex: contact, user>",
	Short: "Create write method",
	Long: "Provides a template to start implementing a write method. " +
		"Manual test will have create/update/delete template for objectName",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		recipe := GetRecipe()
		applyTemplatesFromDirectory("write", recipe,
			filepath.Join(recipe.Output, recipe.Package),
		)
		createManualTest(recipe, strcase.ToCamel(args[0]), "write-delete")
	},
}

var deleteCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "delete objectName <objectName, ex: contact, user>",
	Short: "Create delete method",
	Long: "Provides a template to start implementing a delete method. " +
		"Manual test will have create/update/delete template for objectName",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		recipe := GetRecipe()
		applyTemplatesFromDirectory("delete", recipe,
			filepath.Join(recipe.Output, recipe.Package),
		)
		createManualTest(recipe, strcase.ToCamel(args[0]), "write-delete")
	},
}

var metadataCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "metadata objectName <objectName, ex: contact, user>",
	Short: "Create metadata method",
	Long: "Provides a template to start implementing a list object metadata method. " +
		"Manual test will have template for objectName",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		recipe := GetRecipe()
		applyTemplatesFromDirectory("metadata", recipe,
			filepath.Join(recipe.Output, recipe.Package),
		)
		createManualTest(recipe, strcase.ToCamel(args[0]), "metadata")
	},
}

func createManualTest(recipe *Recipe, objectName string, directory string) {
	applyTemplatesFromDirectory("test", recipe,
		filepath.Join(recipe.Output, "test", recipe.Package),
	)

	type ManualTestParams struct {
		*Recipe
		ObjectName string
	}

	data := &ManualTestParams{
		Recipe:     recipe,
		ObjectName: objectName,
	}
	applyTemplatesFromDirectory(filepath.Join("test", directory), data,
		filepath.Join(recipe.Output, "test", recipe.Package, directory),
	)
}

func init() {
	rootCmd.AddCommand(baseCmd)
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(writeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(metadataCmd)
}
