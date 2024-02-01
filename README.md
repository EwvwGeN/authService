# AuthService

An authorization service that works through a grpc connection.

# Table of Contents

- [Proto](#Proto)
- [Build](#Build)
- [Configuration](#Configuration)
- [Usages](#Usages)
    - [Register user](#Registration)
    - [Login user](#Loginning)
    - [Check user is admin](#Checking-for-admin)

## Proto

First you need to create a protofile.

> If you use standard queries, you can skip this step.

After the change, you need to rebuild the go files. To do this, use the command: <br>
`make proto`.

## Build

If you have changed the protofile and/or are launching the service for the first time, you must first build it.

>in the future, a step-by-step assembly will be written and this will be handled by the docker himself

For building only service: <br>
`make buildServer`

## Configuration

The config is located in the configs directory. 

It should specify the:
- logging level;
- the port on which the service will run;
- validation parameters;
- params for mongodb.

If u will use docker-compose u shoulde change th .env file.

## Startup

If you are not launching the service from docker-compose, you must specify the path to the config at startup (when using it):

`./serverMain -config=./configs/example.yaml`

or make a clean launch (then environment variables will be used):

`./serverMain`

If you are want launch service by docker-compose you can use Makefile command:<br>
`make` or `make runservice`

or run docker-compose by yourself.

## Usages

### Registration

Request for registration:
```json
{
    "email": "",
    "password": ""
}
```

Response for registration:
```json
{
    "user_id": "",
}
```

### Loginning

Request for Loginning:
```json
{
    "email": "",
    "password": "",
    "app_id": ""
}
```

Response for Loginning:
```json
{
    "token": "",
}
```


### Checking for admin

Request for Loginning:
```json
{
    "user_id": "",
}
```

Response for Loginning:
```json
{
    "is_admin": false,
}
```