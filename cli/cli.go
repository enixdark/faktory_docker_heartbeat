package cli


import (
        "os"
        "log"
        "flag"
)


type CliOptions struct {
        LogLevel         string
}


func ParseArguments() CliOptions {
        defaults := CliOptions{"info"}

        flag.Usage = help
        flag.StringVar(&defaults.LogLevel, "l", "info", "Logging level (error, warn, info, debug)")

        versionPtr := flag.Bool("v", false, "Show version")
        flag.Parse()

        if *versionPtr {
                os.Exit(0)
        }

        return defaults
}


func help() {
        log.Println("-l [level]\tSet logging level (warn, info, debug, verbose), default: info")
        log.Println("-v\t\tShow version and license information")
        log.Println("-h\t\tThis help screen")
}

