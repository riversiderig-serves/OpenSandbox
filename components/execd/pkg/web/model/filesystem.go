// Copyright 2025 The OpenSandbox Authors
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

package model

import (
	"encoding/json"
	"time"
)

type FileInfo struct {
	Path       string    `json:"path,omitempty"`
	Type       string    `json:"type,omitempty"`
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modified_at,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	Permission `json:",inline"`
}

type FileMetadata struct {
	Path       string `json:"path,omitempty"`
	Permission `json:",inline"`
	ModeSet    bool `json:"-"
}

// UnmarshalJSON preserves whether mode was supplied so an explicit zero
// permission can be distinguished from an omitted mode.
func (m *FileMetadata) UnmarshalJSON(data []byte) error {
	type metadataAlias FileMetadata
	var decoded metadataAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if rawMode, ok := fields["mode"]; ok && string(rawMode) != "null" {
		decoded.ModeSet = true
	}
	*m = FileMetadata(decoded)
	return nil
}

type Permission struct {
	Owner string `json:"owner"`
	Group string `json:"group"`
	Mode  int    `json:"mode"`
}

type RenameFileItem struct {
	Src  string `json:"src,omitempty"`
	Dest string `json:"dest,omitempty"`
}

type ReplaceFileContentItem struct {
	Old string `json:"old,omitempty"`
	New string `json:"new,omitempty"`
}

type ReplaceFileContentResult struct {
	ReplacedCount int `json:"replacedCount"`
}
