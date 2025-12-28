package wwwrpc

import "fmt"

type Args struct {
	NodeName string
	VecName  string
	DefKeys  []string
	AddItem  map[string]string
	Count    uint
}

func (self Args) String() string {
	nodeName, vecName := self.NodeName, self.VecName
	if nodeName == "" {
		nodeName = "_"
	}

	if vecName == "" {
		vecName = "_"
	}

	return fmt.Sprintf(
		"%s/%s %v %v count %d",
		nodeName, vecName, self.DefKeys, self.AddItem, self.Count)
}
