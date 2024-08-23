// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

#include "cside.h"
#include <stdatomic.h>

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
uintptr CreateLexer(const char* name);

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
	int type;
	int result;
} Mbargs;

DWORD WINAPI message_thread(void* p) {
	Mbargs *mb = (Mbargs*) p;
	mb->result = MessageBox(0, mb->text, mb->caption,
		mb->type | MB_TASKMODAL | MB_TOPMOST | MB_SETFOREGROUND);
	return 0;
}

int msgbox(char* text, char* caption, int type) {
	Mbargs mb;
	mb.text = text;
	mb.caption = caption;
	mb.type = type;
	HANDLE thread =
		CreateThread(NULL, 0, message_thread, (void*) &mb, 0, NULL);
	if (thread)
		WaitForSingleObject(thread, INFINITE);
	free(text);
	return mb.result;
}

int alert(char* msg, int type) {
	return msgbox(msg, "Alert", type);
}

void fatal(char* msg) {
	msgbox(msg, "FATAL", 0);
}

void errmsg(const char* msg) {
	DWORD nw;
	WriteFile(GetStdHandle(STD_ERROR_HANDLE), msg, strlen(msg), &nw, 0);
}

void msgexit(const char* msg) {
	errmsg(msg);
	exit(1);
}

const int CTRL_BREAK_ID = 1; // arbitrary value passed to RegisterHotKey

static int interrupt() {
	MSG msg;
	int n = 0;
	int hotkey = 0;
	if (HIWORD(GetQueueStatus(QS_HOTKEY))) {
		while (PeekMessage(&msg, NULL, WM_HOTKEY, WM_HOTKEY, PM_REMOVE)) {
			if (msg.wParam == CTRL_BREAK_ID)
				hotkey = 1;
			if (++n > 100)
				msgexit("FATAL: interrupt too many loops\r\n");
		}
	}
	return hotkey;
}

// interact does a call *to* the other side with signalAndWait
// which blocks us and unblocks the other side.
// signalAndWait will return when the other side does signalAndWait.
// We will get the result via msg_result and return.
// However in between there may be calls *from* the other side
// which is why there is the loop.
uintptr interact() {
	if (GetCurrentThreadId() != main_threadid)
		msgexit("FATAL: interact called from different thread\r\n");
	for (;;) {
		signalAndWait(); // block us and unblock the other side
		// these are the messages sent from go-side to c-side
		switch (args[0]) {
		case msg_syscall:
			args[1] = syscall();
			// args[2] = GetLastError();
			args[0] = msg_result;
			break;
		case msg_msgloop:
			message_loop(args[1]);
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
        case msg_createlexer:
            args[1] = CreateLexer((const char*)args[1]);
            args[0] = msg_result;
            break;
		case msg_setupconsole:
			args[1] = DeleteMenu(GetSystemMenu(GetConsoleWindow(), 0),
				SC_CLOSE, MF_BYCOMMAND);
			args[0] = msg_result;
		case msg_result:
			args[0] = msg_none;
			return args[1];
		}
	}
}

static CRITICAL_SECTION lock;
static CONDITION_VARIABLE cond = CONDITION_VARIABLE_INIT;

void signalAndWait() {
	WakeConditionVariable(&cond); // allow other side to run
	SleepConditionVariableCS(&cond, &lock, INFINITE);
}

const int stack_size = 16 * 1024; // ???

enum { END_MSG_LOOP = 0xebb };

static _Atomic int ticks = 0;
static _Atomic HCURSOR prev_cursor = 0;
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
		if (msg.hwnd != window && msg.message == WM_KEYDOWN) {
			uintptr ptr = (uintptr) GetWindowLongPtrA(msg.hwnd, GWLP_USERDATA);
			if (ptr && S_OK == traccel(ptr, (uintptr)&msg))
				continue;
		}
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
// But timers are the lowest priority,
// so it will only be called when there are no other messages to process.
static VOID CALLBACK timer(
	HWND hwnd, UINT uMsg, UINT_PTR idEvent, DWORD dwTime) {
	args[0] = msg_runongoside;
	interact();
}

const int timerIntervalMS = 10;
uintptr helperHwnd = 0; // set by setupHelper

const UINT notifyMsg = WM_USER;
const UINT sunappMsg = WM_USER + 1;

static LRESULT CALLBACK helperWndProc(
	HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam) {
	if (uMsg == notifyMsg) {
		args[0] = msg_runongoside;
		interact();
		return 0;
	} else if (uMsg == sunappMsg) {
		buf_t* buf = (buf_t*) lParam;
		args[0] = msg_sunapp;
		args[1] = (uintptr) buf->buf;
		buf->buf = (char*) interact();
		buf->size = args[2];
		return 0;
	}
	return DefWindowProc(hwnd, uMsg, wParam, lParam);
}

static int setupHelper() {
	WNDCLASS wc;
	memset(&wc, 0, sizeof wc);
	wc.lpszClassName = "helper";
	wc.lpfnWndProc = helperWndProc;
	if (!RegisterClass(&wc)) {
		return FALSE;
	}

	HWND hwnd = CreateWindow("helper", "helper", WS_OVERLAPPEDWINDOW,
		0, 0, 0, 0, HWND_MESSAGE, NULL, NULL, NULL);
	if (!hwnd) {
		return FALSE;
	}

	helperHwnd = (uintptr) hwnd;
	return TRUE;
}

int sunapp_register_classes();

#include <stdio.h>

static LONG WINAPI filter(EXCEPTION_POINTERS* p_info) {
	unsigned int e = p_info->ExceptionRecord->ExceptionCode;
	char buf[64];
	snprintf(buf, sizeof buf, "FATAL: unhandled C exception %x\r\n", e);
	msgexit(buf);
	return EXCEPTION_EXECUTE_HANDLER;
}

static DWORD WINAPI thread(LPVOID lpParameter) {
	OleInitialize(NULL);
	sunapp_register_classes();
	RegisterHotKey(0, CTRL_BREAK_ID, MOD_CONTROL, VK_CANCEL);
	main_threadid = GetCurrentThreadId();
	hook = SetWindowsHookExA(WH_GETMESSAGE, message_hook, 0, main_threadid);
	SetUnhandledExceptionFilter(filter);
	CreateThread(NULL, 8192, timer_thread, 0, 0, 0);
	setupHelper();
	args[0] = msg_none;
	interact(); // let start() continue and return
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
	SleepConditionVariableCS(&cond, &lock, INFINITE);
}

// suneidoAPP is called by sunapp.cpp
buf_t suneidoAPP(char* url) {
	buf_t buf;
	if (GetCurrentThreadId() != main_threadid) {
		buf.buf = url;
		SendMessageA((HWND) helperHwnd, sunappMsg, 0, (LPARAM) &buf);
		return buf;
	}
	args[0] = msg_sunapp;
	args[1] = (uintptr) url;
	buf.buf = (char*) interact();
	buf.size = args[2];
	return buf;
}
