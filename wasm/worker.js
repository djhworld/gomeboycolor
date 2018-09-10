importScripts('wasm_exec.js');


var ROM = "";

function init(e) {
    startEmulator(e.data);
}

addEventListener("message", init, false);

function startEmulator(rom) {
    ROM = rom;
    removeEventListener("message", init, false);
    if (!WebAssembly.instantiateStreaming) { // polyfill
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }




    const go = new Go();

    let mod, inst;
    WebAssembly.instantiateStreaming(fetch("gbc.wasm"), go.importObject).then((result) => {
        mod = result.module;
        inst = result.instance;
        go.run(inst);

    });
}