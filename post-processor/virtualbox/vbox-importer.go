package virtualbox

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func importImageViaWebService(imagePath string, vboxAddress string, vboxUser string, vboxPass string){
	req, err := http.NewRequest("GET", vboxAddress, nil)
	req.SetBasicAuth(vboxUser, vboxPass)
	res, err := http.DefaultClient.Do(req)

	fmt.Println("Status:%s", res)

	if err != nil {
		fmt.Sprintf("%s", err)
	}
}

func importImageViaCommandLine(scpKeyPath string, scpUserName string, virtualBoxHost string, remoteImagePath string){
	cmd := exec.Command("ssh", "-q", "-i", scpKeyPath, scpUserName + "@" + virtualBoxHost, "vboxmanage", "import", remoteImagePath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
            fmt.Println(err)
    }
    stderr, err := cmd.StderrPipe()
    if err != nil {
        fmt.Println(err)
    }
    err = cmd.Start()
    if err != nil {
        fmt.Println(err)
    }
    go io.Copy(os.Stdout, stdout) 
    go io.Copy(os.Stderr, stderr) 
    cmd.Wait()
}