#include "runtime.h"

void ·Ref(uintptr ref) {
	ref = (uintptr)g->m;
	FLUSH(&ref);
}
