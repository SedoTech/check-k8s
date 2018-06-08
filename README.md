# check-k8s

The checks are mostly orientated by the offical nagios [guidelines](http://nagios-plugins.org/doc/guidelines.html)

## Default Flags

| flag | short | description |
| -- | -- | -- |
| critical | c | defines a [threshold](#threshold) for a critical return status |
| warning | w | defines a [threshold](#threshold) for a critical return status |
| namespace | n | the namespace of the kubernetes resource |
| result | r | the service status if the check fails |
| duration | d | if you want to specify a duration where something must happen. <https://golang.org/pkg/time/#ParseDuration> |
| string | s | if you want to check if a property has a concrete value |

## Thresholds

<https://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT>

| Range definition | Generate an alert if x... |
| -- | -- |
| 10 | < 0 or > 10, (outside the range of {0 .. 10}) |
| 10: | < 10, (outside {10 .. ∞}) |
| ~:10 | > 10, (outside the range of {-∞ .. 10}) |
| 10:20 | < 10 or > 20, (outside the range of {10 .. 20}) |
| @10:20 | ≥ 10 and ≤ 20, (inside the range of {10 .. 20} |