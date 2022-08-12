package main

import (
	"strings"
)

func GetFunc(funcname string) uintptr {
	var phModule uintptr
	//Get dll module BaseAddr
	phModule, _ = inMemLoads("ntdll.dll")

	//get dll exports
	ex := GetExport(phModule)

	for _, exp := range ex {
		if exp.Name == funcname{
			return uintptr(exp.VirtualAddress)
		}
	}
	return 0
}




//InMemLoads returns a map of loaded dll paths to current process offsets (aka images) in the current process. No syscalls are made.
func inMemLoads(modulename string) (uintptr, uintptr) {
	s, si, p := gMLO(0)
	start := p
	i := 1
	if strings.Contains(strings.ToLower(p), strings.ToLower(modulename)) {
		return s, si
	}
	for {
		s, si, p = gMLO(i)
		if p != "" {
			if strings.Contains(strings.ToLower(p), strings.ToLower(modulename)) {
				return s, si
			}
		}
		if p == start {
			break
		}
		i++
	}
	return 0, 0
}

//GetModuleLoadedOrder returns the start address of module located at i in the load order. This might be useful if there is a function you need that isn't in ntdll, or if some rude individual has loaded themselves before ntdll.
func gMLO(i int) (start uintptr, size uintptr, modulepath string) {
	var badstring *sstring
	start, size, badstring = getMLO(i)
	modulepath = badstring.String()
	return
}


//getModuleLoadedOrder returns the start address of module located at i in the load order. This might be useful if there is a function you need that isn't in ntdll, or if some rude individual has loaded themselves before ntdll.
func getMLO(i int) (start uintptr, size uintptr, modulepath *sstring)