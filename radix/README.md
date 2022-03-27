# go-radix-tree
Radix Tree implementation for go
A 256 fork Radix Tree
####256叉基数树（含删除节点功能）。欢迎引用，欢迎star。
####引用方法：把radix子目录完整地copy到您的项目里，然后进行您希望的修改。
```azure
I Insert插入数据
S Show完整的树
F Find查找数据
D Delete删除数据
E Exit退出
请输入指令，按回车键：

展示树：(父|Child数)路径|string载荷|int载荷...\child-1路径\child-2路径...
 , (nil|2)as|as|[5]\df\vb
 , (as|1)df|asdf|[2 10]\g , (as|0)vb|asvb|[8]
 , (df|2)g|asdfg|[6]\h\n
 , (g|0)h|asdfgh|[5] , (g|0)n|asdfgn|[6]


```
