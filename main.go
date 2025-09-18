package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/cmd"
	"github.com/vogtp/rag/pkg/cfg"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	command := cmd.New()
	var wg sync.WaitGroup
	wg.Go(func() { prof(ctx) })
	cobra.CheckErr(command.ExecuteContext(ctx))
	stop()
	wg.Wait()
}

func fileName(profile, profileType string) string {
	return fmt.Sprintf("ignore_%s_%s.pprof", profile, profileType)
}

func prof(ctx context.Context) {
	cpuprofile := viper.GetString(cfg.PProf)
	if cpuprofile == "" {
		return
	}
	fmt.Printf("STARTING PPROF PROFILING. File=%q\n", cpuprofile)
	cpuFile, err := os.Create(fileName(cpuprofile, "cpu"))
	if err != nil {
		log.Fatal(err)
	}
	memFile, err := os.Create(fileName(cpuprofile, "mem"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		cpuFile.Close()
		memFile.Close()
		fmt.Printf("Wrote pprof files:\n  %s\n  %s\n", cpuFile.Name(), memFile.Name())
	}()
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()
	last := time.Now()
	intervall := time.Minute
	ticker := time.NewTicker(intervall)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			inter := time.Since(last)
			err = pprof.WriteHeapProfile(memFile)
			pprof.StopCPUProfile()
			// err = pprof.StartCPUProfile(f)
			delta := inter - intervall
			fmt.Printf("%s - %s PProf Update: %s - %s = %s err=%v\n", now.Format(time.TimeOnly), time.Now().Format(time.TimeOnly), inter.Truncate(time.Millisecond), intervall, delta.Truncate(time.Millisecond), err)
			if err := cpuFile.Sync(); err != nil {
				fmt.Printf("CPU profile sync error: %s\n", err)
			}
			if err := memFile.Sync(); err != nil {
				fmt.Printf("Memory profile sync error: %s\n", err)
			}
			last = now
			panic("profiles")
		}
	}
}
