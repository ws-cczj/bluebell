# 项目名称: bluebell

## 项目大概流程
### main函数注册流程
![init.png](http://cdn.cczjblog.top/cczjBlog-img/bluebell_init.png-cczjImage)
![user.png](http://cdn.cczjblog.top/cczjBlog-img/bluebell_user.png-cczjImage)
![community](http://cdn.cczjblog.top/cczjBlog-img/bluebell_community.png-cczjImage)
![post](http://cdn.cczjblog.top/cczjBlog-img/bluebell_post.png-cczjImage)
![comment.png](http://cdn.cczjblog.top/cczjBlog-img/bluebell_comment.png-cczjImage)
初始化定时任务,开启过期帖子定期清除任务,清除残余帖子评论信息,定实监控任务执行情况<br>
开启web服务,监听请求,等待优雅关机.

### 项目功能
- 基于`redis`实现帖子的时间排序和投票分数排序。
- 基本实现帖子的`CRUD`的功能，帖子ID使用雪花算法ID生成器生成。
- 实现用户登录和注册功能，用户ID使用雪花算法ID生成器生成。使用jwt中间件对用户认证进行处理
  并实现了同一时间限制唯一用户的功能，jwt采用双token实现，实现续约的同时增加安全性。
- 基于`redis`实现实现帖子的投票功能，对于帖子投票进行分数统计和排名,帖子的按照时间排序和投票分数排序,
- 实现社区的增加和查找,实现根据用户id来查询用户所管理的所有社区,实现根据社区对于帖子的时间排序和投票分数排序。
- 开启定实任务对过期帖子进行排查，将帖子数据移植到数据库中，并将redis中的过期数据进行清除。
  帖子分为三种情况: 0 待审核, 1 已发布, 2 已保存, 3 已过期
- 实现帖子评论区，评论区可以按照时间和分值进行排序。返回前端的数据已经被压缩过，一个父评论
  可以包括多条子评论。实现删除评论功能，点赞功能。基于`redis`存储点赞数据，将`redis`作为缓存数据库
  存储评论的id，通过主键id去查询数据库中的评论。

### 项目技术
- 实现了优雅关机，在程序结束后有五秒时间处理期间请求。
- 使用第三方库`cron`开启定时任务，使用`sync.Once`对简单业务定时任务进行唯一实例化限制，自定义
结构体`Crontab`记录任务执行时间，如果任务执行错误，将时间重置为上一次任务执行时间等待下次继续执行。
开启另一个定时任务对当前开启的定时任务以及执行情况进行指定时间的日志输出，方便观察记录定时任务情况。
- 使用`gin`作为轻量级web服务框架进行开发
- 使用第三方库`zap-logger`更换了go中原生的logger,自定义了日志格式并通过日志分割将日志进
行打印到`bluebell.log`文件中.
- 使用第三方库`viper`来管理配置文件信息
- 使用第三方库`sqlx`作为数据库的操作工具，将查询速度加快.
- 使用第三方库`JWT`作为用户认证管理中间件，并搭配`redis`实现同一时间内单个用户登录的限制
,使用`atoken和rtoken`方式同步刷新token，实现token的续约.
- 使用第三方库`ratelimit`作为限流中间件，限制突发情况下大规模的请求,避免导致服务器崩溃的情况
- 使用错误包对错误进行管理，并且定义专门的响应流处理响应。
并且使用第三方库`validator`对部分字段进行空值校验，确保参数的正确接收以及相应的中英文返回。

### 使用了redis和mysql进行数据的管理
#### 项目的数据库配置文件在conf包下
- 有`json`和`yaml`两种格式进行配置，通过`viper`对配置文件进行读取，并且更改配置文件时会触发`viper`的重新加载。
