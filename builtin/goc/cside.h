typedef unsigned long long uintptr;

enum { maxargs = 20 };
extern uintptr args[maxargs];

extern unsigned long threadid;

enum {
	ncb2s = 32,
	ncb3s = 32,
	ncb4s = 1024,
};

uintptr cb2s[ncb2s];
uintptr cb3s[ncb2s];
uintptr cb4s[ncb4s];

void start();
void signalAndWait();

enum {
	msg_result,
	msg_syscall,
	msg_callback2,
	msg_callback3,
	msg_callback4,
	msg_msgloop,
	msg_timerid,
};

typedef uintptr(__stdcall* cb2_t)(uintptr a, uintptr b);
typedef uintptr(__stdcall* cb3_t)(uintptr a, uintptr b, uintptr c);
typedef uintptr(__stdcall* cb4_t)(uintptr a, uintptr b, uintptr c, uintptr d);

uintptr cb2(uintptr i, uintptr a, uintptr b);
uintptr cb3(uintptr i, uintptr a, uintptr b, uintptr c);
uintptr cb4(uintptr i, uintptr a, uintptr b, uintptr c, uintptr d);
