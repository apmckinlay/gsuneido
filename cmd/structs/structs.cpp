#include <windows.h>
#include <commctrl.h>
#include <cstdio>
#include <Richedit.h>
#include <ShlObj.h>
#include "../cSuneido/vs2019scintilla/include/Scintilla.h"

#if _WIN64
#define printSize(t) printf("_ = x[unsafe.Sizeof(" #t "{}) - %zd]\n", sizeof(t))
#else
#define printSize(t) printf("Assert(" #t ".Size() is: %zd, msg: '" #t "')\n", sizeof(t))
#endif

int main(int argc, char** argv)
{
	//printf("TCITEMA offsetof(pszText) %zd\n", offsetof(TCITEMA, pszText));

	printSize(ACCEL);
	printSize(BITMAP);
	printSize(BITMAPINFOHEADER);
	printSize(BROWSEINFO);
	printSize(CHARRANGE);
	printSize(CHOOSECOLOR);
	printSize(CHOOSEFONTA);
	printSize(CWPRETSTRUCT);
	printSize(DEVMODE);
	printSize(DEVNAMES);
	printSize(DOCINFO);
	printSize(DRAWITEMSTRUCT);
	printSize(DRAWTEXTPARAMS);
	printSize(EDITBALLOONTIP);
	printSize(FIXED);
	printSize(FLASHWINFO);
	printSize(GLYPHMETRICS);
	printSize(GUID);
	printSize(HDHITTESTINFO);
	printSize(HDITEM);
	printSize(IMAGEINFO);
	printSize(INITCOMMONCONTROLSEX);
	printSize(LOGBRUSH);
	printSize(LOGFONT);
	printSize(LVCOLUMN);
	printSize(LVHITTESTINFO);
	printSize(LVITEM);
	printSize(MENUITEMINFO);
	printSize(MINMAXINFO);
	printSize(MSG);
	printSize(NMDAYSTATE);
	printSize(NMHDR);
	printSize(NMLISTVIEW);
	printSize(NMLVDISPINFO);
	printSize(NMTTDISPINFO);
	printSize(NMTVDISPINFO);
	printSize(NOTIFYICONDATA);
	printSize(OPENFILENAME);
	printSize(OSVERSIONINFOEX);
	printSize(PAGESETUPDLG);
	printSize(PAINTSTRUCT);
	printSize(POINT);
	printSize(PRINTDLG);
	printSize(RECT);
	printSize(SCNotification);
	printSize(SCROLLINFO);
	printSize(SHELLEXECUTEINFO);
	printSize(SYSTEMTIME);
	printSize(TCITEM);
	printSize(TEXTMETRIC);
	printSize(TEXTRANGE);
	printSize(TOOLINFO);
	printSize(TPMPARAMS);
	printSize(TRACKMOUSEEVENT);
	printSize(TVINSERTSTRUCT);
	printSize(TVITEM);
	printSize(WINDOWPLACEMENT);
	printSize(WNDCLASS);

	printSize(VARIANT);

	TOOLINFOA ti;
	ti.cbSize = 0x11111111;
	ti.uFlags = 0x22222222;
	ti.hwnd = (HWND)0x3333333333333333;
	ti.uId = 0x4444444444444444;
	ti.rect.left = 0x11111111;
	ti.rect.top = 0x22222222;
	ti.rect.right = 0x33333333;
	ti.rect.bottom = 0x44444444;
	ti.hinst = (HINSTANCE)0x5555555555555555;
	ti.lpszText = (LPSTR)0x6666666666666666;
	ti.lParam = 0x7777777777777777;
	ti.lpReserved = (void*)0x8888888888888888;
	for (auto i = 0Ui64; i < sizeof(TOOLINFOA); i++) {
		printf("%02x\n", ((char*)&ti)[i]);
		if (i % 4 == 3)
			printf("\n");
	}

	return 0;
}