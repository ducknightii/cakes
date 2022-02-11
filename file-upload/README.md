对于大文件的上传，不能全部加载在内存，同时我们需要为文件生成摘要(md5)用于检查是否完整

- http.Request.ParseMultipartForm(memSize) 用于设置内存使用大小，如果超出则会用临时文件存储
- io.Copy 会使用到32K大小的内存
- 计算大文件的md5: 可以在 new后 以append的方式写入，参考main.fileSaveStream()
- 同样disk存储，main.fileSaveBytes 需要将整个文件读取到内存，大约需要2倍文件大小内存(具体大小参照slice扩容规律)，main.fileSaveStream 只需要io.Copy用到的 32K内存大小的额外空间
- bos 简单流式上传 是需要读取全部数据（也即完整加载到内存）后才发送数据，测试时 20Mi 大概需要6s（当然与带宽有关）, stream 方式要比bytes 耗费内存更少
- main.fileSaveStream 实现了 bos分块上传，不过测试 100Mi文件分块上传 采用stream 或者 multi 耗时差不多 30s+, multi 内存占用相对较少
- 本demo 由于bos原因 最终还是整个file加载到了内存（可以改动下multi分批次分块上传，实现部分内存加载），可以在业务层限制下并发，达到内存以及io控制的目的