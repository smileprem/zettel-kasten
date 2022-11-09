package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

const TimeFormat = "20060102150405"
const NewLineChar = '\n'
const DefaultZettelConfigFile = "/Users/smileprem/.zettel.json"
const DefaultZettelKastenLocation = "./"

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
		fmt.Println("Error reading the title. Error: ", err.Error())
	}

	fmt.Print("tags (comma de-limited): ")
	zettelTags, err := reader.ReadString(NewLineChar)
	if err != nil {
		fmt.Println("Error reading the tags. Error: ", err.Error())
	}

	zettelTemplate += "# " + zettelID + " " + zettelTitle
	zettelTemplate += formatZettelTags(zettelTags) + string(NewLineChar)
	zettelTemplate += "## Links" + string(NewLineChar)
	zettelTemplate += "- [[]]" + string(NewLineChar)
	zettelTemplate += "## Source" + string(NewLineChar)

	zettelFileName, err := createZettelFile(zettelID, zettelTitle, zettelTemplate)
	if err != nil {
		fmt.Println("Error creating the zettel file. Error: ", err.Error())
	}

	err = exec.Command("code", zettelFileName).Run()
	if err != nil {
		fmt.Println("Error opening vscode with the zettel file. Error: ", err.Error())
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
	zettelFilename := getZettelFileNameFromConfig(zettelID, zettelTitle)
	zettelFile, err := os.Create(zettelFilename)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	_, err = zettelFile.WriteString(zettelTemplate)
	if err != nil {
		return "", err
	}
	return zettelFilename, nil
}

func getZettelFileNameFromConfig(zettelID string, zettelTitle string) string {
	zettelKastenPath := DefaultZettelKastenLocation
	jsonFile, err := os.Open(DefaultZettelConfigFile)
	if err != nil {
		fmt.Println("Unable to open the zettel config file. Creating zettel in current directory. Error: ", err.Error())
	} else {
		byteValue, _ := io.ReadAll(jsonFile)
		var result map[string]interface{}
		err = json.Unmarshal(byteValue, &result)
		if err != nil {
			fmt.Println("Zettel config file is not a proper json. Creating zettel in current directory. Error: ", err.Error())
		}
		zettelKastenPath = result["zettelkasten"].(string)
	}
	_ = jsonFile.Close()

	return zettelKastenPath + zettelID + "-" + strcase.ToKebab(zettelTitle) + ".md"
}
