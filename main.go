package main

import (
	"fmt"

	"github.com/wamshawn/vip-demo/npuruntime/viplite"
)

func main() {
	fmt.Println("vip")
	vers := viplite.GetVersion()
	fmt.Println(vers)

	//npu := npuruntime.NewNpuUnit()
	//defer npu.Close()
	//defer npu.Destroy()
	//vers := npu.GetDriverVersion()
	//fmt.Println(vers)
}
