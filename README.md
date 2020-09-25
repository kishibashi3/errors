# pkg/errors with grpc errors

grpcアプリケーションのためのエラー基盤。


grpc-goはstatus.Statusを用いてエラーをハンドリングできるが、以下のような問題点があった。

## 1. ErrorLevelの概念が存在しない  

エラーレベルによって、Sentryにアラートを飛ばすかどうかなど、ハンドリングしたいことはある。

## 2. causeがない
 
Pythonでは`raise Error() from e` のような、例外の原因となる別の例外を指定することができるが、golangではこれがない。

そのため、pkg/errorsをForkし、拡張して新しくgrpc errorに対応する例外を作成した。


# GrpcErrorの機能

* pkg/errorsの既存の例外と同じstack traceを持ち、出力することができる。
* grpc error情報（status)を持つ。
* origin (cause)情報をもつ。pkg/errorsもcauseスタックを持っているが、これはエラーのcauseではなく単なるスタックトレース。
* level（エラーレベル）を持つ。


```go
package errors
type grpcError struct {
	status *status.Status
	origin error
	level  ErrorLevel
	stack  *stack
}
```

attr | desc
---|---
status | grpc status
origin | error cause
level | error level
stack | pkg/errors stack





`go get github.com/kishibashi3/grpc/errors`



