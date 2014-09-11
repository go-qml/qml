package qml

/*
#include "capi.h"
*/
import "C"

func init() {
	C.initGoQmlLib()
}
