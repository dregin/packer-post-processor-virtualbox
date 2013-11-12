package virtualbox

import (
	"fmt"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"os/exec"
	"strings"
	"io"
	"os"
)

// Keeping this to leave opportunity for VMWare and AWS Post-Processors
var builtins = map[string]string{
	"mitchellh.virtualbox": "virtualbox",
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	// Username for SCP operation.
	// SSH keys should be used for authentication.
	ScpUserName string `mapstructure:"scp_user_name"`

	// Path to private SSH Key
	ScpKeyPath string  `mapstructure:"scp_key_path"`

	// Path to which the exported VirtualBox image will be transferred.
	RemoteImageDir string `mapstructure:"remote_image_dir"`

	// The VirtualBox Host
	VirtualBoxHost string `mapstructure:"virtual_box_host"`

	// The Address of PHP Virtualbox
	PhpVirtualBoxAddress string `mapstructure:"php_virtualbox_address"`

	// The Admin User for PHP Virtualbox
	PhpVirtualBoxUser string `mapstructure:"php_virtualbox_user"`

	// The Admin Password for PHP Virtualbox
	PhpVirtualBoxPass string `mapstructure:"php_virtualbox_pass"`
}

type PostProcessor struct {
	config  Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	_, err := common.DecodeConfig(&p.config, raws...)
	if err != nil {
		return err
	}

	errors := new (packer.MultiError)

	_, err1 := exec.LookPath("scp")
	if err1 != nil{
		errors = packer.MultiErrorAppend(
			errors, fmt.Errorf("This tool depends on scp.", err1))
	}

	validates := map[string]*string{
		"scp_user_name":	&p.config.ScpUserName,
		"scp_key_path": &p.config.ScpKeyPath,
		"remote_image_dir": &p.config.RemoteImageDir,
		"virtual_box_host": &p.config.VirtualBoxHost,
		"php_virtualbox_address": &p.config.PhpVirtualBoxAddress,
		"php_virtualbox_user": &p.config.PhpVirtualBoxUser,
		"php_virtualbox_pass": &p.config.PhpVirtualBoxPass,
	}

	for n := range validates {
		if *validates[n] == "" {
			errors = packer.MultiErrorAppend(
			errors, fmt.Errorf("Argument not set: %s", n))
		}
	}

	if len(errors.Errors) > 0 {
		return errors
	}
	return nil
}

// Send the Virtual Box Image to the host.
func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	remoteImagePath := ""
	_, ok := builtins[artifact.BuilderId()]

	if !ok {
		return nil, false, fmt.Errorf("Unknown artifact type, can't build box: %s", artifact.BuilderId())
	}

	// Each Image comprises of a .ovf and a .vmdk file
	for _, fileName := range artifact.Files(){
		if strings.HasSuffix(fileName, ".ovf"){
			remoteImagePath = p.config.RemoteImageDir + fileName
		}
		ui.Message(fmt.Sprintf("The Virtualbox Post-Processor is uploading %s to the Virtualbox Host", fileName))
		cmd := exec.Command("scp", "-i", p.config.ScpKeyPath, fileName, p.config.ScpUserName + "@" + p.config.VirtualBoxHost + ":" + p.config.RemoteImageDir)

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

	// Fire off HTTP request to PHPVirtualBox
	//importImageViaWebService(remoteImagePath, p.config.PhpVirtualBoxAddress, p.config.PhpVirtualBoxUser, p.config.PhpVirtualBoxPass)

	// Run command line import over SSH
	importImageViaCommandLine(p.config.ScpKeyPath, p.config.ScpUserName, p.config.VirtualBoxHost, remoteImagePath)

	return artifact, false, nil
}