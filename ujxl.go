package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"golang.org/x/sync/semaphore"
	"gopkg.in/ini.v1"
)

const version = "Jpeg XL universal batch utility, version 0.9.1"

func main() {

	fmt.Println(version)

	args := os.Args
	if len(args) <= 1 {
		fmt.Printf("Usage: %s \"[path]filename|wildcard.ext\" [destination path]\n", filepath.Base(args[0]))
		return
	} else if len(args) > 3 {
		fmt.Printf("Error: Excessive arguments specified.\n")
		fmt.Printf("Usage: %s \"[path]filename|wildcard.ext\" [destination path]\n", filepath.Base(args[0]))
		fmt.Printf(" Note: Did you forget to put first argument in quotation marks?\n")
		return
	}

	ex, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	exPath := filepath.Dir(ex)
	cfg, err := ini.Load(exPath + string(os.PathSeparator) + "ujxl.ini")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	var cfgSec *ini.Section
	var cmdLine []string
	var outExt string
	var Index int
	var doneIndex int
	var mu sync.Mutex

	inPattern := string(args[1])
	destination := ""
	if len(os.Args) > 2 {
		destination = string(args[2])
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
	case "cjxl":
		cfgSec = cfg.Section("cjxl")
		outExt = ".jxl"
		loslessJpeg := cfgSec.Key("lossless").String()
		if loslessJpeg == "0" {
			cmdLine = []string{
				"--lossless_jpeg=" + loslessJpeg,
				"--effort=" + cfgSec.Key("effort").String(),
				"--distance=" + cfgSec.Key("distance").String(),
				"--num_threads=" + numThreads,
				"-v",
			}
		} else {
			cmdLine = []string{
				"--lossless_jpeg=" + loslessJpeg,
				"--effort=" + cfgSec.Key("effort").String(),
				"--num_threads=" + numThreads,
				"-v",
			}
		}
	case "djxl":
		cfgSec = cfg.Section("djxl")
		outExt = cfgSec.Key("out-ext").String()
		cmdLine = []string{
			"--color_space=" + cfgSec.Key("color-space").String(),
			"--jpeg_quality=" + cfgSec.Key("quality").String(),
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
			outFilename := filepath.Clean(dir) + string(os.PathSeparator) + file[:len(file)-len(filepath.Ext(file))] + outExt
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
	fmt.Println()
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
