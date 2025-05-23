// Copyright The HTNN Authors.
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

//go:build envoydev

package api

import (
	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
)

var (
	LogTrace     = api.LogTrace
	LogDebug     = api.LogDebug
	LogInfo      = api.LogInfo
	LogWarn      = api.LogWarn
	LogError     = api.LogError
	LogCritical  = api.LogCritical
	LogTracef    = api.LogTracef
	LogDebugf    = api.LogDebugf
	LogInfof     = api.LogInfof
	LogWarnf     = api.LogWarnf
	LogErrorf    = api.LogErrorf
	LogCriticalf = api.LogCriticalf

	GetLogLevel = api.GetLogLevel
)

// SecretManager exports the managed secret
type SecretManager = api.SecretManager
