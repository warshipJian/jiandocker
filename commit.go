package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func commitContainer(imageName string) {
	mntURL := "/newroot"
	imageTar := "/tmp/" + imageName + ".tar"
	fmt.Printf("%s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error %v", mntURL, err)
	} else {
		log.Infof("commit success")
	}
}
