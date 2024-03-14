# gonja-cosmosSDK-writer

now, you can edit config.toml & app.toml using yaml "key - value" data only as you need!


## config.yaml
```yaml
configToml:
  - instrumentation:
    - namespace: "hellWOrld"

  - goodHelo: "second"
  - third: "!FWEOISDJ"

appToml:
  - state_sync:
    - snapshot_interval: "12"
```
## out/app.toml
```yaml
....
[state-sync]
snapshot-interval = "12"
snapshot-keep-recent = "2"
....
```

## out/config.toml
```yaml
....
[instrumentation]
....
namespace = "hellWOrld"
```# dynamic-cosmosSDK-writer
# dynamic-cosmosSDK-writer
