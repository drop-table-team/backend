Service config example:

```json
{
  "modules": [
    "debian"
  ],
  "module_definitions":  [
    {
      "name": "debian",
      "image": "debian:latest"
    },
    {
      "name": "alpine",
      "image": "alpine:latest"
    }
  ]
}
```

_See [template.config.example](template.config.json) for an example env file._
