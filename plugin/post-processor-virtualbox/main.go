import (
  "github.com/mitchellh/packer/plugin"
)

// Assume this implements packer.Builder
type PostProcessor interface {
    Configure(interface{}) error
    PostProcess(Ui, Artifact) (a Artifact, keep bool, err error)
}

func main() {
    plugin.ServePostProcessor(new(vagrant.PostProcessor))
}