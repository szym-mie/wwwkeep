package wwwkeep

import "fmt"

type Args struct {
	NodeName string
	VecName  string
	DefKeys  []string
	AddItem  map[string]string
	Count    uint
}

func (it Args) String() string {
	nodeName, vecName := it.NodeName, it.VecName
	if nodeName == "" {
		nodeName = "_"
	}

	if vecName == "" {
		vecName = "_"
	}

	return fmt.Sprintf(
		"%s/%s %v %v count %d",
		nodeName, vecName, it.DefKeys, it.AddItem, it.Count)
}
