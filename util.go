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

import "errors"

func recv(method string, res interface{}, err error) (interface{}, error) {

	if err != nil {
		return nil, err
	}

	if res != nil {

		arr, ok := res.([]interface{})
		if !ok {
			return nil, nil
		}

		switch method {

		case "Query":

			return res, nil

		default:

			// Check that we have at least one query response

			if len(arr) == 0 {
				break
			}

			// Check that the first query response is an object

			obj, ok := arr[0].(map[string]interface{})
			if !ok {
				break
			}

			// Check that the first query response has a status

			sta, ok := obj["status"].(string)
			if !ok {
				break
			}

			// Check the return status of the first query response

			switch {
			case sta == "OK":
				if val, ok := obj["result"]; ok {
					return val, nil
				}
			case sta != "OK":
				if val, ok := obj["detail"].(string); ok {
					return nil, errors.New(val)
				}
			}

		}

	}

	return nil, nil

}
