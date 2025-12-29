package wwwkeep

type List struct {
	head *ListNode
	tail *ListNode
	free *ListNode
	Len  uint
	Cap  uint
}

type ListNode struct {
	Val  *string
	next *ListNode
}

func (it *List) Grow(cap uint) {
	for range cap {
		it.Cap++
		node := new(ListNode)
		node.next = it.free
		it.free = node
	}
}

func (it *List) Append(val string) {
	it.Len++
	if it.free == nil {
		it.Cap++
		node := new(ListNode)
		node.Val = &val
		if it.head == nil {
			it.head = node
		} else {
			it.tail.next = node

		}

		it.tail = node
	} else {
		it.free.Val = &val
		it.tail.next = it.free
		it.tail = it.free
		it.free = it.free.next
	}
}

func (it *List) Pop() *string {
	if it.head == nil {
		return nil
	}

	it.Len--
	node := it.head
	val := node.Val
	it.head = node.next
	node.next = it.free
	it.free = node
	node.Val = nil
	return val
}

func (it *List) Shrink() uint {
	count := it.Cap - it.Len
	it.free = nil
	it.Cap = it.Len
	return count
}

func (it *List) Slice() []string {
	out := make([]string, it.Len)
	i := 0
	for node := it.head; node != nil; node = node.next {
		if node.Val == nil {
			break
		}

		out[i] = *node.Val
		i++
	}

	return out
}
