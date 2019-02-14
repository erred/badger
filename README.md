# badger

[![https://img.shields.io/badge/endpoint.svg?url=https://badger.seankhliao.com/r/github_seankhliao_badger]()](https://console.cloud.google.com/cloud-build/builds?project=com-seankhliao&query=source.repo_source.repo_name%20%3D%20%22github_seankhliao_badger%22)
[![License](https://img.shields.io/github/license/seankhliao/badger.svg?style=for-the-badge)](LICENSE)

badges for GCP Cloud Build

## usage

```
badger [-p 8080] [-pr com-seankhliao]
  -p    port to listen on
  -pr   GCP project to query
```

Accepts the following urls:

- `/r/$REPO`: name of Source Repo, github repos are the form `github_$USER_$REPO`
- `/success`: returns a success
- `/failure`: returns a failure
- `/status_unkown`: returns a status_unkown

To be used with the [Shields.io Endpoint API](https://shields.io/endpoint)

```
https://img.shields.io/badge/endpoint.svg?url=https://badger.seankhliao.com/r/$REPO

ex:
https://img.shields.io/badge/endpoint.svg?url=https://badger.seankhliao.com/r/github_seankhliao_badger
```

## Design

steps:

1. Query GCP for builds list fitting repo filter
2. Filter skipping working / queued / cancelled, use first result
3. generate JSON

## TODO

- manage caching
- doc url
- markdown generator
