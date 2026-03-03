package cmd

import (
	"fmt"
	"github.com/direnv/direnv/v2/pkg/dotenv"
	"os"
	"bufio"
	"strings"
	"path/filepath"
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

	if file, err = os.Open(target); err != nil {
		return
	}

	defer file.Close()

    var data strings.Builder

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        data.WriteString(scanner.Text())
        data.WriteString("\n")
    }

    if err := scanner.Err(); err != nil {
        return
    }

	// Set PWD env var to the directory the .env file resides in. This results
	// in the least amount of surprise, as a dotenv file is most often defined
	// in the same directory it's loaded from, so referring to PWD should match
	// the directory of the .env file.
	path, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	if err := os.Setenv("PWD", filepath.Dir(path)); err != nil {
		return err
	}

	newenv, err = dotenv.Parse(data)
	if err != nil {
		return err
	}

	str, err := newenv.ToShell(shell)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return
}
