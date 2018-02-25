package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/andyleap/anfs/anfs/proto"
	"github.com/andyleap/anfs/dsdist"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/coreos/etcd/clientv3"
	"github.com/jessevdk/go-flags"
)

var Options struct {
	DatastoreHosts []string `long:"datastore"`
	EtcdCluster    []string `long:"etcd"`
	Volume         string   `long:"volume"`
	Debug          bool     `long:"debug"`
	AllowOther     bool     `long:"allow-other"`
}

func main() {
	args, _ := flags.Parse(&Options)
	if len(args) != 1 {
		log.Fatal("expected 1 arg")
	}
	mountpoint := args[0]

	if Options.Debug {
		f, err := os.Create(mountpoint + ".log")
		if err == nil {
			defer f.Close()
			fuse.Debug = func(msg interface{}) {
				fmt.Fprintln(f, msg)
			}
		}
	}

	options := []fuse.MountOption{
		fuse.NoAppleDouble(),
		fuse.NoAppleXattr(),
		fuse.FSName("anfs"),
		fuse.Subtype("anfs"),
	}

	if Options.AllowOther {
		options = append(options, fuse.AllowOther())
	}

	c, err := fuse.Mount(mountpoint, options...)

	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	names := Options.DatastoreHosts
	log.Println("datastore", names)
	ds := dsdist.New(names, 2, 1)

	names = Options.EtcdCluster
	log.Println("etcd", names)
	client, _ := clientv3.New(clientv3.Config{Endpoints: names})

	anfs := New(ds, client, Options.Volume)
	err = fs.Serve(c, anfs)
	if err != nil {
		log.Fatal(err)
	}

	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}

type FS struct {
	ds     *dsdist.DSDist
	client *clientv3.Client
	volume string

	open   map[uint64]interface{}
	openMu sync.Mutex
}

func New(ds *dsdist.DSDist, client *clientv3.Client, volume string) *FS {
	return &FS{
		ds:     ds,
		client: client,
		volume: volume,

		open: map[uint64]interface{}{},
	}
}

func (fs *FS) Root() (fs.Node, error) {
	fs.openMu.Lock()
	defer fs.openMu.Unlock()
	resp, err := fs.client.Get(context.Background(), fmt.Sprintf("/anfs/volume/%s", fs.volume))
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		inode := fs.nextInode()
		dir := &Dir{
			inode: inode,
			fs:    fs,
		}
		dir.data.Attr = proto.Attr{
			Mode: uint32(os.ModeDir) | 0777,
		}
		dir.save()
		val := make([]byte, 8)
		binary.LittleEndian.PutUint64(val, inode)
		fs.client.Put(context.Background(), fmt.Sprintf("/anfs/volume/%s", fs.volume), string(val))
		fs.atomicAdd(fmt.Sprintf("/anfs/inode/%d/link", inode), 1)
		return dir, nil
	}
	inode := binary.LittleEndian.Uint64(resp.Kvs[0].Value)
	o := fs.open[inode]
	if o != nil {
		return o.(*Dir), nil
	}
	resp, err = fs.client.Get(context.Background(), fmt.Sprintf("/anfs/inode/%d", inode))
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) != 1 {
		return nil, fuse.ENOENT
	}
	dir := &Dir{
		fs:    fs,
		inode: inode,
	}
	dir.data.Unmarshal(resp.Kvs[0].Value)
	fs.open[inode] = dir
	return dir, nil
}

func (fs *FS) Statfs(ctx context.Context, req *fuse.StatfsRequest, resp *fuse.StatfsResponse) error {
	resp.Bsize = 16 * 1024
	return nil
}

func (fs *FS) nextInode() uint64 {
	return uint64(fs.atomicAdd("/anfs/inode", 1))
}

func (fs *FS) atomicAdd(key string, add int64) int64 {
	ctx := context.Background()
	for {
		resp, _ := fs.client.Get(ctx, key)
		val := int64(0)
		txn := fs.client.Txn(ctx)
		if len(resp.Kvs) > 0 {
			val = int64(binary.LittleEndian.Uint64(resp.Kvs[0].Value))
			txn.If(clientv3.Compare(clientv3.Version(key), "=", resp.Kvs[0].Version))
		} else {
			txn.If(clientv3.Compare(clientv3.Version(key), "=", 0))
		}
		if add == 0 {
			return val
		}
		val += add
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(val))
		txn.Then(clientv3.OpPut(key, string(buf)))
		txnresp, _ := txn.Commit()
		if txnresp.Succeeded {
			return val
		}
	}
}

type Inoder interface {
	Inode() uint64
}

type Typer interface {
	Type() uint16
}

type Dir struct {
	fs    *FS
	inode uint64
	data  proto.Directory
}

func (d *Dir) Inode() uint64 {
	return d.inode
}

func (d *Dir) Type() uint16 {
	return proto.TypeDir
}

func (d *Dir) save() {
	buf, _ := d.data.Marshal(nil)
	d.fs.client.Put(context.Background(), fmt.Sprintf("/anfs/inode/%d", d.inode), string(buf))
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Println("ATTR")
	attr := d.data.Attr
	a.Inode = d.inode
	a.Mode = os.ModeDir | os.FileMode(attr.Mode)
	a.Uid = attr.UID
	a.Gid = attr.GID
	a.Nlink = uint32(d.fs.atomicAdd(fmt.Sprintf("/anfs/inode/%d/link", d.inode), 0))
	return nil
}

func (d *Dir) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	log.Println("SETATTR")
	attr := d.data.Attr
	if req.Valid.Mode() {
		attr.Mode = uint32(req.Mode)
	}
	if req.Valid.Uid() {
		attr.UID = req.Uid
	}
	if req.Valid.Gid() {
		attr.GID = req.Gid
	}
	d.data.Attr = attr
	d.save()
	return nil
}

func (d *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	d.fs.openMu.Lock()
	defer d.fs.openMu.Unlock()
	inode := d.fs.nextInode()
	d.data.Items = append(d.data.Items, proto.Item{Name: req.Name, Inode: inode, Type: proto.TypeDir})
	dir := &Dir{
		fs:    d.fs,
		inode: inode,
	}
	dir.data.Attr = proto.Attr{
		GID:  req.Gid,
		UID:  req.Uid,
		Mode: uint32(req.Mode),
	}
	dir.save()
	d.fs.atomicAdd(fmt.Sprintf("/anfs/inode/%d/link", inode), 1)
	d.fs.open[inode] = dir
	d.save()
	return dir, nil
}

func (d *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	d.fs.openMu.Lock()
	defer d.fs.openMu.Unlock()
	for _, i := range d.data.Items {
		if i.Name == req.Name {
			if req.Flags&fuse.OpenExclusive != 0 {
				return nil, nil, fuse.EEXIST
			} else {
				panic("creating existing file")
			}
		}
	}
	inode := d.fs.nextInode()
	d.data.Items = append(d.data.Items, proto.Item{Name: req.Name, Inode: inode, Type: proto.TypeFile})
	file := &File{
		fs:    d.fs,
		inode: inode,
	}
	file.data.BlockSize = 1024 * 16
	file.data.Attr = proto.Attr{
		Mode: uint32(req.Mode),
		UID:  req.Uid,
		GID:  req.Gid,
	}
	d.save()
	file.save()
	d.fs.atomicAdd(fmt.Sprintf("/anfs/inode/%d/link", inode), 1)
	return file, file, nil
}

func (d *Dir) Link(ctx context.Context, req *fuse.LinkRequest, old fs.Node) (fs.Node, error) {
	d.data.Items = append(d.data.Items, proto.Item{
		Inode: old.(Inoder).Inode(),
		Name:  req.NewName,
		Type:  old.(Typer).Type(),
	})
	d.save()
	d.fs.atomicAdd(fmt.Sprintf("/anfs/inode/%d/link", old.(Inoder).Inode()), 1)
	return old, nil
}

func (d *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	node, err := d.Lookup(ctx, req.Name)
	if err != nil {
		return err
	}
	node.(Unlinker).Unlink()
	for i, item := range d.data.Items {
		if item.Name == req.Name {
			copy(d.data.Items[i:], d.data.Items[i+1:])
			d.data.Items = d.data.Items[:len(d.data.Items)-1]
			d.save()
			return nil
		}
	}
	return nil
}

func (d *Dir) Rename(ctx context.Context, req *fuse.RenameRequest, newNode fs.Node) error {
	newDir := newNode.(*Dir)
	if old, err := newDir.Lookup(ctx, req.NewName); err == nil {
		old.(Unlinker).Unlink()
	renameDeleteLoop:
		for i, item := range newDir.data.Items {
			if item.Name == req.NewName {
				copy(newDir.data.Items[i:], newDir.data.Items[i+1:])
				newDir.data.Items = newDir.data.Items[:len(newDir.data.Items)-1]
				break renameDeleteLoop
			}
		}
	}

	if d.Inode() == newDir.Inode() {
		for i, item := range d.data.Items {
			if item.Name == req.OldName {
				d.data.Items[i].Name = req.NewName
				return nil
			}
		}
	}
	var renameItem proto.Item
renameFromLoop:
	for i, item := range d.data.Items {
		if item.Name == req.OldName {
			renameItem = item
			copy(d.data.Items[i:], d.data.Items[i+1:])
			d.data.Items = d.data.Items[:len(d.data.Items)-1]
			d.save()
			break renameFromLoop
		}
	}

	renameItem.Name = req.NewName
	newDir.data.Items = append(newDir.data.Items, renameItem)
	newDir.save()
	return nil
}

type Unlinker interface {
	Unlink()
}

func (d *Dir) Unlink() {
	linkCount := d.fs.atomicAdd(fmt.Sprintf("/anfs/inode/%d/link", d.inode), -1)
	if linkCount > 0 {
		return
	}
	ctx := context.Background()
	for _, item := range d.data.Items {
		node, err := d.Lookup(ctx, item.Name)
		if err != nil {
			continue
		}
		node.(Unlinker).Unlink()
	}
	d.fs.client.Delete(ctx, fmt.Sprintf("/anfs/inode/%d", d.inode), clientv3.WithPrefix())
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	for _, item := range d.data.Items {
		if item.Name == name {
			d.fs.openMu.Lock()
			defer d.fs.openMu.Unlock()
			switch item.Type {
			case proto.TypeFile:
				o := d.fs.open[item.Inode]
				if o != nil {
					return o.(*File), nil
				}
				resp, err := d.fs.client.Get(ctx, fmt.Sprintf("/anfs/inode/%d", item.Inode))
				if err != nil {
					return nil, err
				}
				if len(resp.Kvs) != 1 {
					return nil, fuse.ENOENT
				}
				file := &File{
					fs:    d.fs,
					inode: item.Inode,
				}
				file.data.Unmarshal(resp.Kvs[0].Value)
				d.fs.open[item.Inode] = file
				return file, nil
			case proto.TypeDir:
				o := d.fs.open[item.Inode]
				if o != nil {
					return o.(*Dir), nil
				}
				resp, err := d.fs.client.Get(ctx, fmt.Sprintf("/anfs/inode/%d", item.Inode))
				if err != nil {
					return nil, err
				}
				if len(resp.Kvs) != 1 {
					return nil, fuse.ENOENT
				}
				dir := &Dir{
					fs:    d.fs,
					inode: item.Inode,
				}
				dir.data.Unmarshal(resp.Kvs[0].Value)
				d.fs.open[item.Inode] = dir
				return dir, nil
			}
		}
	}
	return nil, fuse.ENOENT
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	ret := make([]fuse.Dirent, len(d.data.Items))
	for i, item := range d.data.Items {
		ret[i].Inode = item.Inode
		ret[i].Name = item.Name
		ret[i].Type = fuse.DT_File
		if item.Type == proto.TypeDir {
			ret[i].Type = fuse.DT_Dir
		}
	}
	return ret, nil
}

type File struct {
	fs    *FS
	inode uint64
	data  proto.File
}

type Handle struct {
	file *File
	mode uint32
}

func (f *File) save() {
	buf, _ := f.data.Marshal(nil)
	f.fs.client.Put(context.Background(), fmt.Sprintf("/anfs/inode/%d", f.inode), string(buf))
}

func (f *File) Unlink() {
	linkCount := f.fs.atomicAdd(fmt.Sprintf("/anfs/inode/%d/link", f.inode), -1)
	if linkCount > 0 {
		return
	}
	ctx := context.Background()
	f.truncate(0)
	f.fs.client.Delete(ctx, fmt.Sprintf("/anfs/inode/%d", f.inode), clientv3.WithPrefix())
}

func (f *File) Inode() uint64 {
	return f.inode
}

func (f *File) Type() uint16 {
	return proto.TypeFile
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	attr := f.data.Attr
	a.Mode = os.FileMode(attr.Mode)
	a.Size = f.data.Len
	a.Uid = attr.UID
	a.Gid = attr.GID
	a.Nlink = uint32(f.fs.atomicAdd(fmt.Sprintf("/anfs/inode/%d/link", f.inode), 0))
	return nil
}

func (f *File) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	attr := f.data.Attr
	if req.Valid.Mode() {
		attr.Mode = uint32(req.Mode)
	}
	if req.Valid.Uid() {
		attr.UID = req.Uid
	}
	if req.Valid.Gid() {
		attr.GID = req.Gid
	}
	f.data.Attr = attr
	if req.Valid.Size() {
		f.truncate(req.Size)
	}
	f.save()
	return nil
}

func (f *File) truncate(size uint64) {
	oldSize := f.data.Len
	f.data.Len = size
	if size > oldSize {
		return
	}
	end := ((size - 1) / f.data.BlockSize) + 1
	if size == 0 {
		end = 0
	}
	index, bit := end/8, uint(end%8)
	for index < uint64(len(f.data.SparseMap)) {
		if f.data.SparseMap[index]&(1<<bit) != 0 {
			key := make([]byte, 16)
			binary.LittleEndian.PutUint64(key, f.inode)
			binary.LittleEndian.PutUint64(key[8:], uint64(end))
			f.fs.ds.Delete(key)
			f.data.SparseMap[index] ^= (1 << bit)
		}
		bit++
		end++
		if bit >= 8 {
			bit = 0
			index++
		}
	}
	f.save()
}

func (f *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) (err error) {
	if cap(resp.Data) < req.Size {
		resp.Data = make([]byte, req.Size)
	} else {
		resp.Data = resp.Data[:req.Size]
	}
	off := req.Offset

	if off >= int64(f.data.Len) {
		resp.Data = resp.Data[:0]
		return nil
	}

	end := off + int64(len(resp.Data))
	n := len(resp.Data)
	if end > int64(f.data.Len) {
		end = int64(f.data.Len)
		n = int(f.data.Len) - int(off)
	}
	start := off / int64(f.data.BlockSize)
	woff := off % int64(f.data.BlockSize)
	for l1 := start; l1 < ((end-1)/int64(f.data.BlockSize))+1; l1++ {
		var data []byte

		index, bit := l1/8, uint(l1%8)
		if index >= int64(len(f.data.SparseMap)) || f.data.SparseMap[index]&(1<<bit) == 0 {
			data = make([]byte, f.data.BlockSize)
		} else {
			var err error
			key := make([]byte, 16)
			binary.LittleEndian.PutUint64(key, f.inode)
			binary.LittleEndian.PutUint64(key[8:], uint64(l1))
			data, err = f.fs.ds.Read(key)
			if err != nil {
				resp.Data = resp.Data[:0]
				return err
			}
		}
		boff := ((l1 - start) * int64(f.data.BlockSize)) - woff
		if boff < 0 {
			copy(resp.Data, data[-boff:])
		} else {
			copy(resp.Data[boff:], data)
		}
	}
	return

	resp.Data = resp.Data[:n]
	return err
}

func (f *File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) (err error) {
	off := req.Offset
	end := int64(len(req.Data)) + off
	if uint64(end) > f.data.Len {
		f.data.Len = uint64(end)
	}
	neededBlocks := ((end - 1) / int64(f.data.BlockSize)) + 1
	if int64(len(f.data.SparseMap))/8 <= neededBlocks {
		old := f.data.SparseMap
		f.data.SparseMap = make([]byte, (neededBlocks/8)+1)
		copy(f.data.SparseMap, old)
	}
	n := len(req.Data)
	start := off / int64(f.data.BlockSize)
	woff := off % int64(f.data.BlockSize)
	wg := sync.WaitGroup{}
	for l1 := start; l1 < ((end-1)/int64(f.data.BlockSize))+1; l1++ {
		wg.Add(1)
		go func(l1 int64) {
			defer wg.Done()
			var data []byte
			boff := ((l1 - start) * int64(f.data.BlockSize)) - woff
			if (l1 == start && woff != 0) || (l1 == end/int64(f.data.BlockSize) && end%int64(f.data.BlockSize) != 0) {
				index, bit := l1/8, uint(l1%8)
				if f.data.SparseMap[index]&(1<<bit) == 0 {
					data = make([]byte, f.data.BlockSize)
				} else {
					var err error
					key := make([]byte, 16)
					binary.LittleEndian.PutUint64(key, f.inode)
					binary.LittleEndian.PutUint64(key[8:], uint64(l1))
					data, err = f.fs.ds.Read(key)
					if err != nil {
						//                return 0, err
					}
				}
				if uint64(len(data)) < f.data.BlockSize {
					old := data
					data = make([]byte, f.data.BlockSize)
					copy(data, old)
				}
				if boff < 0 {
					copy(data[-boff:], req.Data)
				} else {
					copy(data, req.Data[boff:])
				}
			} else {
				data = req.Data[boff : boff+int64(f.data.BlockSize)]
			}
			var err error
			key := make([]byte, 16)
			binary.LittleEndian.PutUint64(key, f.inode)
			binary.LittleEndian.PutUint64(key[8:], uint64(l1))
			err = f.fs.ds.Write(key, data)
			index, bit := l1/8, uint(l1%8)
			f.data.SparseMap[index] |= (1 << bit)
			if err != nil {
				//              return 0, err
			}
		}(l1)
	}
	wg.Wait()
	f.save()
	resp.Size = n
	return
}

func (f *File) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	return nil
}

func (f *File) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	return nil
}
