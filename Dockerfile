FROM scratch
ADD /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD init-consul-acl /init-consul-acl
ENTRYPOINT ["/init-vault"]
