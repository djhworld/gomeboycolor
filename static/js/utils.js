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
                var array = new Uint8Array(ab);
                onLoadHandler(file.name, base64js.fromByteArray(array));
            }
            reader.readAsArrayBuffer(file);
        }
    };
}







