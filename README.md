
# SpeedyAuth
![example workflow](https://github.com/kwesidev/authserver/actions/workflows/go.yml/badge.svg)

An API-only standalone authentication server that streamlines the user authentication process.

## Features
- Issuing Access Tokens and Refresh Tokens
- Password Recovery
- User Registration
- Two factor Authentication (EMAIL,TOTP)

## Build Dependencies
- Go >= 1.19 
- PostgreSQL >= 9.x 


## Build
```
make buildserver

```

## Run
If you docker engine installed run the following command :
```
docker compose -f docker-compose-dev.yml up

```
