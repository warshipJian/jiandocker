package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return
	}
	log.Infof("Current location is %s", pwd)

	// 设置根目录为private，实现unshare -m效果
	if err := syscall.Mount("/", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		log.Errorf("mount / private error ", err)
	}

	// pivot_root挂载
	pivotRoot()

	// 挂载proc
	procMount()

	// 挂载tmpfs
	tmpfsMount()
}

func pivotRoot() error {

	// 把root重新mount一次
	root := "/newroot"
	fmt.Println(root)
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	// 创建 rootfs/pivot_root 存储old_root
	pivotDir := filepath.Join(root, "pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		fmt.Println(err)
		os.Remove(pivotDir)
		return err
	}
	// pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/pivot_root
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		fmt.Println("pivot_root error: ", err)
		os.Remove(pivotDir)
		return fmt.Errorf("pivot_root %v", err)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		fmt.Println("chdir / %v", err)
		os.Remove(pivotDir)
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", "pivot_root")
	// umount rootfs/pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		fmt.Println("unmount pivot_root dir %v", err)
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	// 删除临时文件夹
	return os.Remove(pivotDir)
}

func procMount() {
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
}

func tmpfsMount() {
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}
