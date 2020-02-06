 <h1 align="center">Agente</h1>
 <p align="center">
  <img height="150" src="assets/agente.png"/>
 </p>
 <p align="center">
   <a href="https://circleci.com/gh/streetbyters/agente">
    <img src="https://circleci.com/gh/streetbyters/agente.svg?style=svg"/>
   </a>
   <a href="https://github.com/streetbyters/agente/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/streetbyters/agente"/>
   </a>
   <a href="https://codecov.io/gh/streetbyters/agente">
     <img src="https://codecov.io/gh/streetbyters/agente/branch/master/graph/badge.svg" />
   </a>
   <a href="https://goreportcard.com/report/github.com/streetbyters/agente">
    <img src="https://goreportcard.com/badge/github.com/streetbyters/agente"/>
   </a>
 </p>

Distributed simple and robust release management and monitoring system.

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

## Docker Environment
For PostgreSQL
```shell script
docker run --name agente_PostgreSQL -e POSTGRES_PASSWORD=123456 -e POSTGRES_USER=agente -p 5432:5432 -d postgres

docker exec agente_PostgreSQL psql --username=agente -c 'create database agente_dev;'
```
For RabbitMQ
```shell script
docker run --hostname my-rabbit --name agente_RabbitMQ -e RABBITMQ_DEFAULT_USER=local -e RABBITMQ_DEFAULT_PASS=local -p 5672:5672 -d rabbitmq:3-management
```

## Development
```shell script
git clone -b develop https://github.com/streetbyters/agente

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

Copyright 2019 StreetByters Community

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
