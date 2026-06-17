# Author: Ira Yago
# Date: 6/4/2026
# Description: HTTP client

# send $i HTTP requests through same TCP connection
<#
for ($i = 0; $i -le 2; $i++) {
	#curl -v -s http://REPLACE_WITH_IP #will need to change the ip addr to server host ip
	Invoke-WebRequest -Uri "http://REPLACE_WITH_IP" | Select-Object -ExpandProperty Content
}
#>

Invoke-WebRequest -Uri "http://REPLACE_WITH_IP" | Select-Object -ExpandProperty Content
