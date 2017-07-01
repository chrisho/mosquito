package utils

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetGOPATHs returns all paths in GOPATH variable.
func GetGoPATHs() []string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" && strings.Compare(runtime.Version(), "go1.8") >= 0 {
		gopath = defaultGOPATH()
	}
	return filepath.SplitList(gopath)
}

func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		return filepath.Join(home, "go")
	}
	return ""
}

// __FILE__ returns the file name in which the function was invoked
func FILE() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

// __LINE__ returns the line number at which the function was invoked
func LINE() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

// get computer macs
func GetLocalMacs() (macs []string) {
	interfaces, err := net.Interfaces()

	if err != nil {
		log.Println("Error : " + err.Error())
	}

	for _, inter := range interfaces {
		mac := inter.HardwareAddr.String()
		macs = append(macs, mac)
	}
	return
}

// get computer ips
func GetLocalIps() (ips []string) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Println("Error : " + err.Error())
	}

	for _, address := range addrs {

		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				ips = append(ips, ip)
			}
		}
	}
	return
}
