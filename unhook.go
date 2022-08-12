package main

import (
	"fmt"
	"github.com/Binject/debug/pe"
	"golang.org/x/sys/windows"
	"unsafe"
)

func JmpUnhook(ptr uintptr,funcN string)(uintptr,uintptr,error){

	//构造填充
	jmpPre := []byte{0x49, 0xBB}
	//jmpAddr := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF }
	jmpRet := []byte{0x41, 0xFF, 0xE3, 0xC3}


	//获取内存中函数地址
	memPtr := GetFunc(funcN)

	//获取文件内函数
	p, e := pe.Open("C:\\windows\\system32\\ntdll.dll")
	if e!= nil {
		return 0,memPtr,e
	}
	ex, e := p.Exports()
	if e != nil {
		return 0,memPtr,e
	}
	var buff []byte
	for _, exp := range ex {
		if exp.Name == funcN {
			dllOffset := uintptr(rvaToOffset(p, exp.VirtualAddress))
			b, _ := p.Bytes()
			buff = b[dllOffset : dllOffset+21]
		}
	}
	if buff == nil{
		return 0,memPtr,fmt.Errorf("not found in file")
	}

	//对比函数字节，判断是否被hook,若未被hook,返回原地址
	if buff[0] ==  *(*byte)(unsafe.Pointer(memPtr)) && buff[1] ==  *(*byte)(unsafe.Pointer(memPtr+1)){
		return memPtr,memPtr,nil
	}

	//构造unhook函数
	for i := 0;i < len(buff)-2;i++{
		if buff[i] == *(*byte)(unsafe.Pointer(memPtr+uintptr(i))) && buff[i+1] == *(*byte)(unsafe.Pointer(memPtr+uintptr(i+1))){
			addr := memPtr+uintptr(i)
			jmpAddr := uintptrToBytes(&addr)

			windows.WriteProcessMemory(0xffffffffffffffff,ptr,&buff[0],uintptr(i),nil)
			windows.WriteProcessMemory(0xffffffffffffffff,ptr+uintptr(i),&jmpPre[0],uintptr(2),nil)
			windows.WriteProcessMemory(0xffffffffffffffff,ptr+uintptr(i)+2,&jmpAddr[0],uintptr(len(jmpAddr)),nil)
			windows.WriteProcessMemory(0xffffffffffffffff,ptr+uintptr(i)+2+8,&jmpRet[0],uintptr(4),nil)

			return ptr,memPtr,nil
		}
	}

	return 0,memPtr,fmt.Errorf("last err")

}


const sizeOfUintPtr = unsafe.Sizeof(uintptr(0))

func uintptrToBytes(u *uintptr) []byte {
	return (*[sizeOfUintPtr]byte)(unsafe.Pointer(u))[:]
}


//rvaToOffset converts an RVA value from a PE file into the file offset. When using binject/debug, this should work fine even with in-memory files.
func rvaToOffset(pefile *pe.File, rva uint32) uint32 {
	for _, hdr := range pefile.Sections {
		baseoffset := uint64(rva)
		if baseoffset > uint64(hdr.VirtualAddress) &&
			baseoffset < uint64(hdr.VirtualAddress+hdr.VirtualSize) {
			return rva - hdr.VirtualAddress + hdr.Offset
		}
	}
	return rva
}
