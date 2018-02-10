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
	"github.com/rs/xid"
)

type DB struct {
	sock *Socket
}

func New(url string, opts *Config) (*DB, error) {
	sock, err := NewSocket(url, opts)
	if err != nil {
		return nil, err
	}
	return &DB{sock: sock}, nil
}

func (db *DB) Close() {
	db.sock.cli.Close()
}

func (db *DB) Info() (interface{}, error) {
	return db.send("Info")
}

func (db *DB) Auth(to string) (interface{}, error) {
	return db.send("Auth", to)
}

func (db *DB) Live(tb string) (interface{}, error) {
	return db.send("Live", tb)
}

func (db *DB) Kill(id string) (interface{}, error) {
	return db.send("Kill", id)
}

func (db *DB) Query(sql string, vars map[string]interface{}) (interface{}, error) {
	return db.send("Query", sql, vars)
}

func (db *DB) Select(tb string, id interface{}) (interface{}, error) {
	return db.send("Select", tb, id)
}

func (db *DB) Create(tb string, id interface{}, data map[string]interface{}) (interface{}, error) {
	return db.send("Create", tb, id, data)
}

func (db *DB) Update(tb string, id interface{}, data map[string]interface{}) (interface{}, error) {
	return db.send("Update", tb, id, data)
}

func (db *DB) Change(tb string, id interface{}, data map[string]interface{}) (interface{}, error) {
	return db.send("Change", tb, id, data)
}

func (db *DB) Modify(tb string, id interface{}, data map[string]interface{}) (interface{}, error) {
	return db.send("Modify", tb, id, data)
}

func (db *DB) Delete(tb string, id interface{}) (interface{}, error) {
	return db.send("Delete", tb, id)
}

func (db *DB) send(method string, params ...interface{}) (interface{}, error) {

	id := xid.New().String()

	chn, err := db.sock.Once(id, method)

	db.sock.Send(id, method, params)

	for {
		select {
		default:
		case e := <-err:
			return nil, e
		case r := <-chn:
			return r, nil
		}
	}

}
