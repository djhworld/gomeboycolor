class WorkerMessage {
    constructor(msgType, body) {
        this.msgType = msgType;
        this.body = body;
    }
}

function _arrayBufferToBase64(buffer) {
    var binary = '';
    var bytes = new Uint8ClampedArray(buffer);
    var len = bytes.byteLength;
    for (var i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary);
}

function _strToUint8ClampedArray(str) {
    return new Uint8ClampedArray(base64js.toByteArray(str));
}


function sleep(time) {
    return new Promise((resolve) => setTimeout(resolve, time));
}


function handleFileSelect(onLoadHandler) {
    return function(evt) {
        var files = evt.target.files; // FileList object
        if (files.length >= 1) {
            let file = files[0];
            var reader = new FileReader();
            reader.onload = function(e) {
                var ab = reader.result;
                onLoadHandler(file.name, _arrayBufferToBase64(ab));
            }
            reader.readAsArrayBuffer(file);
        }
    };
}







