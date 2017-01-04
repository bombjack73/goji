// Copyright © 2017 Marc Vandenbosch
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// autobuildCmd represents the autobuild command
var autobuildCmd = &cobra.Command{
	Use:   "autobuild [packagename]",
	Short: "Automatically build packages",
	Long: `Automatically build the packages passed as parameter.
As soon as a file is modified in the package, go build is called on it.autobuildCmd	
Example: goji autobuild github.com/myusername/mypackage`,
	Run: func(cmd *cobra.Command, args []string) {

		//todo handle multiple args
		//todo check that directory exist
		prefixPath := os.Getenv("GOPATH") + "/src/" //todo use os path separator
		lastEvent := ""
		lastTime := time.Now()

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		done := make(chan bool)
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					//log.Println("event:", event)
					if event.Op&fsnotify.Write == fsnotify.Write {
						dobuild := false
						now := time.Now()
						//the following lastEvent and lastTime are here because fsnotify notifies twice for the same file
						if event.Name == lastEvent {
							if now.Sub(lastTime).Seconds() > 1 { //skip built if built a seconds ago
								dobuild = true
							}
						} else {
							dobuild = true
						}
						if dobuild {
							packageName := extractPackageName(prefixPath, event.Name)
							log.Println("Building", packageName)
							cmd := exec.Command("go", "install", packageName)
							output, err := cmd.CombinedOutput()
							//todo : better hangling of execution and error
							if err != nil {
								os.Stderr.WriteString(err.Error())
							}
							fmt.Print(string(output))

							lastEvent = event.Name
							lastTime = now
						}
					}
				case err := <-watcher.Errors:
					log.Println("error:", err)
				}
			}
		}()

		//todo could have absolute or relative to gopath
		//todo notify only if go files are modified

		filepath.Walk(prefixPath+args[0], func(path string, info os.FileInfo, err error) error {
			//todo add exclude dir list somewhere
			//handle the following by doing parsing by myself to exclude all .git subdirs directly
			if info.IsDir() && !strings.Contains(path, ".git") {
				fmt.Println("Watching", path)
				err = watcher.Add(path)
				if err != nil {
					log.Fatal(err)
				}
			}
			return nil
		})

		<-done
	},
}

func extractPackageName(prefixPath string, filePath string) (packageName string) {
	if strings.HasPrefix(filePath, prefixPath) {
		sep := string(filepath.Separator)
		dirs := strings.Split(filePath[len(prefixPath):], sep)
		numdirs := len(dirs)
		if numdirs > 3 {
			numdirs = 3
		}
		var buffer bytes.Buffer
		for i := 0; i < numdirs; i++ {
			buffer.WriteString(dirs[i])
			if i < numdirs-1 {
				buffer.WriteString(sep)
			}
		}
		packageName = buffer.String()
	}
	return
}

func init() {
	RootCmd.AddCommand(autobuildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// autobuildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// autobuildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
