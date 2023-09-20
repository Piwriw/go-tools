## cobra-cli
更快速的cobra工具
```bash
# 安装
go install github.com/spf13/cobra-cli@latest 

# 添加命令
cobra-cli add help

```
## 添加flag
```go
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	exportCmd.Flags().StringVarP(&file, "file","f","local","file to out put")
```