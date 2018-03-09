package main

import (
	"log"

	"github.com/dtynn/jmgen/parser"
	"github.com/spf13/cobra"
)

func main() {
	var opt parser.Options
	var ffmt string

	rootCmd := cobra.Command{
		Use:   "jmgen",
		Short: "jmgen",
		Long:  "add tags for go structs",
		Run: func(cmd *cobra.Command, args []string) {
			opt.FieldFormat = parser.FieldFormat(ffmt)
			if err := opt.Parse(); err != nil {
				log.Println(err)
			}
		},
	}

	rootCmd.Flags().StringVarP(&opt.Input, "input", "i", "", "input file path")
	rootCmd.Flags().StringSliceVarP(&opt.Structs, "structs", "s", nil, "specified struct names")
	rootCmd.Flags().IntSliceVarP(&opt.Lines, "lines", "l", nil, "specified lines")
	rootCmd.Flags().StringSliceVarP(&opt.Tags, "tags", "t", nil, "tags to add")
	rootCmd.Flags().StringVarP(&ffmt, "format", "f", "", "field name format type, default empty. \"camel\" or \"snake\" allowed")
	rootCmd.Flags().BoolVarP(&opt.Rewrite, "rewrite", "r", false, "rewrite src file, default false")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
