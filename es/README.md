## 目标

- 一些基本概念
- 单机es docker 搭建
- mysql -> es ---> canal
- es 基本使用

## 基本概念

### 节点 Node

- 主节点： 负责集群范围内的所有变更，如增加/删除索引、增加/删除节点等。主节点 不涉及文档级别的变更/搜索等操作。node.master=true
- 数据节点： 存储数据以及其对应的倒排索引。node.data=true
- 协调节点： 接收请求、转发请求到其它节点、汇总数据、均衡负载等功能。 node.master=false;node.data=false

### 分片 Shard

一个分片便是一个Lucene 的实例，它本身就是一个完整的搜索引擎。

ES实际上就是利用分片来实现分布式。

分片是数据的容器，文档保存在分片内，分片又被分配到集群内的各个节点里。当你的集群规模扩大或者缩小时， ES会自动的在各节点中迁移分片，使得数据仍然均匀分布在集群里。

主分片与副分片。

### 索引 Index/Indices

### 文档 Document

### 类型 Type （7.x 被移除 也就一个index 一个type）

它是虚拟的逻辑分组，用来过滤 Document。

不同的 Type 应该有相似的结构（schema）。

## 单机部署[https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#docker-cli-run-dev-mode]

官方yaml[https://github.com/elastic/elasticsearch/blob/main/docs/reference/setup/install/docker/docker-compose.yml]

```
docker network create elastic

docker run --name es01 --net elastic -p 9200:9200 -p 9300:9300 -v /ssd/leihan/home/es/data:/usr/share/elasticsearch/data -e "discovery.type=single-node" -t docker.elastic.co/elasticsearch/elasticsearch:8.4.3

----
-> Elasticsearch security features have been automatically configured!
-> Authentication is enabled and cluster connections are encrypted.

->  Password for the elastic user (reset with `bin/elasticsearch-reset-password -u elastic`):
  aETaTa=F3lsh8rFt3+e0

->  HTTP CA certificate SHA-256 fingerprint:
  46d6440c9ff966d2bb17e5d08c72009c2e5eb91bcdbb71e7ea88bec9f4c414dd

->  Configure Kibana to use this cluster:
* Run Kibana and click the configuration link in the terminal when Kibana starts.
* Copy the following enrollment token and paste it into Kibana in your browser (valid for the next 30 minutes):
  eyJ2ZXIiOiI4LjQuMyIsImFkciI6WyIxMC4xMDAuMTg1LjI6OTIwMCJdLCJmZ3IiOiI0NmQ2NDQwYzlmZjk2NmQyYmIxN2U1ZDA4YzcyMDA5YzJlNWViOTFiY2RiYjcxZTdlYTg4YmVjOWY0YzQxNGRkIiwia2V5Ijoib29aMjg0TUJxTnUwWFlkOHM0UUI6cExvYUNrRFRTTnV2RnhMQ0RLNjBIZyJ9

-> Configure other nodes to join this cluster:
* Copy the following enrollment token and start new Elasticsearch nodes with `bin/elasticsearch --enrollment-token <token>` (valid for the next 30 minutes):
  eyJ2ZXIiOiI4LjQuMyIsImFkciI6WyIxMC4xMDAuMTg1LjI6OTIwMCJdLCJmZ3IiOiI0NmQ2NDQwYzlmZjk2NmQyYmIxN2U1ZDA4YzcyMDA5YzJlNWViOTFiY2RiYjcxZTdlYTg4YmVjOWY0YzQxNGRkIiwia2V5Ijoib1laMjg0TUJxTnUwWFlkOHM0UUI6VXZvQnNKWUxTYW10UERtT3ZTUlk5dyJ9

  If you're running in Docker, copy the enrollment token and run:
  `docker run -e "ENROLLMENT_TOKEN=<token>" docker.elastic.co/elasticsearch/elasticsearch:8.4.3`
-----
  

然后拷贝出 证书 docker cp es01:/usr/share/elasticsearch/config/certs/http_ca.crt .

curl --cacert http_ca.crt -u elastic https://localhost:9200 验证启动成功
elastic:aETaTa=F3lsh8rFt3+e0
```

kibana

```
docker pull docker.elastic.co/kibana/kibana:8.4.3

docker run --name kibana --net elastic -p 5601:5601 docker.elastic.co/kibana/kibana:8.4.3


```


