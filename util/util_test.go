package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

func Test001(test *testing.T){

	var m = make(map[string]string)
	fmt.Println("========= first =============")
	for _,v := range m{
		fmt.Printf("%s\r\n",v)
	}
	m["a1"]="a2"
	m["b1"]="b2"
	m["c1"]="c2"

	fmt.Println("========= Second =============")
	for _,v := range m{
		fmt.Printf("%s\r\n",v)
	}
}

func sayHello(){
	for i := 0; i < 10; i++{
		fmt.Print("hello ")
		runtime.Gosched()
	}
}
func sayWorld(){
	for i := 0; i < 10; i++ {
		fmt.Println("world")
		runtime.Gosched()
	}
}

func Test002(test *testing.T){

	go sayHello()
	go sayWorld()
	var str string
	fmt.Scan(&str)
}

func TestSelect3(test *testing.T) {
	start := time.Now()
	c := make(chan interface{})

	go func() {
		time.Sleep(2*time.Second)
		close(c)
	}()

	fmt.Println("Blocking on read...")
	select {
	case val,ok := <-c:
		fmt.Printf("val=%v, ok=%v\n",val,ok)
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}


func TestSelect4(test *testing.T) {
	a := "1JQR18000001"
	b := strings.ToUpper(string(a[1:4]))

	fmt.Printf("a=%s, b=%s\n",a,b)

	if "JQR" != b {
		fmt.Printf("JQR != %s\n",b)
	}else{
		fmt.Printf("JQR == %s\n",b)
	}
}

func TestSelect5(test *testing.T) {
	a := []byte{1,2}
	var tmp int16
	binary.Read(bytes.NewBuffer(a), binary.BigEndian, &tmp)

	fmt.Printf("tmp=%d\n",int(tmp))


}

func TestSelect6(test *testing.T) {

	const shortForm = "20060102150405"
	t := time.Now()
	//temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	DATE_LAYOUT := "2006-01-02 15:04:05"
	//str := t.Format(DATE_LAYOUT)
	fmt.Println(t.Format(DATE_LAYOUT))
	a := "20190731"
	b :=[]byte(a)

	bin_buf := bytes.NewBuffer([]byte{0,0})
	binary.Write(bin_buf, binary.BigEndian, len(b)/2)

	fmt.Printf("b=%v,len=%d\n",b,bin_buf.Bytes())

	fmt.Printf("1<<8 = %d",1<<8)


}




