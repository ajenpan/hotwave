package battle

import (
	"fmt"
	"strings"
	"sync"
)

type GameLogicCreator struct {
	Store sync.Map
}

func (c *GameLogicCreator) Add(name string, creator func() GameLogic) error {
	c.Store.Store(name, creator)
	return nil
}

func (c *GameLogicCreator) CreateLogic(name string) (GameLogic, error) {
	v, has := c.Store.Load(name)
	if !has {
		return nil, fmt.Errorf("game logic %s not found", name)
	}
	creator := v.(func() GameLogic)
	return creator(), nil
}

var LogicCreator = &GameLogicCreator{}

func RegisterGame(name, version string, creator func() GameLogic) error {
	return LogicCreator.Add(strings.Join([]string{name, version}, "-"), creator)
}
