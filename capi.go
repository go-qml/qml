package qml

func init() {
	// Install the C++ log handler that diverts calls to the hook below.
	C.installLogHandler()
}

//void initHooks(HookHandlers* h);