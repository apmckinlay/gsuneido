#include "MyWebView2Base.h"
#include "MyContextMenuRequestedEventHandler.h"
#include "MyDOMContentLoadedEventHandler.h"
#include "MyAcceleratorKeyPressedHandler.h"
#include "MyCustomSchemeHandler.h"
#include "MyNewWindowRequestedhandler.h"

// Define the ICoreWebView2CreateCoreWebView2ControllerCompletedHandler callback class
class MyCreateControllerCompletedHandler : public ICoreWebView2CreateCoreWebView2ControllerCompletedHandler {
public:
    // Constructor
    MyCreateControllerCompletedHandler(HWND hwnd, MyBrowserObject* pBrowserObject) : refCount(1), hwnd(hwnd), pBrowserObject(pBrowserObject) {
        pBrowserObject->AddRef();
    }

    ~MyCreateControllerCompletedHandler() {
        pBrowserObject->Release();
    }

    // IUnknown interface methods
    STDMETHODIMP QueryInterface(REFIID riid, void** ppvObject) override {
        if (riid == IID_IUnknown || riid == IID_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler) {
            *ppvObject = static_cast<ICoreWebView2CreateCoreWebView2ControllerCompletedHandler*>(this);
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

    // ICoreWebView2CreateCoreWebView2ControllerCompletedHandler::Invoke method
    STDMETHODIMP Invoke(HRESULT result, ICoreWebView2Controller* controller) override {
        if (FAILED(result)) {
            pBrowserObject->onReady(CREATING_CONTROLLER, result);
            return S_OK;
        }

        ICoreWebView2 *webview2;
        HRESULT rtn = controller->get_CoreWebView2(&webview2);
        if (FAILED(rtn)) {
            pBrowserObject->onReady(CREATING_WEBVIEW, rtn);
            return S_OK;
        }

        rtn = RegisterCustomScheme(webview2);
        if (FAILED(rtn)) {
            webview2->Release();
            pBrowserObject->onReady(REGISTER_CUSTOM_SCHEME, rtn);
            return S_OK;
        }

        ICoreWebView2_2 *webview2_2;
        rtn = webview2->QueryInterface(IID_ICoreWebView2_2, (void**)&webview2_2);
        if (FAILED(rtn) || !webview2_2) {
            webview2->Release();
            pBrowserObject->onReady(ADD_LOAD_EVENT_LISTENER, rtn);
            return S_OK;
        }
        EventRegistrationToken token;
        MyDOMContentLoadedEventHandler* eventListener = new MyDOMContentLoadedEventHandler(pBrowserObject);
        rtn = webview2_2->add_DOMContentLoaded(eventListener, &token);
        if (FAILED(rtn)) {
            webview2->Release();
            webview2_2->Release();
            eventListener->Release();
            pBrowserObject->onReady(ADD_LOAD_EVENT_LISTENER, rtn);
            return S_OK;
        }
        eventListener->Release();

        MyNavigationCompletedEventHandler* eventListener2 = new MyNavigationCompletedEventHandler(pBrowserObject);
        rtn = webview2_2->add_NavigationCompleted(eventListener2, &token);
        if (FAILED(rtn)) {
            webview2->Release();
            webview2_2->Release();
            eventListener2->Release();
            pBrowserObject->onReady(ADD_LOAD_EVENT_LISTENER, rtn);
            return S_OK;
        }
        eventListener2->Release();

        MyAcceleratorKeyPressedHandler* accHandler = new MyAcceleratorKeyPressedHandler(pBrowserObject);
        rtn = controller->add_AcceleratorKeyPressed(accHandler, &token);
        if (FAILED(rtn)) {
            webview2->Release();
            webview2_2->Release();
            accHandler->Release();
            pBrowserObject->onReady(ADD_ACCELERATOR_HANDLER, rtn);
            return S_OK;
        }
        webview2_2->Release();
        accHandler->Release();

        ICoreWebView2_11 *webview2_11;
        rtn = webview2->QueryInterface(IID_ICoreWebView2_11, (void**)&webview2_11);
        if (FAILED(rtn) || !webview2_11) {
            webview2->Release();
            webview2_2->Release();
            pBrowserObject->onReady(ADD_CONTEXTMENU_HANDLER, rtn);
            return S_OK;
        }
        MyContextMenuRequestedEventHandler* contextMenuHandler = new MyContextMenuRequestedEventHandler(pBrowserObject);
        rtn = webview2_11->add_ContextMenuRequested(contextMenuHandler, &token);
        if (FAILED(rtn)) {
            webview2->Release();
            webview2_11->Release();
            contextMenuHandler->Release();
            pBrowserObject->onReady(ADD_CONTEXTMENU_HANDLER, rtn);
            return S_OK;
        }
        webview2_11->Release();
        contextMenuHandler->Release();

        ICoreWebView2_13 *webview2_13;
        rtn = webview2->QueryInterface(IID_ICoreWebView2_13, (void**)&webview2_13);
        if (SUCCEEDED(rtn) && webview2_13 != nullptr) {
            ICoreWebView2Profile *profile;
            if (SUCCEEDED(webview2_13->get_Profile(&profile) && profile != nullptr)) {
                profile->put_PreferredColorScheme(COREWEBVIEW2_PREFERRED_COLOR_SCHEME_LIGHT);
                profile->Release();
            }
            webview2_13->Release();
        }

        ICoreWebView2Settings *settings;
        if (SUCCEEDED(webview2->get_Settings(&settings) && settings != nullptr)) {
            settings->put_IsStatusBarEnabled(FALSE);
            settings->Release();
        }

        webview2->add_NewWindowRequested(new MyNewWindowRequestedHandler(), &token);

        pBrowserObject->SetController(controller);
        pBrowserObject->SetWebView2(webview2);
        pBrowserObject->onReady(FINISH, OK);
        webview2->Release();
        return S_OK;
    }

private:
    ULONG refCount; // Reference count for COM interface
    HWND hwnd;
    MyBrowserObject* pBrowserObject;
};