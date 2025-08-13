## Dependencies

- [Docker](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-22-04)
- [Docker Compose](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-compose-on-ubuntu-22-04)
- Make

```bash
   sudo apt install make
```

- [golang-migrate](https://github.com/golang-migrate/migrate)
  - [Installation guide](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

You should export the `migrate` to the `$PATH` variable

```bash
   export PATH=$PATH:/path/to/migrate
```

- [sqlc](https://sqlc.dev/)

## Run the application locally

### Start the services 

```bash
    docker compose up -d
```

### Run the migrations

```bash
    make migrate-up
```

The server will be available at `http://localhost:8080`.

### TODO 

- [ ] User should be able to create a new category
    - [ ] Category should be associated with the user account
