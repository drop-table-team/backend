1. Setup Environment File

```env
# file: .env
#
# MongoDB Environment
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_USERNAME=root
MONGO_INITDB_PASSWORD_FILE=/run/secrets/mongo_password
MONGO_PORT=27017
QDRANT_PORT=6333
QDRANT_GRPC_PORT=6334
QDRANT_INITDB_ROOT_PASSWORD_FILE=/run/secrets/qdrant_password
BACKEND_PORT=8080
```

2. Setup MongoDB Password

> file: secrets/mongo_password.txt

3. Setup Qdrant Password

> file: secrets/qdrant_password.txt
