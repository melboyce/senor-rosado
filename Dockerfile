FROM scratch
ADD senor-rosado /
ADD ./ca-certificates.crt /etc/ssl/certs/
CMD ["/senor-rosado"]
