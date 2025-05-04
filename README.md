# myestr-diff

[![Go Reference](https://pkg.go.dev/badge/github.com/oy3o/myestr-diff.svg)](https://pkg.go.dev/github.com/oy3o/myestr-diff)
[![Go Report Card](https://goreportcard.com/badge/github.com/oy3o/myestr-diff)](https://goreportcard.com/report/github.com/oy3o/myestr-diff)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`myestr-diff` 是一个 Go 语言包，用于计算两个字节切片（`[]byte`）之间的差异。它基于 Neil Fraser 对 Myer's 差分算法的优化实现。

此包可以生成描述差异的结构化数据，并能将这些差异序列化为补丁字符串，以及将补丁应用到原始数据上以生成目标数据。

它还包含一个内部使用的 `bytesutil` 子包，提供了一些字节切片操作的辅助函数。

## 特性

*   **高效差异计算:** 使用优化的 Myer's 算法计算 `[]byte` 差异。
*   **补丁生成与应用:**
    *   将差异列表 (`diff.Diffs`) 转换为紧凑的补丁字符串 (`diffs.ToPatch()`)。
    *   将补丁字符串应用到源字节切片 (`diff.Patch(src, patch)`).
*   **可定制化:** `diff.Get` 函数接受选项以调整差异计算行为：
    *   `diff.WithDeadline`: 设置计算超时。
    *   `diff.WithEditcost`: 调整效率优化参数。
    *   `diff.WithChecklines`: 启用/禁用行级预处理以提高速度（可能牺牲最优性）。
    *   `diff.WithSemantic`: 启用语义化清理，使差异更易于人类阅读。
*   **字节工具集 (`bytesutil`):** 包含如 `Reverse`, `Itoa`, `Quote`, `Random` 等字节切片辅助函数 (请注意，部分函数，尤其是涉及 `unsafe` 的函数，需谨慎使用)。

## 安装

```bash
go get github.com/oy3o/myestr-diff@latest
```
或者指定一个版本：
```bash
go get github.com/oy3o/myestr-diff@v0.1.0 # 使用你发布的具体版本号
```

## 使用示例

```go
package main

import (
	"fmt"
	diff "github.com/oy3o/myestr-diff/diff"
)

var text1 = `Hamlet: Do you see yonder cloud that's almost in shape of a camel?
Polonius: By the mass, and 'tis like a camel, indeed.
Hamlet: Methinks it is like a weasel.
Polonius: It is backed like a weasel.
Hamlet: Or like a whale?
Polonius: Very like a whale.
-- Shakespeare`
var text2 = `Hamlet: Do you see the cloud over there that's almost the shape of a camel?
Polonius: By golly, it is like a camel, indeed.
Hamlet: I think it looks like a weasel.
Polonius: It is shaped like a weasel.
Hamlet: Or like a whale?
Polonius: It's totally like a whale.
-- Shakespeare`

var b1 = []byte(text1)
var b2 = []byte(text2)

func main() {
    diffs := diff.Get(b1, b2, diff.WithChecklines(true), diff.WithSemantic(true))
    patch := diffs.ToPatch()
    newText, err := diff.Patch(b1, patch)
    if err != nil{
        // log
    }
    // newText == text2
}
```

## 文档

更详细的 API 文档可以在 [pkg.go.dev](https://pkg.go.dev/github.com/oy3o/myestr-diff) 上找到。

## 许可

本项目采用 [MIT 许可证](LICENSE) 授权。

## 贡献

欢迎贡献！如果你发现任何问题或有改进建议，请随时提交 Pull Request 或开启 Issue。

## 致谢

核心的差异计算逻辑基于 Neil Fraser 先生关于 Diff 策略的研究成果：[https://neil.fraser.name/writing/diff/](https://neil.fraser.name/writing/diff/)
