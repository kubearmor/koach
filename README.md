# koach

Koach collects observability data from relay server and save it in database. Koach also provided gRPC method to query this data, which currently used by the KubeArmor CLI.

## Quick Start

To start koach service, build the program first using the `make build` command, then go to the koach directory and run this command

```
./koach
```

If you are starting the koach service for the first time, the database file hasn't been created, to create it, start the service using the `migrate` flag

```
./koach --migrate
```
