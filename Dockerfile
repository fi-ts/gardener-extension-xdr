FROM golang:1.25 AS builder

WORKDIR /go/src/github.com/fi-ts/gardener-extension-xdr
COPY . .
RUN make install \
 && strip /go/bin/gardener-extension-xdr

FROM alpine:3.22
WORKDIR /
COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-xdr /gardener-extension-xdr
CMD ["/gardener-extension-xdr"]
