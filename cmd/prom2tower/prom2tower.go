package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/pja237/prom2tower/internal/configuration"
	"github.com/pja237/prom2tower/internal/pipe"
	"github.com/pja237/prom2tower/internal/version"
)

const defaultConfFile = "conf.yaml"

func getConfFile() *string {
	var cf string

	flag.StringVar(&cf, "c", defaultConfFile, "Path to configuration file.")
	flag.Parse()

	return &cf
}

func main() {
	var (
		logFile io.Writer
		wg      sync.WaitGroup
	)

	cf := getConfFile()
	conf, err := configuration.GetConfig(*cf)
	if err != nil {
		log.Fatalf("getconfig: %s\n", err)
	}
	//log.Printf("Config %v\n", conf)

	if lf, err := configuration.WantString(conf.Globals["logFile"]); err == nil {
		logFile, err = os.OpenFile(*lf, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("can not open configured log file. Exiting.\n")
		}
	} else {
		log.Printf("WARNING: logFile err %s, silently proceeding with logfile=stderr\n", err)
		logFile = os.Stderr
	}

	log := log.New(logFile, "prom2tower:", log.Lshortfile|log.Ldate|log.Lmicroseconds)

	version.DumpVersion(log)
	log.Println("======================== goglu start ===========================================")

	// loop through the glue, spin up pipes
	log.Printf("Got %d pipes configured.\n", len(conf.Glue))
	for i, v := range conf.Glue {
		log.Printf("Pipe %d : %q @ %p\n", i, v, &v)
		pw, err := pipe.NewPipeWorker(i, conf.Globals, v, log)
		if err != nil {
			log.Printf("PipeWorker %s init FAILED: %s\n", v.Name, err)
		} else {
			wg.Add(1)
			go pw.SpinUp(&wg)
		}
	}

	log.Println("Waiting for handlers to register...")
	wg.Wait()
	log.Println("Ready to ListenAndServe()")
	log.Printf("listen %#v\n", conf.Globals["listenAddr"])

	if s, err := configuration.WantString(conf.Globals["listenAddr"]); err != nil {
		log.Printf("ERROR: listenAddr %s", err)
	} else {
		log.Printf("ListenAndServe() starting...\n")
		err := http.ListenAndServe(*s, nil)
		if err != nil {
			log.Fatalf("ERROR: %s\n", err)
		}
	}

	log.Println("======================== goglu end =============================================")
}
