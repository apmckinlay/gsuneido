// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

typedef unsigned long long uintptr;

extern uintptr helperHwnd;

void setup();
int run();
int interrupt();
int alert(char* msg, int type);
void fatal(char* msg);
int message_loop(uintptr hdlg);
uintptr createLexer(char* name);

long EmbedBrowserObject(uintptr hwnd, void* pBrowserObject, void* pPtr);
void UnEmbedBrowserObject(uintptr browserObject, uintptr ptr);

uintptr queryIDispatch(uintptr iunk);
uintptr createInstance(char* progid);
long invoke(uintptr idisp, char* name, uintptr flags, void* args, void* result);
long release(uintptr iunk);

typedef struct {
	char* data;
	int size;
} buf_t;

long WebView2_Create(uintptr hwnd, void* pBrowserObject, char* dllPath, 
    char* userDataFolder, uintptr cb);
long WebView2_Resize(uintptr pBrowserObject, long w, long h);
long WebView2_Navigate(uintptr pBrowserObject, char* s);
long WebView2_NavigateToString(uintptr pBrowserObject, char* s);
long WebView2_ExecuteScript(uintptr pBrowserObject, char* script);
long WebView2_GetSource(uintptr pBrowserObject, char* dst);
long WebView2_Print(uintptr pBrowserObject);
long WebView2_SetFocus(uintptr pBrowserObject);
long WebView2_Close(uintptr pBrowserObject);

// deps last modified 2025-07-18 17:47:02 UTC
