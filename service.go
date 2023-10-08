package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func Update(url string, file string) {
	fmt.Printf("updated %s\n",file)
	resp,err := http.Get(url)
	if err != nil {
		panic(err)
	}

	out,err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	daemon, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	daemon.Truncate(0)
	daemon.WriteAt(out,0)
	daemon.Close()
}

func main() {
	// example: "https://gilded-frangollo-a960a9.netlify.app"
	git := os.Getenv("TM_PACKAGE_URL") 
	if git != "" {
		Update(fmt.Sprintf("%s/hub.exe"   , git),"./hub.exe")
		Update(fmt.Sprintf("%s/daemon.exe", git),"./daemon.exe")
	}

	proxy,err := os.OpenFile("./secret/proxy.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	proxy_cred,err := io.ReadAll(proxy)
	if err != nil {
		panic(err)
	}
	if len(proxy_cred) == 0{
		fmt.Printf("enter proxy account : ")
		var cred string
		fmt.Scanln(&cred)
		proxy.Truncate(0)
		proxy.WriteAt([]byte(cred),0)
	}
	proxy.Close()
	
	cmd := exec.Command("./daemon.exe")
	cmd.Stdout 	= os.Stdout
    cmd.Stderr 	= os.Stderr
	cmd.Stdin 	= os.Stdin
	cmd.Run()
}