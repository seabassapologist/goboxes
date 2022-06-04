package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Box struct {
	Name     string `yaml:"name"`
	Cmd      string `yaml:"cmd"`
	Icon     string `yaml:"icon"`
	Interval int    `yaml:"interval"`
}

type Conf struct {
	Boxes     []Box  `yaml:"boxes,flow"`
	Format    string `yaml:"format,omitempty"`
	Delimeter string `yaml:"delimeter,omitempty"`
	Refresh   int    `yaml:"refresh"`
}

type Buffer struct {
	Mut    sync.Mutex
	Status map[int]string
}

// Run command and write to buffer
func runCmd(cmd []string, idx int, buf *Buffer) {
	// exec cmd and arguments
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buf.Mut.Lock()
	buf.Status[idx] = strings.TrimSpace(string(out))
	buf.Mut.Unlock()
}

func main() {

	// It's just gonna be config.yaml don't complain
	cfgFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var cfg Conf
	err = yaml.Unmarshal(cfgFile, &cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	l := len(cfg.Boxes)

	// Initialize the status buffer
	var buf Buffer
	buf.Status = make(map[int]string)

	// Launch boxes as goroutines, and assign index to thread
	// indexes let the goroutine know which value to update in status map
	for i := 0; i < l; i++ {
		// If interval is negative only run cmd once
		if cfg.Boxes[i].Interval < 0 {
			go runCmd(strings.Fields(cfg.Boxes[i].Cmd), i, &buf)
		} else {
			go func(cmd []string, ref int, idx int, buf *Buffer) {
				for {
					runCmd(cmd, idx, buf)
					time.Sleep(time.Duration(ref) * time.Second)
				}
			}(strings.Fields(cfg.Boxes[i].Cmd), cfg.Boxes[i].Interval, i, &buf)
		}
	}

	// Sleep to let goroutines run at least once before printing status
	time.Sleep(time.Duration(cfg.Refresh) * time.Second)

	for {
		// Sleep and refresh output at specified interval
		buf.Mut.Lock()

		var out string
		if cfg.Format == "" {
			out = buf.Status[0]
			for i := 1; i < l; i++ {
				out += cfg.Delimeter + buf.Status[i]
			}
		} else {
			out = cfg.Format
			for i := 0; i < l; i++ {
				out = strings.Replace(out, "{{"+cfg.Boxes[i].Name+"}}", buf.Status[i], -1)
			}
		}
		fmt.Println(out)
		buf.Mut.Unlock()
		time.Sleep(time.Duration(cfg.Refresh) * time.Second)
	}

	fmt.Println("End")
}
