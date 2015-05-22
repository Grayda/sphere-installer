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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if len(r.URL.Path[1:]) == 0 {
		fmt.Println("No path specified. Defaulting to index.html")
		http.ServeFile(w, r, "html/index.html")
	} else if len(r.Form["deburl"]) > 0 {
		if isDeb(r.Form["deburl"][0]) == true {
			fmt.Println("Path is a Debian URL!")
			fmt.Println(r.Form["deburl"][0])
			installDEB(w, r.Form["deburl"][0])
		} else {
			fmt.Println("Not a .deb. Showing error")
			fmt.Fprintf(w, "<h1>Not a valid DEB file!</h1>")
		}

	} else {
		http.ServeFile(w, r, "html/"+r.URL.Path[1:])
	}

}

func isDeb(deb string) bool {
	return strings.Contains(deb, ".deb")
}

func getDeb(rURL string) string {
	pURL, _ := url.Parse(rURL)
	return path.Base(pURL.Path)[:4]
}

func installDEB(w http.ResponseWriter, deb string) {

	filename := getDeb(deb)
	out, err := os.Create("/tmp/" + filename)
	defer out.Close()
	if err != nil {
		fmt.Fprintf(w, "<h1>Error!</h1><p>Unable to write "+filename+" to /tmp/. Error was:")
		fmt.Println(err)
	}
	resp, err := http.Get(deb)
	if err != nil {
		fmt.Fprintf(w, "<h1>Error!</h1><p>Unable to Get "+deb+". Error was:")
		fmt.Println(err)
	}

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Fprintf(w, "<h1>Error!</h1><p>Unable to output "+filename+" to '/tmp/"+filename+"'. Error was:")
		fmt.Println(err)
	}

	runCommand("nservice", filename+" stop")
	runCommand("dpkg", "-i /tmp/installer.deb")
	runCommand("nservice", filename+" start")
}

func runCommand(cmd string, args string) {
	eCmd := exec.Command(cmd, args)
	eCmd.Run()
}
