const go = new Go();

async function loadWasm() {
  if (!("instantiateStreaming" in WebAssembly)) {
    const response = await fetch("app.wasm");
    const bytes = await response.arrayBuffer();
    const result = await WebAssembly.instantiate(bytes, go.importObject);
    go.run(result.instance);
    return;
  }

  const result = await WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject);
  go.run(result.instance);
}

loadWasm().catch((error) => {
  console.error("Failed to load WebAssembly", error);
});
