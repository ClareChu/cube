package aggregate

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"hidevops.io/cube/console/pkg/command"
	"hidevops.io/hiboot-data/starter/etcd"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"time"
)

type CubeUpdateAggregate interface {
	Add(name string, update *command.CubeUpdate) error
	Get(name string) (update *command.CubeUpdate, err error)
	Delete(name string) (err error)
}

type CubeUpdate struct {
	BuildAggregate
	repository etcd.Repository
}

func init() {
	app.Register(NewCubeUpdateService)
}

func NewCubeUpdateService(repository etcd.Repository) CubeUpdateAggregate {
	return &CubeUpdate{
		repository: repository,
	}
}

func (s *CubeUpdate) Add(name string, update *command.CubeUpdate) error {
	updateBuf, _ := json.Marshal(update)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	res, err := s.repository.Put(ctx, name, string(updateBuf))
	cancel()
	if err != nil {
		log.Errorf("cube update get body err : %v", err)
		return err
	}
	log.Debug(res)
	return err
}

func (s *CubeUpdate) Get(name string) (update *command.CubeUpdate, err error) {
	update = &command.CubeUpdate{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := s.repository.Get(ctx, name)
	cancel()
	if err != nil {
		log.Debugf("failed to get data from etcd, err: %v", err)
		return nil, err
	}

	if resp.Count == 0 {
		return nil, errors.New("record not found")
	}

	if err = json.Unmarshal(resp.Kvs[0].Value, &update); err != nil {
		log.Debugf("failed to unmarshal data, err: %v", err)
		return nil, err
	}
	return
}

func (s *CubeUpdate) Delete(name string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = s.repository.Delete(ctx, name)
	cancel()
	return
}
