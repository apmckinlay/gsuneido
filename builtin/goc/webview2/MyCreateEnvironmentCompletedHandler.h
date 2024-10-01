#include "MyWebview2Base.h"
#include "MyCreateControllerCompletedHandler.h"

// Define the ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler interface
class MyCreateEnvironmentCompletedHandler : public ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler {
public:
    // Constructor
    MyCreateEnvironmentCompletedHandler(HWND hwnd, MyBrowserObject* pBrowserObject) : refCount(1), hwnd(hwnd), pBrowserObject(pBrowserObject) {
        pBrowserObject->AddRef();
    }

    ~MyCreateEnvironmentCompletedHandler() {
        pBrowserObject->Release();
    }

    // IUnknown interface methods
    STDMETHODIMP QueryInterface(REFIID riid, void** ppvObject) override {
        if (riid == IID_IUnknown || riid == IID_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) {
            *ppvObject = static_cast<ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler*>(this);
            AddRef();
            return S_OK;
        }
        *ppvObject = nullptr;
        return E_NOINTERFACE;
    }

    STDMETHODIMP_(ULONG) AddRef() override {
        return InterlockedIncrement(&refCount);
    }

    STDMETHODIMP_(ULONG) Release() override {
        ULONG newCount = InterlockedDecrement(&refCount);
        if (newCount == 0) {
            delete this;
        }
        return newCount;
    }

    // ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler::Invoke method
    STDMETHODIMP Invoke(HRESULT result, ICoreWebView2Environment* created_environment) override {
        // Handle the result of creating the WebView2 environment
        if (SUCCEEDED(result)) {
            MyCreateControllerCompletedHandler* handler = new MyCreateControllerCompletedHandler(hwnd, pBrowserObject);
            HRESULT rtn = created_environment->CreateCoreWebView2Controller(hwnd, handler);
            if (FAILED(rtn)) {
                pBrowserObject->onReady(CREATING_CONTROLLER, rtn);
            }
            handler->Release();
        } else {
            pBrowserObject->onReady(CREATING_ENVIRONMENT, result);
        }

        return S_OK;
    }

private:
    ULONG refCount; // Reference count for COM interface
    HWND hwnd;
    MyBrowserObject* pBrowserObject;
};