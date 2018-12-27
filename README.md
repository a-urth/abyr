#### Installation
1. DB
- `docker-compose up postgres`
- `make port-migrations-up`

2. Port service
- `docker-compose up port-service`

3. Client API service
- `docker-compose up client-api`

#### Future improvements and some comments
- there is ZERO configuration right now, since I didn't have time to implement it properly; i'm not exactly aware on what's the best tool for that, but i'd use viper
- use go-proto-validator
- ports can be sent to port service with grpc streaming
- client service's rest interface can be implemented with grpc + grpc-gateway for consistency wise, but now for clarity its just straight http approach
- error wrapping to preserve call stack
- predefined/precompiled sql statements
- context logger with request id and better logging in general
- there should be better error handling, like sql error should be converted to corresponding grpc/http one

#### Tests
This is whole another story, I'm up for spinning all services and testing with actual grpc calls (possibly mocking some dependencies through interface), and proper db cleanup (not just execute truncate).
But in order to save some time its done in a dead simple manner.
