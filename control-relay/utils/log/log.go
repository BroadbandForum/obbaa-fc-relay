/*
 * Copyright 2020 Broadband Forum
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
* Control Relay log file
*
* Created by Filipe Claudio(Altice Labs) on 01/09/2020
*/ 

package log

import (
	"os"
	
	//"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var LogTerminal = logrus.New()
//var LogFile 	= logrus.New()

func init() {
	
	LogTerminal.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		FullTimestamp: true,
	})

	LogTerminal.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	LogTerminal.SetLevel(logrus.InfoLevel)

/* 	pathMap := lfshook.PathMap{
		logrus.ErrorLevel: "/var/log/error.log",
		logrus.FatalLevel: "/var/log/fatal.log",
	}

	LogFile.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		},
	)) */
}

func Debug(msg...interface{}) {
	LogTerminal.Debug(msg)
} 

func Info(msg...interface{}) {
	LogTerminal.Info(msg)
} 

func Warning(msg...interface{}) {
	LogTerminal.Warn(msg)
}

func Error(msg...interface{}) {
	LogTerminal.Error(msg)
	//LogFile.Error(msg)
}

func Fatal(msg...interface{}) {
	LogTerminal.Fatal(msg)
	//LogFile.Fatal(msg)
}