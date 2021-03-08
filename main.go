package main

import (
	"context"
	"fmt"
	image "github.com/containers/image/v5/image"
	alltransports "github.com/containers/image/v5/transports/alltransports"
	types "github.com/containers/image/v5/types"
	inspect "github.com/containers/skopeo/cmd/skopeo/inspect"
	"github.com/sirupsen/logrus"
	"time"
)



type optionalBool struct {
	present bool
	value   bool
}

type optionalString struct {
	present bool
	value   string
}

type globalOptions struct {
	debug              bool          // Enable debug output
	tlsVerify          optionalBool  // Require HTTPS and verify certificates (for docker: and docker-daemon:)
	policyPath         string        // Path to a signature verification policy file
	insecurePolicy     bool          // Use an "allow everything" signature verification policy
	registriesDirPath  string        // Path to a "registries.d" registry configuration directory
	overrideArch       string        // Architecture to use for choosing images, instead of the runtime one
	overrideOS         string        // OS to use for choosing images, instead of the runtime one
	overrideVariant    string        // Architecture variant to use for choosing images, instead of the runtime one
	commandTimeout     time.Duration // Timeout for the command execution
	registriesConfPath string        // Path to the "registries.conf" file
	tmpDir             string        // Path to use for big temporary files
}

type dockerImageOptions struct {
	global         *globalOptions      // May be shared across several imageOptions instances.
	authFilePath   optionalString      // Path to a */containers/auth.json (prefixed version to override shared image option).
	credsOption    optionalString      // username[:password] for accessing a registry
	registryToken  optionalString      // token to be used directly as a Bearer token when accessing the registry
	dockerCertPath string              // A directory using Docker-like *.{crt,cert,key} files for connecting to a registry or a daemon
	tlsVerify      optionalBool        // Require HTTPS and verify certificates (for docker: and docker-daemon:)
	noCreds        bool                // Access the registry anonymously
}



type imageOptions struct {
	dockerImageOptions
	sharedBlobDir    string // A directory to use for OCI blobs, shared across repositories
	dockerDaemonHost string // docker-daemon: host to connect to
}

func main() {

	ctx := context.Background()

	//imageName := "docker://library/ubuntu:20.04"
	imageName := "docker://quay.io/coreos/kube-state-metrics:v1.9.7"

	sys := &types.SystemContext{}
	var imgInspect *types.ImageInspectInfo
	ref, _ := alltransports.ParseImageName(imageName)
	src, _ := ref.NewImageSource(ctx, sys)

	img, err := image.FromUnparsedImage(ctx, sys, image.UnparsedInstance(src, nil))

	if err != nil {
		logrus.Error(err)
	}
	//manifestData, _, _ := img.Manifest(ctx)
	//fmt.Println(string(manifestData))
	//fmt.Println(err)

	imgInspect, err = img.Inspect(ctx)
	outputData := inspect.Output{
		Name: "", // Set below if DockerReference() is known
		Tag:  imgInspect.Tag,
		// Digest is set below.
		RepoTags:      []string{}, // Possibly overridden for docker.Transport.
		Created:       imgInspect.Created,
		DockerVersion: imgInspect.DockerVersion,
		Labels:        imgInspect.Labels,
		Architecture:  imgInspect.Architecture,
		Os:            imgInspect.Os,
		Layers:        imgInspect.Layers,
		Env:           imgInspect.Env,
	}
	for key, element := range outputData.Labels {
		fmt.Println("Key:", key, "=>", "Element:", element)
	}
}