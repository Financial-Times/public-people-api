FROM golang:1.8-alpine

RUN mkdir -p "$GOPATH/src"

ADD . "$GOPATH/src/github.com/Financial-Times/public-people-api"

WORKDIR "$GOPATH/src/github.com/Financial-Times/public-people-api"

RUN apk --no-cache --virtual .build-dependencies add git \
    && apk --no-cache --upgrade add ca-certificates \
    && update-ca-certificates --fresh \
    && git config --global http.https://gopkg.in.followRedirects true \
    && cd $GOPATH/src/github.com/Financial-Times/public-people-api \
    && BUILDINFO_PACKAGE="github.com/Financial-Times/public-people-api/vendor/github.com/Financial-Times/service-status-go/buildinfo." \
    && VERSION="version=$(git describe --tag --always 2> /dev/null)" \
    && DATETIME="dateTime=$(date -u +%Y%m%d%H%M%S)" \
    && REPOSITORY="repository=$(git config --get remote.origin.url)" \
    && REVISION="revision=$(git rev-parse HEAD)" \
    && BUILDER="builder=$(go version)" \
    && LDFLAGS="-X '"${BUILDINFO_PACKAGE}$VERSION"' -X '"${BUILDINFO_PACKAGE}$DATETIME"' -X '"${BUILDINFO_PACKAGE}$REPOSITORY"' -X '"${BUILDINFO_PACKAGE}$REVISION"' -X '"${BUILDINFO_PACKAGE}$BUILDER"'" \
    && go get -u github.com/kardianos/govendor \
    && $GOPATH/bin/govendor sync \
    && go-wrapper download \
    && go-wrapper install -ldflags="${LDFLAGS}" \
    && apk del .build-dependencies \
    && rm -rf $GOPATH/src $GOPATH/pkg /usr/local/go

EXPOSE 8080

CMD ["go-wrapper", "run"]