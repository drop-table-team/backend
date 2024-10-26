### Environment variables

| env name                         | example                      |
|----------------------------------|------------------------------|
| MODULE_CONFIG_PATH               | /run/config/modules.json     |
| MONGO_INITDB_ROOT_USERNAME       | admin                        |
| MONGO_INITDB_USERNAME            | root                         |
| MONGO_INITDB_PASSWORD_FILE       | /run/secrets/mongo_password  |
| MONGO_PORT                       | 27017                        |
| QDRANT_PORT                      | 6333                         |
| QDRANT_GRPC_PORT                 | 6334                         |
| QDRANT_INITDB_ROOT_PASSWORD_FILE | /run/secrets/qdrant_password |
| BACKEND_PORT                     | 8080                         |

_See [template.env](template.env) for an example env file._

#### Setup MongoDB Password

> file: secrets/mongo_password.txt

#### Setup Qdrant Password

> file: secrets/qdrant_password.txt

### Module config file example

Service config example:

```json
{
  "modules": [
    "nginx"
  ],
  "module_definitions":  [
    {
      "name": "nginx",
      "image": "nginx:latest"
    },
    {
      "name": "alpine",
      "image": "alpine:latest"
    }
  ]
}
```

_See [example.config.example](example.config.json) for an example env file._
