package wwwkeep

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"
)

type OpErr struct {
	Code int
	Err  error
}

func (it Keep) Handle(opCtxIn <-chan *OpCtx) {
	for {
		op := <-opCtxIn
		log.Printf("keep_handle: %s\n", op)
		switch op.Id {
		case NilOp:
			op.errorOut <- nil
		case DefOp:
			op.errorOut <- op.replyWithUint(it.def(&op.Args))
		case AddOp:
			op.errorOut <- op.replyWithUint(it.add(&op.Args))
		case GetOp:
			op.errorOut <- op.replyWithVals(it.get(&op.Args))
		case PopOp:
			op.errorOut <- op.replyWithUint(it.pop(&op.Args))
		case LenOp:
			op.errorOut <- op.replyWithUint(it.len(&op.Args))
		case DirOp:
			op.errorOut <- op.replyWithDirs(it.dir(&op.Args))
		default:
			op.errorOut <- fmt.Errorf(
				"keep_handle: bad op id (%s)",
				op.Id)
		}
	}
}

func (err OpErr) write(w http.ResponseWriter, enc *gob.Encoder) {
	log.Printf("op_err: %d -> %s\n", err.Code, err.Err)
	w.WriteHeader(err.Code)
	if err := enc.Encode(err.Err.Error()); err != nil {
		log.Panicf("op_err_write: %s\n", err)
	}
}

func (it Keep) Serve(addr string) error {
	opCtxIn := make(chan *OpCtx)
	go it.Handle(opCtxIn)

	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		enc, dec := gob.NewEncoder(w), gob.NewDecoder(r.Body)
		if r.Method != "POST" {
			OpErr{
				http.StatusMethodNotAllowed,
				fmt.Errorf("serve: method %s not allowed\n", r.Method)}.write(w, enc)
			return
		}

		opCall := new(OpCall)
		if err := dec.Decode(opCall); err != nil {
			OpErr{
				http.StatusBadRequest,
				fmt.Errorf("serve: %s - decode fail\n", err)}.write(w, enc)
			return
		}

		if err := ctx.Err(); err != nil {
			log.Printf("serve: %s - cancel\n", ctx.Err())
			return
		}

		opCtx := newOpCtx(opCall)
		opCtxIn <- opCtx
		timer := time.NewTimer(2 * time.Second)
		select {
		case err := <-opCtx.errorOut:
			timer.Stop()
			if err != nil {
				OpErr{
					http.StatusInternalServerError,
					fmt.Errorf("serve: %s - call fail\n", err)}.write(w, enc)
				return
			}

			reply, err := opCtx.getOpReply()
			if err != nil {
				OpErr{
					http.StatusInternalServerError,
					fmt.Errorf("serve: %s - no reply\n", err)}.write(w, enc)
				return
			}

			w.WriteHeader(200)
			enc.Encode(*reply)
		case <-timer.C:
			OpErr{
				http.StatusServiceUnavailable,
				fmt.Errorf("serve: call timeout\n")}.write(w, enc)
		case <-ctx.Done():
			timer.Stop()
			log.Printf("serve: %s - cancel\n", ctx.Err())
		}
	})

	return http.ListenAndServe(addr, mux)
}
