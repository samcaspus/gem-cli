package utils

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func ExecuteCommand(command string) {

	parts := strings.Fields(command)

	if len(parts) == 0 {
		fmt.Println("no command provided")
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("failed to start command: %w", err)
	}

	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	go func() {
		for stdoutScanner.Scan() {
			fmt.Println(stdoutScanner.Text())
		}
	}()

	go func() {
		for stderrScanner.Scan() {
			fmt.Println(stderrScanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("command finished with error: %w", err)
	}

}
