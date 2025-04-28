---
title: E2E Testing
---

## Aggressive Push Mode

When running e2e tests, it's often desirable to have configuration changes propagate quickly to the data plane to speed up test execution. HTNN provides an "aggressive push mode" for istiod that reduces the debounce times and increases the push throughput.

### Configuration Parameters

The aggressive push mode modifies the following istiod parameters:

- `PILOT_DEBOUNCE_AFTER`: Reduced from 100ms to 10ms
- `PILOT_DEBOUNCE_MAX`: Reduced from 10s to 1s
- `PILOT_ENABLE_EDS_DEBOUNCE`: Set to false (disabled)
- `PILOT_PUSH_THROTTLE`: Increased from 100 to 200

### Usage

To run e2e tests with aggressive push mode enabled:

```bash
make e2e-aggressive
```

This will deploy istiod with the aggressive push configuration and run the e2e tests.

### Notes

- Aggressive push mode is intended for testing environments only and should not be used in production.
- While this mode speeds up configuration propagation, it may increase CPU usage on the control plane.
