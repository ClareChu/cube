package operator

import (
	"hidevops.io/cube/operator/cr"
	"hidevops.io/cube/operator/crd"
)

func Start() {
	initCRD, err := crd.InitCRD()
	if err != nil {
		panic(err)
	}
	initCRD.Run()

	cube, err := cr.InitCube()
	if err != nil {
		panic(err)
	}
	cube.Run()
}
