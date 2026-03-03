package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/direnv/direnv/v2/pkg/dotenv"
)

// CmdDotEnv is `direnv dotenv [SHELL [PATH_TO_DOTENV]]`
// Transforms a .env file to evaluatable `export KEY=PAIR` statements.
//
// See: https://github.com/bkeepers/dotenv and https://github.com/ddollar/foreman
var CmdDotEnv = &Cmd{
	Name:    "dotenv",
	Desc:    "Transforms a .env file to evaluatable `export KEY=PAIR` statements",
	Args:    []string{"[SHELL]", "[PATH_TO_DOTENV]"},
	Private: true,
	Action:  actionSimple(cmdDotEnvAction),
}

func cmdDotEnvAction(_ Env, args []string) (err error) {
	var shell Shell
	var newenv Env
	var target string

	// NOTE: adjust indexing here if your framework passes args differently.
	if len(args) > 1 {
		shell = DetectShell(args[1])
	} else {
		shell = Bash
	}

	if len(args) > 2 {
		target = args[2]
	}
	if target == "" {
		target = ".env"
	}

	file, err := os.Open(target)
	if err != nil {
		return err
	}
	defer file.Close()

	var data strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data.WriteString(scanner.Text())
		data.WriteString("\n")
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return scanErr
	}

	// Set PWD env var to the directory the .env file resides in.
	path, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	if err := os.Setenv("PWD", filepath.Dir(path)); err != nil {
		return err
	}

	newenv, err = dotenv.Parse(data.String())
	if err != nil {
		return err
	}

	str, err := newenv.ToShell(shell)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}