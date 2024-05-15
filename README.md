# Accounts
## _service for auth operations and admin management_

## Features

- Admin login (endpoint)
- Admin create (endpoint)
- Admin list (endpoint)
- Admin middleware (external lib)
- Token service for admins auth operatins (external lib)



## Installation

```sh
go get github.com/proj-go-5/accounts
```

## Local development

Clone project :

```sh
git clone git@github.com:proj-go-5/accounts.git
```

##### Prepare env envoronment:
create .env file (see .env.example)

Up db server:
```sh
docker-compose -f internal/db/docker-compose.yml up
```

Create db tables:
```sh
make migration_up
```

Up server:
```sh
make start
```
