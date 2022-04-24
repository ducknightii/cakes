可以依照 kube-prom 很方便完成部署: 

https://github.com/prometheus-operator/kube-prometheus
```
简易开放web port：

 kubectl --namespace monitoring port-forward --address 0.0.0.0 svc/prometheus-k8s 9090
 
 kubectl --namespace monitoring port-forward --address 0.0.0.0 svc/grafana 3000
```

一些说明 以及 pod monitor加入:
https://cloud.tencent.com/developer/article/1802679

```
kubectl create secret generic additional-configs --from-file=prometheus-additional.yaml -n monitoring
```

loki 是单机部署的一个demo，正式使用 promtail应该是采用DaemonSet的方式部署，或者直接参考官网 helm方式部署

https://grafana.com/docs/loki/latest/installation/simple-scalable-helm/

一个k8s 中的promtail配置 （主要是看下如何与k8s结合获取label）

https://blog.csdn.net/u011943534/article/details/120871258
