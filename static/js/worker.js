if (!WebAssembly.instantiateStreaming) { // polyfill
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

importScripts('wasm_exec.js');
const GOMEBOY_COLOR_WASM = "../wasm/gbc.wasm";

const go = new Go();
let mod, inst;
WebAssembly.instantiateStreaming(fetch(GOMEBOY_COLOR_WASM), go.importObject).then((result) => {
    mod = result.module;
    inst = result.instance;
    go.run(inst);
});