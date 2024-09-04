package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome Github Tracker, use github-activity <username>")

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.SplitN(input, " ", 2)
		command := parts[0]

		switch command {
		case "github-activity":
			name := parts[1]

			url := fmt.Sprintf("https://api.github.com/users/%s/events", name)
			response, err := http.Get(url)

			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			defer response.Body.Close()

			var jsonData interface{}
			err = json.Unmarshal(responseData, &jsonData)

			if err != nil {
				fmt.Println("Error Parsing JSON:", err)
				return
			}

			prettyJSON, err := json.MarshalIndent(jsonData, "", " ")
			if err != nil {
				fmt.Println("Error formating JSON:", err)
				return
			}

			res := string(prettyJSON)
			var events []map[string]interface{}
			err = json.Unmarshal([]byte(res), &events)
			if err != nil {
				fmt.Println("Error parsing JSON events:", err)
				return
			}

			var s []string

			for _, event := range events {
				if eventStar, validEventStar := event["type"].(string); validEventStar && eventStar == "WatchEvent" {
					repo, _ := event["repo"].(map[string]interface{})
					repoName, _ := repo["name"].(string)
					fmt.Println("Found a Starred:")
					res := fmt.Sprintf("- Starred %s\n", repoName)
					s = append(s, res)
				}

				if eventType, ok := event["type"].(string); ok && eventType == "PushEvent" {
					// Handle PushEvent
					fmt.Println("Found a PushEvent:")

					// Safely access commits and repo name
					commits, _ := event["payload"].(map[string]interface{})["commits"].([]interface{})
					repo, _ := event["repo"].(map[string]interface{})

					if len(commits) > 0 {
						repoName, _ := repo["name"].(string)
						if repoName != "" {
							res := fmt.Sprintf("- Pushed %d commits to %s\n", len(commits), repoName)
							s = append(s, res)
						} else {
							fmt.Println("Error: Unable to retrieve repository name")
						}
					} else {
						fmt.Println("Error: Unable to retrieve commits or repository information")
					}
				}
			}

			for _, action := range s {
				fmt.Println(action)
			}
		}

	}

}
