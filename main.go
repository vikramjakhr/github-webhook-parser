package main

import (
	"fmt"

	"net/http"

	"flag"
	"net/url"
	"os"

	"github.com/yosida95/golang-jenkins"
	"gopkg.in/go-playground/webhooks.v5/github"
)

var port = flag.String("port", "6686", "port on which application will run")
var secret = flag.String("secret", "", "github webhook secret")
var endpoint = flag.String("endpoint", "/github-webhook/", "github webhook endpoint to configure")
var jenkinsUrl = flag.String("jenkinsUrl", "", "jenkins base URL")
var jenkinsUsername = flag.String("jenkinsUsername", "", "jenkins username whom REST api token is associated")
var jenkinsApiToken = flag.String("jenkinsApiToken", "", "jenkins REST API token")
var jenkinsJobToTrigger = flag.String("jenkinsJobToTrigger", "", "jenkins job to trigger")

const usage = `Github webhook parser.

Usage:

  githubwebhookparser [commands|flags]

The commands & flags are:

  help              prints help

  --port              port on which application will run (default: 6686)
  --secret            github webhook secret
  --endpoint          github webhook endpoint to configure (default: /github-webhook/)
  --jenkinsUrl          jenkins base URL
  --jenkinsUsername          jenkins username whom REST api token is associated
  --jenkinsApiToken          jenkins REST API token
  --jenkinsJobToTrigger          jenkins job to trigger

Examples:

  # prints help:
  githubwebhookparser help

  # sample usage
  githubwebhookparser --port 6686 --secret MyGitHubSuperSecretSecrect --endpoint /github-webhook/
  --jenkinsUrl http://jenkins.icmarkets.com/ --jenkinsUsername admin --jenkinsApiToken 11ed0505476f3313215e0b3b12bf179e9e
  --jenkinsJobToTrigger webhook-test
`

// Jenkins parameter Definitions
var actions = []gojenkins.Action{{
	ParameterDefinitions: []gojenkins.ParameterDefinition{
		{
			Name: "REF",
		},
		{
			Name: "REF_TYPE",
		},
		{
			Name: "REPOSITORY",
		},
		{
			Name: "SENDER",
		},
	},
}}

func usageExit(rc int) {
	fmt.Println(usage)
	os.Exit(rc)
}

func validate() {
	if *port == "" || *endpoint == "" || *secret == "" {
		usageExit(1)
	}

	if *jenkinsJobToTrigger != "" && (*jenkinsUrl == "" || *jenkinsUsername == "" || *jenkinsApiToken == "") {
		usageExit(1)
	}
}

func main() {
	flag.Usage = func() { usageExit(0) }
	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		switch args[0] {
		case "help":
			usageExit(0)
			return
		}
	}

	validate()

	// Github authentication
	hook, _ := github.New(github.Options.Secret(*secret))

	// Webhook handler
	http.HandleFunc(*endpoint, func(w http.ResponseWriter, r *http.Request) {
		event, err := hook.Parse(r, github.CreateEvent, github.DeleteEvent, github.PushEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				fmt.Println("E! Specified event not found")
				fmt.Println(err)
			}
		}

		switch event.(type) {

		case github.CreatePayload:
			if *jenkinsJobToTrigger != "" {
				payload := event.(github.CreatePayload)
				job := gojenkins.Job{
					Name:    *jenkinsJobToTrigger,
					Actions: actions,
				}
				checkError(jenkins().Build(job, withParams(payload.Ref, payload.RefType, payload.Repository.Name, payload.Sender.Login)))
			}

		case github.DeletePayload:
			if *jenkinsJobToTrigger != "" {
				payload := event.(github.DeletePayload)
				job := gojenkins.Job{
					Name:    *jenkinsJobToTrigger,
					Actions: actions,
				}
				checkError(jenkins().Build(job, withParams(payload.Ref, payload.RefType, payload.Repository.Name, payload.Sender.Login)))
			}

		case github.PushPayload:
			if *jenkinsJobToTrigger != "" {
				payload := event.(github.PushPayload)
				job := gojenkins.Job{
					Name:    *jenkinsJobToTrigger,
					Actions: actions,
				}
				checkError(jenkins().Build(job, withParams(payload.Ref, "", payload.Repository.Name, payload.Sender.Login)))
			}

		}
	})
	fmt.Printf("I! Starting the server on port %s\n", *port)
	http.ListenAndServe(":"+*port, nil)
}

func jenkins() *gojenkins.Jenkins {
	auth := &gojenkins.Auth{
		Username: *jenkinsUsername,
		ApiToken: *jenkinsApiToken,
	}

	return gojenkins.NewJenkins(auth, *jenkinsUrl)
}

func withParams(ref, refType, repository, sender string) url.Values {
	params := url.Values{}
	params.Set("REF", ref)
	params.Set("REF_TYPE", refType)
	params.Set("REPOSITORY", repository)
	params.Set("SENDER", sender)
	return params
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("E! Error occured while executing job %s\n", *jenkinsJobToTrigger)
		fmt.Println(err)
	} else {
		fmt.Printf("I! Successfully executed job %s\n", *jenkinsJobToTrigger)
	}
}
