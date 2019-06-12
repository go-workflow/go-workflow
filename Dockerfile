FROM scratch
ADD /go-workflow //
ADD /config.json //
EXPOSE 8080
ENTRYPOINT [ "/go-workflow" ]