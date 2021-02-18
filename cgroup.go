package main

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

const cgroup = "jiandocker"
const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func MemoryLimit(pid int,limit string)  {
	os.Mkdir(path.Join(cgroupMemoryHierarchyMount, cgroup), 0755)
	ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, cgroup, "tasks"),[]byte(strconv.Itoa(pid)), 0644)
	ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, cgroup, "memory.limit_in_bytes"),[]byte(limit), 0644)
}
