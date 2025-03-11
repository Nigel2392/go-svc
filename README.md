# go-svc

Simpler in- and deinstallation of your windows services.

Custom services must adhere to the following interface:

```go
type Service interface {
	Name() string
	Config() *mgr.Config
	Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32)
	AcceptedCommands() svc.Accepted
	EventLog() services.EventLog
}
```

Example included in [the examples directory](./example/service).