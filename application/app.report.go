package application

import (
	"fmt"
	"os"
)

func traceReport(epc string,total int,succ int,duration int64){
	file, err := os.OpenFile("report_result.csv", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil {
		return
	}

	defer file.Close()

	file.WriteString(fmt.Sprintf("%s,%d,%d,%d\r\n",epc,duration,total,succ))
}
