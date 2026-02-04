# Personal Accounting

## Pre-requisites

1. yarn (+node)
1. go
1. docker

## Usage

1. make sure docker is running
1. create .env file
   1. create file from example `cp .env.example .env`
   1. update entries as needed
1. install dependencies
   1. `yarn tidy`
1. serve all projects
   1. `yarn start`

## Commands

_(for full list, see [package.json])_

| Use                    | Commands     |
| ---------------------- | ------------ |
| start dev servers      | `yarn start` |
| run all tests          | `yarn test`  |
| lint everything        | `yarn lint`  |
| update go dependencies | `yarn tidy`  |
