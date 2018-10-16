# 使用方法

## 一、二进制直接运行
    1. 去 [release界面](https://github.com/orestonce/autocert/releases/) 下载对应版本的二进制到域名所在主机
    2. 确保主机80端口未被占用, 并且可被外部访问
    3. 运行autocert命令, 生成证书 `./autocert -ServerName example.com`
## 二、源代码安装
    1. 下载源码: ` go get github.com/orestonce/autocert `
    2. 将编译出来的二进制放到域名所在主机上, 确保80端口未被占用
    3. 运行autocert命令, 生成证书 `./autocert -ServerName example.com `
