// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

typedef unsigned long long uintptr;

enum { maxargs = 20 };
extern uintptr args[maxargs];

extern unsigned long threadid;
extern uintptr helperHwnd;

enum {
	ncb2s = 32,
	ncb3s = 32,
	ncb4s = 1024,
};

extern uintptr cb2s[ncb2s];
extern uintptr cb3s[ncb2s];
extern uintptr cb4s[ncb4s];

void start();
void signalAndWait();
int alert(char* msg, int type);
void fatal(char* msg);

enum {
	msg_none,
	msg_result,
	msg_syscall,
	msg_callback2,
	msg_callback3,
	msg_callback4,
	msg_msgloop,
	msg_timer,
	msg_sunapp,
	msg_queryidispatch,
	msg_createinstance,
	msg_invoke,
	msg_release,
	msg_interrupt,
	msg_embedbrowserobject,
	msg_unembedbrowserobject,
	msg_webview2,
    msg_createlexer,
	msg_shutdown,
	msg_setupconsole,
};

typedef uintptr(__stdcall* cb2_t)(uintptr a, uintptr b);
typedef uintptr(__stdcall* cb3_t)(uintptr a, uintptr b, uintptr c);
typedef uintptr(__stdcall* cb4_t)(uintptr a, uintptr b, uintptr c, uintptr d);

uintptr cb2(uintptr i, uintptr a, uintptr b);
uintptr cb3(uintptr i, uintptr a, uintptr b, uintptr c);
uintptr cb4(uintptr i, uintptr a, uintptr b, uintptr c, uintptr d);

typedef struct {
	char* buf;
	int size;
} buf_t;

// deps last modified 2025-02-09 23:57:44 UTC
