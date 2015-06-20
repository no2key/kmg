package kmgView

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestNewHtmlRendererListFromList(ot *testing.T) {
	NewHtmlRendererListFromList([]String{})
	out := NewHtmlRendererListFromList([]String{"1"})
	kmgTest.Equal(len(out), 1)
	kmgTest.Equal(out[0].(String), String("1"))
}
