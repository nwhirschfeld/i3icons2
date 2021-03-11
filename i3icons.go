package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/mdirkse/i3ipc"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// handle command line arguments
	var configFileName = flag.String("c", usr.HomeDir+"/.config/i3icons/i3icons2.config", "config file")
	var verbose = flag.Bool("v", false, "verbose")
	flag.Parse()

	// Open our configFile
	configFile, err := os.Open(*configFileName)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}
	defer configFile.Close()

	// read our config File and write to hash map
	byteValue, _ := ioutil.ReadAll(configFile)
	configLines := strings.Split(string(byteValue), "\n")
	config := make(map[string]string)
	for _, ci := range configLines {
		p := strings.Split(string(ci), "=")
		if len(p) == 2 {
			config[p[0]] = p[1]
		}
	}

	// open I3IPC socket and subscribe to window events
	ipcsocket, _ := i3ipc.GetIPCSocket()
	i3ipc.StartEventListener()
	channel, err := i3ipc.Subscribe(i3ipc.I3WindowEvent)
	EventLoop(channel, ipcsocket, config, *verbose)

}

// EventLoop - main event loop
func EventLoop(events chan i3ipc.Event, ipcsocket *i3ipc.IPCSocket, config map[string]string, verbose bool) {
	for range events {
		tree, _ := ipcsocket.GetTree()
		wss := tree.Workspaces()
		for _, ws := range wss {
			name := ws.Name
			number := strings.Split(name, " ")[0]
			windows := ws.Leaves()
			// empty workspace - leave it
			if len(windows) == 0 {
				continue
			}
			newname := number
			windownames := make([]string, len(windows))
			for i, win := range windows {
				winname := strings.ToLower(win.Window_Properties.Class)
				if verbose {
					fmt.Println(winname)
				}
				// rename window to config item, if present
				if val, ok := config[winname]; ok {
					winname = val
				} else if len(winname) > 7 {
					winname = winname[:4] + ".." + winname[len(winname)-3:]
				}
				// check if workspace name already contains window title
				choose := true
				for _, n := range windownames {
					if strings.Compare(n, winname) == 0 {
						choose = false
					}
				}
				if choose {
					windownames[i] = winname
				}
			}
			// rename workspace
			for _, windowname := range windownames {
				if len(windowname) > 0 {
					newname = fmt.Sprintf("%s %s", newname, windowname)
				}
			}
			ipcsocket.Command("rename workspace \"" + name + "\" to " + newname)
		}
	}
}
