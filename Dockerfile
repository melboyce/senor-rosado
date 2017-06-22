FROM scratch
ADD senor-rosado-static senor-rosado
ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD zoneinfo.zip /usr/lib/go/lib/time/zoneinfo.zip
ENTRYPOINT ["/senor-rosado"]
