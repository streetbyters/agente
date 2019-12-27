 <h1 align="center">Agente</h1>
 <p align="center">
  <img height="150" src="assets/agente.png"/>
 </p>
 <p align="center">
   <a href="https://travis-ci.org/akdilsiz/agente">
    <img src="https://travis-ci.org/akdilsiz/agente.svg?branch=master"/>
   </a>
   <a href="https://github.com/akdilsiz/agente/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/akdilsiz/agente"/>
   </a>
   <a href="https://codecov.io/gh/akdilsiz/agente">
     <img src="https://codecov.io/gh/akdilsiz/agente/branch/master/graph/badge.svg" />
   </a>
   <a href="https://goreportcard.com/report/github.com/akdilsiz/agente">
    <img src="https://goreportcard.com/badge/github.com/akdilsiz/agente"/>
   </a>
 </p>

Distributed simple and robust release management and monitoring system.

This project currently maintained by **[@TransferChain](https://github.com/TransferChain)**

***This project on going work.*

### Road map
 - [ ] Core system
 - [ ] First worker agent
 - [ ] Management dashboard
 - [ ] Jenkins vs CI tool extensions
 - [ ] Management dashboard
 - [ ] First master agent
 - [ ] All relevant third-party system integrations (version control, CI, database, queuing etc.)

## Requirements
 - Go > 1.11
 - Redis or RabbitMQ
 - PostgreSQL

## Development
```shell script
git clone -b develop https://github.com/akdilsiz/agente

go mod vendor

# Development Mode
go run ./cmd -mode dev -migrate -reset
go run ./cmd -mode dev

# Test Mode
go run ./cmd -mode test -migrate -reset
go run ./cmd -mode test
```

## Build
We will release firstly Agente for Linux environment.\
[See detail](docs/build.md)

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
