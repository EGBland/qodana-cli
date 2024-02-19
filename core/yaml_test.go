/*
 * Copyright 2021-2024 JetBrains s.r.o.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"github.com/JetBrains/qodana-cli/v2024/platform"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func setupTestFile(fileName string, content string) {
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, fileName)

	// create a test file with provided content and filename
	file, err := os.Create(tempFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func TestLoadQodanaYaml(t *testing.T) {
	testCases := []struct {
		description string
		setup       func(name string)
		project     string
		filename    string
		expected    *platform.QodanaYaml
	}{
		{
			description: "file exists but is empty",
			setup: func(name string) {
				setupTestFile(name, "")
			},
			project:  os.TempDir(),
			filename: "empty.yaml",
			expected: &platform.QodanaYaml{},
		},
		{
			description: "file exists with valid content",
			setup: func(name string) {
				content := `version: 1.0`
				setupTestFile(name, content)
			},
			project:  os.TempDir(),
			filename: "valid.yaml",
			expected: &platform.QodanaYaml{
				Version: "1.0",
			},
		},
		{
			description: "file exists with .net section",
			setup: func(name string) {
				content := `version: 1.0
dotnet:
  project: test.csproj
  frameworks: "!netstandard2.0;!netstandard2.1"`
				setupTestFile(name, content)
			},
			project:  os.TempDir(),
			filename: "dotnet.yaml",
			expected: &platform.QodanaYaml{
				Version: "1.0",
				DotNet: platform.DotNet{
					Project:    "test.csproj",
					Frameworks: "!netstandard2.0;!netstandard2.1",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.setup(tc.filename)
			actual := platform.LoadQodanaYaml(tc.project, tc.filename)
			_ = os.Remove(filepath.Join(tc.project, tc.filename))
			assert.Equal(t, tc.expected, actual)
		})
	}
}
