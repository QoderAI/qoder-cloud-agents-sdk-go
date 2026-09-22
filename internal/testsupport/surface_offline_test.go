package testsupport

// 反向守护：若任一离线符号被误挪到 //go:build live 文件，此文件离线编译即 RED。
var (
	_ Endpoint
	_ Transport
	_ TurnResult
	_ CleanupFailure
	_ = FailureContracts
	_ = CheckFields
	_ = InvokeJSON
	_ = SafeError
	_ = ResourceAlreadyGone
)
