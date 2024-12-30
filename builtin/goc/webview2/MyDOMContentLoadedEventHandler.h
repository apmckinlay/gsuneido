#include "MyWebView2Base.h"

class MyDOMContentLoadedEventHandler : public ICoreWebView2DOMContentLoadedEventHandler {
public:
    // Constructor
    MyDOMContentLoadedEventHandler(MyBrowserObject* pBrowserObject) : m_refCount(1), pBrowserObject(pBrowserObject) {
        pBrowserObject->AddRef();
    }

    ~MyDOMContentLoadedEventHandler() {
        pBrowserObject->Release();
    }

    // IUnknown methods
    STDMETHODIMP QueryInterface(REFIID riid, void** ppvObject) override {
        if (riid == IID_IUnknown || riid == IID_ICoreWebView2DOMContentLoadedEventHandler) {
            *ppvObject = static_cast<ICoreWebView2DOMContentLoadedEventHandler*>(this);
            AddRef();
            return S_OK;
        }
        *ppvObject = nullptr;
        return E_NOINTERFACE;
    }

    STDMETHODIMP_(ULONG) AddRef() override {
        return InterlockedIncrement(&m_refCount);
    }

    STDMETHODIMP_(ULONG) Release() override {
        ULONG refCount = InterlockedDecrement(&m_refCount);
        if (refCount == 0) {
            delete this;
        }
        return refCount;
    }

    // ICoreWebView2DOMContentLoadedEventHandler method
    STDMETHODIMP Invoke(ICoreWebView2* webview, ICoreWebView2DOMContentLoadedEventArgs* args) override {
        pBrowserObject->onLoaded();
        return S_OK;
    }

private:
    // Reference count for COM object
    ULONG m_refCount;
    MyBrowserObject* pBrowserObject;
};

class MyNavigationCompletedEventHandler : public ICoreWebView2NavigationCompletedEventHandler {
public:
    // Constructor
    MyNavigationCompletedEventHandler(MyBrowserObject* pBrowserObject) : m_refCount(1), pBrowserObject(pBrowserObject) {
        pBrowserObject->AddRef();
    }

    ~MyNavigationCompletedEventHandler() {
        pBrowserObject->Release();
    }

    // IUnknown methods
    STDMETHODIMP QueryInterface(REFIID riid, void** ppvObject) override {
        if (riid == IID_IUnknown || riid == IID_ICoreWebView2DOMContentLoadedEventHandler) {
            *ppvObject = static_cast<ICoreWebView2NavigationCompletedEventHandler *>(this);
            AddRef();
            return S_OK;
        }
        *ppvObject = nullptr;
        return E_NOINTERFACE;
    }

    STDMETHODIMP_(ULONG) AddRef() override {
        return InterlockedIncrement(&m_refCount);
    }

    STDMETHODIMP_(ULONG) Release() override {
        ULONG refCount = InterlockedDecrement(&m_refCount);
        if (refCount == 0) {
            delete this;
        }
        return refCount;
    }

    // ICoreWebView2DOMContentLoadedEventHandler method
    STDMETHODIMP Invoke(ICoreWebView2* webview, ICoreWebView2NavigationCompletedEventArgs* args) override {
        pBrowserObject->onNavCompleted();
        return S_OK;
    }

private:
    // Reference count for COM object
    ULONG m_refCount;
    MyBrowserObject* pBrowserObject;
};