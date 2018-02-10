// Copyright © 2016 Abcum Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package surreal

import (
	"sync"

	"github.com/abcum/fibre"
)

type Socket struct {
	cli  *fibre.Client
	quit chan error
	send chan<- *fibre.RPCRequest
	recv <-chan *fibre.RPCResponse
	emit struct {
		lock sync.Mutex
		once map[interface{}][]func(error, interface{})
		when map[interface{}][]func(error, interface{})
	}
}

func NewSocket(url string, opts *Config) (*Socket, error) {

	c, err := fibre.NewClient(url, opts.parse())
	if err != nil {
		return nil, err
	}

	cli := &Socket{cli: c}

	cli.send, cli.recv, cli.quit = c.Rpc()

	return cli.init(), nil

}

func (c *Socket) init() *Socket {

	go func() {
		for {
			select {
			case <-c.quit:
				break
			case res := <-c.recv:
				switch {
				case res.Error == nil:
					c.done(res.ID, nil, res.Result)
				case res.Error != nil:
					c.done(res.ID, res.Error.Message, res.Result)
				}
			}
		}
	}()

	return c

}

func (c *Socket) once(id interface{}, fn func(error, interface{})) {

	c.emit.lock.Lock()
	defer c.emit.lock.Unlock()

	if c.emit.once == nil {
		c.emit.once = make(map[interface{}][]func(error, interface{}))
	}

	c.emit.once[id] = append(c.emit.once[id], fn)

}

func (c *Socket) when(id interface{}, fn func(error, interface{})) {

	c.emit.lock.Lock()
	defer c.emit.lock.Unlock()

	if c.emit.when == nil {
		c.emit.when = make(map[interface{}][]func(error, interface{}))
	}

	c.emit.when[id] = append(c.emit.when[id], fn)

}

func (c *Socket) done(id interface{}, err error, res interface{}) {

	c.emit.lock.Lock()
	defer c.emit.lock.Unlock()

	if c.emit.when != nil {
		if _, ok := c.emit.when[id]; ok {
			for i := len(c.emit.when[id]) - 1; i >= 0; i-- {
				c.emit.when[id][i](err, res)
			}
		}
	}

	if c.emit.once != nil {
		if _, ok := c.emit.once[id]; ok {
			for i := len(c.emit.once[id]) - 1; i >= 0; i-- {
				c.emit.once[id][i](err, res)
				c.emit.once[id][i] = nil
				c.emit.once[id] = c.emit.once[id][:i]
			}
		}
	}

}

func (c *Socket) Send(id string, method string, params []interface{}) {

	go func() {
		c.send <- &fibre.RPCRequest{
			ID:     id,
			Async:  true,
			Method: method,
			Params: params,
		}
	}()

}

func (c *Socket) Once(id, method string) (<-chan interface{}, <-chan error) {

	err := make(chan error)
	res := make(chan interface{})

	c.once(id, func(e error, r interface{}) {
		r, e = recv(method, r, e)
		switch {
		case e != nil:
			err <- e
			close(err)
			close(res)
		case e == nil:
			res <- r
			close(err)
			close(res)
		}
	})

	return res, err

}

func (c *Socket) When(id, method string) (<-chan interface{}, <-chan error) {

	err := make(chan error)
	res := make(chan interface{})

	c.when(id, func(e error, r interface{}) {
		r, e = recv(method, r, e)
		switch {
		case e != nil:
			err <- e
		case e == nil:
			res <- r
		}
	})

	return res, err

}
