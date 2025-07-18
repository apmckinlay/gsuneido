// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

#include "cside.h"
#include <stdatomic.h>

int Scintilla_RegisterClasses(void* hInstance);

typedef unsigned int uint32;
long traccel(uintptr ob, uintptr msg);
uintptr CreateLexer(char* name);

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

int interrupt() {
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

enum { END_MSG_LOOP = 0xebb };

static _Atomic int ticks = 0;
static _Atomic HCURSOR prev_cursor = 0;
static HHOOK hook = 0;

static DWORD WINAPI cursor_thread(LPVOID lpParameter) {
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

int message_loop(uintptr hdlg) { // hdlg is 0 for main message loop
	MSG msg;
	for (;;) {
		ticks = 0; // stop timer
		if (prev_cursor) {
			// restore cursor
			SetCursor(prev_cursor);
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
	EnumThreadWindows(main_threadid, destroy_func, 0);
}

uintptr helperHwnd = 0; // set by setupHelper

const UINT sunappMsg = WM_USER + 1;

extern void SuneidoAPP(buf_t* buf);

static LRESULT CALLBACK helperWndProc(
	HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam) {
	if (uMsg == sunappMsg) {
		buf_t* buf = (buf_t*) lParam;
		SuneidoAPP(buf);
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

void setup() {
	Scintilla_RegisterClasses(GetModuleHandle(NULL));
	OleInitialize(NULL);
	sunapp_register_classes();
	RegisterHotKey(0, CTRL_BREAK_ID, MOD_CONTROL, VK_CANCEL);
	main_threadid = GetCurrentThreadId();
	hook = SetWindowsHookExA(WH_GETMESSAGE, message_hook, 0, main_threadid);
	CreateThread(NULL, 8192, cursor_thread, 0, 0, 0);
	setupHelper();
	SetUnhandledExceptionFilter(filter);
}

int run() {
	int exitcode = message_loop(0);
	destroy_windows();
	return exitcode;
}

uintptr createLexer(char* name) {
	return CreateLexer(name);
}

// suneidoAPP is called by sunapp.cpp
buf_t suneidoAPP(char* url) {
	buf_t buf;
	if (GetCurrentThreadId() != main_threadid) {
		buf.data = url;
		SendMessageA((HWND) helperHwnd, sunappMsg, 0, (LPARAM) &buf);
		return buf;
	}
	buf.data = url;
	SuneidoAPP(&buf);
	return buf;
}
