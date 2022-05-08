package traversal

import "github.com/rytrose/pixelsound/api"

// TraverseFuncs contains all of the traversal functions.
var TraverseFuncs = map[string]api.TraverseFunc{
	"Random":   Random,
	"TtoBLtoR": TtoBLtoR,
}
