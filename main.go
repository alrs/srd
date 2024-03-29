/*
srd - generates structured home source directory
Copyright (C) 2023 Lars Lehtonen

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type ErrShortURL struct{}

func (e ErrShortURL) Error() string {
	return fmt.Sprint("a forge URL should have at least a user and a project")
}

func paths(u *url.URL) (*url.URL, string, error) {
	cleanPath := path.Clean(u.Path)
	pathSlice := strings.Split(cleanPath, "/")
	if len(pathSlice) < 3 {
		return u, "", ErrShortURL{}
	}
	user := pathSlice[1]
	project := pathSlice[2]
	cleanUser := strings.ToLower(strings.TrimPrefix(user, "~"))
	cleanProject := strings.ToLower(project)
	gitDir := path.Join(u.Host, cleanUser, cleanProject)
	var nu url.URL
	nu.Scheme = u.Scheme
	nu.Host = u.Host
	nu.Path = path.Join(user, project)
	return &nu, gitDir, nil
}

func main() {
	var root string
	flag.StringVar(
		&root, "root",
		path.Join(os.Getenv("HOME"), "src"),
		"root path to clone projects",
	)
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("requires a git URL as an argument")
	}

	arg := flag.Args()[0]
	u, err := url.Parse(arg)
	if err != nil {
		log.Fatalf("error %T parsing url: %v", err, err)
	}

	gitUrl, relPath, err := paths(u)
	if err != nil {
		log.Fatal(err)
	}

	dir := path.Join(root, relPath)

	err = os.MkdirAll(filepath.Dir(dir), 0755)
	if err != nil {
		log.Fatalf("error %T creating directory: %v", err, err)
	}
	err = os.Chdir(filepath.Dir(dir))
	if err != nil {
		log.Fatalf("error changing to directory %q: %v", dir, err)
	}
	clone := exec.Command("git", "clone", gitUrl.String(), strings.ToLower(path.Base(dir)))
	// the only thing we want to go to stdout is the full path of
	// the git repo
	clone.Stdout = os.Stderr
	clone.Stderr = os.Stderr
	err = clone.Run()
	if err != nil {
		log.Fatalf("error %T cloning repo %q: %v", err, arg, err)
	}
	fmt.Println(dir)
}
