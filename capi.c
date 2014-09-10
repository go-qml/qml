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
	p->hookIdleTimer = hookIdleTimer;
	p->hookLogHandler = hookLogHandler;
	p->hookGoValueReadField = hookGoValueReadField;
	p->hookGoValueWriteField = hookGoValueWriteField;
	p->hookGoValueCallMethod = hookGoValueCallMethod;
	p->hookGoValueDestroyed = hookGoValueDestroyed;
	p->hookGoValuePaint = hookGoValuePaint;
	p->hookRequestImage = hookRequestImage;
	p->hookGoValueTypeNew = hookGoValueTypeNew;
	p->hookWindowHidden = hookWindowHidden;
	p->hookSignalCall = hookSignalCall;
	p->hookSignalDisconnect = hookSignalDisconnect;
	p->hookPanic = hookPanic;
	p->hookListPropertyCount = hookListPropertyCount;
	p->hookListPropertyAt = hookListPropertyAt;
	p->hookListPropertyAppend = hookListPropertyAppend;
	p->hookListPropertyClear = hookListPropertyClear;
	initHooks(p);
}
