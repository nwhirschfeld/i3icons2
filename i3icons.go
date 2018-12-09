package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nwhirschfeld/i3ipc"
)

func main() {
	// handle command line arguments
	var configFileName = flag.String("c", "/etc/i3icons2.config", "config file")
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
	channel, err := i3ipc.Subscribe(i3ipc.I3WindowEvent)
	EventLoop(channel, ipcsocket, config)
}

// get the subnode of an I3Node by name
func SubNodeByName(node *i3ipc.I3Node, name string) (root i3ipc.I3Node, err error) {
	if strings.Compare(node.Name, name) == 0 {
		return *node, nil
	}
	for _, arm := range node.Nodes {
		res, err := SubNodeByName(&arm, name)
		if err == nil {
			return res, err
		}
	}
	return i3ipc.I3Node{}, errors.New("no such Node")
}

// get the subnodes of an I3Node which doesn't match to a specific name
func SubNodesWithoutName(node *i3ipc.I3Node, name string) (nodes []i3ipc.I3Node, err error) {
	result := make([]i3ipc.I3Node, len(node.Nodes))
	for i, item := range node.Nodes {
		if strings.Compare(item.Name, name) != 0 {
			result[i] = item
		}
	}
	return result, nil
}

// get the ends of the trees
func FlattenNode(node *i3ipc.I3Node) (nodes []i3ipc.I3Node, err error) {
	if len(node.Nodes) == 0 {
		result := make([]i3ipc.I3Node, 1)
		result[0] = *node
		return result, nil
	}

	result := make([]i3ipc.I3Node, 0)
	for _, item := range node.Nodes {
		oldnodes := result
		newnodes, _ := FlattenNode(&item)
		result = make([]i3ipc.I3Node, len(oldnodes)+len(newnodes))
		for i, item := range oldnodes {
			result[i] = item
		}
		for i, item := range newnodes {
			result[i+len(oldnodes)] = item
		}
	}
	return result, nil
}

// main event loop
func EventLoop(events chan i3ipc.Event, ipcsocket *i3ipc.IPCSocket, config map[string]string) {
	for _ = range events {
		tree, _ := ipcsocket.GetTree()
		screens, _ := SubNodesWithoutName(&tree, "__i3")
		for _, screen := range screens {
			wss, _ := SubNodeByName(&screen, "content")
			for _, ws := range wss.Nodes {
				name := ws.Name
				number := strings.Split(name, ":")[0]
				windows, _ := FlattenNode(&ws)
				newname := number + ":"
				windownames := make([]string, len(windows))
				for i, win := range windows {
					winname := win.WindowProperties.Class
					// rename window to config item, if present
					if val, ok := config[winname]; ok {
						winname = val
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
					newname = newname + windowname
				}
				ipcsocket.Command("rename workspace \"" + name + "\" to " + newname)
			}
		}
	}
}
