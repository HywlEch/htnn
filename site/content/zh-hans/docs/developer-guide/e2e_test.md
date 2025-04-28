---
title: E2E 测试
---

## 积极推送模式

在运行 e2e 测试时，通常希望配置变更能够快速传播到数据平面，以加快测试执行速度。HTNN 为 istiod 提供了一个"积极推送模式"，该模式减少了延迟时间并增加了推送吞吐量。

### 配置参数

积极推送模式修改了以下 istiod 参数：

- `PILOT_DEBOUNCE_AFTER`：从 100ms 减少到 10ms
- `PILOT_DEBOUNCE_MAX`：从 10s 减少到 1s
- `PILOT_ENABLE_EDS_DEBOUNCE`：设置为 false（禁用）
- `PILOT_PUSH_THROTTLE`：从 100 增加到 200

### 使用方法

要启用积极推送模式运行 e2e 测试：

```bash
make e2e-aggressive
```

这将部署带有积极推送配置的 istiod 并运行 e2e 测试。

### 注意事项

- 积极推送模式仅适用于测试环境，不应在生产环境中使用。
- 虽然此模式加快了配置传播速度，但可能会增加控制平面的 CPU 使用率。
