# Api-example

This repository contains example of simple API service.

You can clone this repo and run a service. Just run following commands:

1. Clone the repo:

```sh
git clone https://github.com/gibiw/api-example && cd /api-example
```

2. Run a database in docker:

```sh
make database_up
```

3. Apply migrations to database:

```sh
make migration_up
```

4. Run the service:

```sh
make run
```
