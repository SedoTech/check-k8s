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


# How to develop

## Requirements

- Go [https://golang.org/dl/]
- Dep [https://github.com/golang/dep] 
  ```go get -u github.com/golang/dep/cmd/dep```
- Delve [https://github.com/go-delve/delve] (Optional for debugging)
`go get github.com/go-delve/delve/cmd/dlv`
  
## Update/Install dependencies

Before you can compile the package you need to make sure required libraries are available locally

- `dep ensure`                             install the project's dependencies
- `dep ensure -update`                     update the locked versions of all dependencies

This will create `vendor` derectory and put them all there

## Compile executable

- Build final binary  
  `scripts/build.sh 1.8.4` will create ./check-k8s executable binary with version 1.8.4
  
- Run some checks
  `./check-k8s statefulset --kube-context integration --namespace integration-infrastructure availableReplicas trading-queue`


### Debug locally using Intellij

- Compile with debug flags
`go build -gcflags "all=-N -l" -o check-k8s cmd/*.go`
- Start script with debug wrapper 
  `dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./check-k8s -- statefulset --kube-context integration --namespace integration-infrastructure availableReplicas trading-queue` 
- Create Run configuration in Intellij using `Go Remote` template
- Start the created configuration