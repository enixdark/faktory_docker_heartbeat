package main

import (
        "fmt"
        "log"
        "os"
        "net/url"
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
                        return err
                }
                handleEvent(Shutdown, mgr)
                return nil
        })

        if err != nil {
                panic(err)
                os.Exit(0)
        }

        return nil
}



func main() {
        var err error
        // var client *faktory.Client
        log.SetFlags(0)
        opts := cli.ParseArguments()
        util.InitLogger(opts.LogLevel)

        _, ok := os.LookupEnv("FAKTORY_URL")

        if ok {
                _, err = faktory.Open()
                if err != nil {
                        fmt.Errorf("error")
                        os.Exit(0)
                }
        } else {
                
                util.Debugf("Options: %v", opts) 
                uri, err := url.Parse(opts.Uri)

                if err != nil {
                        fmt.Errorf("error")
                        os.Exit(0)
                }
                
                server := faktory.Server{}

                server.Network = uri.Scheme
                server.Address = fmt.Sprintf("%s:%s", uri.Hostname(), uri.Port())
                if uri.User != nil {
                        server.Password, _ = uri.User.Password()
                }
                _, err = server.Open()  
                if err != nil {
                        fmt.Errorf("error")
                        os.Exit(0)
                }

        }

        mgr := worker.NewManager()
        var quit bool
	mgr.On(worker.Shutdown, func() {
		quit = true
        })

        go func() {
		for {
			if quit {
                                fmt.Println("Success")
				os.Exit(0)
                        }
                        heartbeat(mgr)
			time.Sleep(1 * time.Second)
		}
	}()
        
	mgr.Run()

}