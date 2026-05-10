package main

import "fmt"

const version = "0.1.0"

func main() {
	printBanner()
}

func printBanner() {
	fmt.Println(`
  ██████  ██  ██████  ██████  ██      ███████ ██   ██     ██████  ██
 ██       ██ ██       ██   ██ ██      ██       ██ ██      ██   ██ ██
 ██   ███ ██ ██   ███ ██████  ██      █████     ███       ███████ ██
 ██    ██ ██ ██    ██ ██      ██      ██       ██ ██      ██   ██ ██
  ██████  ██  ██████  ██      ███████ ███████ ██   ██     ██   ██ ██
	`)
	fmt.Printf("  gigplex.ai v%s — background jobs with agentic AI observability\n\n", version)
}
