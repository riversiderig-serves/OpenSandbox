// Copyright 2026 The OpenSandbox Authors
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

package controller

import (
	"encoding/json"
	"os"
	"testing"
	"strconv"

	"github.com/alibaba/opensandbox/execd/pkg/web/model"
)

func TestFileMetadataTracksModePresence(t *testing.T) {
	tests := []struct {
		name string
		json string
		want bool
		mode int
	}{
		{name: "omitted", json: `{"path":"/tmp/file"}`, want: false, mode: 0},
		{name: "explicit zero", json: `{"path":"/tmp/file","mode":0}`, want: true, mode: 0},
		{name: "nonzero", json: `{"path":"/tmp/file","mode":755}`, want: true, mode: 755},
		{name: "null treated as omitted", json: `{"path":"/tmp/file","mode":null}`, want: false, mode: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var metadata model.FileMetadata
			if err := json.Unmarshal([]byte(tt.json), &metadata); err != nil {
				t.Fatalf("unmarshal metadata: %v", err)
			}
			if metadata.ModeSet != tt.want {
				t.Fatalf("ModeSet = %v, want %v", metadata.ModeSet, tt.want)
			}
			if metadata.Mode != tt.mode {
				t.Fatalf("Mode = %d, want %d", metadata.Mode, tt.mode)
			}
		})
	}
}

func TestParseUploadMode(t *testing.T) {
	tests := []struct {
		value int
		want  os.FileMode
	}{
		{value: 0, want: 0},
		{value: 644, want: 0644},
		{value: 755, want: 0755},
		{value: 4755, want: 0755 | os.ModeSetuid},
		{value: 2755, want: 0755 | os.ModeSetgid},
		{value: 1777, want: 0777 | os.ModeSticky},
		{value: 7777, want: 0777 | os.ModeSetuid | os.ModeSetgid | os.ModeSticky},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.value), func(t *testing.T) {
			got, err := parseUploadMode(tt.value)
			if err != nil {
				t.Fatalf("parseUploadMode(%d): %v", tt.value, err)
			}
			if got != tt.want {
				t.Errorf("parseUploadMode(%d) = %#o, want %#o", tt.value, got, tt.want)
			}
		})
	}
}

func TestParseUploadModeRejectsInvalidValues(t *testing.T) {
	for _, value := range []int{-1, 8, 99, 10000, 100644} {
		if _, err := parseUploadMode(value); err == nil {
			t.Errorf("parseUploadMode(%d) unexpectedly succeeded", value)
		}
	}
}
