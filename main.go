package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func gitDir(u *url.URL) (string, error) {
	cleanPath := path.Clean(u.Path)
	pathSlice := strings.Split(cleanPath, "/")
	if len(pathSlice) < 3 {
		return "", errors.New("a forge URL should have at least a user and a project")
	}
	user := strings.ToLower(strings.TrimPrefix(pathSlice[1], "~"))
	project := strings.ToLower(pathSlice[2])
	return path.Join(os.Getenv("HOME"), "src", u.Host, user, project), nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("requires a git URL as an argument")
	}

	arg := os.Args[1]
	u, err := url.Parse(arg)
	if err != nil {
		log.Fatalf("error %T parsing url: %v", err, err)
	}
	dir, err := gitDir(u)
	if err != nil {
		log.Fatalf("error parsing URL: %v", err)
	}
	err = os.MkdirAll(filepath.Dir(dir), 0755)
	if err != nil {
		log.Fatalf("error %T creating directory: %v", err, err)
	}
	err = os.Chdir(filepath.Dir(dir))
	if err != nil {
		log.Fatalf("error changing to directory %q: %v", dir, err)
	}
	clone := exec.Command("git", "clone", u.String())
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
