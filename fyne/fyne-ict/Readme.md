# 利用fyne 实现的SGU的业绩点计算
## 环境准备
web URL：https://jmeubank.github.io/tdm-gcc/download/

下载来链接（64位操作系统，windows）：https://github.com/jmeubank/tdm-gcc/releases/download/v10.3.0-tdm64-2/tdm64-gcc-10.3.0-2.exe
需要安装gcc （windows需要，mac自带）
## 功能介绍
1. 支持fyne 中文支持（Done）
2. 读取某一个文件夹下所有的表格（一个学生一个表格），输出每一个学生的每个学年的绩点平均值（通过excelize 实现支持表格操作）和大学四年的绩点（如果是有四年的就输出，没有就输出目前获取的比例）
   - 需要考虑：学生补考、重修的问题（参照SGU绩点计算规则处理）
3. 记得需要支持GUI的界面，界面也需要变动

## fyne不支持中文解决方案
1. 下载对应的支持中文的字体
2. 安装fyne工具和命令工具
    ```bash
   go get  github.com/flopp/go-findfont 支持中文
   go install fyne.io/fyne/v2/cmd/fyne@latest
   ```
4. 输出文件
    ```bash
   fyne bundle ShangShouJianSongXianXiTi-2.ttf >> bundle.go
    ```
5. 组织Mytheme
    ```go
    
    package gui
    
    import (
        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/theme"
        "image/color"
    )
    
    //example 
    type MyTheme struct{}
    
    var _ fyne.Theme = (*MyTheme)(nil)
    
    // Font 会爆红 无所谓 能正常运行
    func (*MyTheme) Font(s fyne.TextStyle) fyne.Resource {
        if s.Monospace {
            return theme.DefaultTheme().Font(s)
        }
        return resourceSimkaiTtf
    }
    
    func (*MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
        return theme.DefaultTheme().Color(n, v)
    }
    
    func (*MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
        return theme.DefaultTheme().Icon(n)
    }
    
    func (*MyTheme) Size(n fyne.ThemeSizeName) float32 {
        return theme.DefaultTheme().Size(n)
    }
    ```
5. 使用MYtheme
    ```go
        myApp := app.New()
        myApp.Settings().SetTheme(&MyTheme{})
    ```