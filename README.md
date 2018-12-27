#### Installation
``

2. Test
`make test`

#### Running
1. DB
- `docker-compose up postgres`
- `make port-migrations-up`

2. Port service
- `docker-compose up port-service`

3. Client API service
- `docker-compose up clientapi`
- NOTE since there is no configuration service relies to have "ports.json" file in project root

#### Future improvements and some comments
- there is ZERO configuration right now, since I didn't have time to implement it properly; i'm not exactly aware on what's the best tool for that, but i'd use viper
- use go-proto-validator
- ports can be sent to port service with grpc streaming
- client service's rest interface can be implemented with grpc + grpc-gateway for consistency wise, but now for clarity its just straight http approach
- error wrapping to preserve call stack
- predefined/precompiled sql statements
- context logger with request id and better logging in general
- there should be better error handling, like sql error should be converted to corresponding grpc/http one
- i should install modules into container and not link them from outside, or at least vendor them into repository
- i'm pretty sure that I messed in some way http part - it was a while for me since last time I've implemented something like that without any framework or external library
- proper shutdown
- json reader writer should be moved to its own abstration layer and be better structured, this was something i did last, so i've spent only time i had =(

#### Tests
I'm down to appoach spinning all services and testing with actual grpc calls (possibly mocking some dependencies through interface), and proper db cleanup (currently its a bit fragile).
There are no tests for client api service, but I imagine them being end-to-end, spinning port service and starting consuming some test data.
Also tests should be better structured and more isolated, instead I tried to test everything in one test.

#### Disclaimer
Huge mistake on my behalf was to use go modules, which I've never used before, and my editor was really messed up because of that.
