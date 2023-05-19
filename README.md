# dapr-ambient-examples

This repository aims to show how to use Dapr Ambient and Dapr building blocks (State management and Pub/Sub) with multiples services into a cluster kubernetes.

## Architecture
Below, you can see a high-level and simple architecture used on this example.

![architecture](/docs/img/architecture.png)

### subscriber

Subscriber just listen by notifications sent from [write-values](#write-values). This component receives all notifications and requests from `dapr` through `dapr-ambient` proxy.

### write-values

Write-values is responsible for save values into `redis` through `dapr-ambient`.

```
curl -X POST http://<host>:<port>?value=90
```

### read-values

Read-values reads all values created by `write-values` and returns an average.

```
curl http://<host>:<port>
```
