# leveldbmanager

#### 介绍
该软件为一工具软件，使用golang语言，基于[goleveldb](github.com/syndtr/gole    veldb/leveldb)的命令行工具，能够实现读取leveldb数据库目录，计算数据量，读取和写入数值
#### 软件架构
命令行 跨平台运行于windows,linux,macos


#### 安装教程

1.  git clone https://gitee.com/dadigang/leveldbmanager
2.  cd leveldbmanager
3.  go build

#### 使用说明

1.  ./leveldbmanager
2.  ./leveldbmanager dbpath

#### 具体使用
1.	count
	计算数据库中所有数据数量
2.	list
	枚举所有数据库中信息
3.	put
	在数据库中写入信息，如put key value
4.	get
	获取数据中心信息，get key
5.	exit
	退出系统
6.	q
	退出系统
7.	rm
	删除数据库中数据 rm key

