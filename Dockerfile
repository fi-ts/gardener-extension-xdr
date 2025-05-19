FROM golang:1.24 AS builder

WORKDIR /go/src/github.com/metal-stack/gardener-extension-xdr
COPY . .
RUN make install \
 && strip /go/bin/gardener-extension-xdr

FROM alpine:3.21
WORKDIR /
COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-xdr /gardener-extension-xdr
CMD ["/gardener-extension-xdr"]
