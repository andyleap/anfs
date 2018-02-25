package proto

//go:generate gencode go -package proto -schema proto.schema

type Request interface {
	request()
}

type Response interface {
	response()
}

func (ReadReq) request()  {}
func (WriteReq) request() {}
func (DeleteReq) request() {}
func (ClearReq) request() {}
func (CheckReq) request() {}
func (SweepReq) request() {}

func (ReadResp) response()  {}
func (WriteResp) response() {}
func (DeleteResp) response() {}
func (ClearResp) response() {}
func (CheckResp) response() {}
func (SweepResp) response() {}
