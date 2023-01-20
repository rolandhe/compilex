# compilex
配合daog组件的工具，可以把create table语句编译成daog需要的go代码

# 原理
基于mysql sql语句解析器解析出create table语句，遍历出表名、字段名、字段类型，字段是否为自增长，然后生成对应的go代码文件。
生成两个文件：
* 主文件，开发者不能修改，包括表所对应的go struct、元数据、记录表字段名的匿名struct对象
* 扩展文件，开发者可以修改、扩展，尤其是分表规则函数需要在init函数中设置

sql解析使用tidb/parser包
# 编译代码
运行build文件夹下的脚本

# 编译sql
```
./compilex -i="sql file" -pkg packageName -o xxx/xx
```
