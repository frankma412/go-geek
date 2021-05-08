基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

说明：
三种方式退出
1. stop 信号
2. httpserver  http://localhost:8080/stop
3. mock error`
