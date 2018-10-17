package cli


import (
        "os"
        "log"
        "flag"
)


type CliOptions struct {
        Uri             string
        LogLevel         string
}


func ParseArguments() CliOptions {
        defaults := CliOptions{"tcp://localhost:7419", "info"}

        flag.Usage = help
        flag.StringVar(&defaults.Uri, "b", "localhost:7419", "Network binding")
        flag.StringVar(&defaults.LogLevel, "l", "info", "Logging level (error, warn, info, debug)")

        versionPtr := flag.Bool("v", false, "Show version")
        flag.Parse()

        if *versionPtr {
                os.Exit(0)
        }

        return defaults
}


func help() {
        log.Println("-b [binding]\tHost binding (use :7419 to listen on all interfaces), default: tcp://localhost:7419")
        log.Println("-l [level]\tSet logging level (warn, info, debug, verbose), default: info")
        log.Println("-v\t\tShow version and license information")
        log.Println("-h\t\tThis help screen")
}

