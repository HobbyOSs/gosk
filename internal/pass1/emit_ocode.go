package pass1

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

func emitCommand(env *Pass1, command string, ocodes []int32) {
	strSlice := lo.Map(ocodes, func(n int32, _ int) string {
		return strconv.FormatInt(int64(n), 10)
	})
	args := strings.Join(strSlice, ",")
	env.Client.Emit(fmt.Sprintf("%s %s\n", command, args))
}

func stringToOcodes(str string) []int32 {
	ocodes := []int32{}
	ocodes = append(ocodes, lo.Map([]byte(str), func(b byte, index int) int32 {
		return int32(b)
	})...)
	return ocodes
}
