# Agente
 [![License](https://img.shields.io/github/license/akdilsiz/agente)](https://opensource.org/licenses/Apache-2.0)
 \
 \
Distributed simple and robust release management and system monitoring system.

This project currently maintained by **@TransferChain**

***This project on going work.*

### Road map
 - [ ] Core system
 - [ ] First worker agent
 - [ ] Jenkins vs CI tool extensions
 - [ ] Management dashboard
 - [ ] First master agent
 - [ ] All relevant third-party system integrations (version control, CI, database, queuing etc.)

## Requirements
 - Go > 1.9
 - Redis or RabbitMQ

## Development
```shell script
go mod vendor
go run ./cmd -dev
```

## Contribution
I would like to accept any contributions to make Agente better and feature rich. So feel free to contribute your features(i.e. more 3rd-party(version control, CI, database, queuing etc.) tools), improvements and fixes.\
[See detail](docs/contribution.md)
## LICENSE

Copyright 2019 Abdulkadir DILSIZ - TransferChain

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
