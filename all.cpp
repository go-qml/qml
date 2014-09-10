
// windows use shared library
#if !defined(__MINGW32__) && !defined(__MINGW64__)
#	include "cpp/capi.cpp"
#	include "cpp/govalue.cpp"
#	include "cpp/govaluetype.cpp"
#	include "cpp/idletimer.cpp"
#	include "cpp/connector.cpp"
#	include "cpp/moc_all.cpp"
#endif
