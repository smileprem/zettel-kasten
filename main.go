package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

const TimeFormat = "20060102150405"
const NewLineChar = '\n'
const ZettelKastenLocation = "~/Library/Mobile Documents/iCloud~md~obsidian/Documents/smileprem/zettelkasten/"

func main() {
	fmt.Println("------------------------------------------")
	fmt.Println("Create a new zettel (note) in kasten (box)")
	fmt.Println("------------------------------------------")

	zettelID := time.Now().Format(TimeFormat)
	var zettelTemplate string
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("title: ")
	zettelTitle, err := reader.ReadString(NewLineChar)
	if err != nil {
		fmt.Println("Error reading the title", err.Error())
	}

	fmt.Print("tags (comma de-limited): ")
	zettelTags, err := reader.ReadString(NewLineChar)
	if err != nil {
		fmt.Println("Error reading the tags", err.Error())
	}

	zettelTemplate += "# " + zettelID + " " + zettelTitle
	zettelTemplate += formatZettelTags(zettelTags) + string(NewLineChar)
	zettelTemplate += "## Links" + string(NewLineChar)
	zettelTemplate += "- [[]]" + string(NewLineChar)
	zettelTemplate += "## Source" + string(NewLineChar)

	zettelFileName, err := createZettelFile(zettelID, zettelTitle, zettelTemplate)
	if err != nil {
		fmt.Println("Error creating the zettel file", err.Error())
	}

	err = exec.Command("code", zettelFileName).Run()
	if err != nil {
		fmt.Println("Error opening vscode with the zettel file", err.Error())
	}
}

func formatZettelTags(zettelTags string) string {
	var formattedZettelTags string
	for _, tag := range strings.Split(zettelTags, ",") {
		formattedZettelTags += "#" + strings.Trim(tag, "") + " "
	}
	return formattedZettelTags
}

func createZettelFile(zettelID string, zettelTitle string, zettelTemplate string) (string, error) {
	zettelFilename := ZettelKastenLocation + zettelID + "-" + strcase.ToKebab(zettelTitle) + ".md"
	zettelFile, err := os.Create(zettelFilename)
	defer zettelFile.Close()
	if err != nil {
		return "", err
	}
	_, err = zettelFile.WriteString(zettelTemplate)
	if err != nil {
		return "", err
	}
	return zettelFilename, nil
}
