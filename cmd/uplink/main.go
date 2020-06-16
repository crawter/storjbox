package main

import (
	"fmt"
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
		Use:   "cp SOURCE DESTINATION",
		Short: "Copies a local file or Storj object to another location locally or in Storj",
		RunE:  copyMain,
		Args:  cobra.ExactArgs(2),
	})
}

func main() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
	}
}

// copyMain is the function executed when cpCmd is called.
func copyMain(_ *cobra.Command, args []string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("no object specified for copy")
	}
	if len(args) == 1 {
		return fmt.Errorf("no destination specified")
	}

	src, err := fpath.New(args[0])
	if err != nil {
		return err
	}

	dst, err := fpath.New(args[1])
	if err != nil {
		return err
	}

	return upload(src, dst)
}

func upload(src fpath.FPath, dst fpath.FPath) (err error) {
	if !src.IsLocal() {
		return fmt.Errorf("source must be local path: %s", src)
	}

	if dst.IsLocal() {
		return fmt.Errorf("destination must be Storj URL: %s", dst)
	}

	uplinkExecutable, _ := exec.LookPath( "uplink" )

	cmdUplinkCopy := &exec.Cmd {
		Path: uplinkExecutable,
		Args: []string{ uplinkExecutable, "cp", src.String(), dst.String(), },
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	err = cmdUplinkCopy.Run()
	if err != nil {
		return err
	}

	cmdUplinkShare := &exec.Cmd {
		Path: uplinkExecutable,
		Args: []string{ uplinkExecutable, "share", dst.String()+"/"+src.Base(), "--readonly", },
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	err = cmdUplinkShare.Run()
	if err != nil {
		return err
	}

	return nil
}
