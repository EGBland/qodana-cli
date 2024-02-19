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
	"os"
	"path/filepath"
)

type QodanaOptions struct {
	*platform.QodanaOptions
}

func (o *QodanaOptions) fixesSupported() bool {
	return o.guessProduct() != platform.QDNET && o.guessProduct() != platform.QDNETC && o.guessProduct() != platform.QDCL
}

func (o *QodanaOptions) appInfoXmlPath(ideBinDir string) string {
	appInfoPath := filepath.Join(ideBinDir, qodanaAppInfoFilename)
	if _, err := os.Stat(appInfoPath); err != nil && o.AnalysisId != "FAKE" {
		log.Fatalf("%s should exist in IDE directory %s. Unsupported IDE detected, exiting.", qodanaAppInfoFilename, ideBinDir)
	}
	return appInfoPath
}

func (o *QodanaOptions) vmOptionsPath() string {
	return filepath.Join(o.ConfDirPath(), "ide.vmoptions")
}
