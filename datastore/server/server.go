package main

import (
	"encoding/hex"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/andyleap/anfs/datastore/proto"
)

type Server struct {
	l       net.Listener
	baseDir string
	inUse   *Marker
}

type Marker struct {
	subMarkers map[string]*Marker
	items      map[string]struct{}
	mu         sync.RWMutex
}

func NewMarker() *Marker {
	return &Marker{
		subMarkers: map[string]*Marker{},
		items:      map[string]struct{}{},
	}
}

func (m *Marker) Mark(item ...string) {
	if len(item) > 1 {
		m.mu.RLock()
		if sm, ok := m.subMarkers[item[0]]; ok {
			m.mu.RUnlock()
			sm.Mark(item[1:]...)
			return
		}
		m.mu.RUnlock()
		m.mu.Lock()
		if sm, ok := m.subMarkers[item[0]]; ok {
			m.mu.Unlock()
			sm.Mark(item[1:]...)
		} else {
			m.subMarkers[item[0]] = NewMarker()
			m.mu.Unlock()
		}
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[item[0]] = struct{}{}
}

func (m *Marker) Check(item ...string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(item) > 1 {
		if sm, ok := m.subMarkers[item[0]]; ok {
			return sm.Check(item[1:]...)
		}
		return false
	}
	_, ok := m.items[item[0]]
	return ok
}

func (m *Marker) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sm := range m.subMarkers {
		sm.Clear()
	}
	m.items = map[string]struct{}{}
}

func New(addr string, baseDir string) (*Server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		l:       l,
		baseDir: baseDir,
		inUse:   NewMarker(),
	}, nil
}

func (s *Server) Serve() error {
	for {
		conn, err := s.l.Accept()
		if err != nil {
			return err
		}
		go s.Handle(conn)
	}
}

func (s *Server) Handle(conn net.Conn) {
	for {
		var req proto.Req
		var resp proto.Resp
		err := req.Deserialize(conn)
		if err != nil {
			return
		}
		go func() {
			resp.ID = req.ID
			switch r := req.Request.(type) {
			case proto.ReadReq:
				hexKey := hex.EncodeToString(r.Key)

				dataPath := filepath.Join(s.baseDir, hexKey[:2], hexKey[2:4], hexKey[4:])

				data, err := ioutil.ReadFile(dataPath)
				errStr := ""
				if err != nil {
					errStr = err.Error()
				}
				resp.Response = proto.ReadResp{
					Data:  data,
					Error: errStr,
				}
			case proto.WriteReq:
				hexKey := hex.EncodeToString(r.Key[:])

				dataPath := filepath.Join(s.baseDir, hexKey[:2], hexKey[2:4], hexKey[4:])

				os.MkdirAll(filepath.Join(s.baseDir, hexKey[:2], hexKey[2:4]), 0777)
				err := ioutil.WriteFile(dataPath, r.Data, 0666)
				if err != nil {
					resp.Response = proto.WriteResp{
						Error: err.Error(),
					}
				} else {
					resp.Response = proto.WriteResp{}
				}
			case proto.DeleteReq:
				hexKey := hex.EncodeToString(r.Key[:])

				dataPath := filepath.Join(s.baseDir, hexKey[:2], hexKey[2:4], hexKey[4:])

				err := os.Remove(dataPath)
				if err != nil {
					resp.Response = proto.DeleteResp{
						Error: err.Error(),
					}
				} else {
					resp.Response = proto.DeleteResp{}
				}

			case proto.ClearReq:
				s.inUse.Clear()
				resp.Response = proto.ClearResp{}
			case proto.CheckReq:
				hexKey := hex.EncodeToString(r.Key)
				path := []string{hexKey[:2], hexKey[2:4], hexKey[4:]}
				dataPath := filepath.Join(s.baseDir, hexKey[:2], hexKey[2:4], hexKey[4:])
				_, err := os.Stat(dataPath)
				if err != nil {
					s.inUse.Mark(path...)
				}
				resp.Response = proto.CheckResp{
					Found: err == nil,
				}
			case proto.SweepReq:
				s.GCSweep(make([]string, 0, 10)...)
				resp.Response = proto.SweepResp{}
			}
			resp.Serialize(conn)
		}()
	}
}

func (s *Server) GCSweep(path ...string) {
	dir, _ := os.Open(filepath.Join(s.baseDir, filepath.Join(path...)))
	defer dir.Close()
	items, _ := dir.Readdir(0)
	for _, i := range items {
		if i.IsDir() {
			s.GCSweep(append(path, i.Name())...)
		} else {
			if !s.inUse.Check(append(path, i.Name())...) {
				os.Remove(filepath.Join(s.baseDir, filepath.Join(path...), i.Name()))
			}
		}
	}
}
