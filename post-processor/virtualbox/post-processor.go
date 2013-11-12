package virtualbox

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"os/exec"
	"strings"
)

// Keeping this to leave opportunity for VMWare and AWS Post-Processors
var builtins = map[string]string{
	"mitchellh.virtualbox": "virtualbox",
}

type Config struct {
	// Username for SCP operation.
	// SSH keys should be used for authentication.
	scpUserName string `mapstructure:"scp_user_name"`

	// Path to private SSH Key
	scpKeyPath string  `mapstructure:"scp_key_path"`

	// Path to which the exported VirtualBox image will be transferred.
	remoteImagePath string `mapstructure:"remote_image_path"`

	// The VirtualBox Host
	virtualBoxHost string `mapstructure:"virtual_box_host"`

	// The Address of PHP Virtualbox
	phpVirtualBoxAddress string `mapstructure:"php_virtualbox_address"`

	// The Admin User for PHP Virtualbox
	phpVirtualBoxUser string `mapstructure:"php_virtualbox_user"`

	// The Admin Password for PHP Virtualbox
	phpVirtualBoxPass string `mapstructure:"php_virtualbox_pass"`
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
		"scp_user_name":	&p.config.scpUserName,
		"scp_key_path": &p.config.scpKeyPath,
		"remote_image_path": &p.config.remoteImagePath,
		"virtual_box_host": &p.config.virtualBoxHost,
		"php_virtualbox_address": &p.config.phpVirtualBoxAddress,
		"php_virtualbox_user": &p.config.phpVirtualBoxUser,
		"php_virtualbox_pass": &p.config.phpVirtualBoxPass,
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
	//remoteImagePath := ""
	_, ok := builtins[artifact.BuilderId()]

	if !ok {
		return nil, false, fmt.Errorf("Unknown artifact type, can't build box: %s", artifact.BuilderId())
	}

	// Each Image comprises of a .ovf and a .vmdk file
	for _, fileName := range artifact.Files(){
		if strings.HasSuffix(fileName, ".ovf"){
			//remoteImagePath = p.config.remoteImagePath + fileName
		}
		ui.Message(fmt.Sprintf("The Virtualbox Post-Processor is uploading %s to the Virtualbox Host", fileName))
		cmd := exec.Command("scp", "-i", p.config.scpKeyPath, fileName, p.config.scpUserName + "@" + p.config.virtualBoxHost + ":" + p.config.remoteImagePath)

		var out bytes.Buffer
		cmd.Stdout = &out
		err_run := cmd.Run()
		if err_run != nil {
			return nil, false, fmt.Errorf("%s", out.String())
		}
		ui.Message(fmt.Sprintf("%s", out.String()))
	}

	// Fire off HTTP request to PHPVirtualBox
	//importImageViaWebService(remoteImagePath, p.config.phpVirtualBoxAddress, p.config.phpVirtualBoxUser, p.config.phpVirtualBoxPass)

	// Run command line import over SSH
	importImageViaCommandLine(p.config.scpKeyPath, p.config.scpUserName, p.config.virtualBoxHost, p.config.remoteImagePath)

	return artifact, false, nil
}