# logrus-multihandler
用于logrus的多路输出器，每个Logger可以向多个Writer以不同的Fromatter和不同的Level输出日志

## 介绍

`MultiHandler`结构实现了`Formatter`接口，可使用`SetFormatter`配置进`logrus.Logger`

使用`NewMultiHandler`函数创建该接口时，传入任意使用`NewHandler`函数创建的`Handler`结构，这些结构相当于一个个`logrus.Logger`，当在原本的`logrus.Logger`上写入时，将会把同样的`Entry`分发到这些`Handler`，以配置的参数分别解析并写入。

在原本的`logrus.Logger`上写入时，`MultiHandler`作为一个`Formatter`，会返回空字节数组，即不再向本来的`Logger`的`Out`写入日志

## 示例

该示例得到一个logger，该logger可以同时向标准输出和文件写入日志，并且标准输出中的日志为彩色，文件中的日志没有彩色。

```go
package main

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gedoy9793/logrus-multihandler"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var logPath = "/var/log/examle"

func main() {
	stdLogFormatter := &nested.Formatter{
		TimestampFormat: time.RFC3339,
	}
	fileLogFormatter := &nested.Formatter{
		TimestampFormat: time.RFC3339,
		NoColors:        true,
	}

	logFile, err := rotatelogs.New(
		logPath+".%Y%m%d",
		rotatelogs.WithRotationTime(time.Hour*time.Duration(24)),
		rotatelogs.WithLinkName(logPath),
	)

	if err != nil {
		panic(err)
	}

	multiFormatter := multihandler.NewMultiHandler(multihandler.NewHandler(
		stdLogFormatter,
		logrus.DebugLevel,
		os.Stdout,
		&multihandler.HandlerConfig{},
	), multihandler.NewHandler(
		fileLogFormatter,
		logrus.DebugLevel,
		logFile,
		&multihandler.HandlerConfig{},
	))
	
	logger := logrus.New()
	logger.SetFormatter(multiFormatter)
}
```