package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
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
		Use:   "cp source destination",
		Short: "copies a local file or Storj object to another location locally or in Storj",
		RunE:  copyMain,
		Args:  cobra.ExactArgs(2),
	})

	RootCmd.AddCommand(&cobra.Command{
		Use:   "cp folder",
		Short: "creates a new bucket from a local folder and place all files from folder to the bucket",
		RunE:  makeBucketAndUpload,
		Args:  cobra.ExactArgs(1),
	})

	RootCmd.AddCommand(&cobra.Command{
		Use:   "setup",
		Short: "setups folder watcher",
		RunE:  watcherSetup,
		Args:  cobra.ExactArgs(0),
	})

	RootCmd.AddCommand(&cobra.Command{
		Use:   "share via Storj",
		Short: "generates link to a file and copy it to the clipboard",
		RunE:  generateLink,
		Args:  cobra.ExactArgs(1),
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

	uplinkExecutable, err := exec.LookPath("uplink")
	if err != nil {
		return err
	}

	cmdUplinkCopy := &exec.Cmd{
		Path:   uplinkExecutable,
		Args:   []string{uplinkExecutable, "cp", src.String(), dst.String()},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = cmdUplinkCopy.Run()
	if err != nil {
		return err
	}

	cmdUplinkShare := &exec.Cmd{
		Path:   uplinkExecutable,
		Args:   []string{uplinkExecutable, "share", dst.String() + "/" + src.Base(), "--readonly"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = cmdUplinkShare.Run()
	if err != nil {
		return err
	}

	return nil
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

	uplinkExecutable, err := exec.LookPath("uplink")
	if err != nil {
		return err
	}

	cmdMakeBucket := &exec.Cmd{
		Path:   uplinkExecutable,
		Args:   []string{uplinkExecutable, "mb", "sj://" + info.Name()},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = cmdMakeBucket.Run()
	if err != nil {
		return err
	}

	goExecutable, err := exec.LookPath("go")
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src.String())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i := 0; i < len(files); i++ {
		if files[i].IsDir() {
			fmt.Println("folders are being ignored")
			continue
		}
		srcPath := src.String() + "/" + files[i].Name()

		cmdGoRun := &exec.Cmd{
			Path:   goExecutable,
			// TODO: take path from config
			Args:   []string{goExecutable, "run", "/Users/vitalii/Work/storjbox/core/cmd/main.go", "cp", srcPath, "sj://" + info.Name()},
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}

		err = cmdGoRun.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func generateLink(_ *cobra.Command, args []string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("no object name to generate link")
	}

	name, err := fpath.New(args[0])
	if err != nil {
		return err
	}

	uplinkExecutable, err := exec.LookPath("uplink")
	if err != nil {
		return err
	}

	cmdUplinkShare := &exec.Cmd{
		Path:   uplinkExecutable,
		Args:   []string{uplinkExecutable, "share", "distConfigValue"+name.String(), "--readonly"},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	err = cmdUplinkShare.Run()
	if err != nil {
		return err
	}

	return nil
}

func watcherSetup(_ *cobra.Command, _ []string) (err error) {
	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
		return err
	}
	defer func() {
		err = watcher.Close()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println(event.String())
				if event.Op == 1 {
					goExecutable, _ := exec.LookPath("go")

					// TODO ignore .DS_store like files
					cmdGoRun := &exec.Cmd{
						Path:   goExecutable,
						Args:   []string{goExecutable, "run", "/Users/vitalii/Work/storjbox/core/cmd/main.go", "cp", event.Name, "sj://bucket"},
						Stdout: os.Stdout,
						Stderr: os.Stderr,
					}

					err = cmdGoRun.Run()
					if err != nil {
						fmt.Println("ERROR", err)
					}
				}

			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add("/Users/vitalii/stoprjbox"); err != nil {
		fmt.Println("ERROR", err)
		return err
	}

	<-done

	return nil
}
