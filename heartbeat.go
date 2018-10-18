package main

import (
        "fmt"
        "log"
        "os"
        "time"
        util "github.com/enixdark/faktory_docker_heartbeat/util"
        cli "github.com/enixdark/faktory_docker_heartbeat/cli"
        faktory "github.com/contribsys/faktory/client"
        worker "github.com/contribsys/faktory_worker_go"
)


type eventType int

const (
	Quiet    eventType = 1
	Shutdown eventType = 2
)

var (
	RandomProcessWid = ""
)


func handleEvent(sig eventType, mgr *worker.Manager) {
	switch sig {
	case Shutdown:
		go func() {
			mgr.Terminate()
		}()
	case Quiet:
		go func() {
			mgr.Quiet()
		}()
	}
}


func with(mgr *worker.Manager, fn func(fky *faktory.Client) error) error {
	conn, err := mgr.Pool.Get()
	if err != nil {
		return err
	}
	pc := conn.(*worker.PoolConn)
	f, ok := pc.Closeable.(*faktory.Client)
	if !ok {
		return fmt.Errorf("Connection is not a Faktory client instance: %+v", conn)
	}
	err = fn(f)
	if err != nil {
		pc.MarkUnusable()
	}
	conn.Close()
	return err
}


func Beat(c* faktory.Client) (string, error) {
        val, err := c.Generic("BEAT " + fmt.Sprintf(`{"wid":"%s"}`, c.Options.Wid))
	if val == "OK" {
		return "", nil
	}
	return val, err
}


func heartbeat(mgr *worker.Manager) error {
	
        err := with(mgr, func(c *faktory.Client) error {
                
                _, err := Beat(c)

                if err != nil {
                        panic(err)
                        util.Debugf("error")
                        return err
                }
                return nil
        })

        if err != nil {
                panic(err)
                util.Debugf("error")
                os.Exit(0)
        }
        handleEvent(Shutdown, mgr)
        return nil
}



func main() {

        log.SetFlags(0)
        opts := cli.ParseArguments()
        util.InitLogger(opts.LogLevel)

        mgr := worker.NewManager()
        var quit bool
	mgr.On(worker.Shutdown, func() {
		quit = true
        })

        go func() {
		for {
                        heartbeat(mgr)
			time.Sleep(15 * time.Second)
		}
	}()
	mgr.Run()

}