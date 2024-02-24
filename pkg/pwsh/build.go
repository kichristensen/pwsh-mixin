package pwsh

import (
	"context"
	"fmt"
	"text/template"

	"get.porter.sh/porter/pkg/exec/builder"
	"gopkg.in/yaml.v2"
)

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the pwsh mixin in porter.yaml
// mixins:
// - pwsh:
//	  clientVersion: "v0.0.0"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`
}

const buildTemplate string = `RUN --mount=type=cache,target=/var/cache/apt --mount=type=cache,target=/var/lib/apt \
	apt-get update && apt-get install -y curl
{{- if eq .ClientVersion "" }}
RUN URL=$(curl -s https://api.github.com/repos/powershell/powershell/releases/latest | grep 'browser_download_url.*deb' | cut -d : -f 2,3 | tr -d \" | head -n 1) && \
	curl -L -o powershell.deb $URL
{{- else }}
RUN curl -L -o powershell.deb https://github.com/PowerShell/PowerShell/releases/download/v{{.ClientVersion}}/powershell_{{.ClientVersion}}-1.deb_amd64.deb
{{- end }}
RUN --mount=type=cache,target=/var/cache/apt --mount=type=cache,target=/var/lib/apt \
	dpkg -i powershell.deb || apt-get install -f -y 
RUN rm powershell.deb`

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build(ctx context.Context) error {

	// Create new Builder.
	var input BuildInput

	err := builder.LoadAction(ctx, m.RuntimeConfig, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	tmpl, err := template.New("dockerfile").Parse(buildTemplate)
	if err != nil {
		return fmt.Errorf("error parsing Dockerfile template for the az mixin: %w", err)
	}

	cfg := MixinConfig{ClientVersion: input.Config.ClientVersion}
	if err = tmpl.Execute(m.Out, cfg); err != nil {
		return fmt.Errorf("error generating Dockerfile lines for the az mixin: %w", err)
	}

	return nil
}
