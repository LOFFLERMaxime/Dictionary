package dictionary

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Entry struct {
	Definition string `json:"definition"`
}

func (e Entry) String() string {

	return e.Definition
}

type Dictionary struct {
	entries  map[string]Entry
	filePath string
}

func New(filePath string) *Dictionary {
	d := &Dictionary{
		entries:  make(map[string]Entry),
		filePath: filePath,
	}
	d.load()
	return d
}

func (d *Dictionary) Save() error {
	data, err := json.Marshal(d.entries)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(d.filePath, data, 0644)
	if err != nil {
		log.Printf("Error saving dictionary to file: %v\n", err)
		return err
	}

	log.Println("Dictionary saved to file.")
	return nil
}

func (d *Dictionary) load() {
	if _, err := os.Stat(d.filePath); os.IsNotExist(err) {
		log.Printf("Dictionary file not found. Creating a new file: %s\n", d.filePath)
		return
	}

	data, err := ioutil.ReadFile(d.filePath)
	if err != nil {
		log.Printf("Error reading dictionary file: %v\n", err)
		return
	}

	err = json.Unmarshal(data, &d.entries)
	if err != nil {
		log.Printf("Error unmarshalling dictionary data: %v\n", err)
		return
	}
	log.Println("Dictionary loaded from file.")
}

func (d *Dictionary) Add(word string, definition string) {
	entry := Entry{Definition: definition}
	d.entries[word] = entry
}

func (d *Dictionary) Get(word string) (Entry, error) {

	entry, found := d.entries[word]
	if !found {
		return Entry{}, errors.New("Word not found in the dictionary")
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) {
	delete(d.entries, word)
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	words := make([]string, 0, len(d.entries))
	for word := range d.entries {
		words = append(words, word)
	}
	return words, d.entries
}
