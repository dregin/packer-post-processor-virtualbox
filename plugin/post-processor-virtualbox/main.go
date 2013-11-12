package main

import (
	"github.com/dregin/packer-post-processor-virtualbox/post-processor/virtualbox"
	"github.com/mitchellh/packer/packer/plugin"
)

func main() {
	plugin.ServePostProcessor(new(virtualbox.PostProcessor))
}
