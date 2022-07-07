To start the server please
1. Run postgres docker container with `make postgres` command
2. Run `make db` command to create database
3. Run `migrate up` command to run migration
4. Run `make sqlc` command to create db queries
5. Run `make server` command to start the server on port 8080 (this can be configured in app.env file)

To run the tests
1. Create db mocks with `make mock` command
2. Start tests with `make test` command