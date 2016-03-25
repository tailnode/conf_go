// parse configure from a directory or a file
// configuer file must end with .conf subffix
// call conf.Load() once to parse and then use conf.GetConf() to get config
package conf

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const MAX_DEPTH = 10
const SUFFIX = ".conf"

var confRoot Node

type Node map[string]interface{}

func copyNode(src Node) Node {
	desc := newNode()
	for k, v := range src {
		switch v.(type) {
		case Node:
			desc[k] = copyNode(v.(Node))
		case string:
			desc[k] = v
		}
	}
	return desc
}

func newNode() Node {
	return make(map[string]interface{})
}

func Load(path string) {
	log.Println("start load config", path)
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		// path is a directory
		confRoot = parseDir(path, 0)
	} else if info, err = os.Stat(path + SUFFIX); err == nil && !info.IsDir() {
		// path.conf is a configure file
		confRoot = parseFile(path + SUFFIX)
	}
	log.Println("finish load config", path)
}

func GetConf(path string) (value interface{}) {
	value = ""
	splitPath := strings.Split(strings.Trim(path, "/"), "/")

	tmpConfig := confRoot
	for i := 0; i < len(splitPath); i++ {
		if node, ok := tmpConfig[splitPath[i]]; ok {
			switch v := node.(type) {
			case Node:
				if i == len(splitPath)-1 {
					// make a Node copy to avoid the caller change inner data by GetConf()
					return copyNode(v)
				}
				tmpConfig = v
			case string:
				if i == len(splitPath)-1 {
					value = v
				}
				return
			default:
			}
		}
	}
	return
}

func parseDir(path string, depth int) Node {
	log.Printf("parse dir  [%v] current depth[%v]\n", path, depth)
	if depth > MAX_DEPTH {
		return newNode()
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	node := newNode()
	for _, file := range files {
		if file.IsDir() {
			node[file.Name()] = parseDir(path+"/"+file.Name(), depth+1)
		} else {
			// only parse files have .conf suffix
			if strings.HasSuffix(file.Name(), SUFFIX) {
				base := strings.TrimRight(file.Name(), SUFFIX)
				node[base] = parseFile(path + "/" + file.Name())
			}
		}
	}
	return node
}

func parseFile(path string) Node {
	log.Printf("parse file [%v]\n", path)
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return newNode()
	}
	defer f.Close()
	scanner := bufio.NewScanner(bufio.NewReader(f))
	node := newNode()
	currentNode := node
	for scanner.Scan() {
		str := scanner.Text()
		commentIndex := strings.Index(str, "#")
		if commentIndex > -1 {
			str = str[:commentIndex]
		}
		str = strings.TrimSpace(str)
		if len(str) == 0 {
			continue
		}
		// [group]
		// key1:value1
		// key2: value2
		//  key3 : value3
		if len(str) > 2 && str[0] == '[' && str[len(str)-1] == ']' {
			str = str[1 : len(str)-1]
			node[str] = newNode()
			currentNode = (node[str]).(Node)
			continue
		}
		pairs := strings.SplitN(str, ":", 2)
		// key:value
		if len(pairs) == 2 && len(pairs[0]) > 0 && len(pairs[1]) > 0 {
			currentNode[strings.TrimSpace(pairs[0])] = strings.TrimSpace(pairs[1])
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return node
}
