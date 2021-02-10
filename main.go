package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func run() {
	rpt := GenerateReport()
	rpt.SendLog()
	rpt.PublishMetric()
}

func main() {
	daemon := false
	runEvery := int64(120)

	flag.BoolVar(&daemon, "daemon", false, "run as a daemon")
	flag.Int64Var(&runEvery, "run-every", 120, "when in daemon mode, run every X seconds")
	flag.Parse()

	run()

	if daemon {
		ticker := time.NewTicker(time.Duration(runEvery) * time.Second)
		quit := make(chan os.Signal)

		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		signal.Notify(quit, os.Interrupt, syscall.SIGHUP)
		signal.Notify(quit, os.Interrupt, syscall.SIGQUIT)
		signal.Notify(quit, os.Interrupt, syscall.SIGKILL)

		for {
			select {
			case <-ticker.C:
				run()
			case <-quit:
				return
			}
		}
	}
}
