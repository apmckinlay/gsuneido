#include "cside.h"

extern void timerId();

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

static void message_loop(uintptr hdlg);

uintptr interact() {
	for (;;) {
		switch (args[0]) {
		case msg_callback2:
		case msg_callback3:
		case msg_callback4:
		case msg_timerid:
			break;
		case msg_syscall:
			args[1] = syscall();
			args[0] = msg_result;
			break;
		case msg_msgloop:
			message_loop(args[1]);
			args[0] = msg_result;
			break;
		case msg_result:
			return args[1];
		}
		signalAndWait();
	}
}

#undef UNICODE
#undef _UNICODE
#include <windows.h>
#include <synchapi.h>

static CRITICAL_SECTION lock;
static CONDITION_VARIABLE cond = CONDITION_VARIABLE_INIT;

void signalAndWait() {
	WakeConditionVariable(&cond);
	SleepConditionVariableCS(&cond, &lock, INFINITE);
}

const int stack_size = 16 * 1024; // ???

enum { END_MSG_LOOP = 0xebb };

static void message_loop(uintptr hdlg) {
	MSG msg;
	for (;;) {
		if (-1 == GetMessageA(&msg, 0, 0, 0))
			continue; // ignore error ???
		if (msg.message == WM_QUIT) {
			if (hdlg != 0)
				PostQuitMessage(msg.wParam);
			return;
		}
		if (hdlg != 0 && (uintptr)(msg.hwnd) == hdlg &&
			msg.message == WM_NULL && msg.wParam == END_MSG_LOOP &&
			msg.lParam == END_MSG_LOOP)
			return;
		if (msg.message == WM_USER && msg.hwnd == 0) {
			// from SetTimer in another thread
			uintptr id =
				SetTimer(0, 0, (UINT)(msg.wParam), (TIMERPROC)(msg.lParam));
			// timerId(id);
			args[0] = msg_timerid;
			args[1] = id;
			interact();
			continue;
		}
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

static DWORD WINAPI thread(LPVOID lpParameter) {
	signalAndWait();
	interact(); // allow go side to run init, finishing with result
	message_loop(0);
	destroy_windows();
	exit(0);
}

DWORD threadid;

void start() {
	InitializeCriticalSection(&lock);
	EnterCriticalSection(&lock);
	CreateThread(NULL, stack_size, thread, 0, 0, &threadid);
	EnterCriticalSection(&lock); // wait for thread to be in signalAndWait
}
