// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"fmt"
	"net/http"
	"strings"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	// TODO replace or remove
	BaseURL  string   `envconfig:"PLUGIN_URL" envDefault:"https://ntfy.sh"`
	Topic    string   `envconfig:"PLUGIN_TOPIC", required:"true"`
	Username string   `envconfig:"PLUGIN_USERNAME"`
	Password string   `envconfig:"PLUGIN_PASSWORD"`
	Token    string   `envconfig:"PLUGIN_TOKEN"`
	Title    string   `envconfig:"PLUGIN_TITLE"`
	Priority string   `envconfig:"PLUGIN_PRIORITY" default:"default"`
	Tags     []string `envconfig:"PLUGIN_TAGS" envSeparator:","`
	//DefaultTags []string `envconfig:"PLUGIN_DEFAULT_TAGS" default:"drone"`
	Message string `envconfig:"PLUGIN_MESSAGE"`
}

func buildAppMessage(args *Args) {
	args.Title = fmt.Sprintf("Build #%d %s", args.Build.Number, args.Build.Status)

	if strings.Contains(args.Commit.Ref, "refs/tags/") {
		args.Tags = append(args.Tags, args.Tag.Name)
		args.Message = "Tag " + args.Tag.Name + " created"
	} else {
		args.Tags = append(args.Tags, args.Repo.Name+"/"+args.Commit.Branch)
		args.Message = "[" + args.Commit.Rev[0:8] + "] " + args.Commit.Message
	}
}

func addResultToTags(args *Args) {
	if args.Build.Status == "success" {
		args.Tags = append(args.Tags, "white_check_mark")
	} else if args.Build.Status == "failure" {
		args.Tags = append(args.Tags, "x")
	} else {
		args.Tags = append(args.Tags, "grey_question")
	}
}

func getActions(args *Args) string {
	var buildLink = "view, Build, " + args.Build.Link

	if strings.Contains(args.Commit.Ref, "refs/tags/") {
		return buildLink
	}

	return buildLink + "; view, Changes, " + args.Commit.Link
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			return ""
		}
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

// Exec executes the plugin.
func Exec(args *Args) (string, error) {
	buildAppMessage(args)
	addResultToTags(args)

	//args.Tags = append(args.DefaultTags, args.Tags...)

	req, _ := http.NewRequest("POST",
		args.BaseURL+"/"+args.Topic,
		strings.NewReader(args.Message))

	if args.Token != "" {
		req.Header.Add("Authorization", "Bearer "+args.Token)
	} else {
		req.SetBasicAuth(args.Username, args.Password)
	}

	req.Header.Set("Title", args.Title)
	req.Header.Set("Priority", args.Priority)
	req.Header.Set("Tags", strings.Join(args.Tags, ","))
	req.Header.Set("Actions", getActions(args))

	//fmt.Printf("--> %s\n\n", formatRequest(req))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("error trying to notify the result. Error: %+v", err)
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("error from server. HTTP status: %d. Error: %e", res.StatusCode, err)
	}

	return "[SUCCESS] Notification sent", nil
}
