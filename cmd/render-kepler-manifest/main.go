package main

import (
	"fmt"

	"github.com/liquid-reply/gardener-extension-shoot-kepler/kepler"
	apisconfig "github.com/liquid-reply/gardener-extension-shoot-kepler/pkg/apis/config"
)

func main() {
	out, err := kepler.Render(&apisconfig.Configuration{}, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
