#include "capi.h"

#ifdef __cplusplus
extern "C" {
#	include "_cgo_export.h"
}
#else
#	include "_cgo_export.h"
#endif

void initGoQmlLib() {
	HookHandlers p[1];
	p->hookIdleTimer = cgoHookIdleTimer;
	p->hookLogHandler = cgoHookLogHandler;
	p->hookGoValueReadField = cgoHookGoValueReadField;
	p->hookGoValueWriteField = cgoHookGoValueWriteField;
	p->hookGoValueCallMethod = cgoHookGoValueCallMethod;
	p->hookGoValueDestroyed = cgoHookGoValueDestroyed;
	p->hookGoValuePaint = cgoHookGoValuePaint;
	p->hookRequestImage = cgoHookRequestImage;
	p->hookGoValueTypeNew = cgoHookGoValueTypeNew;
	p->hookWindowHidden = cgoHookWindowHidden;
	p->hookSignalCall = cgoHookSignalCall;
	p->hookSignalDisconnect = cgoHookSignalDisconnect;
	p->hookPanic = cgoHookPanic;
	p->hookListPropertyCount = cgoHookListPropertyCount;
	p->hookListPropertyAt = cgoHookListPropertyAt;
	p->hookListPropertyAppend = cgoHookListPropertyAppend;
	p->hookListPropertyClear = cgoHookListPropertyClear;
	initHooks(p);
}
