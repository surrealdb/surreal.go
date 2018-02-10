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
	"encoding/base64"
	"fmt"
)

type Config struct {
	NS   string
	DB   string
	User string
	Pass string
}

func (c *Config) parse() []string {

	var prots = []string{"json", "pack"}

	if len(c.NS) > 0 {
		prots = append(prots, fmt.Sprintf("ns-%s", c.NS))
	}

	if len(c.DB) > 0 {
		prots = append(prots, fmt.Sprintf("db-%s", c.DB))
	}

	if len(c.User) > 0 && len(c.Pass) > 0 {
		str := fmt.Sprintf("%s:%s", c.User, c.Pass)
		val := base64.StdEncoding.EncodeToString([]byte(str))
		prots = append(prots, fmt.Sprintf("auth-%s", val))
	}

	return prots

}
