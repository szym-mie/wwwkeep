package wwwkeep

import (
	"fmt"
)

const (
	NilOp = "nil"
	DefOp = "def"
	AddOp = "add"
	GetOp = "get"
	PopOp = "pop"
	LenOp = "len"
	DirOp = "dir"
	OptOp = "opt"
)

type OpId string

// OpCall represents remote call originating from a client.
type OpCall struct {
	Id   OpId
	Args Args
}

// OpReply is a union type of all possible reply types. Normally, only one of
// the pointers will be non-nil.
type OpReply struct {
	Uint *uint
	Vals *Vals
	Dirs *Dirs
}

// OpCtx represents server-side operation state.
type OpCtx struct {
	OpCall
	OpReply
	errorOut chan error
}

func newOpCtx(opCall *OpCall) *OpCtx {
	opCtx := new(OpCtx)
	opCtx.Id = opCall.Id
	opCtx.Args = opCall.Args
	opCtx.errorOut = make(chan error)
	return opCtx
}

func (it *OpCtx) String() string {
	return fmt.Sprintf("%s: %s", it.Id, it.Args)
}

func (it *OpReply) replyWithUint(val uint, err error) error {
	if err != nil {
		return err
	}

	it.Uint = &val
	return nil
}

func (it *OpReply) replyWithVals(val Vals, err error) error {
	if err != nil {
		return err
	}

	it.Vals = &val
	return nil
}

func (it *OpReply) replyWithDirs(val Dirs, err error) error {
	if err != nil {
		return err
	}

	it.Dirs = &val
	return nil
}

func (it *OpCtx) getOpReply() (*OpReply, error) {
	if it.Uint == nil && it.Vals == nil && it.Dirs == nil {
		return nil, fmt.Errorf("get_op_reply: no value in %p", it)
	}

	return &it.OpReply, nil
}
