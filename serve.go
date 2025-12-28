package wwwkeep

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
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

func getInfoMsg(id, addr string) string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		errMsg := "index: read build info fail"
		log.Println(errMsg)
		return errMsg
	}

	goVer := buildInfo.GoVersion
	return fmt.Sprintf("wwwkeep\n%s %s\n%s\n", id, addr, goVer)
}

func (err OpErr) write(w http.ResponseWriter, enc *gob.Encoder) {
	log.Println(fmt.Errorf("op_err: %d -> %w", err.Code, err.Err))
	w.WriteHeader(err.Code)
	if err := enc.Encode(err.Err.Error()); err != nil {
		log.Panicln(fmt.Errorf("op_err_write: %w", err))
	}
}

func (it Keep) Serve(id, addr string) error {
	if strings.Contains(id, " ") {
		return fmt.Errorf("serve: id contains whitespace")
	}

	infoMsgBytes := []byte(getInfoMsg(id, addr))

	opCtxIn := make(chan *OpCtx)
	go it.Handle(opCtxIn)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write(infoMsgBytes); err != nil {
			log.Println(fmt.Errorf("index: %w - write fail", err))
		}
	})
	mux.HandleFunc("POST /rpc", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		enc, dec := gob.NewEncoder(w), gob.NewDecoder(r.Body)
		if r.Method != "POST" {
			OpErr{
				http.StatusMethodNotAllowed,
				fmt.Errorf("serve: method %s not allowed", r.Method)}.write(w, enc)
			return
		}

		opCall := new(OpCall)
		if err := dec.Decode(opCall); err != nil {
			OpErr{
				http.StatusBadRequest,
				fmt.Errorf("serve: %w - decode fail", err)}.write(w, enc)
			return
		}

		if err := ctx.Err(); err != nil {
			log.Println(fmt.Errorf("serve: %w - cancel", ctx.Err()))
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
					fmt.Errorf("serve: %w - call fail\n", err)}.write(w, enc)
				return
			}

			reply, err := opCtx.getOpReply()
			if err != nil {
				OpErr{
					http.StatusInternalServerError,
					fmt.Errorf("serve: %w - no reply\n", err)}.write(w, enc)
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
			log.Println(fmt.Errorf("serve: %w - cancel\n", ctx.Err()))
		}
	})

	return http.ListenAndServe(addr, mux)
}
