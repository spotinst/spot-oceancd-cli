package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type FilesFlag struct {
	files []string
}

func (f *FilesFlag) String() string {
	s := ""
	for _, v := range f.files {
		s = s + fmt.Sprintf(",%s", v)
	}
	return s
}
func (f *FilesFlag) Set(s string) error {

	f.files = append(f.files, s)

	return nil
}
func (f *FilesFlag) Type() string {

	return "string"
}

func AddFileFlags(cmd *cobra.Command) {

	//cmd.PersistentFlags().VarP(&filesVar, "file", "f", "manifest file with resource definition")

}
