package virtualbox

import (
	"fmt"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"os"
	"path/filepath"
)

// Keeping this to leave opportunity for VMWare and AWS Post-Processors
var builtins = map[string]string{
	"dregin.virtualbox": "virtualbox",
}

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
	config  Config
	premade map[string]packer.PostProcessor
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	_, err := common.DecodeConfig(&p.config, raws...)
	if err != nil {
		return err
	}

	tpl, err := packer.NewConfigTemplate()
	if err != nil {
		return err
	}
	return nil
}

// Send the OVF file (The artifact) to the virtual box host.
func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	ppName, ok := builtins[artifact.BuilderId()]
	if !ok {
		return nil, false, fmt.Errorf("Unknown artifact type, can't build box: %s", artifact.BuilderId())
	}

	// Use the premade PostProcessor if we have one. Otherwise, we
	// create it and configure it here.
	pp, ok := p.premade[ppName]
	return pp.PostProcess(ui, artifact)
}

func (p *PostProcessor) subPostProcessor(key string, specific interface{}, extra map[string]interface{}) (packer.PostProcessor, error) {
	pp := keyToPostProcessor(key)
	if pp == nil {
		return nil, nil
	}

	if err := pp.Configure(extra, specific); err != nil {
		return nil, err
	}

	return pp, nil
}

// keyToPostProcessor maps a configuration key to the actual post-processor
// it will be configuring. This returns a new instance of that post-processor.
func keyToPostProcessor(key string) packer.PostProcessor {
	switch key {
	case "virtualbox":
		return new(VBoxBoxPostProcessor)
	//case "aws":
	//	return new(AWSBoxPostProcessor)
	//case "vmware":
	//	return new(VMwareBoxPostProcessor)
	default:
		return nil
	}
}
