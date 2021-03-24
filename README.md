
![](img/domi128x128.png)

# domi

[![Go Report Card](https://goreportcard.com/badge/github.com/devops-kung-fu/domi)](https://goreportcard.com/report/github.com/devops-kung-fu/domi) [![codecov](https://codecov.io/gh/devops-kung-fu/domi/branch/main/graph/badge.svg?token=R1TFX89WFQ)](https://codecov.io/gh/devops-kung-fu/domi) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/devops-kung-fu/domi)


__NOTICE: This documentation and project is under active development. Consider this pre-alpha.__

## Overview

Enforces policy in CI/CD pipelines using conftest and Open Policy Agent.

Features:

* Integrates with Github, Bitbucket (To Do), Gitlab (To Do), and other services (To Do) via webhooks
* Supports Policy-as-Code written in rego.
* May be run against specific repositories without integration to determine policy

## Building from Source

After cloning the repository, utilize make to build the application.  Simply typing _make_ at the command line will display help for the defined targets in the Makefile.

### Quickstart Targets

| Target | Function                                                |
| ------ | ------------------------------------------------------- |
| help   | Displays all targets and what they do                   |
| docker | Builds the docker file                                  |
| run    | Starts domi locally on the port defined in _config.env_ |

## What's with the name Domi?

The name domi is a shout out to one of the greatest enforcers of all time in the NHL - [Tie Domi](https://en.wikipedia.org/wiki/Tie_Domi). He's best known as a fighter, and holds the all-time record of the most fighting majors (333).

## Contributing
Pull requests are welcome. For more information see [CONTRIBUTING.md](contributing.md)

If the changes being proposed or requested are breaking changes, please create an issue for discussion.

## License
Distributed under the MPL V2 License, please see the [LICENSE](LICENSE]) file for more details.

## Credits
Logo courtesy of [FlatIcon](https://flaticon.com)