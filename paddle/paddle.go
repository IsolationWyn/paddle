/**
 * @author [wyn]
 * @email [isolationwyn@gmail.com]
 * @create date 2018-11-19 15:48:25
 * @modify date 2018-11-19 15:48:25
 * @desc [description]
*/
package main

import (
	"log"
	"os"
	"syscall"
	"os/exec"
)

func main() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}