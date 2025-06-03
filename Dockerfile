FROM gcr.io/distroless/static-debian12:latest
COPY --chmod=555 build/housesitter-z2m-linux-amd64 .
COPY lib/server/http/index.gohtml lib/server/http/index.gohtml
COPY lib/server/http/device.gohtml lib/server/http/device.gohtml
COPY public public
USER       nobody
EXPOSE     6724
ENTRYPOINT [ "/housesitter-z2m-linux-amd64" ]