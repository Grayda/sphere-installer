package main

import (
	"fmt"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/support"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
)

var password = "temppwd"

var info = ninja.LoadModuleInfo("./package.json")

var installerConfig InstallerConfig

type InstallerConfig struct {
	Port int
}

func init() {

}

type RuntimeConfig struct {
}

type App struct {
	support.AppSupport
	config *InstallerConfig // This is how we save and load IR codes and such. Call this by using driver.config
	conn   *ninja.Connection
}

func (a *App) Start(cfg *RuntimeConfig) error {

	// This tells the API that we're going to expose a UI, and to run GetActions() in configuration.go
	a.Conn.MustExportService(&configService{&installerConfig}, "$driver/"+info.ID+"/configure", &model.ServiceAnnouncement{
		Schema: "/protocol/configuration",
	})

	return nil
}

// Stop the security light app.
func (a *App) Stop() error {

	return nil
}

func isDeb(deb string) bool {
	return strings.Contains(deb, ".deb")
}

func getDeb(rURL string) string {
	pURL, _ := url.Parse(rURL)
	return path.Base(pURL.Path)
}

func installDEB(deb string) error {

	filename := getDeb(deb)

	out, err := os.Create("/tmp/" + filename)
	defer out.Close()
	if err != nil {
		fmt.Println("Unable to create temp file. Error was: ", err, ". Not continuing")
		return err
	}
	resp, err := http.Get(deb)
	if err != nil {
		fmt.Println("Unable to download ", deb, ". Error was: ", err, ". Not continuing")
		return err
	}

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Unable to save temp file. Error was: ", err, ". Not continuing")
		return err

	}

	err = runCommand("/bin/sh", "-c", "echo "+password+" | sudo -S with-rw dpkg -i /tmp/"+filename)
	os.Remove("/tmp/" + filename)
	if err != nil {
		fmt.Println("Cannot run command. Error was", err)
		return err
	}

	return nil
}

func runCommand(cmd string, args1 string, args2 string) error {
	fmt.Println("=========================== Running", cmd, "with args", args1, args2)
	eCmd := exec.Command(cmd, args1, args2)
	return eCmd.Run()
}
