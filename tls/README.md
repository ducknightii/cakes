https://zhuanlan.zhihu.com/p/133375078

自签证书流程: https://www.myfreax.com/creating-a-self-signed-ssl-certificate/


openssl req -newkey rsa:4096 \
-x509 \
-sha256 \
-days 3650 \
-nodes \
-out example.crt \
-keyout example.key


验证配置：
nginx -t

重启：
nginx -s reload