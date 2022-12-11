/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package undo

import (
	"flag"
)

type Compress struct {
	Enable    bool   `yaml:"enable" json:"enable,omitempty" property:"enable"`
	Type      string `yaml:"type" json:"type,omitempty" property:"type"`
	Threshold int    `yaml:"threshold" json:"threshold,omitempty" property:"threshold"`
}

type Config struct {
	DataValidation        bool     `yaml:"data-validation" json:"data-validation,omitempty" property:"data-validation"`
	LogSerialization      string   `yaml:"log-serialization" json:"log-serialization,omitempty" property:"log-serialization"`
	LogTable              string   `yaml:"log-table" json:"log-table,omitempty" property:"log-table"`
	OnlyCareUpdateColumns bool     `yaml:"only-care-update-columns" json:"only-care-update-columns,omitempty" property:"only-care-update-columns"`
	Compress              Compress `yaml:"compress" json:"compress,omitempty" property:"compress"`
}

func (ufg *Config) RegisterFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	f.BoolVar(&ufg.DataValidation, prefix+".data-validation", true, "Judge whether the before image and after image are the same，If it is the same, undo will not be recorded")
	f.StringVar(&ufg.LogSerialization, prefix+".log-serialization", "jackson", "Serialization method.")
	f.StringVar(&ufg.LogTable, prefix+".log-table", "undo_log", "undo log table name.")
	f.BoolVar(&ufg.OnlyCareUpdateColumns, prefix+".only-care-update-columns", true, "The switch for degrade check.")

}

// RegisterFlagsWithPrefix for Compress.
func (cfg *Compress) RegisterFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	f.BoolVar(&cfg.Enable, prefix+".log-table-name", true, "Whether compression is required.")
	f.StringVar(&cfg.Type, prefix+".clean-period", "zip", "Compression type")
	f.IntVar(&cfg.Threshold, prefix+".clean-period", 64, "Compression threshold Unit: k")
}
