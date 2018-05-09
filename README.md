# README
TR Logic LLC <contact@itprojects.management>

# Task

Implement in Go a simple REST API with only single method that uploads images. 

Requirements:
- Ability to accept multiple files.
- Ability to accept multipart/form-data requests.
- Ability to accept JSON requests with BASE64 encoded images.
- Ability to upload image by its URL (hosted somewhere in Internet).
- Create thumb square preview (100 x 100 px) for every uploaded image.

There is no restrictions for time to be spent on test task implementing, or tools/libraries to be used for implementation. Any other aspects of test task which are not described in requirements may be implemented on your own decision.

The following will be appreciated:
- Graceful shutdown of application.
- `Dockerfile` and `docker-compose.yml` which allow to boot up application in a single `docker-compose up` command. 
- Unit tests, functional tests, CI integration (Travis CI, Circle CI, etc).
