package battle

import (
	"fmt"
	"sync"
)

type GameLogicCreator struct {
	store sync.Map
}

func (c *GameLogicCreator) Store(name string, creator func() GameLogic) error {
	c.store.Store(name, creator)
	return nil
}

func (c *GameLogicCreator) CreateLogic(name, conf string) (GameLogic, error) {
	v, has := c.store.Load(name)
	if !has {
		return nil, fmt.Errorf("game logic %s not found", name)
	}
	creator := v.(func() GameLogic)
	return creator(), nil
}
