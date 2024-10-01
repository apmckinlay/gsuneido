#undef UNICODE

#include "webview2\MyWebView2Base.h"
#include "webview2\helpers.h"
#include "webview2\MyCoreWebView2CustomSchemeRegistration.h"
#include "webview2\MyCoreWebView2EnvironmentOptions.h"
#include "webview2\MyCreateEnvironmentCompletedHandler.h"

// Define the function pointer type for CreateWebViewEnvironmentWithOptionsInternal
typedef HRESULT (STDMETHODCALLTYPE *CreateWebViewEnvironmentWithOptionsInternalFunc)(
    bool checkRunningInstance, 
    int runtimeType, 
    PCWSTR userDataFolder, 
    ICoreWebView2EnvironmentOptions *environmentOptions, 
    ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler *webViewEnvironmentCreatedHandler
);

long CreateWebView2(HWND hwnd, MyBrowserObject** pBrowserObject, const char* dllPath, const char* userDataFolder, uintptr cb) {
    HRESULT result;

    // Load EmbeddedBrowserWebView.dll
    HMODULE hWebView2Loader = LoadLibrary(dllPath);
    if (hWebView2Loader == NULL) {
        return DLL_NOT_FOUND;
    }

    // Get the address of CreateCoreWebView2Environment
    CreateWebViewEnvironmentWithOptionsInternalFunc CreateCoreWebView2Environment =
        (CreateWebViewEnvironmentWithOptionsInternalFunc)GetProcAddress(hWebView2Loader, "CreateWebViewEnvironmentWithOptionsInternal");
    if (CreateCoreWebView2Environment == NULL) {
        FreeLibrary(hWebView2Loader);
        return PROC_NOT_FOUND;
    }
    
    MyCoreWebView2EnvironmentOptions* options = new MyCoreWebView2EnvironmentOptions();

    const WCHAR* allowedOrigins[1] = {L"*"};
    auto customSchemeRegistration = new MyCoreWebView2CustomSchemeRegistration(L"suneido");
    ICoreWebView2CustomSchemeRegistration* registrations[1] = { customSchemeRegistration };
    result = options->SetCustomSchemeRegistrations(1, registrations);
    if (FAILED(result)) {
        options->Release();
        customSchemeRegistration->Release();
        FreeLibrary(hWebView2Loader);
        return result;
    }
    customSchemeRegistration->Release();

    *pBrowserObject = new MyBrowserObject(hwnd, (GenericCallback)cb);
    MyCreateEnvironmentCompletedHandler* handler = new MyCreateEnvironmentCompletedHandler(hwnd, *pBrowserObject);
    // Call CreateCoreWebView2Environment function

    auto userDataFolder_w = ConvertCharToWChar(userDataFolder);
    result = CreateCoreWebView2Environment(true, 1, userDataFolder_w, options, handler);
    delete[] userDataFolder_w;

    handler->Release();
    options->Release();
    // Clean up: free the DLL
    FreeLibrary(hWebView2Loader);
    return result;
}

long Close(MyBrowserObject* pBrowserObject) {
    return pBrowserObject->Release();
}

long Resize(MyBrowserObject* pBrowserObject, long w, long h) {
    if (pBrowserObject->controller == 0) {
        return NOT_READY;
    }

    RECT bounds{};
    bounds.left = bounds.top = 0;
    bounds.right = w;
    bounds.bottom = h;
    HRESULT result = pBrowserObject->controller->put_Bounds(bounds);
    return result;
}

long Navigate(MyBrowserObject* pBrowserObject, const char* s) {
    if (pBrowserObject->controller == 0) {
        return NOT_READY;
    }
    wchar_t* convertedChars = ConvertCharToWChar(s);
    HRESULT result = pBrowserObject->webview2->Navigate(convertedChars);
    delete[] convertedChars;
    return result;
}

long NavigateToString(MyBrowserObject* pBrowserObject, const char* s) {
    if (pBrowserObject->controller == 0) {
        return NOT_READY;
    }
    wchar_t* convertedChars = ConvertCharToWChar(s);
    HRESULT result = pBrowserObject->webview2->NavigateToString(convertedChars);
    delete[] convertedChars;
    return result;
}



long ExecuteScript(MyBrowserObject* pBrowserObject, const char* script) {
    if (pBrowserObject->controller == 0) {
        return NOT_READY;
    }
    wchar_t* convertedChars = ConvertCharToWChar(script);
    HRESULT result = pBrowserObject->webview2->ExecuteScript(convertedChars, nullptr);
    delete[] convertedChars;
    return result;
}

long GetSource(MyBrowserObject* pBrowserObject, char* dst) {
    if (pBrowserObject->controller == 0) {
        return NOT_READY;
    }
    LPWSTR uri = nullptr;
    HRESULT result = pBrowserObject->webview2->get_Source(&uri);
    if (SUCCEEDED(result)) {
        char* converted = ConvertWcharToChar(uri);
        CoTaskMemFree(uri);
        strcpy(dst, converted);
        delete[] converted;
    }
    return result;
}

long DoPrint(MyBrowserObject* pBrowserObject) {
    if (pBrowserObject->controller == 0) {
        return NOT_READY;
    }
    
    ICoreWebView2_16* webview2_16 = nullptr;
    HRESULT result = pBrowserObject->webview2->QueryInterface(IID_ICoreWebView2_16, (void**)&webview2_16);
    if (FAILED(result) || !webview2_16) {
        return result;
    }

    result = webview2_16->ShowPrintUI(COREWEBVIEW2_PRINT_DIALOG_KIND_BROWSER);
    webview2_16->Release();
    return result;
}

extern "C"
long WebView2(uintptr op, uintptr arg1, uintptr arg2, uintptr arg3, uintptr arg4, uintptr arg5) {
    switch (op) {
        case webview2_create:
            return CreateWebView2((HWND)arg1, (MyBrowserObject**)arg2, (const char*)arg3, (const char*)arg4, arg5);
        case webview2_close:
            return Close((MyBrowserObject*)arg1);
        case webview2_resize:
            return Resize((MyBrowserObject*)arg1, (long)arg2, (long)arg3);
        case webview2_navigate:
            return Navigate((MyBrowserObject*)arg1, (const char*)arg2);
        case webview2_navigate_to_string:
            return NavigateToString((MyBrowserObject*)arg1, (const char*)arg2);
        case webview2_execute_script:
            return ExecuteScript((MyBrowserObject*)arg1, (const char*)arg2);
        case webview2_get_source:
            return GetSource((MyBrowserObject*)arg1, (char*)arg2);
        case webview2_print:
            return DoPrint((MyBrowserObject*)arg1);
    }
    return INVALID_OP;
}

