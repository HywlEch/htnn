#!/usr/bin/env bash
# Copyright The HTNN Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -eo pipefail
set -x

HELM="${LOCALBIN}/helm"
E2E_DIR="$(pwd)"
# 默认为false，可以通过环境变量覆盖
AGGRESSIVE_PUSH="${AGGRESSIVE_PUSH:-false}"

install() {
    OPTIONS_DISABLED="$1"

    pushd ../manifests/charts

    $HELM dependency update htnn-controller
    $HELM dependency update htnn-gateway
    $HELM package htnn-controller htnn-controller
    $HELM package htnn-gateway htnn-gateway

    CONTROLLER_VALUES_OPT="-f $E2E_DIR/htnn_controller_values.yaml"
    GATEWAY_VALUES_OPT="-f $E2E_DIR/htnn_gateway_values.yaml"

    # 如果启用了积极推送模式，生成包含相应配置的istio-values.yaml
    if [ "$AGGRESSIVE_PUSH" = "true" ]; then
        cat > $E2E_DIR/istio-values.yaml <<EOF
pilot:
  env:
    # 减少推送延迟时间
    PILOT_DEBOUNCE_AFTER: "10ms"
    # 减少最大延迟时间
    PILOT_DEBOUNCE_MAX: "1s"
    # 禁用EDS推送延迟
    PILOT_ENABLE_EDS_DEBOUNCE: "false"
    # 增加并发推送数量
    PILOT_PUSH_THROTTLE: "200"
EOF
        CONTROLLER_VALUES_OPT="$CONTROLLER_VALUES_OPT -f $E2E_DIR/istio-values.yaml"
        echo "Enabled aggressive push mode for istiod"
    fi

    if [ -n "$OPTIONS_DISABLED" ]; then
        CONTROLLER_VALUES_OPT=
        GATEWAY_VALUES_OPT=
    fi

    # shellcheck disable=SC2086
    $HELM install htnn-controller htnn-controller --namespace istio-system --create-namespace --wait $CONTROLLER_VALUES_OPT \
        || exitWithAnalysis

    # shellcheck disable=SC2086
    $HELM install htnn-gateway htnn-gateway --namespace istio-system --create-namespace $GATEWAY_VALUES_OPT \
        && \
        (kubectl wait --timeout=5m -n istio-system deployment/istio-ingressgateway --for=condition=Available \
        || exitWithAnalysis)

    popd
}

installWithoutOptions() {
    install WithoutOptions
}

exitWithAnalysis() {
    kubectl get pods -n istio-system -o yaml
    for pod in $(kubectl get pods -n istio-system | grep 'istiod-' | awk '{print $1}'); do
        kubectl -n istio-system logs --tail=1000 "$pod"
        echo
    done
    for pod in $(kubectl get pods -n istio-system | grep 'istio-ingressgateway' | awk '{print $1}'); do
        kubectl -n istio-system logs --tail=1000 "$pod"
        echo
    done
    exit 1
}

uninstall() {
    $HELM uninstall htnn-controller -n istio-system && $HELM uninstall htnn-gateway -n istio-system && kubectl delete ns istio-system
}

opt=$1
shift

${opt} "$@"
