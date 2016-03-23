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

var config *configTree

type configTree struct {
	node map[string]interface{}
}

func newConfigTree() *configTree {
	node := make(map[string]interface{})
	return &configTree{node}
}

func Load(path string) {
	config = parse(path)
}

func GetConf(path string) (value string) {
	splitPath := strings.Split(strings.Trim(path, "/"), "/")

	tmpConfig := config
	for i := 0; i < len(splitPath); i++ {
		if node, ok := tmpConfig.node[splitPath[i]]; ok {
			switch v := node.(type) {
			case *configTree:
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

func parse(path string) *configTree {
	log.Println("start parse", path)
	info, err := os.Stat(path)
	var config *configTree
	if err == nil && info.IsDir() {
		// path is a directory
		config = parseDir(path, 0)
	} else if info, err = os.Stat(path + SUFFIX); err == nil && !info.IsDir() {
		// path.conf is a configure file
		config = parseFile(path + SUFFIX)
	}
	log.Println("finish parse", path)
	return config
}

func parseDir(path string, depth int) *configTree {
	log.Printf("parse dir  [%v] current depth[%v]\n", path, depth)
	if depth > MAX_DEPTH {
		return nil
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	config := newConfigTree()
	for _, file := range files {
		if file.IsDir() {
			config.node[file.Name()] = parseDir(path+"/"+file.Name(), depth+1)
		} else {
			// only parse files have .conf suffix
			if strings.HasSuffix(file.Name(), SUFFIX) {
				base := strings.TrimRight(file.Name(), SUFFIX)
				config.node[base] = parseFile(path + "/" + file.Name())
			}
		}
	}
	return config
}

func parseFile(path string) *configTree {
	log.Printf("parse file [%v]\n", path)
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer f.Close()
	scanner := bufio.NewScanner(bufio.NewReader(f))
	config := newConfigTree()
	currentNode := config
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
			config.node[str] = newConfigTree()
			currentNode = (config.node[str]).(*configTree)
			continue
		}
		pairs := strings.SplitN(str, ":", 2)
		// key:value
		if len(pairs) == 2 && len(pairs[0]) > 0 && len(pairs[1]) > 0 {
			currentNode.node[strings.TrimSpace(pairs[0])] = strings.TrimSpace(pairs[1])
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return config
}
