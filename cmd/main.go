/**
 * Copyright 2025 OceanBase
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

package main

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/oceanbase/configserver/config"
	"github.com/oceanbase/configserver/logger"
	"github.com/oceanbase/configserver/server"
)

var (
	configserverCommand = &cobra.Command{
		Use:   "configserver",
		Short: "configserver is used to store and query ob rs_list",
		Long:  "configserver is used to store and query ob rs_list, used by observer, obproxy and other tools",
		Run: func(cmd *cobra.Command, args []string) {
			err := runConfigServer()
			if err != nil {
				log.WithField("args:", args).Errorf("start configserver failed: %v", err)
			}
		},
	}
)

func init() {
	configserverCommand.PersistentFlags().StringP("config", "c", "etc/config.yaml", "config file")
	_ = viper.BindPFlag("config", configserverCommand.PersistentFlags().Lookup("config"))
}

func main() {
	if err := configserverCommand.Execute(); err != nil {
		log.WithField("args", os.Args).Errorf("configserver execute failed %v", err)
	}
}

func runConfigServer() error {
	configFilePath := viper.GetString("config")
	configServerConfig, err := config.ParseConfigServerConfig(configFilePath)
	if err != nil {
		return errors.Wrap(err, "read and parse configserver config")
	}

	// init logger
	logger.InitLogger(logger.LoggerConfig{
		Level:      configServerConfig.Log.Level,
		Filename:   configServerConfig.Log.Filename,
		MaxSize:    configServerConfig.Log.MaxSize,
		MaxAge:     configServerConfig.Log.MaxAge,
		MaxBackups: configServerConfig.Log.MaxBackups,
		LocalTime:  configServerConfig.Log.LocalTime,
		Compress:   configServerConfig.Log.Compress,
	})

	// init config server
	configServer := server.NewConfigServer(configServerConfig)

	err = configServer.Run()
	if err != nil {
		return errors.Wrap(err, "start config server")
	}

	return nil
}
