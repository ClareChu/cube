package dispatch

import "hidevops.io/hiboot/pkg/app"

type BuildPack struct {
	BuildPackInterface
	Status map[string]BuildPackAggregate
}

type BuildPackInterface interface {
	Add(key string, value BuildPackAggregate)
	Get(key string) BuildPackAggregate
}

func init() {
	app.Register(NewBuildPack)
}

func NewBuildPack() BuildPackInterface {
	return &BuildPack{
		Status: map[string]BuildPackAggregate{},
	}
}

func (b *BuildPack) Add(key string, value BuildPackAggregate) {
	b.Status[key] = value
}

func (b *BuildPack) Get(key string) BuildPackAggregate {
	for k, v := range b.Status {
		if key == k {
			return v
		}
	}
	return nil
}
