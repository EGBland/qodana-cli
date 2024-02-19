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

package platform

import (
	"bufio"
	"fmt"
	"github.com/JetBrains/qodana-cli/v2024/cloud"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	nugetConfigName           = "nuget.config"
	nugetConfigNamePascalCase = "NuGet.Config"
)

func UnsetNugetVariables() {
	variables := []string{qodanaNugetUser, qodanaNugetPassword, qodanaNugetName, qodanaNugetUrl}
	for _, variable := range variables {
		if err := os.Unsetenv(variable); err != nil {
			log.Fatal("couldn't unset env variable ", err.Error())
		}
	}
}

func WarnIfPrivateFeedDetected(prodCode string, projectPath string) {
	if prodCode != QDNET && prodCode != QDNETC || qodanaNugetVarsSet() {
		return
	}
	configFileNames := []string{nugetConfigName, nugetConfigNamePascalCase}
	for _, fileName := range configFileNames {
		if _, err := os.Stat(filepath.Join(projectPath, fileName)); err == nil {
			if checkForPrivateFeed(filepath.Join(projectPath, fileName)) {
				_, _ = fmt.Fprintf(os.Stderr, "\nWarning: private NuGet feed detected. Please set %s, %s, %s and %s (optional) environment variables to provide credentials for the private feed.\n",
					qodanaNugetUser, qodanaNugetPassword, qodanaNugetUrl, qodanaNugetName)
			}
		}
	}
}

func checkForPrivateFeed(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Failed to open file:", fileName)
		return false
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("couldn't close file ", err.Error())
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "<add ") {
			match, _ := regexp.MatchString(`http(s)?://`, line)
			return match
		}
	}
	return false
}

func IsNugetConfigNeeded() bool {
	return IsContainer() && qodanaNugetVarsSet()
}

func qodanaNugetVarsSet() bool {
	return os.Getenv(qodanaNugetUrl) != "" && os.Getenv(qodanaNugetUser) != "" && os.Getenv(qodanaNugetPassword) != ""
}

func PrepareNugetConfig(userPath string) {
	nugetConfig := filepath.Join(userPath, ".nuget", "NuGet")
	if _, err := os.Stat(nugetConfig); err != nil {
		// mkdir -p ~/.nuget/NuGet
		if err := os.MkdirAll(nugetConfig, os.ModePerm); err != nil {
			log.Fatal("couldn't create a directory ", err.Error())
		}
	}
	nugetConfig = filepath.Join(nugetConfig, "NuGet.Config")
	config := nugetWithPrivateFeed(cloud.GetEnvWithDefault(qodanaNugetName, "qodana"), os.Getenv(qodanaNugetUrl), os.Getenv(qodanaNugetUser), os.Getenv(qodanaNugetPassword))
	if err := os.WriteFile(nugetConfig, []byte(config), 0644); err != nil {
		log.Fatal("couldn't create a file ", err.Error())
	}
}

func nugetWithPrivateFeed(nugetSourceName string, nugetUrl string, nugetUser string, nugetPassword string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <clear />
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
    <add key="%s" value="%s" />
  </packageSources>
  <packageSourceCredentials>
    <%s>
      <add key="Username" value="%s" />
      <add key="ClearTextPassword" value="%s" />
    </%s>
  </packageSourceCredentials>
</configuration>`,
		nugetSourceName,
		nugetUrl,
		nugetSourceName,
		nugetUser,
		nugetPassword,
		nugetSourceName)
}
