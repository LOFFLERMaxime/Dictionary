package main

import (
	"bufio"
	"estiam/dictionary"
	"fmt"
	"os"
	"strings"
)

func main() {
	dict := dictionary.New()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter action (add/define/remove/list/exit): ")
		action, _ := reader.ReadString('\n')
		action = strings.TrimSpace(action)

		switch action {
		case "add":
			actionAdd(dict, reader)
		case "define":
			actionDefine(dict, reader)
		case "remove":
			actionRemove(dict, reader)
		case "list":
			actionList(dict)
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Invalid action. Please try again.")
		}
	}
}

func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	fmt.Print("Enter definition: ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	d.Add(word, definition)
	fmt.Printf("Word '%s' added to the dictionary.\n", word)
}

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word to define: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	entry, err := d.Get(word)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("%s: %s\n", word, entry)
	}
}

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word to remove: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	d.Remove(word)
	fmt.Printf("Word '%s' removed from the dictionary.\n", word)
}

func actionList(d *dictionary.Dictionary) {
	words, entries := d.List()

	fmt.Println("Dictionary:")
	for _, word := range words {
		fmt.Printf("%s: %s\n", word, entries[word])
	}
}
