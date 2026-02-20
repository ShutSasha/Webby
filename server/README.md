# Webby Backend

This project contains multiple backend services.
You can run all services together or start only one.

## Requirements

* .NET 8 SDK
* Task runner installed (`task` command)

## Run All Services

From the `server` folder, run:

```bash
task run_all
```

## Run a Single Service

Example for the Auth service:

```bash
task run_auth_service
```

Other services follow the same pattern:

```bash
task run_{service_name}
```

## Stop Services

Press:

```
Ctrl + C
```

## Task List

To see all available commands:

```bash
task --list
```
