# domi

__NOTICE: This documentation and project is under massive development. Consider this volatile__

## Overview

Enforces policy in CI/CD pipelines using conftest and Open Policy Agent.

Features:

* Integrates with Github, Bitbucket, Gitlab, and other services via webhooks
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