package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eswar-7116/k8s-reconciler/internal/apiserver"
)

func main() {
	fmt.Println("k8s reconciler simulation")
	fmt.Println("commands: set <n>, status, exit")
	fmt.Println()

	api := apiserver.NewAPIServer(0)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "set":
			if len(parts) != 2 {
				fmt.Println("usage: set <replicas>")
				continue
			}
			n, err := strconv.Atoi(parts[1])
			if err != nil || n < 0 {
				fmt.Println("replicas must be a non-negative integer")
				continue
			}
			api.SetReplicas(n)
			time.Sleep(5 * time.Millisecond)

		case "status":
			pods, replicas := api.State()
			fmt.Printf("Desired no. of replicas: %d\nNo. of running pods: %d\n", replicas, pods)

		case "exit":
			api.Close()
			return

		default:
			fmt.Println("unknown command, try: set <n>, status, exit")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
	}

	api.Close()
}
