// Copyright The OpenTelemetry Authors
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

package inspect

import (
	"encoding/json"
	"os"

	"github.com/go-logr/logr"
	"github.com/hashicorp/go-version"

	"go.opentelemetry.io/auto/internal/pkg/structfield"
)

// Cache is a cache of struct field offsets.
type Cache struct {
	log  logr.Logger
	data *structfield.Index
}

// NewCache loads struct field offsets from offsetFile and returns them as a
// new Cache.
func NewCache(l logr.Logger, offsetFile string) (*Cache, error) {
	c := newCache(l)

	f, err := os.Open(offsetFile)
	if err != nil {
		return c, err
	}
	defer f.Close()

	c.data = structfield.NewIndex()
	err = json.NewDecoder(f).Decode(&c.data)
	return c, err
}

func newCache(l logr.Logger) *Cache {
	return &Cache{log: l.WithName("cache")}
}

// GetOffset returns the cached offset key and true for the id at the specified
// version is found in the cache. If the cache does not contain a valid offset for the provided
// values, 0 and false are returned.
func (c *Cache) GetOffset(ver *version.Version, id structfield.ID) (structfield.OffsetKey, bool) {
	if c.data == nil {
		return structfield.OffsetKey{}, false
	}

	off, ok := c.data.GetOffset(id, ver)
	msg := "cache "
	if ok {
		msg += "hit"
	} else {
		msg += "miss"
	}
	c.log.V(1).Info(msg, "version", ver, "id", id)
	return off, ok
}
