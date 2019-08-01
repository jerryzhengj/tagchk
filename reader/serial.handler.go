package reader

import (
	"bytes"
	"github.com/jerryzhengj/tagchk/util"
	log "github.com/jeanphorn/log4go"
	"io"
	"time"
)

func (session Reader) readPacketHead() int {

	//log.LOGGER("App").Debug("read packet head")
	data := readBlock(session, 2)
	packetSize := data[1]
	return int(packetSize)
}

func (session Reader) readPacketBody(packetSize int) []byte {

	//log.LOGGER("App").Debug("read packet body")
	return readBlock(session, packetSize)
}

func readBlock(session Reader, size int) []byte {

	hasRead := 0
	buffer := bytes.Buffer{}
	for {
		p := make([]byte, size - hasRead)
		readSize, err := session.conn.Read(p)
		if err != nil {
			if err != io.EOF {
				log.LOGGER("App").Error(err)
				panic(err)
			}
		}
		if readSize > 0 {
			hasRead += readSize
			buffer.Write(p[:readSize])
		}

		if hasRead >= size {
			break
		}
		time.Sleep(time.Millisecond * 30)
	}

	data := buffer.Bytes()

	//log.LOGGER("App").Debug("<<%s", util.BytesToHex(data))
	return data
}

// head(0)/len(1)/body(2 : n-2)/check(n - 1)
func (session Reader) writeSerialBody(body []byte) int {

	buffer := bytes.Buffer{}
	bodySize := len(body)
	lenFlag := bodySize + 1 // len的值

	buffer.WriteByte(0xA0) // head
	buffer.WriteByte(byte(lenFlag)) // len
	buffer.Write(body)
	crc := util.CRC(buffer.Bytes())
	buffer.WriteByte(crc)
	return session.writeSerial(buffer.Bytes())
}

func (session Reader) writeSerial(data []byte) int {

	//log.LOGGER("App").Debug(">>%s", util.BytesToHex(data))
	writeSize, err := session.conn.Write(data)
	if err != nil {
		log.LOGGER("App").Error(err)
		panic(err)
	}
	time.Sleep(50 * time.Millisecond)
	return writeSize
}