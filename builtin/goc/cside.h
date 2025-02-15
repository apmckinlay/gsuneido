// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

typedef unsigned long long uintptr;

extern unsigned long threadid;
extern uintptr helperHwnd;

void setup();
int run();
int interrupt();
void signalAndWait();
int alert(char* msg, int type);
void fatal(char* msg);
int message_loop(uintptr hdlg);
uintptr createLexer(uintptr name);

long EmbedBrowserObject(uintptr hwnd, uintptr pBrowserObject, uintptr pPtr);
void UnEmbedBrowserObject(uintptr browserObject, uintptr ptr);
long WebView2(uintptr op, uintptr arg1, uintptr arg2, uintptr arg3, uintptr arg4, uintptr arg5);

uintptr queryIDispatch(uintptr iunk);
uintptr createInstance(uintptr progid);
long invoke(uintptr idisp, uintptr name, uintptr flags, uintptr args, uintptr result);
long release(uintptr iunk);

typedef struct {
	char* data;
	int size;
} buf_t;

// deps last modified 2025-02-14 22:06:58 UTC
