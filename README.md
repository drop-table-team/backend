## Environment variables

### Backend

| env name                   | description                                                   |
|----------------------------|---------------------------------------------------------------|
| PUBLIC_BACKEND_PORT        | Port to access the backend on the host                        |
| BACKEND_MODULE_CONFIG_PATH | Path to the [module config file](#module-config-file-example) |
| BACKEND_MINIO_BUCKET       | Name of the minio bucket                                      |
| BACKEND_MINIO_ACCESS_KEY   | Minio access key                                              |
| BACKEND_MINIO_SECRET_KEY   | Minio secret key                                              |

### Minio

| env name                  | description                                  |
|---------------------------|----------------------------------------------|
| MINIO_ROOT_USER           | Minio root user                              |
| MINIO_ROOT_PASSWORD       | Minio root password                          |
| PUBLIC_MINIO_CONSOLE_PORT | Port to access the minio console on the host |
| PUBLIC_MINIO_PORT         | Port to access minio on the host             |

### MongoDB

| env name          | description                        |
|-------------------|------------------------------------|
| PUBLIC_MONGO_PORT | Port to access mongodb on the host |

### Qdrant

| env name                | description                                         |
|-------------------------|-----------------------------------------------------|
| PUBLIC_QDRANT_PORT      | Port to access qdrant on the host                   |
| PUBLIC_QDRANT_GRPC_PORT | Port to access the qdrant grpc endpoint on the host |

_See [.env.example](.env.example) for an example env file._

## Module config file example

Service config example:

```json
{
  "modules": [
    "nginx"
  ],
  "input_module_definitions": [
    {
      "name": "nginx",
      "image": "nginx:latest",
      "types": ["application/pdf"]
    },
    {
      "name": "alpine",
      "image": "alpine:latest",
      "types": ["audio/ogg"]
    }
  ]
}
```

_See [example.config.json](example.config.json) for an example env file._

## Endpoints

`GET /modules/input`
```json
[
  {
    "name": <string>,
    "types": [<mime types as string>]
  }
]
```

`POST /modules/input/<module name>` - Send data to an output module
```
Sent data is module/mime specific
```

## Dev setup

All services besides the backend itself can be started by running `docker compose -f docker-compose.dev.yml`
