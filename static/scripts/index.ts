// let functions: DispatchFunctions;
let conn_id: string | undefined = undefined;
let base_url: string | undefined = undefined;
let verbose = false;

type DispatchFunctions = {
    [key: string]: (data: Dispatch) => Dispatch | void;
};

type FnCustom = {
    function: string;
    data: Object;
};

type FnEventListener = {
    id: string;
    target_id: string;
    on: string;
    action: string;
    method: string;
    form_data: string;
    data: Object;
};

type FnRender = {
    target_id: string;
    tag: string;
    inner: boolean;
    outer: boolean;
    append: boolean;
    prepend: boolean;
    html: string;
    event_listeners: FnEventListener[];
};

type FnRedirect = {
    url: string;
};

type FnError = {
    message: string;
};

type Dispatch = {
    function: "render" | "class" | "redirect" | "event" | "error" | "custom";
    id: string;
    key: string;
    conn_id: string;
    handler_id: string;
    action: string;
    label: string;
    event: FnEventListener;
    render: FnRender;
    redirect: FnRedirect;
    custom: FnCustom;
    error: FnError;
};

class Socket {
    private ws: WebSocket | null = null;
    private addr: string | undefined = undefined;
    private key: string | undefined = undefined;

    constructor() {
        let key = localStorage.getItem("fncmp_key");
        if (!key) {
            key = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(
                /[xy]/g,
                function (c) {
                    var r = (Math.random() * 16) | 0,
                        v = c == "x" ? r : (r & 0x3) | 0x8;
                    return v.toString(16);
                }
            );
            localStorage.setItem("fncmp_key", key);
        }
        this.key = key;
        let path = window.location.pathname.split("");
        let path_parsed = "";
        if (path[-1] == "/" || (path.length == 1 && path[0] == "/")) {
            path.pop();
        }
        path_parsed = path.join("");
        
        if (path_parsed == "") {
            path_parsed = "/";
        }
        this.addr = "ws://" + window.location.host + path_parsed + "?fncmp_id=" + this.key;
        this.connect()
    }

    private connect() {
        try {
            this.ws = new WebSocket(this.addr);
        } catch {
            throw new Error("ws: failed to connect to server...");
        }

        this.ws.onopen = function () {};
        this.ws.onclose = function () {};
        this.ws.onerror = function () {};

        this.ws.onmessage = function (event) {
            let d = JSON.parse(event.data) as Dispatch;
            api.Process(this, d);
        };
    }
}

class API {
    private ws: WebSocket | null = null;
    constructor() {
    }

    public Process(ws: WebSocket, d: Dispatch) {
        if (!this.ws) {
            this.ws = ws;
        }
        switch (d.function) {
            case "redirect":
                window.location.href = d.redirect.url;
                break;
            default:
                const result = this.funs[d.function](d);
                this.Dispatch(result);
                break;
        }
    }

    private Dispatch = (data: Dispatch | void) => {
        if (!data) return;
        if (!this.ws) {
            throw new Error("ws: not connected to server...");
        }
        this.ws.send(JSON.stringify(data));
    };

    private funs: DispatchFunctions = {
        render: (d: Dispatch) => {
            let elem: Element | null = null;
            const parsed = new DOMParser().parseFromString(
                d.render.html,
                "text/html"
            ).firstChild as HTMLElement;
            const html = parsed.getElementsByTagName("body")[0].innerHTML;

            if (d.render.tag != "") {
                elem = document.getElementsByTagName(d.render.tag)[0];
                if (!elem) {
                    return this.Error(
                        d,
                        "element with tag not found: " + d.render.tag
                    );
                }
            } else if (d.render.target_id != "") {
                elem = document.getElementById(d.render.target_id);
                if (!elem) {
                    return this.Error(
                        d,
                        "element with target_id not found: " +
                            d.render.target_id
                    );
                }
            } else {
                return this.Error(d, "no target or tag specified");
            }

            if (d.render.inner) {
                elem.innerHTML = html;
            }
            if (d.render.outer) {
                elem.outerHTML = html;
            }
            if (d.render.append) {
                elem.innerHTML += html;
            }
            if (d.render.prepend) {
                elem.innerHTML = html + elem.innerHTML;
            }

            d = this.utils.parseEventListeners(elem, d);
            this.Dispatch(this.utils.addEventListeners(d));
            return;
        },
        class: (d: Dispatch) => {
            return;
        },
        custom: (d: Dispatch) => {
            const result = window[d.custom.function](d.custom.data)
            return;
        },
    };

    private utils = {
        parseEventListeners: (element: Element, d: Dispatch): Dispatch => {
            const events = this.utils.getAttributes(element, "events")
            const listeners = events.map((e) => {
                const event = JSON.parse(e);
                if (!event) return
                return event as FnEventListener[];
            });
            const listeners_flat = listeners.flat();
            const listeners_filtered = listeners_flat.filter((e) => e != null);
            d.render.event_listeners = listeners_filtered;
            return d;
        },
        // Element selectors
        parseFormData: (ev: Event, d: Dispatch) => {
            const form = ev.target as HTMLFormElement;
            const formData = new FormData(form);
            d.event.data = Object.fromEntries(formData.entries());
            return d;
        },
        getAttributes: (elem: Element, attribute: string): string[] => {
            const elems = elem.querySelectorAll(`[${attribute}]`);
            return Array.from(elems).map((el) => el.getAttribute(attribute));
        },
        addEventListeners: (d: Dispatch) => {
            if (!d.render.event_listeners) return;
            // Event listeners
            d.render.event_listeners.forEach((listener: FnEventListener) => {
                let elem = document.getElementById(listener.target_id);
                if (!elem) {
                    this.Error(d, "element not found");
                    return;
                }
                if (elem.firstChild) {
                    elem = elem.firstChild as HTMLElement;
                }
                elem.addEventListener(listener.on, (ev) => {
                    ev.preventDefault();
                    d.function = "event";
                    d.event = listener;
                    switch (listener.on) {
                        case "submit":
                            d = this.utils.parseFormData(ev, d);
                            break;
                        case "pointerdown" || "pointerup" || "pointermove"||  "click" || "contextmenu" || "dblclick" :
                            d.event.data = ParsePointerEvent(ev as PointerEvent);
                            break;
                        case "drag" || "dragend" || "dragenter" || "dragexitcapture" || "dragleave" || "dragover" || "dragstart" || "drop":
                            d.event.data = ParseDragEvent(ev as DragEvent);
                            break;
                        case  "mousedown" || "mouseup" || "mousemove":
                            d.event.data = ParseMouseEvent(ev as MouseEvent);
                            break;
                        case "keydown" || "keyup" || "keypress":
                            d.event.data = ParseKeyboardEvent(ev as KeyboardEvent);
                            break;
                        case "change" || "input" || "invalid" || "reset" || "search" || "select" || "focus" || "blur" || "copy" || "cut" || "paste":
                            d.event.data = ParseEventTarget(ev.target);
                            break;
                        case "touchstart" || "touchend" || "touchmove" || "touchcancel":
                            d.event.data = ParseTouchEvent(ev as TouchEvent & { layerX: number; layerY: number; pageX: number; pageY: number });
                            break;
                        default:
                            d.event.data = ParseEventTarget(ev.target);       
                    }
                    this.Dispatch(d);
                });
            });
        },
    };

    private Error = (d: Dispatch, message: string) => {
        d.function = "error";
        d.error.message = message;
        this.Dispatch(d);
    };
}

function ParseEventTarget(ev: any)  {
    return {
        id: ev.id || "",
        name: ev.name || "",
        tagName: ev.tagName || "",
        innerHTML: ev.innerHTML || "",
        outerHTML: ev.outerHTML || "",
        value: ev.value || "",
    } as Partial<EventTarget>;
}

function ParsePointerEvent(ev: PointerEvent): PointerEventProperties {
    return {
        isTrusted: ev.isTrusted,
        altKey: ev.altKey,
        bubbles: ev.bubbles,
        button: ev.button,
        buttons: ev.buttons,
        cancelable: ev.cancelable,
        clientX: ev.clientX,
        clientY: ev.clientY,
        composed: ev.composed,
        ctrlKey: ev.ctrlKey,
        currentTarget: ParseEventTarget(ev.currentTarget),
        defaultPrevented: ev.defaultPrevented,
        detail: ev.detail,
        eventPhase: ev.eventPhase,
        height: ev.height,
        isPrimary: ev.isPrimary,
        metaKey: ev.metaKey,
        movementX: ev.movementX,
        movementY: ev.movementY,
        offsetX: ev.offsetX,
        offsetY: ev.offsetY,
        pageX: ev.pageX,
        pageY: ev.pageY,
        pointerId: ev.pointerId,
        pointerType: ev.pointerType,
        pressure: ev.pressure,
        relatedTarget: ParseEventTarget(ev.relatedTarget),

    };
}

function ParseTouchEvent(ev: TouchEvent & { layerX: number; layerY: number; pageX: number; pageY: number}): TouchEventProperties & { layerX: number; layerY: number; pageX: number; pageY: number }{
    return {
        changedTouches: Array.from(ev.changedTouches).map((t) => ParseTouch(t)),
        targetTouches: Array.from(ev.targetTouches).map((t) => ParseTouch(t)),
        touches: Array.from(ev.touches).map((t) => ParseTouch(t)),
        layerX: ev.layerX,
        layerY: ev.layerY,
        pageX: ev.pageX,
        pageY: ev.pageY,
    };
}

function ParseTouch(ev: Touch): TouchProperties {
    return {
        clientX: ev.clientX,
        clientY: ev.clientY,
        identifier: ev.identifier,
        pageX: ev.pageX,
        pageY: ev.pageY,
        radiusX: ev.radiusX,
        radiusY: ev.radiusY,
        rotationAngle: ev.rotationAngle,
        screenX: ev.screenX,
        screenY: ev.screenY,
        target: ParseEventTarget(ev.target),
    };
}

function ParseDragEvent(ev: DragEvent): DragEventProperties {
    return {
        isTrusted: ev.isTrusted,
        altKey: ev.altKey,
        bubbles: ev.bubbles,
        button: ev.button,
        buttons: ev.buttons,
        cancelable: ev.cancelable,
        clientX: ev.clientX,
        clientY: ev.clientY,
        composed: ev.composed,
        ctrlKey: ev.ctrlKey,
        currentTarget: ParseEventTarget(ev.currentTarget),
        defaultPrevented: ev.defaultPrevented,
        detail: ev.detail,
        eventPhase: ev.eventPhase,
        metaKey: ev.metaKey,
        movementX: ev.movementX,
        movementY: ev.movementY,
        offsetX: ev.offsetX,
        offsetY: ev.offsetY,
        pageX: ev.pageX,
        pageY: ev.pageY,
        relatedTarget: ParseEventTarget(ev.relatedTarget),
    };
}

function ParseMouseEvent(ev: MouseEvent): MouseEventProperties {
    return {
        isTrusted: ev.isTrusted,
        altKey: ev.altKey,
        bubbles: ev.bubbles,
        button: ev.button,
        buttons: ev.buttons,
        cancelable: ev.cancelable,
        clientX: ev.clientX,
        clientY: ev.clientY,
        composed: ev.composed,
        ctrlKey: ev.ctrlKey,
        currentTarget: ParseEventTarget(ev.currentTarget),
        defaultPrevented: ev.defaultPrevented,
        detail: ev.detail,
        eventPhase: ev.eventPhase,
        metaKey: ev.metaKey,
        movementX: ev.movementX,
        movementY: ev.movementY,
        offsetX: ev.offsetX,
        offsetY: ev.offsetY,
        pageX: ev.pageX,
        pageY: ev.pageY,
        relatedTarget: ParseEventTarget(ev.relatedTarget),
    };
}

function ParseKeyboardEvent(ev: KeyboardEvent): KeyboardEventProperties {
    return {
        isTrusted: ev.isTrusted,
        altKey: ev.altKey,
        bubbles: ev.bubbles,
        cancelable: ev.cancelable,
        code: ev.code,
        composed: ev.composed,
        ctrlKey: ev.ctrlKey,
        currentTarget: ParseEventTarget(ev.currentTarget),
        defaultPrevented: ev.defaultPrevented,
        detail: ev.detail,
        eventPhase: ev.eventPhase,
        isComposing: ev.isComposing,
        key: ev.key,
        location: ev.location,
        metaKey: ev.metaKey,
        repeat: ev.repeat,
        shiftKey: ev.shiftKey,
    };
}

function ParseFormData(ev: SubmitEvent) {
    const form = ev.target as HTMLFormElement;
    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());
    return data;
}

// Event types
type PointerEventProperties = {
    isTrusted: boolean;
    altKey: boolean;
    bubbles: boolean;
    button: number;
    buttons: number;
    cancelable: boolean;
    clientX: number;
    clientY: number;
    composed: boolean;
    ctrlKey: boolean;
    currentTarget: Partial<EventTarget> | null;
    defaultPrevented: boolean;
    detail: number;
    eventPhase: number;
    height: number;
    isPrimary: boolean;
    metaKey: boolean;
    movementX: number;
    movementY: number;
    offsetX: number;
    offsetY: number;
    pageX: number;
    pageY: number;
    pointerId: number;
    pointerType: string;
    pressure: number;
    relatedTarget: Partial<EventTarget> | null;
};

type TouchEventProperties = {
    changedTouches: TouchProperties[];
    targetTouches: TouchProperties[];
    touches: TouchProperties[];
    layerX: number;
    layerY: number;
    pageX: number;
    pageY: number;
};

type TouchProperties = {
    clientX: number;
    clientY: number;
    identifier: number;
    pageX: number;
    pageY: number;
    radiusX: number;
    radiusY: number;
    rotationAngle: number;
    screenX: number;
    screenY: number;
    target: Partial<EventTarget> | null;
};

type DragEventProperties = {
    isTrusted: boolean;
    altKey: boolean;
    bubbles: boolean;
    button: number;
    buttons: number;
    cancelable: boolean;
    clientX: number;
    clientY: number;
    composed: boolean;
    ctrlKey: boolean;
    currentTarget: Partial<EventTarget> | null;
    defaultPrevented: boolean;
    detail: number;
    eventPhase: number;
    metaKey: boolean;
    movementX: number;
    movementY: number;
    offsetX: number;
    offsetY: number;
    pageX: number;
    pageY: number;
    relatedTarget: Partial<EventTarget> | null;
};

type MouseEventProperties = {
    isTrusted: boolean;
    altKey: boolean;
    bubbles: boolean;
    button: number;
    buttons: number;
    cancelable: boolean;
    clientX: number;
    clientY: number;
    composed: boolean;
    ctrlKey: boolean;
    currentTarget: Partial<EventTarget> | null;
    defaultPrevented: boolean;
    detail: number;
    eventPhase: number;
    metaKey: boolean;
    movementX: number;
    movementY: number;
    offsetX: number;
    offsetY: number;
    pageX: number;
    pageY: number;
    relatedTarget: Partial<EventTarget> | null;
};

type KeyboardEventProperties = {
    isTrusted: boolean;
    altKey: boolean;
    bubbles: boolean;
    cancelable: boolean;
    code: string;
    composed: boolean;
    ctrlKey: boolean;
    currentTarget: Partial<EventTarget> | null;
    defaultPrevented: boolean;
    detail: number;
    eventPhase: number;
    isComposing: boolean;
    key: string;
    location: number;
    metaKey: boolean;
    repeat: boolean;
    shiftKey: boolean;
};

type EventProperties = {
    pointer: PointerEventProperties;
    drag: DragEventProperties;
    mouse: MouseEventProperties;
    keyboard: KeyboardEventProperties;
};

type EventTargetProperties = {
    id: string;
    name: string;
    tagName: string;
    innerHTML: string;
    outerHTML: string;
    value: string;
};

const api = new API();
new Socket();