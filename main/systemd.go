package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "path/filepath"
    "errors"
)

func contains(distros []string, distro string) bool {
    if distro != "" {
        for _, v := range distros {
            if v == distro {
                return true
            }
        }
    }
    return false
}

func getDistro() (string, error) {
    path, _ := os.Getwd()
    output, err := exec.Command("sudo", filepath.Join(path, "scripts/distro_script.sh")).Output()
    if err != nil {
        return "", errors.Wrapf(err, "error running shell script")
    }
    distroName := strings.ToLower(string(output))
    return distroName, nil
}

func GetSystemdUnitFileInstallPath() string {

    var systemdPath string

    distro, err := getDistro()

    if (err != nil) {
        errorMessage := fmt.Sprintf("Error getting distro name: %+v : %s", err, err.Error())
	    errorLogger.Println(errorMessage)
    } else {
        infoLogger.Println(fmt.Sprintf("distro name: %s", distro))
    }

    distros1 :=[]string{"ubuntu", "debian"}
    distros2 :=[]string{"suse", "sle_hpc", "sles", "opensuse", "redhat", "rhel", "centos", "oracle"}

    if contains(distros1, distro) {
        systemdPath = "/lib/systemd/system"
    } else if contains(distros2, distro) {
        systemdPath = "/usr/lib/systemd/system"
    } else {
        systemdPath = "/lib/systemd/system"
    }

    return systemdPath
}