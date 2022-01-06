package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "path/filepath"
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

func get_distro() string {
    path, _ := os.Getwd()
    output, err := exec.Command("sudo", filepath.Join(path, "scripts/distro_script.sh")).Output()
    if err != nil {
        errorMessage := fmt.Sprintf("Error getting distro name: %+v", err)
		errorLogger.Println(errorMessage)
    }
    distroName := strings.ToLower(string(output))
    infoLogger.Println(fmt.Sprintf("distro name: %s", distroName))
    return distroName
}

func getSystemdUnitFileInstallPath() string {

    var systemdPath string

    distro := get_distro()

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