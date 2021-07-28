// Copyright 2021 Juan Pablo Tosso
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package actions

import (
	"fmt"

	engine "github.com/jptosso/coraza-waf/v1"
	"github.com/jptosso/coraza-waf/v1/utils"
)

type Exec struct {
	cachedScript string
}

func (a *Exec) Init(r *engine.Rule, data string) error {
	fdata, err := utils.OpenFile(data)
	if err != nil {
		return fmt.Errorf("Cannot load file %s", data)
	}
	a.cachedScript = string(fdata)
	return nil
}

func (a *Exec) Evaluate(r *engine.Rule, tx *engine.Transaction) {
	// Not implemented
}

func (a *Exec) GetType() int {
	return engine.ACTION_TYPE_NONDISRUPTIVE
}
