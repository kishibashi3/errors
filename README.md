# pkg/errors with grpc errors

grpcアプリケーションのためのエラー基盤。


grpc-goはstatus.Statusのようなエラーをハンドリングするためのstructを持っている。
単にエラーを返したいだけなら、`status.Err()`を返してやればいい。

だが、扱う際に以下のような問題点があった。

1. ErrorLevelの概念が存在しない  
    > エラーレベルはStatusと必ずしも紐づくわけではない。

2. causeがない
    > causeは、


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

grpcエラーを返却するための材料としてのstatusを持つ。
また、エラーをレベルに応じたハンドリングをするためのlevel、それからpkg/errorsの機能であるstackを持つ。

多言語だとerror causeに相当するエラー同士の紐づけが大抵あるが、golangでは見当たらないので、追加。




`go get github.com/kishibashi3/pkg/errors`



