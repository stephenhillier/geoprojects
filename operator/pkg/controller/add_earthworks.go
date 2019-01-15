package controller

import (
	"github.com/stephenhillier/geoprojects/operator/pkg/controller/earthworks"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, earthworks.Add)
}
