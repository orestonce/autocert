# 自动从let's encrypt签发免费SSL证书的工具
-----
## 使用方法
### 一、直接使用二进制
- 去 [release界面](https://github.com/orestonce/autocert/releases/)下载对应版本的二进制到域名所在主机
- 确保主机80端口未被占用, 并且可被外部访问
- 运行autocert命令, 生成证书 ` ./autocert -ServerName example.com `
- 命令运行完毕后, 当前目录会自动生成1个证书文件example.com.crt, 1个私钥文件 example.com.key
### 二、源代码安装
- 下载源码: ` go get github.com/orestonce/autocert `
- 将编译出来的二进制放到域名所在主机上, 确保80端口未被占用
- 运行autocert命令, 生成证书 ` ./autocert -ServerName example.com `
