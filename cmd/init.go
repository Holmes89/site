/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"github.com/Holmes89/personal-site/site/internal"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize site",
	Long: `Create basic structure required for the tool to compile source and run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.BuildDirStruct("./content"); err != nil {
			return err
		}
		return CreateIndexFiles()
	},
}

func CreateIndexFiles() error {
	d := []byte("# Change Me")
	err := ioutil.WriteFile("./content/posts/index.md", d, 0644)
	if err != nil && err != os.ErrExist {
		return errors.New("unable to create directory")
	}
	err = ioutil.WriteFile("./content/projects/index.md", d, 0644)
	if err != nil && err != os.ErrExist {
		return errors.New("unable to create directory")
	}
	err = ioutil.WriteFile("./content/index.md", d, 0644)
	if err != nil && err != os.ErrExist {
		return errors.New("unable to create directory")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
