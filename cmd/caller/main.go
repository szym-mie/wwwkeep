package main

import (
	"log"

	"smbsd.local/wwwkeep"
)

func main() {
	cl, err := wwwkeep.Dial("127.0.0.1:4400")
	if err != nil {
		panic(err)
	}

	log.Println("dial: ok")
	log.Println("dial_remote_id:", cl.RemoteId)
	log.Println("dial_remote_meta:", cl.RemoteMeta)

	var a, b *wwwkeep.Vals
	var aCount, bCount *uint
	var dirs *wwwkeep.Dirs

	log.Println("def: node/a, node/b, node/c")
	_, err = cl.Def("node", []string{"a", "b", "c"}, 0)
	if err != nil {
		panic(err)
	}

	log.Println("add: node/_: a: 0, b: 1")
	_, err = cl.Add("node", map[string]string{"a": "0", "b": "1"})
	if err != nil {
		panic(err)
	}

	log.Println("add: node/_: a: 2, b: 3")
	_, err = cl.Add("node", map[string]string{"a": "2", "b": "3"})
	if err != nil {
		panic(err)
	}

	log.Println("add: node/_: b: 5")
	_, err = cl.Add("node", map[string]string{"b": "5"})
	if err != nil {
		panic(err)
	}

	log.Println("get: node/a")
	a, err = cl.Get("node", "a")
	if err != nil {
		panic(err)
	}

	log.Printf("a: %v\n", *a)

	log.Println("pop: node/b -2")
	bCount, err = cl.Pop("node", "b", 2)
	if err != nil {
		panic(err)
	}

	log.Printf("b count: %d\n", *bCount)

	log.Println("add: node/_: b: 7")
	_, err = cl.Add("node", map[string]string{"b": "7"})
	if err != nil {
		panic(err)
	}

	log.Println("get: node/b")
	b, err = cl.Get("node", "b")
	if err != nil {
		panic(err)
	}

	log.Printf("b: %v\n", *b)

	log.Println("len: node/a")
	aCount, err = cl.Len("node", "a")
	if err != nil {
		panic(err)
	}

	log.Printf("a count: %d\n", *aCount)

	log.Println("dir: *")
	dirs, err = cl.Dir("*")
	if err != nil {
		panic(err)
	}

	log.Printf("dir: %v at %v:\n", dirs.EntType, dirs.EntName)
	for k, v := range dirs.Sizes {
		log.Printf("| %8s - %d vecs\n", k, v.Len)
	}

	log.Println("dir: node")
	dirs, err = cl.Dir("node")
	if err != nil {
		panic(err)
	}

	log.Printf("dir: %v at %v:\n", dirs.EntType, dirs.EntName)
	for k, v := range dirs.Sizes {
		log.Printf("| %8s - len %d cap %d\n", k, v.Len, v.Cap)
	}
}
