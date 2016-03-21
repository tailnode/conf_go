package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const MAX_DEPTH = 10
const SEPRATOR = "  "
const SUFFIX = ".conf"

type configTree struct {
	node  map[string]*configTree
	leave map[string]string
}

func newConfigTree() *configTree {
	node := make(map[string]*configTree)
	leave := make(map[string]string)
	return &configTree{node, leave}
}

func (c *configTree) String() string {
	var str string
	for k, v := range c.node {
		str += "[" + k + "]\n"
		str += v.String()
	}
	for k, v := range c.leave {
		str += k + ":" + v + "\n"
	}
	return str
}

func main() {
	//getAllFile("/home/ming", 0)
	test()
}

func test() {
	//conf := parse("/home/ming/work")
	conf := parse("testcase")
	log.Println(conf)
}

func Load(path string) {
}

func parse(path string) *configTree {
	log.Println("start parse", path)
	info, err := os.Stat(path)
	var config *configTree
	if err == nil && info.IsDir() {
		// path is a directory
		config = parseDir(path, MAX_DEPTH)
	} else if info, err = os.Stat(path + SUFFIX); err == nil && !info.IsDir() {
		// path.conf is a configure file
		config = parseFile(path + SUFFIX)
	}
	log.Println("finish parse", path)
	return config
}
func parseDir(path string, depth int) *configTree {
	if depth > MAX_DEPTH {
		return nil
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	config := newConfigTree()
	log.Println("parse dir", path)
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
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer f.Close()
	log.Println("parse file", path)
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
			currentNode = config.node[str]
			continue
		}
		pairs := strings.SplitN(str, ":", 2)
		// key:value
		if len(pairs) == 2 && len(pairs[0]) > 0 && len(pairs[1]) > 0 {
			currentNode.leave[strings.TrimSpace(pairs[0])] = strings.TrimSpace(pairs[1])
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return config
}
