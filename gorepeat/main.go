package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var stages = []int64{0, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 6765, 10946, 17711, 28657, 46368}

type association struct {
	ID    uint64
	Time  int64
	Stage int
}

type node struct {
	ID            uint64
	Name          string
	Associations  []association
	filePath      string
	directoryPath string
}

type settings struct {
	Root string
	file *os.File
}

const (
	nodeFileName     = "Node.json"
	stageTimeScatter = 15 * 60
	secondsInDay     = 86400
	settingsName     = "GoRepeat.json"
)

func now() int64 {
	return time.Now().Unix()
}

func stageTime(stage int) int64 {
	return now() + stages[stage]*secondsInDay + rand.Int63n(stageTimeScatter)
}

func nextStage(stage int) int {
	stage--
	if stage >= len(stages) {
		return stage + 1
	}
	return stage + 2
}

func previousStage(stage int) int {
	return 0
}

func makeID() uint64 {
	return rand.Uint64()
}

func unmarshalNode(nodeBytes []byte) (node, error) {

	decoder := json.NewDecoder(bytes.NewReader(nodeBytes))
	decoder.DisallowUnknownFields()

	n := node{}

	e := decoder.Decode(&n)
	if e != nil {
		return node{}, e
	}

	return n, nil
}

func readNode(path string) (node, error) {

	data, e := ioutil.ReadFile(path)
	if e != nil {
		return node{}, e
	}

	n, e := unmarshalNode(data)
	if e != nil {
		return node{}, e
	}

	n.filePath = path

	directoryPath, _ := filepath.Split(path)

	n.directoryPath = filepath.Clean(directoryPath)

	return n, nil
}

func isNode(info os.FileInfo) bool {
	return !info.IsDir() || info.Name() == nodeFileName
}

func findNodes() []node {

	nodes := make([]node, 0)

	filepath.Walk(".", func(path string, info os.FileInfo, e error) error {

		if !isNode(info) {
			return nil
		}

		n, e := readNode(path)
		if e != nil {
			return nil
		}

		nodes = append(nodes, n)

		return nil
	})

	return nodes
}

func (n *node) update() error {

	data, e := json.Marshal(*n)
	if e != nil {
		return e
	}

	return ioutil.WriteFile(n.filePath, data, os.ModePerm)
}

func nodeWithPath(nodes []node, path string) (node, bool) {

	for _, n := range nodes {
		if n.directoryPath == path {
			return n, true
		}
	}

	return node{}, false
}

func nodeWithID(nodes []node, id uint64) (node, bool) {

	for _, n := range nodes {
		if n.ID == id {
			return n, true
		}
	}

	return node{}, false
}

func writeNewNode(path string, name string) error {

	if len(name) < 1 {
		name = filepath.Clean(filepath.Base(path))
	}

	n := node{
		ID:   makeID(),
		Name: name,
	}

	nodePath := filepath.Join(path, nodeFileName)

	data, e := json.Marshal(n)
	if e != nil {
		return e
	}

	e = os.MkdirAll(path, os.ModePerm)
	if e != nil {
		return e
	}

	return ioutil.WriteFile(nodePath, data, os.ModePerm)
}

func nodeFiles(n node) ([]string, error) {

	infos, e := ioutil.ReadDir(n.directoryPath)
	if e != nil {
		return []string{}, e
	}

	filePaths := make([]string, 0)

	for _, info := range infos {

		name := info.Name()

		if !info.IsDir() && name != nodeFileName {

			filePath := filepath.Join(n.directoryPath, name)
			filePaths = append(filePaths, filePath)
		}
	}

	return filePaths, nil
}

func associationWithLeastTime(nodes []node) (node, association, int, int) {

	resultNode := node{}
	resultAssociation := association{}

	resultNodeI := -1
	resultAssociationI := -1

	if len(nodes) < 1 {
		return resultNode, resultAssociation, resultNodeI, resultAssociationI
	}

	minimum := int64(math.MaxInt64)

	for i, n := range nodes {
		for j, a := range n.Associations {

			if a.Time < minimum {

				minimum = a.Time

				resultNode = n
				resultAssociation = a

				resultNodeI = i
				resultAssociationI = j
			}
		}
	}

	return resultNode, resultAssociation, resultNodeI, resultAssociationI
}

func timeIsInFuture(t int64) bool {
	if t > now() {
		fmt.Printf("No ready nodes, next at %s", time.Unix(t, 0).Format(time.RFC1123))
		return true
	}
	return false
}

func deleteAssociation(node1 node, node2 node) (node, bool) {

	index := -1

	for i, a := range node1.Associations {
		if a.ID == node2.ID {
			index = i
			break
		}
	}

	if index == -1 {
		return node1, false
	}

	node1.Associations[index] = node1.Associations[len(node1.Associations)-1]
	node1.Associations = node1.Associations[:len(node1.Associations)-1]

	return node1, true
}

func addAssociation(node1 node, node2 node) (node, bool) {

	for _, a := range node1.Associations {
		if a.ID == node2.ID {
			return node1, false
		}
	}

	a := association{
		ID:    node2.ID,
		Stage: 0,
		Time:  stageTime(0),
	}

	node1.Associations = append(node1.Associations, a)

	return node1, true
}

func associationName(node1 node, node2 node) string {
	return node1.Name + " -> " + node2.Name
}

func prepareName(name string) string {
	return strings.Trim(name, " \n\r\t")
}

func main() {

	rand.Seed(time.Now().UnixNano())

	associate := flag.String("associate", "", "Use with `-with` flag to associate two nodes")
	with := flag.String("with", "", "Use with `-associate` flag to associate two nodes")
	listAssociations := flag.String("list-associations", "", "Show all associations for a node")
	newNode := flag.String("new-node", "", "Create a new node")
	name := flag.String("name", "", "Use with `-new-node` to set new node name. Use alone to see node name")
	question := flag.Bool("question", false, "Show question")
	answer := flag.Bool("answer", false, "Show answer")
	yes := flag.Bool("yes", false, "Correct answer")
	no := flag.Bool("no", false, "Incorrect answer")
	unassociate := flag.String("unassociate", "", "Use with `-with` flag to unassociate two nodes")
	uni := flag.Bool("uni", false, "Associate / unassociate both ways")
	rename := flag.String("rename", "", "Rename a node")
	to := flag.String("to", "", "Use with `-rename` to set new node name")
	text := flag.String("text", "", "Create text file in the unit with `-new-node` flag, or in currect directory if alone")
	is := flag.String("is", "", "Use with `-text` flag to set text file content")

	flag.Parse()

	nodes := findNodes()

	if len(*associate) > 0 && len(*with) > 0 {

		*associate = filepath.Clean(*associate)
		*with = filepath.Clean(*with)

		node1, ok := nodeWithPath(nodes, *associate)
		if !ok {
			fmt.Println("Node is not found for `-associate` flag")
			return
		}

		node2, ok := nodeWithPath(nodes, *with)
		if !ok {
			fmt.Println("Node is not found for `-with` flag")
			return
		}

		node1, ok = addAssociation(node1, node2)
		if ok {
			e := node1.update()
			if e != nil {
				fmt.Println("Could not update first node")
			}
		} else {
			fmt.Println("First node is already associated with the second")
		}

		if *uni {
			node2, ok = addAssociation(node2, node1)
			if ok {
				e := node2.update()
				if e != nil {
					fmt.Println("Could not update second node")
				}
			} else {
				fmt.Println("Second node is already associated with the first")
			}
		}

		return

	} else if len(*listAssociations) > 0 {

		*listAssociations = filepath.Clean(*listAssociations)

		n, ok := nodeWithPath(nodes, *listAssociations)
		if !ok {
			fmt.Println("Node is not found for `-list-associations` flag")
			return
		}

		for i, a := range n.Associations {
			nn, ok := nodeWithID(nodes, a.ID)
			if !ok {
				continue
			}

			fmt.Printf("%d) %s || %s\n", i+1, associationName(n, nn), nn.directoryPath)
		}

		return

	} else if len(*newNode) > 0 {

		*newNode = filepath.Clean(*newNode)
		*name = prepareName(*name)

		e := writeNewNode(*newNode, *name)
		if e != nil {
			fmt.Println("Could not create a node")
			return
		}

		if len(*text) > 0 && len(*is) > 0 {
			e := ioutil.WriteFile(filepath.Join(*newNode, *text), []byte(*is), os.ModePerm)
			if e != nil {
				fmt.Println("Could not create text file")
				return
			}
		}

		return

	} else if *question {

		n, a, nI, _ := associationWithLeastTime(nodes)
		if nI == -1 {
			fmt.Println("No nodes are found")
			return
		}

		if timeIsInFuture(a.Time) {
			return
		}

		filePaths, e := nodeFiles(n)
		if e != nil {
			fmt.Println("Could not read question node files")
			return
		}

		if len(filePaths) < 1 {
			fmt.Printf("No files in the question node at \"%s\"\n", n.directoryPath)
			return
		}

		nn, ok := nodeWithID(nodes, a.ID)
		if !ok {
			fmt.Println("Could not find answer node")
			return
		}

		fmt.Println(associationName(n, nn))

		for _, filePath := range filePaths {
			e := open(filePath)
			if e != nil {
				fmt.Println("Could not open question node file")
				continue
			}
		}

		return

	} else if *answer {

		_, a, nI, _ := associationWithLeastTime(nodes)
		if nI == -1 {
			fmt.Println("No nodes are found")
			return
		}

		if timeIsInFuture(a.Time) {
			return
		}

		n, ok := nodeWithID(nodes, a.ID)
		if !ok {
			fmt.Println("Could not find answer node")
			return
		}

		filePaths, e := nodeFiles(n)
		if e != nil {
			fmt.Println("Could not read answer node files")
			return
		}

		if len(filePaths) < 1 {
			fmt.Printf("No files in the answer node at \"%s\"\n", n.directoryPath)
			return
		}

		for _, filePath := range filePaths {
			e := open(filePath)
			if e != nil {
				fmt.Println("Could not open answer node file")
				continue
			}
		}

		return

	} else if *yes {

		n, a, nI, aI := associationWithLeastTime(nodes)
		if nI == -1 {
			fmt.Println("No nodes are found")
			return
		}

		a.Stage = nextStage(a.Stage)
		a.Time = stageTime(a.Stage)

		n.Associations[aI] = a

		e := n.update()
		if e != nil {
			fmt.Println("Could not update node")
			return
		}

		return

	} else if *no {

		n, a, nI, aI := associationWithLeastTime(nodes)
		if nI == -1 {
			fmt.Println("No nodes are found")
			return
		}

		a.Stage = previousStage(a.Stage)
		a.Time = stageTime(a.Stage)

		n.Associations[aI] = a

		e := n.update()
		if e != nil {
			fmt.Println("Could not update node")
			return
		}

		return

	} else if len(*unassociate) > 0 && len(*with) > 0 {

		*unassociate = filepath.Clean(*unassociate)
		*with = filepath.Clean(*with)

		node1, ok := nodeWithPath(nodes, *unassociate)
		if !ok {
			fmt.Println("Node is not found for `-unassociate` flag")
			return
		}

		node2, ok := nodeWithPath(nodes, *with)
		if !ok {
			fmt.Println("Node is not found for `-with` flag")
			return
		}

		n, ok := deleteAssociation(node1, node2)
		if ok {
			e := n.update()
			if e != nil {
				fmt.Println("Could not update first node")
			}
		} else {
			fmt.Println("First node is already unassociated with the second")
		}

		if *uni {
			n, ok := deleteAssociation(node2, node1)
			if ok {
				e := n.update()
				if e != nil {
					fmt.Println("Could not update second node")
				}
			} else {
				fmt.Println("Second node is already unassociated with the first")
			}
		}

		return

	} else if len(*name) > 0 {

		*name = filepath.Clean(*name)

		n, ok := nodeWithPath(nodes, *name)
		if !ok {
			fmt.Println("Could not find node")
			return
		}

		fmt.Println(n.Name)

		return

	} else if len(*rename) > 0 && len(*to) > 0 {

		*rename = filepath.Clean(*rename)

		n, ok := nodeWithPath(nodes, *rename)
		if !ok {
			fmt.Println("Node is not found")
			return
		}

		n.Name = prepareName(*to)

		e := n.update()
		if e != nil {
			fmt.Println("Could not update node")
			return
		}

		return

	} else if len(*text) > 0 && len(*is) > 0 {

		*text = filepath.Clean(*text)

		e := ioutil.WriteFile(*text, []byte(*is), os.ModePerm)
		if e != nil {
			fmt.Println("Could not create text file")
			return
		}

		return

	} else {
		fmt.Println("Unknown flags")
		return
	}
}
