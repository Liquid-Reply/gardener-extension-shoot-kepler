package kepler

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"

	yttcmd "carvel.dev/ytt/pkg/cmd/template"
	yttui "carvel.dev/ytt/pkg/cmd/ui"
	yttfiles "carvel.dev/ytt/pkg/files"
	"github.com/andybalholm/brotli"

	api "github.com/liquid-reply/gardener-extension-shoot-kepler/pkg/apis/config"
)

//go:embed kepler.yaml
var manifest string

//go:embed seed.yaml
var seedManifest string

func Render(config *api.Configuration, compress bool) ([]byte, error) {
	opts := yttcmd.NewOptions()
	noopUI := yttui.NewCustomWriterTTY(false, os.Stderr, os.Stderr)

	var files []*yttfiles.File
	files = append(files, templateAsFile("manifest.yaml", manifest))
	inputs := yttcmd.Input{Files: yttfiles.NewSortedFiles(files)}

	output := opts.RunWithFiles(inputs, noopUI)
	if output.Err != nil {
		return nil, output.Err
	}
	manifest, err := output.DocSet.AsBytes()
	if err != nil {
		return nil, err
	}
	if compress {
		var buf bytes.Buffer
		w := brotli.NewWriterV2(&buf, 7)
		if _, err := w.Write(manifest); err != nil {
			return nil, err
		}
		if err := w.Close(); err != nil {
			return nil, err
		}
		manifest = buf.Bytes()
	}
	return manifest, nil
}

func namespaceOverlay(namespace string) string {
	return fmt.Sprintf(`#@ load("@ytt:overlay", "overlay")

#@overlay/match by=overlay.all, expects="1+"
---
metadata:
  #@overlay/match missing_ok=True
  namespace: %q
`, namespace)
}

func RenderSeed(config *api.Configuration, shootNamespace string, compress bool) ([]byte, error) {
	opts := yttcmd.NewOptions()
	noopUI := yttui.NewCustomWriterTTY(false, os.Stderr, os.Stderr)

	var files []*yttfiles.File
	files = append(files, templateAsFile("manifest.yaml", seedManifest))
	files = append(files, templateAsFile("namespace.yaml", namespaceOverlay(shootNamespace)))
	inputs := yttcmd.Input{Files: yttfiles.NewSortedFiles(files)}

	output := opts.RunWithFiles(inputs, noopUI)
	if output.Err != nil {
		return nil, output.Err
	}
	manifest, err := output.DocSet.AsBytes()
	if err != nil {
		return nil, err
	}
	if compress {
		var buf bytes.Buffer
		w := brotli.NewWriterV2(&buf, 7)
		if _, err := w.Write(manifest); err != nil {
			return nil, err
		}
		if err := w.Close(); err != nil {
			return nil, err
		}
		manifest = buf.Bytes()
	}
	return manifest, nil
}

type noopWriter struct{}

func (w noopWriter) Write(data []byte) (int, error) { return len(data), nil }

func templateAsFile(name, tpl string) *yttfiles.File {
	file, err := yttfiles.NewFileFromSource(yttfiles.NewBytesSource(name, []byte(tpl)))
	if err != nil {
		// should not happen
		panic(err)
	}

	return file
}
