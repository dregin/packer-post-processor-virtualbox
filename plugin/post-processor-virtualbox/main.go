package main

import (
	"github.com/dregin/packer-post-processor-virtualbox/post-processor/virtualbox"
	"github.com/mitchellh/packer/packer/plugin"
)

// Assume this implements packer.Builder
type PostProcessor interface {
	Configure(interface{}) error
	PostProcess(Ui, Artifact) (a Artifact, keep bool, err error)
}

func main() {
	plugin.ServePostProcessor(new(virtualbox.PostProcessor))
}
