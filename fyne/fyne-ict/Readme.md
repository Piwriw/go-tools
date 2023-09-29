# 利用fyne 实现的SGU的业绩点计算
## 功能介绍
1. 支持fyne 中文支持
2. 通过excelize 实现支持表格操作，读取学生成绩表 ｜ 输出学生绩点折线图

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