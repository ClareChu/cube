package aggregate

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"hidevops.io/hiboot-data/starter/etcd"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/mio/console/pkg/command"
	"time"
)

type MioUpdateAggregate interface {
	Add(name string, update *command.MioUpdate) error
	Get(name string) (update *command.MioUpdate, err error)
	Delete(name string) (err error)
}

type MioUpdate struct {
	BuildAggregate
	repository etcd.Repository
}

func init() {
	app.Register(NewMioUpdateService)
}

func NewMioUpdateService(repository etcd.Repository) MioUpdateAggregate {
	return &MioUpdate{
		repository: repository,
	}
}

func (s *MioUpdate) Add(name string, update *command.MioUpdate) error {
	updateBuf, _ := json.Marshal(update)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	res, err := s.repository.Put(ctx, name, string(updateBuf))
	cancel()
	if err != nil {
		log.Errorf("mio update get body err : %v", err)
		return err
	}
	log.Debug(res)
	return err
}

func (s *MioUpdate) Get(name string) (update *command.MioUpdate, err error) {
	update = &command.MioUpdate{}
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

func (s *MioUpdate) Delete(name string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = s.repository.Delete(ctx, name)
	cancel()
	return
}
