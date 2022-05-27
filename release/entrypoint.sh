#!/bin/sh
./pgrokd -domain $DOMAIN -httpAddr=:80 -httpsAddr=:443 -tunnelAddr=:4443 -tlsCrt=certs/tls.crt -tlsKey=certs/tls.key