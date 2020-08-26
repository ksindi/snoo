# SNOO

[![Build Status](https://travis-ci.org/ksindi/owlet.svg?branch=main)](https://travis-ci.org/ksindi/owlet)

Unofficial client for Owlet.

## Legal

This code is in no way affiliated with, authorized, maintained, sponsored or endorsed by Owlet Baby Care or any of its affiliates or subsidiaries. This is an independent and unofficial API. Use at your own risk.

Enable GitHub vulnerability alerts for all repositories.

## Usage

```bash
# default usage: enable for all repositories with automated security fixes
owlet -org=myorg -alerts=true -fixes=true

# enable for single respository
owlet -org=myorg -alerts=true -fixes=true -repo=myrepo

# enable for all repositories but disable automated security fixes
owlet -org=myorg -alerts=true -fixes=false

# enable for all repositories but do nothing with automated security fixes
owlet -org=myorg -alerts=true


owlet -help

  -alerts
      Boolean to enable/disable alerts (GITHUB_VUL_ALERTS)
  -dry
      Dry run (GITHUB_VUL_DRY)
  -fixes
      [Optional] Boolean to enable/disable automated (GITHUB_VUL_FIXES)
  -org string
      GitHub org (GITHUB_VUL_ORG)
  -repo string
      [Optional] Specify a repository
  -token string
      GitHub API token (GITHUB_VUL_TOKEN)
```

## Requirements

[Generate a personal access token](https://github.com/settings/tokens) with `repo` and `read:org` permissions.

## Installation

### Releases

Download the binary for your platform from the [releases](https://github.com/ksindi/owlet/releases) page.

### Docker

```sh
docker pull ksindi/owlet
docker run -it -e $GITHUB_VUL_TOKEN ksindi/owlet -alert=true -org=ksindi -dry=true
```

### Go

```sh
go get -u github.com/ksindi/owlet
```

## License

GitHub Vul is provided under the [Apache License v2.0](./LICENSE).
