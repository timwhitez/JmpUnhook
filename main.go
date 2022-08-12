package main

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
)

func main(){

	//构造函数(申请可执行内存空间)
	a := func()uintptr{return 0}
	ptr := syscall.NewCallback(a)

	//获取脱钩的函数地址
	p,raw,e := JmpUnhook(ptr,"NtWriteVirtualMemory")
	if e != nil{
		panic(e)
	}

	//打印原函数地址
	fmt.Printf("Raw: 0x%x\n",raw)
	//打印新函数地址
	fmt.Printf("Addr: 0x%x\n",p)


	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')


}

