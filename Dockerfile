# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

############# builder
FROM golang:1.23 AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-shoot-kepler
COPY . .
RUN make install

############# gardener-extension-shoot-kepler
FROM alpine:3.20.3 AS gardener-extension-shoot-kepler

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-shoot-kepler /gardener-extension-shoot-kepler
ENTRYPOINT ["/gardener-extension-shoot-kepler"]
