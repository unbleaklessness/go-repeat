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
	stageTimeScatter = 20 * 60
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

	result := node{}

	e := decoder.Decode(&result)
	if e != nil {
		return node{}, e
	}

	return result, nil
}

func readNode(filePath string) (node, error) {

	data, e := ioutil.ReadFile(filePath)
	if e != nil {
		return node{}, e
	}

	result, e := unmarshalNode(data)
	if e != nil {
		return node{}, e
	}

	result.filePath = filePath

	directoryPath, _ := filepath.Split(filePath)

	result.directoryPath = filepath.Clean(directoryPath)

	return result, nil
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

		node, e := readNode(path)
		if e != nil {
			return nil
		}

		nodes = append(nodes, node)

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

func writeNewNode(directoryPath string) error {

	n := node{
		ID: makeID(),
	}

	nodePath := filepath.Join(directoryPath, nodeFileName)

	data, e := json.Marshal(n)
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

func isNodesAssociated(node1 node, node2 node) bool {
	for _, a := range node1.Associations {
		if a.ID == node2.ID {
			return true
		}
	}
	return false
}

func isNodesUniassociated(node1 node, node2 node) bool {
	return isNodesAssociated(node1, node2) && isNodesAssociated(node2, node1)
}

func timeIsInFuture(t int64) bool {
	if t > now() {
		fmt.Printf("No ready nodes, next at %s", time.Unix(t, 0).Format(time.RFC1123))
		return true
	}
	return false
}

func removeAssociation(n node, id uint64) (node, bool) {

	index := -1

	for i, a := range n.Associations {
		if a.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return node{}, false
	}

	n.Associations[index] = n.Associations[len(n.Associations)-1]
	n.Associations = n.Associations[:len(n.Associations)-1]

	return n, true
}

func main() {

	rand.Seed(now())

	associate := flag.String("associate", "", "Use with `-with` flag to associate two nodes")
	with := flag.String("with", "", "Use with `-associate` flag to associate two nodes")
	listAssociations := flag.String("list-associations", "", "Show all associations for a node")
	newNode := flag.String("new-node", "", "Create a new node")
	question := flag.Bool("question", false, "Show question")
	answer := flag.Bool("answer", false, "Show answer")
	yes := flag.Bool("yes", false, "Correct answer")
	no := flag.Bool("no", false, "Incorrect answer")
	unassociate := flag.String("unassociate", "", "Use with `-with` flag to unassociate two nodes")
	uni := flag.Bool("uni", false, "Associate / unassociate both ways")

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

		if isNodesAssociated(node1, node2) {
			fmt.Println("First node is already associated with the second")
			return
		}

		a := association{
			ID:    node2.ID,
			Stage: 0,
			Time:  now(),
		}

		node1.Associations = append(node1.Associations, a)

		e := node1.update()
		if e != nil {
			fmt.Println("Could not update first node")
			return
		}

		if *uni {
			if isNodesAssociated(node2, node1) {
				fmt.Println("Second node is already associated with the first")
				return
			}

			a := association{
				ID:    node1.ID,
				Stage: 0,
				Time:  now(),
			}

			node2.Associations = append(node2.Associations, a)

			e := node2.update()
			if e != nil {
				fmt.Println("Could not update second node")
				return
			}
		}

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

			fmt.Printf("%d) %s\n", i+1, nn.directoryPath)
		}

	} else if len(*newNode) > 0 {

		*newNode = filepath.Clean(*newNode)

		e := writeNewNode(*newNode)
		if e != nil {
			fmt.Println("Could not create a node")
			return
		}

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
			fmt.Println("Could not read node files")
			return
		}

		for _, filePath := range filePaths {
			e := open(filePath)
			if e != nil {
				fmt.Println("Could not open node file")
				continue
			}
		}

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
			fmt.Println("Could not find node by ID")
			return
		}

		filePaths, e := nodeFiles(n)
		if e != nil {
			fmt.Println("Could not read node files")
			return
		}

		for _, filePath := range filePaths {
			e := open(filePath)
			if e != nil {
				fmt.Println("Could not open node file")
				continue
			}
		}

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

		n, ok := removeAssociation(node1, node2.ID)
		if !ok {
			fmt.Println("First node is already unassociated with the second")
			return
		}

		e := n.update()
		if e != nil {
			fmt.Println("Could not update first node")
			return
		}

		if *uni {
			n, ok := removeAssociation(node2, node1.ID)
			if !ok {
				fmt.Println("Second node is already unassociated with the first")
				return
			}

			e := n.update()
			if e != nil {
				fmt.Println("Could not update second node")
				return
			}
		}

	} else {
		fmt.Println("Unknown flags")
		return
	}

	fmt.Println("Ok!")
}
