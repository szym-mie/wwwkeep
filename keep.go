package wwwrpc

import "fmt"

type Vals []string
type node map[string]Vals
type Keep map[string]node

type Dirs struct {
	EntName string
	EntType string
	Counts  map[string]uint
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
		node[k] = make(Vals, args.Count)
	}

	it[nn] = node
	return uint(len(args.DefKeys)), nil
}

func (it Keep) add(args *Args) (uint, error) {
	nn := args.NodeName
	node, err := it.getNode(nn)
	if err != nil {
		return 0, err
	}

	for k, v := range args.AddItem {
		if node[k] == nil {
			return 0, fmt.Errorf(
				"keep_add: no %s/%s vec in %p",
				nn, k, it)
		}
		node[k] = append(node[k], v)
	}

	return uint(len(args.AddItem)), nil
}

func (it Keep) get(args *Args) (Vals, error) {
	nn := args.NodeName
	node, err := it.getNode(nn)
	if err != nil {
		return nil, err
	}

	k := args.VecName
	vs := node[k]
	if vs == nil {
		return nil, fmt.Errorf(
			"keep_get: no %s/%s vec in %p",
			nn, k, it)
	}

	return vs, nil
}

func (it Keep) pop(args *Args) (uint, error) {
	nn := args.NodeName
	node, err := it.getNode(nn)
	if err != nil {
		return 0, err
	}

	k := args.VecName
	vs := node[k]
	if vs == nil {
		return 0, fmt.Errorf(
			"keep_poplen: no %s/%s vec in %p",
			nn, k, it)
	}

	if args.Count > 0 {
		vs = vs[args.Count:]
		node[k] = vs
	}

	return uint(len(vs)), nil
}

func (it Keep) len(args *Args) (uint, error) {
	return it.pop(&Args{args.NodeName, args.VecName, nil, nil, 0})
}

func (it Keep) dir(args *Args) (Dirs, error) {
	en := args.NodeName
	et := ""
	counts := make(map[string]uint)

	if en == "*" {
		for k, node := range it {
			counts[k] = uint(len(node))
		}

		et = "nodes"
	} else {
		node := it[en]
		if node == nil {
			return Dirs{"", "", nil}, fmt.Errorf(
				"keep_dir: no %s node in %p",
				en, it)
		}

		for k, vs := range node {
			counts[k] = uint(len(vs))
		}

		et = "vecs"
	}

	return Dirs{en, et, counts}, nil
}
