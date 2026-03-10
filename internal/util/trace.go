package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

// MakeTraceHeaders generates Datadog-compatible trace headers.
func MakeTraceHeaders() map[string]string {
	traceID := make([]byte, 16)
	rand.Read(traceID)
	spanID := make([]byte, 8)
	rand.Read(spanID)

	traceIDHex := hex.EncodeToString(traceID)
	spanIDHex := hex.EncodeToString(spanID)

	// traceparent: 00-{traceID}-{spanID}-01
	traceparent := fmt.Sprintf("00-%s-%s-01", traceIDHex, spanIDHex)

	// tracestate: dd=t.dm:-1;t.tid:{first16ofTraceID};s:-1
	tracestate := fmt.Sprintf("dd=t.dm:-1;t.tid:%s;s:-1", traceIDHex[:16])

	// x-datadog-trace-id: decimal conversion of last 16 hex chars (last 8 bytes)
	traceIDInt := new(big.Int).SetBytes(traceID[8:])
	// x-datadog-parent-id: decimal conversion of spanID
	spanIDInt := new(big.Int).SetBytes(spanID)

	return map[string]string{
		"traceparent":                traceparent,
		"tracestate":                 tracestate,
		"x-datadog-trace-id":         traceIDInt.String(),
		"x-datadog-parent-id":        spanIDInt.String(),
		"x-datadog-sampling-priority": "-1",
	}
}