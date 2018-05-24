package connx

import (
	"io"
	"fmt"
	"net"
	"sync/atomic"
	"sync"
	"bytes"
	"encoding/binary"
)

const (
	HeadLenght = 16
	HeadFlagLength = 4
)

var HeadFlag uint32

var connectCreateCount uint64
var connectIndex uint64

var (
	connctionPool sync.Pool
)

func init() {
	HeadFlag = 0x20180618
	connctionPool = sync.Pool{
		New: func() interface{} {
			atomic.AddUint64(&connectCreateCount, 1)
			return &Connection{lock: new(sync.RWMutex)}
		},
	}
}

type HeadInfo struct {
	head_id   uint16 //头标识
	data_type uint16 //数据类型
	data_id   int32  //数据功能ID
	data_len  uint64 //数据长度
}

// GetBytes get bytes
func (h *HeadInfo) GetBytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, h.head_id)
	binary.Write(buf, binary.LittleEndian, h.data_type)
	binary.Write(buf, binary.LittleEndian, h.data_id)
	binary.Write(buf, binary.LittleEndian, h.data_len)
	return buf.Bytes()
}

// FromBytes convert from bytes
func (h *HeadInfo) FromBytes(b []byte) {
	buf := bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &h.head_id)
	binary.Read(buf, binary.LittleEndian, &h.data_type)
	binary.Read(buf, binary.LittleEndian, &h.data_id)
	binary.Read(buf, binary.LittleEndian, &h.data_len)
}

type Connection struct {
	ConnIndex  int64
	lock       *sync.RWMutex
	conn       net.Conn
	Head       *HeadInfo
	flagBuf	    []byte
	headBuf    []byte
	BodyString string
}

//readHeadFlag read head flag message with HeadFlagLength
func readHeadFlag(conn net.Conn) (uint32, error){
	var flagBuf []byte
	err := readSize(conn, HeadFlagLength, &flagBuf)
	if err != nil {
		return 0, err
	}
	var headFlag uint32
	buf := bytes.NewReader(flagBuf)
	err = binary.Read(buf, binary.LittleEndian, &headFlag)
	if err != nil{
		return  0, err
	}else{
		return headFlag, nil
	}
}

//readHead read head message with HeadLenght
func readHead(conn net.Conn) (*HeadInfo, error) {
	var headBuf []byte
	err := readSize(conn, HeadLenght, &headBuf)
	if err != nil {
		return nil, err
	}
	//c.lock.Lock()
	//defer c.lock.Unlock()
	head := &HeadInfo{}
	head.FromBytes(headBuf)
	return head, nil
}

// readSize read message with size
func readSize(conn net.Conn, size int64, buf *[]byte) error {
	*buf = make([]byte, 0)
	var err error
	leftSize := size
	for {

		bufinner := make([]byte, leftSize)
		var n int
		n, err = conn.Read(bufinner)
		leftSize -= int64(n)
		if err == nil {
			*buf = slice_merge(*buf, bufinner)
			if leftSize <= 0 {
				//read end
				break
			}
		} else {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}
	return err
}

func Write(conn net.Conn, p []byte) (int, error) {
	return conn.Write(p)
}

func WriteMerge(conn net.Conn, head []byte, body []byte) (int, error) {
	return conn.Write(slice_merge(head, body))
}

func slice_merge(slice1, slice2 []byte) (c []byte) {
	c = append(slice1, slice2...)
	return
}

