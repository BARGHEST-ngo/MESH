#!/usr/bin/env bash

# Licensed to Elasticsearch B.V. under one or more contributor
# license agreements. See the NOTICE file distributed with
# this work for additional information regarding copyright
# ownership. Elasticsearch B.V. licenses this file to you under
# the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

# Script to test the output of crd-ref-docs

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMP_DIR=$(mktemp -d -t crd-ref-docs-XXXXX)
DEFAULT_ARGS=(--log-level=ERROR --source-path="${SCRIPT_DIR}/test" --output-path="${TEMP_DIR}/out" --config="${SCRIPT_DIR}/test/config.yaml")

trap '[[ $TEMP_DIR ]] && rm -rf "$TEMP_DIR"' EXIT

run_test() {
    local actual="${TEMP_DIR}/out"
    rm -f "$actual"

    local renderer=asciidoctor
    local templates_dir=

    while :; do
        case "${1:-}" in
            --renderer)
                if [[ -n "${2:-}" ]]; then
                    renderer="$2"
                    shift
                else
                    printf "ERROR: '--renderer' cannot be empty.\n\n" >&2
                    exit 1
                fi
                ;;
            --templates-dir)
                if [[ -n "${2:-}" ]]; then
                    templates_dir="$2"
                    shift
                else
                    printf "ERROR: '--templates-dir' cannot be empty.\n\n" >&2
                    exit 1
                fi
                ;;
            *)
                break
                ;;
        esac

        shift
    done

    local args=("${DEFAULT_ARGS[@]}" --renderer="$renderer")
    if [[ -n "$templates_dir" ]]; then
        args+=(--templates-dir="$templates_dir")
    fi

    local expected
    if [[ "$renderer" == "asciidoctor" ]]; then
        expected=expected.asciidoc
    else
        expected=expected.md
    fi

    (
        cd "$SCRIPT_DIR"
        cmd=(go run main.go "${args[@]}")
        echo "${cmd[@]}"

        "${cmd[@]}"

        local diff
        if diff=$(diff -a -y --suppress-common-lines "${SCRIPT_DIR}/test/${expected}" "$actual"); then
            echo "OK"
        else
            echo "ERROR: outputs differ"
            echo ""
            echo "$diff"
            exit 1
        fi
    )
}

run_test
run_test --renderer asciidoctor --templates-dir templates/asciidoctor
run_test --renderer markdown
run_test --renderer markdown --templates-dir templates/markdown
