# kcpServer

>+ golang version = 1.21.0

### kcp
使用kcp作为前后端传输协议，提高稳定性和传输速率。

### proto
使用proto作为消息载体，定制协议和通知等。  
并且proto消息对actor支持友好，直接使用proto数据结构作为actor路由。

### actor
为什么使用actor？  
1、使用actor可以隔离服务，支持一个进程跑多个服务的情况;  
2、使用actor可以解耦消息收发和服务处理，耗时的服务不会阻塞其他服务的消息处理;  
3、使用actor降低锁的消耗，对后续动态路由创造条件;  
4、actor天然支持顺序处理消息，为同步创造条件。

### lockstep