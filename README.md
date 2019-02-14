# badger

[![License](https://img.shields.io/github/license/seankhliao/badger.svg?style=for-the-badge)](githib.com/seankhliao/badger)

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
https://img.shields.io/badge/endpoint.svg?url=badger.seankhliao.com/r/$REPO
```

## Design

steps:

1. Query GCP for builds list fitting repo filter
2. Filter skipping working / queued / cancelled, use first result
3. generate JSON
