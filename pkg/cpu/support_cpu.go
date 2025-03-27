package cpu

// SupportCPU は、CPUが特定の機能をサポートしているかどうかを判定するインターフェースです。
type SupportCPU interface {
	// IsSupported は、指定された機能がサポートされているかどうかを返します。
	IsSupported(feature string) bool
}
