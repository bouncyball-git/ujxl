package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	"golang.org/x/sync/semaphore"
	"gopkg.in/ini.v1"
)

func main() {

	args := os.Args
	if len(args) <= 1 {
		fmt.Printf("Usage: %s [path/]pattern [destination]", filepath.Base(args[0]))
		return
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	cfg, err := ini.Load(exPath + "\\" + "ujxl.ini")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	var cmdLine []string
	var cfgSec *ini.Section
	var Index int
	var doneIndex int
	var mu sync.Mutex

	inPattern := args[1]
	destination := ""
	if len(os.Args) > 2 {
		destination = args[2]
	}
	defaultSec := cfg.Section("default")
	appFilename := defaultSec.Key("app-filename").String()
	maxWorkers, err := defaultSec.Key("max-workers").Int64()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	numThreads := defaultSec.Key("num-threads").String()
	silentMode, err := defaultSec.Key("silent-mode").Bool()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	filesArray, err := getFileList(filepath.Dir(inPattern), filepath.Base(inPattern))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch appFilename {
	case "cjxl.exe":
		cfgSec = cfg.Section("cjxl.exe")
		cmdLine = []string{
			"--distance=" + cfgSec.Key("distance").String(),
			"--effort=" + cfgSec.Key("effort").String(),
			"--num_threads=" + numThreads,
			"-v",
		}
	case "djxl.exe":
		cfgSec = cfg.Section("djxl.exe")
		cmdLine = []string{
			"--jpeg_quality=" + cfgSec.Key("quality").String(),
			"--color_space=" + cfgSec.Key("color-space").String(),
			"--num_threads=" + numThreads,
			"-v",
		}
	default:
		fmt.Println("Wrong executable:", appFilename)
	}

	pool := semaphore.NewWeighted(maxWorkers)
	ctx := context.TODO()
	for index, filename := range filesArray {
		if err := pool.Acquire(ctx, 1); err != nil {
			fmt.Println("Error:", err)
		}
		go func(filename string) {
			cmdFailed := "ok"
			dir, file := filepath.Split(filename)
			if len(destination) > 0 {
				dir = destination
			}
			outFilename := filepath.Clean(dir) + "\\" + file[0:len(file)-len(path.Ext(file))] + cfgSec.Key("out-ext").String()
			argList := append([]string{filename, outFilename}, cmdLine...)
			output, err := exec.Command(appFilename, argList...).CombinedOutput()
			if err != nil {
				cmdFailed = "failed"
			}
			if !(silentMode) {
				fmt.Printf("\rConverting %s, Max %d workers\n%s\n", filename, maxWorkers, output)
			}
			mu.Lock()
			doneIndex++
			mu.Unlock()
			fmt.Printf("\r%d of %d files attempted. Status: %s", doneIndex, Index, cmdFailed)
			pool.Release(1)
		}(filename)
		Index = index + 1
	}
	if err := pool.Acquire(ctx, maxWorkers); err != nil {
		fmt.Println("Error:", err)
	}
}

func getFileList(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
