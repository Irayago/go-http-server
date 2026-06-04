# Author: Ira Yago
# Date: 6/4/2026
# Description: HTTP client

# send 2 HTTP requests through same TCP connection
for ($i = 0; $i -le 1; $i++) {
	curl http://10.124.48.77:9999 #will need to change the ip addr to server host ip
}
