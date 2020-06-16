package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"storj.io/common/fpath"
)

var RootCmd = &cobra.Command{
	Use:   "uplink",
	Short: "The Storj client-side CLI",
	Args:  cobra.OnlyValidArgs,
}

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "mb SOURCE",
		Short: "Create a new bucket",
		RunE:  makeBucketAndUpload,
		Args:  cobra.ExactArgs(1),
	})
}

func makeBucketAndUpload(_ *cobra.Command, args []string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("no object specified for copy")
	}

	src, err := fpath.New(args[0])
	if err != nil {
		return err
	}
	return createAndUpload(src)
}

func main() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
	}
}


func createAndUpload(src fpath.FPath) (err error) {
	if !src.IsLocal() {
		return fmt.Errorf("source must be local path: %s", src)
	}

	info, err := os.Stat(src.String())
	if err != nil {
		return fmt.Errorf("unadble to get os.Stat")
	}
	if !info.IsDir() {
		return fmt.Errorf("you can't from %s", info)
	}

	uplinkExecutable, _ := exec.LookPath( "uplink" )

	cmdMakeBucket := &exec.Cmd {
		Path: uplinkExecutable,
		Args: []string{ uplinkExecutable, "mb", "sj://"+info.Name(), },
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	err = cmdMakeBucket.Run()
	if err != nil {
		return err
	}

	goExecutable, _ := exec.LookPath( "go" )
	files, err := ioutil.ReadDir(src.String())
	for i := 0; i < len(files); i++ {
		if files[i].IsDir() {
			fmt.Println("folders are being ignored")
			continue
		}
		srcPath := src.String()+"/"+files[i].Name()

		cmdGoRun := &exec.Cmd {
			Path: goExecutable,
			Args: []string{ goExecutable, "run", "/Users/nikolaisiedov/Workspace/storjbox/cmd/uplink/main.go", "cp", srcPath, "sj://"+info.Name(), },
			Stdout: os.Stdout,
			Stderr: os.Stdout,
		}

		err = cmdGoRun.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
