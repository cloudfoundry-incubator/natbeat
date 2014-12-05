# NATBEAT

Natbeat registers a route with the gorouter.  It conforms to the gorouter 
registration protocol, described here: https://github.com/cloudfoundry/gorouter#usage

Natbeat implements the ifrit.Runner interface: http://godoc.org/github.com/tedsuo/ifrit

## Example Usage

The suggested usage is to create a backround heartbeat, which will silently wait 
for NATS to be accesible, and resart itself if it loses the connection.

```go
// define the routes you would like to register
registration := natsbeat.RegistryMessage{
   URIs: []string{"foo.example.com","bar.example.com"},
   Host: "123.1.2.3",
   Port: 80,
}

// create a NATS client
natsClient := diegonats.NewClient()

// create a heartbeater that will restart itself whenever the nats connection is lost
heartbeater := natsbeat.NewBackgroundHeartBeat(natsClient, natsAddresses, natsUsername, natsPassword, logger, registration)

// begin heartbeating in the background
process := ifrit.Invoke(heartbeater)

// do your thing....

// signal the heartbeater with `os.Interrupt` to unregister and exit once you are done
process.Signal(os.Interrupt)
err := <-process.Wait()

// alternatively, signal with `os.Kill` to stop heartbeating immediately without unregistering
process.Signal(os.Kill)
err := <-process.Wait()
```
