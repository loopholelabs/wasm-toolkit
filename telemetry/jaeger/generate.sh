rm *.pem 

echo "CA"
openssl req -x509 -newkey rsa:4096 -nodes -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=GB/ST=London/L=London/O=Loophole/OU=Loophole/CN=*.loopholelabs.io/emailAddress=help@loopholelabs.io"
openssl x509 -in ca-cert.pem -noout -text 

echo "Server"
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=GB/ST=London/L=London/O=Loophole/OU=Loophole/CN=server.loopholelabs.io/emailAddress=help@loopholelabs.io"
openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.conf
openssl x509 -in server-cert.pem -noout -text

openssl verify -CAfile ca-cert.pem server-cert.pem

echo "Client"
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/C=GB/ST=London/L=London/O=Loophole/OU=Loophole/CN=client.loopholelabs.io/emailAddress=help@loopholelabs.io"
openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.conf
openssl x509 -in client-cert.pem -noout -text


