package protocol

import "go/types"

// DriverResponse is the GOPACKAGESDRIVER response object.
type DriverResponse struct {
	// NotHandled is returned if the request can't be handled by the current
	// driver. If an external driver returns a response with NotHandled, the
	// rest of the driverResponse is ignored, and go/packages will fallback
	// to the next driver. If go/packages is extended in the future to support
	// lists of multiple drivers, go/packages will fall back to the next driver.
	NotHandled bool

	// Sizes, if not nil, is the types.Sizes to use when type checking.
	Sizes *types.StdSizes

	// Roots is the set of package IDs that make up the root packages.
	// We have to encode this separately because when we encode a single package
	// we cannot know if it is one of the roots as that requires knowledge of the
	// graph it is part of.
	Roots []string `json:",omitempty"`

	// Packages is the full set of packages in the graph.
	// The packages are not connected into a graph.
	// The Imports if populated will be stubs that only have their ID set.
	// Imports will be connected and then type and syntax information added in a
	// later pass (see refine).
	Packages []*FlatPackage
}
