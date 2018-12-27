#### Installation
1. DB
- `docker-compose up postgres`
- `make port-migrations-up`

2. 

#### Future improvements and some comments
- use go-proto-validator
- ports can be sent to port service with grpc streaming
- client service's rest interface can be implemented with grpc + grpc-gateway consistency wise, but now for clarity its just straight http approach
- error wrapping to preserve call stack
- predefined/precompiled sql statements
- context logger with request id
- there should be better error handling, like sql error should be converted to corresponding grpc/http one

#### Tests
This is whole another story, I'm up for spinning all services and testing with actual grpc calls (possibly mocking some dependencies through interface), and proper db cleanup (not just execute truncate).
But in order to save some time its done in a dead simple manner.
