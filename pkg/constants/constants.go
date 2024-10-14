package constants

const (
	// ExtensionType is the name of the extension type.
	ExtensionType = "shoot-kepler"
	// ServiceName is the name of the service.
	ServiceName = ExtensionType

	extensionServiceName = "extension-" + ServiceName

	ManagedResourceNameKeplerConfig     = extensionServiceName + "-config"
	ManagedResourceNameKeplerSeedConfig = extensionServiceName + "-seed-config"
)
