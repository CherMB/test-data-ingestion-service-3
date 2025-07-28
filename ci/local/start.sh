#/bin/sh
set -x

serviceName="reports-service"
repoBasePath="$GOPATH/src/github.com/calculi-corp"
certBasePath="$repoBasePath/grpc-testutil/creds"

export SERVER_TLS_CA="$certBasePath/CalculiCA-cert.pem"
export SERVER_TLS_PRIVATEKEY="$certBasePath/localhost.key.pem"
export SERVER_TLS_CERTIFICATE="$certBasePath/localhost.cert.pem"
export CLIENT_TLS_PRIVATEKEY="$certBasePath/localhost.key.pem"
export CLIENT_TLS_CERTIFICATE="$certBasePath/localhost.cert.pem"
export CLIENT_TLS_CA="$certBasePath/CalculiCA-cert.pem"
export CLIENT_TLS_CA="$certBasePath/localhost.cert.pem"

certsRuntimeArgs="--server.tls.ca=$SERVER_TLS_CA --server.tls.certificate=$SERVER_TLS_CERTIFICATE --server.tls.privateKey=$SERVER_TLS_PRIVATEKEY --client.tls.ca=$CLIENT_TLS_CA --client.tls.certificate=$CLIENT_TLS_CERTIFICATE --client.tls.privateKey=$CLIENT_TLS_PRIVATEKEY"

go build -o "$serviceName" .

./"$serviceName" --configfile=ci/local/local-config.json $certsRuntimeArgs