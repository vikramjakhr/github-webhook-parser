# github-webhook-parser
[![GoDoc](https://godoc.org/github.com/vikramjakhr/github-webhook-parser?status.svg)](https://godoc.org/github.com/vikramjakhr/github-webhook-parser)

#### (Github Event --> github-webhook-parser --> Jenkins)
Github webhook parser cli listens to Github events and allows easy receiving and parsing of GitHub events, It can also 
calls the specified jenkins job on receiving the event.

Features:

 * Listens to Push to repository, Create and Delete a tag or branch
 * Parses the entire REF, REF_TYPE, REPOSITORY and SENDER
 * Forwards REF, REF_TYPE, REPOSITORY and SENDER to the Jenkins job as parameters


# Installation
##### Step 1: Download the [latest release tar](https://github.com/vikramjakhr/github-webhook-parser/releases/latest). Example command below.
```
wget https://github.com/vikramjakhr/github-webhook-parser/releases/download/v1.0.0/githubwebhookparser
```

##### Step 2: Make the binary executable and Move it to /usr/bin
```
chmod +x githubwebhookparser
mv githubwebhookparser /usr/bin
```

##### Step 3: Start the server using cli
```
githubwebhookparser --port <port> 
    --secret <github-secret> 
    --endpoint <webhook-endpoint)  
    --jenkinsUrl <jenkins-url> 
    --jenkinsUsername <jenkins-user> 
    --jenkinsApiToken <jenkins-api-token>
    --jenkinsJobToTrigger <jenkins-job-name-to-trigger>
```

# Build from source
```
go get -u github.com/vikramjakhr/github-webhook-parser.v1
cd $GOPATH/src/github.com/vikramjakhr/github-webhook-parser
go build -o githubwebhookparser
```

# Contributing
Pull requests and suggestions are welcome!

If the changes being proposed or requested are breaking changes, please create an issue for discussion.

Happy coding :-)