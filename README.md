# 项目名称: bluebell







### 项目技术
- zap-logger更换了go中原生的logger ，并通过日志分割将日志进
行打印到`bluebell.log`文件中.
- 使用第三方库sqlx作为数据库的操作工具，将查询速度加快.


### 使用了redis和mysql进行数据的管理
#### 项目的数据库配置文件在conf包下
- 有`json`和`yaml`两种格式进行配置
- 导入了第三方库viper，直接将配置自动装配然后通过结构体进行调用
