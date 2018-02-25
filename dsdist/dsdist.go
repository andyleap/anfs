package dsdist

import (
	"encoding/hex"
	"fmt"
	"net"
	"sync"

	"github.com/andyleap/anfs/datastore"
	"github.com/klauspost/reedsolomon"
	"github.com/serialx/hashring"
)

type DSDist struct {
	nodes        []string
	hr           *hashring.HashRing
	dsClients    map[string]*datastore.Client
	dsClientsMu  sync.RWMutex
	data, parity int
	rsCache      map[dp]reedsolomon.Encoder
	rsCacheMu    sync.RWMutex
}

func New(names []string, data, parity int) *DSDist {
	nodes := []string{}
	for _, name := range names {
		host, port, err := net.SplitHostPort(name)
		if err != nil {
			host = name
			port = "3124"
		}
		ip := net.ParseIP(host)
		if ip == nil {
			ips, _ := net.LookupIP(host)
			for _, ip := range ips {
				nodes = append(nodes, fmt.Sprintf("%s:%s", ip.String(), port))
			}
		} else {
			nodes = append(nodes, fmt.Sprintf("%s:%s", ip.String(), port))
		}
	}

	return &DSDist{
		nodes:     nodes,
		hr:        hashring.New(nodes),
		dsClients: map[string]*datastore.Client{},
		data:      data,
		parity:    parity,
		rsCache:   map[dp]reedsolomon.Encoder{},
	}
}

type dp struct{ Data, Parity int }

func (dsd *DSDist) getRS(data, parity int) reedsolomon.Encoder {
	dsd.rsCacheMu.RLock()
	rs := dsd.rsCache[dp{data, parity}]
	if rs != nil {
		return rs
	}
	dsd.rsCacheMu.RUnlock()
	rs, _ = reedsolomon.New(data, parity)
	dsd.rsCacheMu.Lock()
	dsd.rsCache[dp{data, parity}] = rs
	dsd.rsCacheMu.Unlock()
	return rs
}

func (dsd *DSDist) call(node string, f func(ds *datastore.Client) error) error {
	dsd.dsClientsMu.RLock()
	dsClient := dsd.dsClients[node]
	dsd.dsClientsMu.RUnlock()
	for {
		if dsClient == nil {
			var err error
			dsClient, err = datastore.Dial(node)
			if err != nil {
				return err
			}
			dsd.dsClientsMu.Lock()
			dsd.dsClients[node] = dsClient
			dsd.dsClientsMu.Unlock()
		}
		err := f(dsClient)
		if err != datastore.ErrLostConn {
			return err
		}
	}
}

func (dsd *DSDist) Write(key []byte, data []byte) error {
	nodes, _ := dsd.hr.GetNodes(hex.EncodeToString(key), dsd.data+dsd.parity)
	ds, ps := dsd.data, dsd.parity
	if len(nodes) < dsd.data+dsd.parity {
		ratio := float64(dsd.data) / float64(dsd.parity+dsd.data)
		ds = int(float64(len(nodes)) * ratio)
		if ds < 1 {
			ds = 1
		}
		ps = len(nodes) - ds
	}
	rs := dsd.getRS(ds, ps)
	shards, _ := rs.Split(data)
	rs.Encode(shards)
	trim := len(data) % ds
	for i, shard := range shards {
		node := nodes[i%len(nodes)]
		shard = append([]byte{byte(i), byte(ds), byte(ps), byte(trim)}, shard...)

		dsd.call(node, func(ds *datastore.Client) error {
			return ds.Write(key, shard)
		})
	}
	return nil
}

type Choice struct {
	shards [][]byte
	trim   byte
	found  int
}

func (dsd *DSDist) Read(key []byte) ([]byte, error) {
	nodes, _ := dsd.hr.GetNodes(hex.EncodeToString(key), len(dsd.nodes))
	choices := map[dp]*Choice{}
	found := false
	data := []byte{}
	for _, node := range nodes {
		dsd.call(node, func(ds *datastore.Client) error {
			shard, err := ds.Read(key)
			if err != nil {
				return err
			}
			choice, ok := choices[dp{int(shard[1]), int(shard[2])}]
			if !ok {
				choice = &Choice{
					shards: make([][]byte, int(shard[1])+int(shard[2])),
					trim:   shard[3],
				}
				choices[dp{int(shard[1]), int(shard[2])}] = choice
			}
			choice.shards[shard[0]] = shard[4:]
			choice.found++
			if choice.found >= int(shard[1]) {
				rs := dsd.getRS(int(shard[1]), int(shard[2]))
				err := rs.ReconstructData(choice.shards)
				if err != nil {
					return err
				}
				data = make([]byte, (int(shard[1])*len(choice.shards[0]))-int(choice.trim))
				for i, dshard := range choice.shards[:int(shard[1])] {
					copy(data[i*len(dshard):], dshard)
				}
				found = true
			}
			return nil
		})
		if found {
			return data, nil
		}
	}
	return nil, fmt.Errorf("Could not find data")
}

func (dsd *DSDist) Delete(key []byte) error {
	nodes, _ := dsd.hr.GetNodes(hex.EncodeToString(key), len(dsd.nodes))
	for _, node := range nodes {
		dsd.call(node, func(ds *datastore.Client) error {
			ds.Delete(key)
			return nil
		})
	}
	return nil
}
