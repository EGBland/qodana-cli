//go:build windows
// +build windows

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
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// prepareWinCmd fixes cmd, so it can be run with a space in the batch or args path,
// fix taken from here https://github.com/golang/go/issues/17149
func prepareWinCmd(args ...string) *exec.Cmd {
	var commandLine = strings.Join(args, " ")
	var comSpec = os.Getenv("COMSPEC")
	if comSpec == "" {
		comSpec = os.Getenv("SystemRoot") + "\\System32\\cmd.exe"
	}
	var cmd = exec.Command(comSpec)
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: "/C \"" + commandLine + "\""}
	return cmd
}
