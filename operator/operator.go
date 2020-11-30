package operator

import "hidevops.io/cube/operator/crd"

func Start() {
	initCRD, err := crd.InitCRD()
	if err != nil {
		panic(err)
	}
	initCRD.Run()
}
