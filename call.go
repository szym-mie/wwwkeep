package wwwkeep

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Caller struct {
	Addr       string
	RemoteId   string
	RemoteMeta map[string]string
}

func (it *Caller) fetchInfo() error {
	// TODO: add version checking
	url := fmt.Sprintf("http://%s/", it.Addr)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	i := 0
	sc := bufio.NewScanner(resp.Body)
	for ; sc.Scan(); i++ {
		line := sc.Text()
		switch i {
		case 0:
			if line != "wwwkeep" {
				return fmt.Errorf("caller_fetch_info: bad signature %s", line)
			}
		case 1:
			id, addr, found := strings.Cut(line, "@")
			if !found {
				return fmt.Errorf("caller_fetch_info: bad id field %s", line)
			}

			if it.Addr != addr {
				log.Printf("caller_fetch_info: addr mismatch %s\n", addr)
			}

			it.RemoteId = id
		case 2:
			// goX.Y.Z version - don't bother
		default:
			key, field, found := strings.Cut(line, "=")
			if !found {
				return fmt.Errorf("caller_fetch_info: bad meta %s", line)
			}

			it.RemoteMeta[key] = field
		}
	}

	if i < 3 {
		return fmt.Errorf("caller_fetch_info: info msg is too short")
	}

	return nil
}

func (it *Caller) rpcUrl() string {
	return fmt.Sprintf("http://%s/rpc", it.Addr)
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
		return nil, fmt.Errorf("op_call: %w - encode fail\n", err)
	}

	resp, err := http.Post(cl.rpcUrl(), "application/octet-stream", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := gob.NewDecoder(resp.Body)
	status := resp.StatusCode
	if status != 200 {
		remoteErrMsg := new(string)
		if err := dec.Decode(remoteErrMsg); err != nil {
			return nil, fmt.Errorf("op_call: %w - decode fail\n", err)
		}
		return nil, fmt.Errorf("op_call: %d -> %s\n", status, *remoteErrMsg)
	}

	reply := new(OpReply)
	if err := dec.Decode(reply); err != nil {
		return nil, fmt.Errorf("op_call: %w - decode fail\n", err)
	}

	return reply, nil
}

func Dial(addr string) (*Caller, error) {
	it := new(Caller)
	it.Addr = addr
	it.RemoteMeta = make(map[string]string)
	if err := it.fetchInfo(); err != nil {
		return nil, fmt.Errorf("caller_dial: %w - info fail", err)
	}

	return it, nil
}
