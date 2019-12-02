package dispatch

import "hidevops.io/hiboot/pkg/app"

type PipelineFactory struct {
	PipelineFactoryInterface
	Status map[string]Aggregate
}

type PipelineFactoryInterface interface {
	Add(key string, value Aggregate)
	Get(key string) Aggregate
}

func init() {
	app.Register(NewPipelineFactory)
}

func NewPipelineFactory() PipelineFactoryInterface {
	return &PipelineFactory{
		Status: map[string]Aggregate{},
	}
}

func (p *PipelineFactory) Add(key string, value Aggregate) {
	p.Status[key] = value
}

func (p *PipelineFactory) Get(key string) Aggregate {
	for k, v := range p.Status {
		if key == k {
			return v
		}
	}
	return nil
}
