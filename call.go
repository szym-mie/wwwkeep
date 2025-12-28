package wwwkeep

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
)

type Caller struct {
	addr string
}

func getUintReply(reply *OpReply, err error) (*uint, error) {
	if err != nil {
		return nil, err
	}

	return reply.Uint, nil
}

func (it *Caller) Def(nodeName string, keys []string, initCap uint) (*uint, error) {
	args := Args{nodeName, "", keys, nil, initCap}
	opCall := &OpCall{DefOp, args}
	return getUintReply(opCall.call(it))
}

func (it *Caller) Add(nodeName string, tuple map[string]string) (*uint, error) {
	args := Args{nodeName, "", nil, tuple, 0}
	opCall := &OpCall{AddOp, args}
	return getUintReply(opCall.call(it))
}

func (it *Caller) Get(nodeName string, vecName string) (*Vals, error) {
	args := Args{nodeName, vecName, nil, nil, 0}
	opCall := &OpCall{GetOp, args}
	reply, err := opCall.call(it)
	if err != nil {
		return nil, err
	}

	return reply.Vals, nil
}

func (it *Caller) Pop(nodeName string, vecName string, count uint) (*uint, error) {
	args := Args{nodeName, vecName, nil, nil, count}
	opCall := &OpCall{PopOp, args}
	return getUintReply(opCall.call(it))
}

func (it *Caller) Len(nodeName string, vecName string) (*uint, error) {
	args := Args{nodeName, vecName, nil, nil, 0}
	opCall := &OpCall{LenOp, args}
	return getUintReply(opCall.call(it))
}

func (it *Caller) Dir(nodeName string) (*Dirs, error) {
	args := Args{nodeName, "", nil, nil, 0}
	opCall := &OpCall{DirOp, args}
	reply, err := opCall.call(it)
	if err != nil {
		return nil, err
	}

	return reply.Dirs, nil
}

func (it *OpCall) call(cl *Caller) (*OpReply, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(*it); err != nil {
		return nil, fmt.Errorf("op_call: %s - encode fail\n", err)
	}

	resp, err := http.Post(cl.addr+"/rpc", "application/octet-stream", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := gob.NewDecoder(resp.Body)
	status := resp.StatusCode
	if status != 200 {
		remoteErrMsg := new(string)
		if err := dec.Decode(remoteErrMsg); err != nil {
			return nil, fmt.Errorf("op_call: %s - decode fail\n", err)
		}
		return nil, fmt.Errorf("op_call: %d -> %s\n", status, *remoteErrMsg)
	}

	reply := new(OpReply)
	if err := dec.Decode(reply); err != nil {
		return nil, fmt.Errorf("op_call: %s - decode fail\n", err)
	}

	return reply, nil
}

func Dial(addr string) (*Caller, error) {
	// TODO use persistent connection - declare gob enc/dec here etc.
	return &Caller{addr}, nil
}
