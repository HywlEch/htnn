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

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"mosn.io/htnn/e2e/pkg/k8s"
	"mosn.io/htnn/e2e/pkg/suite"
	mosniov1 "mosn.io/htnn/types/apis/v1"
)

func init() {
	suite.Register(suite.Test{
		Manifests: []string{"base/httproute.yml"},
		Run: func(t *testing.T, suite *suite.Suite) {
			// 测试配置变更传播速度
			c := suite.K8sClient()
			ctx := context.Background()

			// 测量配置变更传播所需的时间
			startTime := time.Now()

			// 创建FilterPolicy
			policy := &mosniov1.FilterPolicy{}

			// 创建新的FilterPolicy
			policy = &mosniov1.FilterPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aggressive-push-test",
					Namespace: k8s.DefaultNamespace,
				},
				Spec: mosniov1.FilterPolicySpec{
					TargetRef: &gwapiv1a2.PolicyTargetReferenceWithSectionName{
						PolicyTargetReference: gwapiv1a2.PolicyTargetReference{
							Group: gwapiv1.Group("gateway.networking.k8s.io"),
							Kind:  gwapiv1.Kind("HTTPRoute"),
							Name:  gwapiv1.ObjectName("default"),
						},
					},
					Filters: map[string]mosniov1.Plugin{
						"demo": {
							Config: runtime.RawExtension{
								Raw: []byte(`{"hostName":"aggressive-push-test"}`),
							},
						},
					},
				},
			}
			err := c.Create(ctx, policy)
			require.NoError(t, err)

			// 等待配置生效
			// 在积极推送模式下，这个时间应该很短
			maxWaitTime := 5 * time.Second
			interval := 100 * time.Millisecond

			var configApplied bool
			for start := time.Now(); time.Since(start) < maxWaitTime; {
				// 发送请求检查配置是否已生效
				rsp, err := suite.Get("/", nil)
				if err == nil {
					req, _, err := suite.Capture(rsp)
					if err == nil {
						// 检查请求头中是否包含更新后的值
						if len(req.Headers["Aggressive-Push-Test"]) > 0 {
							configApplied = true
							break
						}
					}
				}
				time.Sleep(interval)
			}

			propagationTime := time.Since(startTime)
			require.True(t, configApplied, "Configuration was not applied within the expected time")

			// 在积极推送模式下，配置传播时间应该小于2秒
			t.Logf("Configuration propagation took %v", propagationTime)
			require.Less(t, propagationTime.Milliseconds(), int64(2000),
				"Configuration propagation took too long, aggressive push mode may not be working")

			// 清理
			err = c.Delete(ctx, policy)
			require.NoError(t, err)
		},
	})
}
