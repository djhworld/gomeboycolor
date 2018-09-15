if (!WebAssembly.instantiateStreaming) { // polyfill
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

importScripts('wasm_exec.js');
importScripts('base64js.min.js');
const GOMEBOY_COLOR_WASM = "../wasm/gbc.wasm";
const SCREEN_UPDATE = "screen-update";

// uses transferable on post message
function sendScreenUpdate(bs64) {
    var buf =  new Uint8ClampedArray(base64js.toByteArray(bs64)).buffer;
    postMessage([SCREEN_UPDATE, buf], [buf]);
}

const go = new Go();
let mod, inst;
WebAssembly.instantiateStreaming(fetch(GOMEBOY_COLOR_WASM), go.importObject).then((result) => {
    mod = result.module;
    inst = result.instance;
    go.run(inst);
});