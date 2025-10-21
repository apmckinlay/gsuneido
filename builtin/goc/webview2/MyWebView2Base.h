#ifndef suneido_webview2_base
#define suneido_webview2_base

#include "WebView2.h"
extern "C" {
#include "../cside.h"
}

const LPCWSTR suneidoPrefix = L"suneido:";
const size_t suneidoPrefixSize = 8;

enum CALLBACK_TYPE {
    ON_READY,
    ON_LOADED,
    ON_ACCEL_KEY_PRESSED,
    ON_CONTEXT_MENU_REQUESTED,
    ON_NAVCOMPLETED,
};

enum webview2_ops {
	webview2_create,
	webview2_close,
	webview2_resize,
	webview2_navigate,
    webview2_navigate_to_string,
    webview2_execute_script,
    webview2_get_source,
    webview2_print,
    webview2_set_focus,
};

enum WEBVIEW_RESULT {
    OK,
    NOT_READY,
    INVALID_OP,
    DLL_NOT_FOUND,
    PROC_NOT_FOUND,
};

enum CREATE_WEBVIEW_STAGE {
    FINISH,
    CREATING_ENVIRONMENT,
    CREATING_CONTROLLER,
    CREATING_WEBVIEW,
    REGISTER_CUSTOM_SCHEME,
    ADD_LOAD_EVENT_LISTENER,
    ADD_ACCELERATOR_HANDLER,
    ADD_CONTEXTMENU_HANDLER,
};

typedef uintptr (__stdcall *GenericCallback)(CALLBACK_TYPE type, uintptr a, uintptr b, uintptr c);

class MyBrowserObject {
public:
    ICoreWebView2Controller* controller;
    ICoreWebView2* webview2;
    HWND hwnd;

    MyBrowserObject(HWND hwnd, GenericCallback cb) : hwnd(hwnd), onCallBackFn(cb) {
        controller = nullptr;
        webview2 = nullptr;
    }

    ~MyBrowserObject() {
        if (webview2 != nullptr) {
            webview2->Release();
        }
        if (controller != nullptr) {
            controller->Close();
            controller->Release();
        }
    }

    void onReady(int stage, int result) {
        if (onCallBackFn != nullptr) {
            onCallBackFn(ON_READY, stage, result, 0);
        }
    }

    void onLoaded() {
        if (onCallBackFn != nullptr) {
            onCallBackFn(ON_LOADED, 0, 0, 0);
        }
    }

    void onNavCompleted() {
        if (onCallBackFn != nullptr) {
            onCallBackFn(ON_NAVCOMPLETED, 0, 0, 0);
        }
    }

    bool onAccelPressed(UINT key) {
        if (onCallBackFn != nullptr) {
            return onCallBackFn(ON_ACCEL_KEY_PRESSED, key, 0, 0);
        }
        return false;
    }

    bool onContextMenuRequested() {
        if (onCallBackFn != nullptr) {
            return onCallBackFn(ON_CONTEXT_MENU_REQUESTED, 0, 0, 0);
        }
        return false;
    }

    void SetController(ICoreWebView2Controller* value) {
        controller = value;
        controller->AddRef();
    }

    void SetWebView2(ICoreWebView2* value) {
        webview2 = value;
        webview2->AddRef();
    }

    ULONG AddRef() {
        return InterlockedIncrement(&refCount);
    }

    ULONG Release() {
        ULONG result = InterlockedDecrement(&refCount);
        if (result == 0) {
            delete this;
        }
        return result;
    }

    long Close() {
        onCallBackFn = nullptr;
        return Release();
    }

private:
    GenericCallback onCallBackFn;
    LONG refCount = 1;
};

#endif