package wwwkeep

import "fmt"

type Vals []string
type node map[string]*List
type Keep map[string]node

type DirSize struct {
	Len, Cap uint
}

type Dirs struct {
	EntName string
	EntType string
	Sizes   map[string]DirSize
}

func (it Keep) getNode(nodeName string) (node, error) {
	node := it[nodeName]
	if node == nil {
		return nil, fmt.Errorf(
			"keep_get_node: no %s node in %p",
			nodeName, it)
	}

	return node, nil
}

func (it Keep) def(args *Args) (uint, error) {
	nn := args.NodeName
	if it[nn] != nil {
		return 0, fmt.Errorf(
			"keep_def: node %s already exists in %p",
			nn, it)
	}

	node := make(node)
	for _, k := range args.DefKeys {
		l := new(List)
		l.Grow(args.Count)
		node[k] = l
	}

	it[nn] = node
	return uint(len(args.DefKeys)), nil
}

func (it Keep) add(args *Args) (uint, error) {
	node, err := it.getNode(args.NodeName)
	if err != nil {
		return 0, err
	}

	count := 0
	for k, v := range args.AddItem {
		if node[k] == nil {
			return 0, fmt.Errorf(
				"keep_add: no %s/%s vec in %p",
				args.NodeName, k, it)
		}
		node[k].Append(v)
		count++
	}

	return uint(count), nil
}

func (it Keep) get(args *Args) (Vals, error) {
	node, err := it.getNode(args.NodeName)
	if err != nil {
		return nil, err
	}

	k := args.VecName
	vs := node[k]
	if vs == nil {
		return nil, fmt.Errorf(
			"keep_get: no %s/%s vec in %p",
			args.NodeName, k, it)
	}

	return vs.Slice(), nil
}

func (it Keep) pop(args *Args) (uint, error) {
	node, err := it.getNode(args.NodeName)
	if err != nil {
		return 0, err
	}

	k := args.VecName
	vs := node[k]
	if vs == nil {
		return 0, fmt.Errorf(
			"keep_poplen: no %s/%s vec in %p",
			args.NodeName, k, it)
	}

	for range args.Count {
		vs.Pop()
	}

	return uint(vs.Len), nil
}

func (it Keep) len(args *Args) (uint, error) {
	return it.pop(&Args{args.NodeName, args.VecName, nil, nil, 0})
}

func (it Keep) dir(args *Args) (Dirs, error) {
	en := args.NodeName
	et := ""
	sizes := make(map[string]DirSize, len(it))

	if en == "*" {
		for k, node := range it {
			ln := uint(len(node))
			sizes[k] = DirSize{ln, ln}
		}

		et = "nodes"
	} else {
		node := it[en]
		if node == nil {
			return Dirs{}, fmt.Errorf(
				"keep_dir: no %s node in %p",
				en, it)
		}

		for k, vs := range node {
			sizes[k] = DirSize{vs.Len, vs.Cap}
		}

		et = "vecs"
	}

	return Dirs{en, et, sizes}, nil
}

func (it Keep) opt(args *Args) (uint, error) {
	node, err := it.getNode(args.NodeName)
	if err != nil {
		return 0, err
	}

	count := uint(0)
	for _, v := range node {
		count += v.Shrink()
	}

	return count, nil
}
