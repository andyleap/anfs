package datastore

import (
	"net"
	"sync"
	"errors"

	"github.com/andyleap/anfs/datastore/proto"
)

var (
	ErrLostConn = errors.New("Connection Lost")
)

type Client struct {
	conn       net.Conn
	requests   map[uint64]chan proto.Response
	requestsMu sync.Mutex
	nextID     uint64
}

func Dial(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:     conn,
		requests: map[uint64]chan proto.Response{},
	}
	go c.run()
	return c, nil
}

func (c *Client) run() {
	for {
		var resp proto.Resp
		err := resp.Deserialize(c.conn)
		c.requestsMu.Lock()
		if err != nil {
			for _, v := range c.requests {
				close(v)
			}
			c.requests = nil
			c.requestsMu.Unlock()
			return
		}
		v, ok := c.requests[resp.ID]
		if ok {
			v <- resp.Response
		}
		c.requestsMu.Unlock()
	}
}

func (c *Client) call(req proto.Request) chan proto.Response {
	c.requestsMu.Lock()
	defer c.requestsMu.Unlock()
	if c.requests == nil {
		return nil
	}
	c.nextID++
	pReq := proto.Req{
		ID:      c.nextID,
		Request: req,
	}
	respC := make(chan proto.Response)
	c.requests[pReq.ID] = respC
	err := pReq.Serialize(c.conn)
	if err != nil {
		return nil
	}
	return respC
}

type ServerError string

func (se ServerError) Error() string {
	return string(se)
}

func (c *Client) Read(key []byte) (data []byte, err error) {
	req := proto.ReadReq{
		Key: key,
	}
	resp := <-c.call(req)
	readResp, ok := resp.(proto.ReadResp)
	if !ok {
		return nil, ErrLostConn
	}
	if len(readResp.Error) != 0 {
		return nil, ServerError(readResp.Error)
	}
	return readResp.Data, nil
}

func (c *Client) Write(key []byte, data []byte) (err error) {
	req := proto.WriteReq{
		Key: key,
		Data: data,
	}
	resp := <-c.call(req)
	writeResp, ok := resp.(proto.WriteResp)
	if !ok {
		return ErrLostConn
	}
	if len(writeResp.Error) != 0 {
		return ServerError(writeResp.Error)
	}
	return nil
}

func (c *Client) Delete(key []byte) (err error) {
	req := proto.DeleteReq{
		Key: key,
	}
	resp := <-c.call(req)
	deleteResp, ok := resp.(proto.DeleteResp)
	if !ok {
		return ErrLostConn
	}
	if len(deleteResp.Error) != 0 {
		return ServerError(deleteResp.Error)
	}
	return nil
}

func (c *Client) Clear() (err error) {
	req := proto.ClearReq{}
	resp := <-c.call(req)
	_, ok := resp.(proto.ClearResp)
	if !ok {
		return ErrLostConn
	}
	return nil
}

func (c *Client) Check(key []byte) (found bool, err error) {
	req := proto.CheckReq{Key: key}
	resp := <-c.call(req)
	checkResp, ok := resp.(proto.CheckResp)
	if !ok {
		return false, ErrLostConn
	}
	return checkResp.Found, nil
}

func (c *Client) Sweep() (err error) {
	req := proto.SweepReq{}
	resp := <-c.call(req)
	_, ok := resp.(proto.SweepResp)
	if !ok {
		return ErrLostConn
	}
	return nil
}
