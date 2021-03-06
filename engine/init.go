package engine

// File: initialization of command line flags.
// Author: Arne Vansteenkiste

import (
	"code.google.com/p/mx3/cuda"
	"code.google.com/p/mx3/prof"
	"flag"
	"log"
	"runtime"
)

var (
	Flag_version  = flag.Bool("v", true, "print version")
	Flag_silent   = flag.Bool("s", false, "Don't generate any log info")
	Flag_od       = flag.String("o", "", "set output directory")
	Flag_force    = flag.Bool("f", false, "force start, clean existing output directory")
	Flag_maxprocs = flag.Int("threads", 0, "maximum number of CPU threads, 0=auto")
)

func Init() {
	flag.Parse()

	log.SetPrefix("")
	log.SetFlags(0)
	if *Flag_silent {
		log.SetOutput(devnul{})
	}

	if *Flag_version {
		log.Print("Mumax Cubed 0.1α ", runtime.GOOS, "_", runtime.GOARCH, " ", runtime.Version(), "(", runtime.Compiler, ")", "\n")
	}

	if *Flag_od != "" {
		SetOD(*Flag_od, *Flag_force)
	}

	if *Flag_maxprocs == 0 {
		*Flag_maxprocs = runtime.NumCPU()
	}
	procs := runtime.GOMAXPROCS(*Flag_maxprocs) // sets it
	log.Println("gomaxprocs:", procs)

	prof.Init(OD)
	cuda.Init()
	cuda.LockThread()
}

type devnul struct{}

func (d devnul) Write(b []byte) (int, error) {
	return len(b), nil
}

func Close() {
	log.Println("shutting down")
	drainOutput()
	if Table != nil {
		Table.(*dataTable).flush()
	}
	prof.Cleanup()
}
