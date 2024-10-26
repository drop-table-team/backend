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
