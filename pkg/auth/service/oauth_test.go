package service

import (
	"hidevops.io/hiboot/pkg/log"
	"testing"
)

func TestClient(t *testing.T)  {

	var p *[3]int


	var p1 [3]**int


	p4 := [3]int{1, 2, 3}
	p = &p4
	log.Info("p: ", p)
	p2 := 3
	p5 := &p2
	p1[0] = &p5
	p1[1] = &p5
	p1[2] = &p5

	log.Info("p1: ", p1)

	for i, j := range p {
		log.Info("p ", i, j)
	}

	for i, j := range p1 {
		log.Info("p1 ", i, *j)
	}




















}


