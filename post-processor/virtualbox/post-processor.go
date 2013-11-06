package virtualbox

import (
	"fmt"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"os"
	"path/filepath"
)

type Config struct {
	// Username for SCP operation.
	// SSH keys should be used for authentication.
	scpUserName string `mapstructure:"scp_user_name"`

	// Path to which the exported VirtualBox image will be transferred.
	remoteOVFPath string `mapstructure:"remote_ovf_path"`

	// The VirtualBox Host
	virtualBoxHost string `mapstructure:"virtual_box_host"`
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	_, err := common.DecodeConfig(&p.config, raws...)
	if err != nill {
		return err
	}

	tpl, err := packer.NewConfigTemplate()
	if err != nil {
		return err
	}
}

// Send the OVF file (The artifact) to the virtual box host.
func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {

}

func (p *PostProcessor) subPostProcessor(key string, specific interface{}, extra map[string]interface{}) (packer.PostProcessor, error) {

}

// keyToPostProcessor maps a configuration key to the actual post-processor
// it will be configuring. This returns a new instance of that post-processor.
func keyToPostProcessor(key string) packer.PostProcessor {
	switch key {
	case "aws":
		return new(AWSBoxPostProcessor)
	case "virtualbox":
		return new(VBoxBoxPostProcessor)
	case "vmware":
		return new(VMwareBoxPostProcessor)
	default:
		return nil
	}
}
