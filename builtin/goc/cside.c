// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

#include "cside.h"

extern void timerId();
int Scintilla_RegisterClasses(void* hInstance);

uintptr interact();

uintptr args[maxargs];

uintptr cb2(uintptr i, uintptr a, uintptr b) {
	args[0] = msg_callback2;
	args[1] = i;
	args[2] = a;
	args[3] = b;
	return interact();
}

uintptr cb3(uintptr i, uintptr a, uintptr b, uintptr c) {
	args[0] = msg_callback3;
	args[1] = i;
	args[2] = a;
	args[3] = b;
	args[4] = c;
	return interact();
}

uintptr cb4(uintptr i, uintptr a, uintptr b, uintptr c, uintptr d) {
	args[0] = msg_callback4;
	args[1] = i;
	args[2] = a;
	args[3] = b;
	args[4] = c;
	args[5] = d;
	return interact();
}

typedef uintptr(__stdcall* p0)();
typedef uintptr(__stdcall* p1)(uintptr a);
typedef uintptr(__stdcall* p2)(uintptr a, uintptr b);
typedef uintptr(__stdcall* p3)(uintptr a, uintptr b, uintptr c);
typedef uintptr(__stdcall* p4)(uintptr a, uintptr b, uintptr c, uintptr d);
typedef uintptr(__stdcall* p5)(
	uintptr a, uintptr b, uintptr c, uintptr d, uintptr e);
typedef uintptr(__stdcall* p6)(
	uintptr a, uintptr b, uintptr c, uintptr d, uintptr e, uintptr f);
typedef uintptr(__stdcall* p7)(uintptr a, uintptr b, uintptr c, uintptr d,
	uintptr e, uintptr f, uintptr g);
typedef uintptr(__stdcall* p8)(uintptr a, uintptr b, uintptr c, uintptr d,
	uintptr e, uintptr f, uintptr g, uintptr h);
typedef uintptr(__stdcall* p9)(uintptr a, uintptr b, uintptr c, uintptr d,
	uintptr e, uintptr f, uintptr g, uintptr h, uintptr i);
typedef uintptr(__stdcall* p10)(uintptr a, uintptr b, uintptr c, uintptr d,
	uintptr e, uintptr f, uintptr g, uintptr h, uintptr i, uintptr j);
typedef uintptr(__stdcall* p11)(uintptr a, uintptr b, uintptr c, uintptr d,
	uintptr e, uintptr f, uintptr g, uintptr h, uintptr i, uintptr j,
	uintptr k);
typedef uintptr(__stdcall* p12)(uintptr a, uintptr b, uintptr c, uintptr d,
	uintptr e, uintptr f, uintptr g, uintptr h, uintptr i, uintptr j, uintptr k,
	uintptr l);

uintptr syscall() {
	void* p = (void*) (args[1]);
	switch (args[2]) {
	case 0:
		return ((p0)(p))();
	case 1:
		return ((p1)(p))(args[3]);
	case 2:
		return ((p2)(p))(args[3], args[4]);
	case 3:
		return ((p3)(p))(args[3], args[4], args[5]);
	case 4:
		return ((p4)(p))(args[3], args[4], args[5], args[6]);
	case 5:
		return ((p5)(p))(args[3], args[4], args[5], args[6], args[7]);
	case 6:
		return ((p6)(p))(args[3], args[4], args[5], args[6], args[7], args[8]);
	case 7:
		return ((p7)(p))(
			args[3], args[4], args[5], args[6], args[7], args[8], args[9]);
	case 8:
		return ((p8)(p))(args[3], args[4], args[5], args[6], args[7], args[8],
			args[9], args[10]);
	case 9:
		return ((p9)(p))(args[3], args[4], args[5], args[6], args[7], args[8],
			args[9], args[10], args[11]);
	case 10:
		return ((p10)(p))(args[3], args[4], args[5], args[6], args[7], args[8],
			args[9], args[10], args[11], args[12]);
	case 11:
		return ((p11)(p))(args[3], args[4], args[5], args[6], args[7], args[8],
			args[9], args[10], args[11], args[12], args[13]);
	case 12:
		return ((p12)(p))(args[3], args[4], args[5], args[6], args[7], args[8],
			args[9], args[10], args[11], args[12], args[13], args[14]);
	}
}

static int message_loop(uintptr hdlg);

typedef unsigned int uint32;
long traccel(uintptr ob, uintptr msg);
uintptr queryIDispatch(uintptr iunk);
uintptr createInstance(char* progid);
uintptr invoke(
	uintptr idisp, uintptr name, uintptr flags, uintptr args, uintptr result);
long release(uintptr iunk);
long EmbedBrowserObject(uintptr hwnd, uintptr pBrowserObject, uintptr pPtr);
void UnEmbedBrowserObject(uintptr browserObject, uintptr ptr);

#undef UNICODE
#undef _UNICODE
#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <synchapi.h>
#include <ole2.h>

static DWORD main_threadid = 0;

typedef struct {
	char* text;
	char* caption;
} Mbargs;

DWORD WINAPI message_thread(void* p) {
	MessageBox(0, ((Mbargs*) p)->text, ((Mbargs*) p)->caption,
		MB_OK | MB_TASKMODAL | MB_TOPMOST | MB_SETFOREGROUND);
	return 0;
}

void msgbox(char* text, char* caption) {
	Mbargs args;
	args.text = text;
	args.caption = caption;
	HANDLE thread =
		CreateThread(NULL, 0, message_thread, (void*) &args, 0, NULL);
	if (thread)
		WaitForSingleObject(thread, INFINITE);
	free(text);
}

void alert(char* msg) {
	msgbox(msg, "Alert");
}

void fatal(char* msg) {
	msgbox(msg, "FATAL");
}

const int CTRL_BREAK_ID = 1; // arbitrary value passed to RegisterHotKey

static int interrupt() {
	MSG msg;

	int hotkey = 0;
	if (HIWORD(GetQueueStatus(QS_HOTKEY))) {
		while (PeekMessage(&msg, NULL, WM_HOTKEY, WM_HOTKEY, PM_REMOVE))
			if (msg.wParam == CTRL_BREAK_ID)
				hotkey = 1;
	}
	return hotkey;
}

#include <io.h>
const int STDERR = 2;

uintptr interact() {
	if (GetCurrentThreadId() != main_threadid) {
		const char* msg = "FATAL: interact called from different thread\r\n";
		write(STDERR, msg, strlen(msg));
		exit(1);
	}
	for (;;) {
		// these are the messages sent from go-side to c-side
		switch (args[0]) {
		case msg_syscall:
			args[1] = syscall();
			args[0] = msg_result;
			break;
		case msg_msgloop:
			message_loop(args[1]);
			args[0] = msg_result;
			break;
		case msg_traccel:
			args[1] = traccel(args[1], args[2]);
			args[0] = msg_result;
			break;
		case msg_queryidispatch:
			args[1] = queryIDispatch(args[1]);
			args[0] = msg_result;
			break;
		case msg_createinstance:
			args[1] = createInstance((char*) args[1]);
			args[0] = msg_result;
			break;
		case msg_invoke:
			args[1] = invoke(args[1], args[2], args[3], args[4], args[5]);
			args[0] = msg_result;
			break;
		case msg_release:
			args[1] = release(args[1]);
			args[0] = msg_result;
			break;
		case msg_interrupt:
			args[1] = interrupt();
			args[0] = msg_result;
			break;
		case msg_embedbrowserobject:
			args[1] = EmbedBrowserObject(args[1], args[2], args[3]);
			args[0] = msg_result;
			break;
		case msg_unembedbrowserobject:
			UnEmbedBrowserObject(args[1], args[2]);
			args[0] = msg_result;
			break;
		case msg_result:
			return args[1];
		}
		signalAndWait();
	}
}

static CRITICAL_SECTION lock;
static CONDITION_VARIABLE cond = CONDITION_VARIABLE_INIT;

void signalAndWait() {
	WakeConditionVariable(&cond);
	SleepConditionVariableCS(&cond, &lock, INFINITE);
}

const int stack_size = 16 * 1024; // ???

enum { END_MSG_LOOP = 0xebb };

static int volatile ticks = 0;
static HCURSOR volatile prev_cursor = 0;
static HHOOK hook = 0;

static DWORD WINAPI timer_thread(LPVOID lpParameter) {
	HCURSOR wait_cursor = LoadCursor(NULL, IDC_WAIT);
	AttachThreadInput(GetCurrentThreadId(), main_threadid, TRUE);
	for (;;) {
		Sleep(25); // milliseconds
		if (ticks > 0 && --ticks == 0)
			prev_cursor = SetCursor(wait_cursor);
	}
}
static LRESULT CALLBACK message_hook(int code, WPARAM wParam, LPARAM lParam) {
	ticks = 0; // stop timer
	return CallNextHookEx(hook, code, wParam, lParam);
}

static int message_loop(uintptr hdlg) {
	MSG msg;
	for (;;) {
		ticks = 0; // stop timer
		if (prev_cursor) {
			// restore cursor
			POINT pt;
			GetCursorPos(&pt);
			SetCursorPos(pt.x, pt.y);
		}
		if (-1 == GetMessageA(&msg, 0, 0, 0))
			continue; // ignore error ???
		ticks = 5;    // start timer
		prev_cursor = 0;
		if (msg.message == WM_QUIT) {
			if (hdlg != 0)
				PostQuitMessage(msg.wParam);
			return msg.wParam;
		}
		if (hdlg != 0 && (uintptr)(msg.hwnd) == hdlg &&
			msg.message == WM_NULL && msg.wParam == END_MSG_LOOP &&
			msg.lParam == END_MSG_LOOP)
			return 0;
		HWND window = GetAncestor(msg.hwnd, GA_ROOT);
		if (window != 0) {
			HACCEL haccel = (HACCEL) GetWindowLongPtrA(window, GWLP_USERDATA);
			if (haccel != 0 && TranslateAcceleratorA(window, haccel, &msg))
				continue;
			if (IsDialogMessage(window, &msg))
				continue;
		}
		TranslateMessage(&msg);
		DispatchMessageA(&msg);
	}
}

static BOOL CALLBACK destroy_func(HWND hwnd, LPARAM lParam) {
	DestroyWindow(hwnd);
	return TRUE; // continue enumeration
}

static void destroy_windows() {
	EnumWindows(destroy_func, (LPARAM) NULL);
}

// timer is called by a Windows timer so it will get called
// even if a Windows message loop is running e.g. in MessageBox
static VOID CALLBACK timer(
	HWND hwnd, UINT uMsg, UINT_PTR idEvent, DWORD dwTime) {
	if (ticks <= 3) {
		// our message loop isnt running
		args[0] = msg_runongoside;
		interact();
	}
}

const int timerIntervalMS = 50;
uintptr helperHwnd = 0; // set by setupHelper
const WPARAM notifyWparam = 0xffffffff;
const WPARAM sunappWparam = 0xeeeeeeee;

static LRESULT CALLBACK helperWndProc(
	HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam) {
	if (uMsg == WM_USER && wParam == notifyWparam) {
		args[0] = msg_runongoside;
		interact();
	} else if (uMsg == WM_USER && wParam == sunappWparam) {
		buf_t* buf = (buf_t*) lParam;
		args[0] = msg_sunapp;
		args[1] = (uintptr) buf->buf;
		buf->buf = (char*) interact();
		buf->size = args[2];
	} else {
		return DefWindowProc(hwnd, uMsg, wParam, lParam);
	}
	return 0;
}

static int setupHelper() {
	WNDCLASS wc;
	memset(&wc, 0, sizeof wc);
	wc.lpszClassName = "helper";
	wc.lpfnWndProc = helperWndProc;
	if (!RegisterClass(&wc)) {
		return FALSE;
	}

	HWND hwnd = CreateWindow("helper", "helper", WS_OVERLAPPEDWINDOW, 0, 0, 0,
		0, HWND_MESSAGE, NULL, NULL, NULL);
	if (!hwnd) {
		return FALSE;
	}

	helperHwnd = (uintptr) hwnd;
	return TRUE;
}

int sunapp_register_classes();

static DWORD WINAPI thread(LPVOID lpParameter) {
	OleInitialize(NULL);
	sunapp_register_classes();
	RegisterHotKey(0, CTRL_BREAK_ID, MOD_CONTROL, VK_CANCEL);
	main_threadid = GetCurrentThreadId();
	hook = SetWindowsHookExA(WH_GETMESSAGE, message_hook, 0, main_threadid);
	CreateThread(NULL, 8192, timer_thread, 0, 0, 0);
	setupHelper();
	signalAndWait();
	interact(); // allow go side to run init, finishing with result
	SetTimer(0, 0, timerIntervalMS, timer);
	int exitcode = message_loop(0);
	destroy_windows();
	args[0] = msg_shutdown;
	args[1] = exitcode;
	interact();
}

DWORD threadid;

void start() {
	Scintilla_RegisterClasses(GetModuleHandle(NULL));
	InitializeCriticalSection(&lock);
	EnterCriticalSection(&lock);
	CreateThread(NULL, stack_size, thread, 0, 0, &threadid);
	EnterCriticalSection(&lock); // wait for thread to be in signalAndWait
}

// suneidoAPP is called by sunapp.cpp
buf_t suneidoAPP(char* url) {
	buf_t buf;
	if (GetCurrentThreadId() != main_threadid) {
		buf.buf = url;
		SendMessageA((HWND) helperHwnd, WM_USER, sunappWparam, (LPARAM) &buf);
		return buf;
	}
	args[0] = msg_sunapp;
	args[1] = (uintptr) url;
	buf.buf = (char*) interact();
	buf.size = args[2];
	return buf;
}
