# SimpleLib

一些常见的工具库

|名称|包名|功能|备注|
|-|-|-|-|
|lru|simple-lru|lru结构实现|-|
|hashring|simple-hashring|一致性哈希实现|-|
|go-logger|go-logger|日志库封装|封装zap|
|kv-storage|kv-storage|缓存库封装|本地、redis库|

## 测试

```
go test -v -cover ./pkg/sip/
```