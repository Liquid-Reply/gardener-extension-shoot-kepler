#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

function camelCase {
  sed -r 's/(.)-+(.)/\1\U\2/g;s/^[a-z]/\U&/' <<< "$1"
}

package_name="${1:-imagevector}"
images_yaml="${2:-images.yaml}"
constant_prefix="${3:-}"

out="
$(cat "$(dirname $0)/LICENSE_BOILERPLATE.txt" | sed "s/YEAR/$(date +%Y)/g")

// Code generated by $(basename $0). DO NOT EDIT.

package $package_name

const ("

for image_name in $(yq -r '[.images[].name] | sort | unique | .[]' $images_yaml); do
  variable_name="$(camelCase "$image_name")"

  out="
$out
	// ${constant_prefix}ImageName${variable_name} is a constant for an image in the image vector with name '$image_name'.
	${constant_prefix}ImageName${variable_name} = \"$image_name\""
done

out="
$out
)
"

out_file="${images_yaml%.yaml}.go"
echo "$out" > "$out_file"
goimports -l -w "$out_file"
