const go = new Go();

function downloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  URL.revokeObjectURL(url);
}

function setupControlsToggle() {
  const toggle = document.getElementById("toggle-controls");
  if (!toggle) return;

  toggle.addEventListener("click", () => {
    document.body.classList.toggle("controls-hidden");
    const hidden = document.body.classList.contains("controls-hidden");
    toggle.textContent = hidden ? "SHOW CONTROLS" : "HIDE CONTROLS";
  });
}

function setupThemeToggle() {
  const toggle = document.getElementById("toggle-theme");
  if (!toggle) return;

  const lightBg = "#fffdfb";
  const darkBg = "#141312";

  const stored = localStorage.getItem("visum-theme");
  const prefersDark = window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;

  if (stored === "dark" || (!stored && prefersDark)) {
    document.body.classList.add("theme-dark");
    document.body.classList.remove("theme-light");
  } else if (stored === "light") {
    document.body.classList.add("theme-light");
    document.body.classList.remove("theme-dark");
  }

  const updateBackground = (dark) => {
    const input = document.getElementById("bg-color");
    if (!input) return;
    const current = (input.value || "").toLowerCase();
    const next = dark ? darkBg : lightBg;
    const previous = dark ? lightBg : darkBg;
    input.defaultValue = next;
    if (current === previous) {
      input.value = next;
      input.dispatchEvent(new Event("input", { bubbles: true }));
    }
  };

  updateBackground(document.body.classList.contains("theme-dark"));
  toggle.addEventListener("click", () => {
    const dark = document.body.classList.toggle("theme-dark");
    document.body.classList.toggle("theme-light", !dark);
    localStorage.setItem("visum-theme", dark ? "dark" : "light");
    updateBackground(dark);
  });
}

function setupHeaderMenu() {
  const toggle = document.getElementById("menu-toggle");
  if (!toggle) return;
  const menu = document.getElementById("header-actions-list");

  toggle.addEventListener("click", (event) => {
    event.preventDefault();
    document.body.classList.toggle("menu-open");
  });

  if (menu) {
    menu.addEventListener("click", (event) => {
      const target = event.target;
      if (target && target.closest("button")) {
        document.body.classList.remove("menu-open");
      }
    });
  }

  document.addEventListener("click", (event) => {
    if (!document.body.classList.contains("menu-open")) return;
    const target = event.target;
    if (!menu || menu.contains(target) || toggle.contains(target)) return;
    document.body.classList.remove("menu-open");
  });
}

function setupStickyHeader() {
  const sticky = document.querySelector(".sticky-eyebrow");
  if (!sticky) return;
  const subtitle = sticky.querySelector(".sticky-subtitle");
  const word = sticky.querySelector(".sticky-word");
  const actions = sticky.querySelector(".header-actions");

  const update = () => {
    const max = 160;
    const progress = Math.min(window.scrollY / max, 1);
    sticky.style.setProperty("--subtitle-opacity", (1 - progress).toFixed(2));
    sticky.style.setProperty("--subtitle-translate", `${-6 * progress}px`);
    const stickyWidth = sticky.getBoundingClientRect().width;
    const actionsWidth = actions ? actions.getBoundingClientRect().width : 0;
    const wordWidth = word ? word.getBoundingClientRect().width : 0;
    const available = Math.max(0, stickyWidth - actionsWidth - wordWidth - 24);
    const collapse = available * (1 - progress);
    sticky.style.setProperty("--subtitle-max", `${collapse}px`);
    if (subtitle) {
      const baseWidth = subtitle.scrollWidth || 1;
      const scale = Math.min(1, collapse / baseWidth);
      sticky.style.setProperty("--subtitle-scale", scale.toFixed(3));
    }
    const height = sticky.getBoundingClientRect().height;
    document.documentElement.style.setProperty("--sticky-offset", `${height}px`);
  };

  update();
  window.addEventListener("scroll", update, { passive: true });
  window.addEventListener("resize", update);
}

function setupExportControls() {
  const canvas = document.getElementById("visum-canvas");
  if (!canvas) return;

  const scaleInput = document.getElementById("export-scale");
  const fpsInput = document.getElementById("export-fps");
  const bitrateInput = document.getElementById("export-bitrate");
  const resolutionInput = document.getElementById("export-resolution");
  const includeReadoutInput = document.getElementById("export-include-readout");
  const loopsInput = document.getElementById("export-loops");
  const exportPng = document.getElementById("export-png");
  const exportWebp = document.getElementById("export-webp");
  const exportSvg = document.getElementById("export-svg");
  const exportVideo = document.getElementById("export-video");
  const recordButton = document.getElementById("export-record");
  const stopButton = document.getElementById("export-stop");
  const cancelButton = document.getElementById("export-cancel");
  const statusEl = document.getElementById("export-status");
  const progressWrap = document.querySelector(".export-progress");
  const progressBar = document.querySelector(".export-progress-bar");
  const guardOverlay = document.getElementById("export-guard");
  const guardContinue = document.getElementById("guard-continue");
  const guardStop = document.getElementById("guard-stop");
  if (exportVideo && !exportVideo.dataset.label) {
    exportVideo.dataset.label = exportVideo.textContent || "EXPORT VIDEO (REAL TIME)";
  }

  const clampNumber = (value, min, max) => {
    let next = value;
    if (Number.isFinite(min)) {
      next = Math.max(min, next);
    }
    if (Number.isFinite(max)) {
      next = Math.min(max, next);
    }
    return next;
  };

  const readNumber = (input, fallback) => {
    if (!input) return fallback;
    const parsed = Number.parseFloat(input.value);
    const value = Number.isFinite(parsed) ? parsed : fallback;
    const min = Number.parseFloat(input.min);
    const max = Number.parseFloat(input.max);
    const clamped = clampNumber(value, Number.isFinite(min) ? min : undefined, Number.isFinite(max) ? max : undefined);
    if (Number.isFinite(clamped) && input.value !== String(clamped)) {
      input.value = clamped;
    }
    return clamped;
  };

  const snapToOptions = (value, options) => {
    if (!options || options.length === 0) return value;
    return options.reduce((closest, current) => {
      return Math.abs(current - value) < Math.abs(closest - value) ? current : closest;
    }, options[0]);
  };

  const fpsOptions = [12, 15, 24, 25, 30, 48, 50, 60, 72, 90, 96, 120];
  const bitrateOptions = [2, 4, 8, 12, 20, 35, 50, 80];
  const scaleOptions = [1, 1.5, 2, 3, 4];

  const setStatus = (message, tone) => {
    if (statusEl) {
      statusEl.textContent = message;
      statusEl.classList.toggle("is-success", tone === "success");
    }
  };

  const clearStatus = () => {
    if (!statusEl) return;
    statusEl.textContent = "";
    statusEl.classList.remove("is-success");
  };

  const pulseHaptic = (pattern) => {
    if (typeof navigator !== "undefined" && typeof navigator.vibrate === "function") {
      navigator.vibrate(pattern);
    }
  };

  const readoutText = () => {
    if (!includeReadoutInput || !includeReadoutInput.checked) return "";
    const el = document.getElementById("live-readout");
    if (!el) return "";
    const text = el.textContent || "";
    const match = text.match(/k=[^|]+/);
    return match ? match[0].trim() : "";
  };

  const drawReadout = (ctx, width, height, override) => {
    const readout = override || readoutText();
    if (!readout) return;
    const fontSize = Math.max(12, width * 0.02);
    const ink = getComputedStyle(document.body).getPropertyValue("--ink-soft").trim() || "#cfc7bb";
    ctx.fillStyle = ink;
    ctx.font = `300 ${fontSize}px \"Source Serif 4\", \"Iowan Old Style\", \"Palatino Linotype\", serif`;
    ctx.textAlign = "left";
    ctx.textBaseline = "alphabetic";
    ctx.fillText(readout, 14, height - 14);
  };

  const scaledCanvas = () => {
    const scale = snapToOptions(readNumber(scaleInput, 1), scaleOptions);
    if (scaleInput) scaleInput.value = scale;
    if (scale === 1) {
      return canvas;
    }
    const offscreen = document.createElement("canvas");
    offscreen.width = Math.max(1, Math.round(canvas.width * scale));
    offscreen.height = Math.max(1, Math.round(canvas.height * scale));
    const ctx = offscreen.getContext("2d");
    if (ctx) {
      ctx.imageSmoothingEnabled = true;
      ctx.imageSmoothingQuality = "high";
      ctx.drawImage(canvas, 0, 0, offscreen.width, offscreen.height);
    }
    return offscreen;
  };

  const exportRaster = (type, ext) => {
    const source = scaledCanvas();
    let target = source;
    const readout = readoutText();
    if (readout) {
      const offscreen = document.createElement("canvas");
      offscreen.width = source.width;
      offscreen.height = source.height;
      const ctx = offscreen.getContext("2d");
      if (ctx) {
        ctx.drawImage(source, 0, 0);
        drawReadout(ctx, offscreen.width, offscreen.height);
      }
      target = offscreen;
    }
    const filename = `visum-${Date.now()}.${ext}`;
    const quality = type === "image/webp" ? 0.95 : undefined;
    target.toBlob((blob) => {
      if (!blob) return;
      downloadBlob(blob, filename);
      setStatus("Image saved.", "success");
    }, type, quality);
  };

  const exportSvgFile = () => {
    const rect = canvas.getBoundingClientRect();
    const scale = snapToOptions(readNumber(scaleInput, 1), scaleOptions);
    const width = rect.width * scale;
    const height = rect.height * scale;
    if (typeof window.visumExportSVG !== "function") {
      return;
    }
    const svg = window.visumExportSVG(width, height, Boolean(includeReadoutInput && includeReadoutInput.checked));
    const blob = new Blob([svg], { type: "image/svg+xml" });
    downloadBlob(blob, `visum-${Date.now()}.svg`);
    setStatus("Image saved.", "success");
  };

  if (exportPng) {
    exportPng.addEventListener("click", () => exportRaster("image/png", "png"));
  }
  if (exportWebp) {
    exportWebp.addEventListener("click", () => exportRaster("image/webp", "webp"));
  }
  if (exportSvg) {
    exportSvg.addEventListener("click", exportSvgFile);
  }

  let recorder = null;
  let recordStream = null;
  let chunks = [];
  let recordCanvas = null;
  let recordCtx = null;
  let recordFrame = 0;
  let recordingMode = "";
  let recordingStart = 0;
  let recordingDuration = 0;
  let progressRaf = 0;
  let discardRecording = false;
  let recordFileHandle = null;
  let recordFilename = "";
  let guardTarget = null;
  let guardBypass = false;
  let restoreRunning = null;

  const isRunning = () => {
    const playToggle = document.getElementById("play-toggle");
    if (!playToggle) return false;
    return playToggle.textContent.trim().toUpperCase() === "PAUSE";
  };

  const parseResolution = () => {
    const fallback = { width: canvas.width, height: canvas.height };
    if (fallback.width === 0 || fallback.height === 0) {
      const rect = canvas.getBoundingClientRect();
      fallback.width = Math.max(1, Math.round(rect.width));
      fallback.height = Math.max(1, Math.round(rect.height));
    }
    if (!resolutionInput) return fallback;
    if (resolutionInput.value === "auto") return fallback;
    const parts = resolutionInput.value.split("x");
    if (parts.length !== 2) return fallback;
    const width = Number.parseInt(parts[0], 10);
    const height = Number.parseInt(parts[1], 10);
    if (!Number.isFinite(width) || !Number.isFinite(height)) return fallback;
    return { width, height };
  };

  const drawRecordingFrame = () => {
    if (!recordCtx || !recordCanvas) return;
    const { width, height } = recordCanvas;
    recordCtx.clearRect(0, 0, width, height);
    const bgInput = document.getElementById("bg-color");
    const bg = bgInput ? bgInput.value : "#000000";
    recordCtx.fillStyle = bg;
    recordCtx.fillRect(0, 0, width, height);
    const scale = Math.min(width / canvas.width, height / canvas.height);
    const drawWidth = canvas.width * scale;
    const drawHeight = canvas.height * scale;
    const offsetX = (width - drawWidth) / 2;
    const offsetY = (height - drawHeight) / 2;
    recordCtx.drawImage(canvas, offsetX, offsetY, drawWidth, drawHeight);
    drawReadout(recordCtx, width, height);
    recordFrame = window.requestAnimationFrame(drawRecordingFrame);
  };

  const pickMimeType = () => {
    const preferred = [
      "video/mp4;codecs=avc1.42E01E",
      "video/mp4;codecs=avc1",
      "video/mp4",
    ];
    return preferred.find((type) => window.MediaRecorder && MediaRecorder.isTypeSupported(type)) || "";
  };

  const stopProgressLoop = () => {
    if (progressRaf) {
      window.cancelAnimationFrame(progressRaf);
      progressRaf = 0;
    }
  };

  const updateCaptureState = () => {
    const recording = Boolean(recorder);
    if (recordButton) recordButton.disabled = recording;
    if (exportVideo) exportVideo.disabled = recording;
    if (cancelButton) cancelButton.disabled = !recording;
    if (stopButton) {
      stopButton.disabled = !recording;
      stopButton.classList.toggle("is-active", recording);
    }
    if (exportVideo) {
      exportVideo.classList.toggle("is-working", recordingMode === "export");
      if (recording) {
        exportVideo.textContent = "EXPORTING...";
      } else {
        exportVideo.textContent = exportVideo.dataset.label || "EXPORT VIDEO (REAL TIME)";
      }
    }
    if (progressWrap) {
      progressWrap.classList.toggle("is-active", recording);
      progressWrap.classList.toggle("is-indeterminate", recording && recordingMode === "record");
    }
    if (progressBar && !recording) {
      progressBar.style.width = "0%";
    }
    if (!recording) {
      recordingMode = "";
      recordingStart = 0;
      recordingDuration = 0;
      discardRecording = false;
      guardTarget = null;
      if (guardOverlay) {
        guardOverlay.classList.remove("is-active");
      }
      window.onbeforeunload = null;
      stopProgressLoop();
    } else {
      window.onbeforeunload = () => "Export in progress.";
    }
  };

  const formatDuration = (seconds) => {
    const total = Math.max(0, Math.floor(seconds));
    const mins = Math.floor(total / 60);
    const secs = total % 60;
    return mins > 0 ? `${mins}:${String(secs).padStart(2, "0")}` : `${secs}s`;
  };

  const startProgressLoop = () => {
    stopProgressLoop();
    const tick = () => {
      if (!recorder) return;
      const elapsed = (performance.now() - recordingStart) / 1000;
      if (recordingMode === "export" && recordingDuration > 0) {
        const total = recordingDuration / 1000;
        const percent = Math.min(100, Math.round((elapsed / total) * 100));
        const remaining = Math.max(0, total - elapsed);
        const eta = formatDuration(remaining);
        setStatus(`Exporting... ${percent}% • ETA ${eta}`);
        if (statusEl) {
          statusEl.title = "Real-time export. Keep this tab open to finish.";
        }
        if (progressBar) {
          progressBar.style.width = `${percent}%`;
        }
        if (progressWrap) {
          progressWrap.title = "Real-time export. Keep this tab open to finish.";
        }
      } else {
        setStatus(`Recording... ${formatDuration(elapsed)}`);
        if (statusEl) {
          statusEl.title = "Manual recording. Click stop to download.";
        }
        if (progressWrap) {
          progressWrap.title = "Manual recording. Click stop to download.";
        }
      }
      progressRaf = window.requestAnimationFrame(tick);
    };
    progressRaf = window.requestAnimationFrame(tick);
  };

  const startRecording = (mode, durationMs, restoreState, fileHandle) => {
    if (!recordButton || !stopButton || recorder) return;
    recordingMode = mode;
    recordingDuration = durationMs || 0;
    recordingStart = performance.now();
    discardRecording = false;
    recordFileHandle = fileHandle || null;
    recordFilename = `visum-${Date.now()}.mp4`;
    if (mode === "export" && typeof window.visumSetRunning === "function") {
      const wasRunning = typeof restoreState === "boolean" ? restoreState : isRunning();
      restoreRunning = wasRunning;
      window.visumSetRunning(true);
    } else {
      restoreRunning = null;
    }
    const fps = snapToOptions(readNumber(fpsInput, 30), fpsOptions);
    if (fpsInput) fpsInput.value = fps;
    const bitrateMbps = snapToOptions(readNumber(bitrateInput, 12), bitrateOptions);
    if (bitrateInput) bitrateInput.value = bitrateMbps;
    const { width, height } = parseResolution();
    recordCanvas = document.createElement("canvas");
    recordCanvas.width = width;
    recordCanvas.height = height;
    recordCtx = recordCanvas.getContext("2d");
    drawRecordingFrame();
    const stream = recordCanvas.captureStream(fps);
    const mimeType = pickMimeType();
    if (!mimeType) {
      window.alert("MP4 recording is not supported in this browser.");
      setStatus("Recording is not supported in this browser.");
      pulseHaptic([10, 30, 10]);
      if (recordFrame) {
        window.cancelAnimationFrame(recordFrame);
        recordFrame = 0;
      }
      recordCanvas = null;
      recordCtx = null;
      if (restoreRunning !== null && typeof window.visumSetRunning === "function") {
        window.visumSetRunning(restoreRunning);
      }
      restoreRunning = null;
      recordingMode = "";
      recordingDuration = 0;
      recordingStart = 0;
      updateCaptureState();
      return;
    }
    const options = {
      mimeType,
      videoBitsPerSecond: Math.round(bitrateMbps * 1_000_000),
    };
    try {
      recorder = new MediaRecorder(stream, options);
    } catch (error) {
      window.alert("Recording failed to start in this browser.");
      setStatus("Recording failed to start.");
      pulseHaptic([10, 30, 10]);
      if (recordFrame) {
        window.cancelAnimationFrame(recordFrame);
        recordFrame = 0;
      }
      recordCanvas = null;
      recordCtx = null;
      if (restoreRunning !== null && typeof window.visumSetRunning === "function") {
        window.visumSetRunning(restoreRunning);
      }
      restoreRunning = null;
      recordingMode = "";
      recordingDuration = 0;
      recordingStart = 0;
      updateCaptureState();
      console.error("Failed to start MediaRecorder", error);
      return;
    }
    recordStream = stream;
    chunks = [];
    recorder.addEventListener("dataavailable", (event) => {
      if (event.data && event.data.size > 0) {
        chunks.push(event.data);
      }
    });
    recorder.addEventListener("stop", async () => {
      const savedToFile = Boolean(recordFileHandle);
      if (!discardRecording) {
        const blob = new Blob(chunks, { type: recorder.mimeType || "video/mp4" });
        if (recordFileHandle && recordFileHandle.createWritable) {
          setStatus("Saving...");
          try {
            const writable = await recordFileHandle.createWritable();
            await writable.write(blob);
            await writable.close();
          } catch (error) {
            console.error("Failed to save video", error);
            downloadBlob(blob, recordFilename);
          }
        } else {
          downloadBlob(blob, recordFilename);
        }
      }
      chunks = [];
      recorder = null;
      if (recordStream) {
        recordStream.getTracks().forEach((track) => track.stop());
        recordStream = null;
      }
      if (recordFrame) {
        window.cancelAnimationFrame(recordFrame);
        recordFrame = 0;
      }
      recordCanvas = null;
      recordCtx = null;
      if (restoreRunning !== null && typeof window.visumSetRunning === "function") {
        window.visumSetRunning(restoreRunning);
      }
      restoreRunning = null;
      recordFileHandle = null;
      recordFilename = "";
      updateCaptureState();
      if (discardRecording) {
        setStatus("Export canceled.");
      } else if (savedToFile) {
        setStatus("Video saved.", "success");
      } else {
        setStatus("Video download started.");
      }
      pulseHaptic([15, 30, 15]);
    });
    recorder.start();
    updateCaptureState();
    if (progressBar) {
      progressBar.style.width = "0%";
    }
    setStatus(mode === "export" ? "Exporting... 0%" : "Recording... 0s");
    pulseHaptic(20);
    startProgressLoop();
    if (durationMs && durationMs > 0) {
      window.setTimeout(() => {
        stopRecording();
      }, durationMs);
    }
  };

  const stopRecording = () => {
    if (!recorder) return;
    recorder.stop();
  };

  if (recordButton) {
    recordButton.addEventListener("click", () => startRecording("record", 0));
  }
  if (stopButton) {
    stopButton.addEventListener("click", stopRecording);
  }

  const animationTimings = () => {
    const readSettings = (prefix) => {
      const enable = document.getElementById(`${prefix}-enable`);
      if (!enable || !enable.checked) return null;
      const start = readNumber(document.getElementById(`${prefix}-start`), 0);
      const end = readNumber(document.getElementById(`${prefix}-end`), 0);
      const speed = Math.abs(readNumber(document.getElementById(`${prefix}-speed`), 0));
      const loop = document.getElementById(`${prefix}-loop`);
      const pingpong = document.getElementById(`${prefix}-pingpong`);
      const range = Math.abs(end - start);
      if (speed <= 0 || range <= 0) return null;
      const base = range / speed;
      const isPingPong = pingpong && pingpong.checked;
      const isLoop = loop && loop.checked;
      const cycle = isPingPong ? base * 2 : base;
      return { base, cycle, isLoop };
    };
    return readSettings("mult-anim") || readSettings("line-anim") || readSettings("points-anim");
  };

  if (exportVideo) {
    exportVideo.addEventListener("click", async () => {
      if (recorder) return;
      const fps = snapToOptions(readNumber(fpsInput, 30), fpsOptions);
      if (fpsInput) fpsInput.value = fps;
      const bitrateMbps = snapToOptions(readNumber(bitrateInput, 12), bitrateOptions);
      if (bitrateInput) bitrateInput.value = bitrateMbps;
      const loops = readNumber(loopsInput, 0);
      if (!Number.isFinite(loops) || loops < 0) {
        setStatus("Loop count must be zero or positive.");
        pulseHaptic([10, 30, 10]);
        return;
      }

      if (!exportVideo.dataset.label) {
        exportVideo.dataset.label = exportVideo.textContent || "EXPORT VIDEO (REAL TIME)";
      }
      pulseHaptic([25, 40, 25]);
      if (progressBar) {
        progressBar.style.width = "0%";
      }

      const timing = animationTimings();
      if (!timing || timing.base <= 0) {
        setStatus("Enable an animation to export a timed clip.");
        pulseHaptic([10, 30, 10]);
        return;
      }
      const singleCycle = loops === 0;
      const loopCount = singleCycle ? 1 : loops;
      const durationMs = Math.max(0, loopCount) * (singleCycle ? timing.base : timing.cycle) * 1000;
      if (durationMs <= 0) {
        setStatus("Set a positive loop count to export a video.");
        pulseHaptic([10, 30, 10]);
        return;
      }
      let fileHandle = null;
      if (window.showSaveFilePicker) {
        setStatus("Choose where to save the export...");
        try {
          fileHandle = await window.showSaveFilePicker({
            suggestedName: `visum-${Date.now()}.mp4`,
            types: [
              {
                description: "MP4 Video",
                accept: { "video/mp4": [".mp4"] },
              },
            ],
          });
        } catch (error) {
          if (error && error.name === "AbortError") {
            setStatus("Export canceled.");
            return;
          }
          console.error("Failed to open file picker", error);
          setStatus("Starting export...");
        }
      }
      const wasRunning = isRunning();
      if (typeof window.visumSetRunning === "function") {
        window.visumSetRunning(false);
      }
      if (typeof window.visumResetAnimations === "function") {
        window.visumResetAnimations();
      }
      const kickoff = () => startRecording("export", durationMs, wasRunning, fileHandle);
      window.requestAnimationFrame(() => {
        window.requestAnimationFrame(kickoff);
      });
    });
  }

  if (cancelButton) {
    cancelButton.addEventListener("click", () => {
      if (!recorder) return;
      discardRecording = true;
      stopRecording();
    });
  }

  if (guardContinue && guardOverlay) {
    guardContinue.addEventListener("click", () => {
      guardOverlay.classList.remove("is-active");
      if (guardTarget) {
        guardBypass = true;
        guardTarget.click();
        guardBypass = false;
        guardTarget = null;
      }
    });
  }

  if (guardStop && guardOverlay) {
    guardStop.addEventListener("click", () => {
      guardOverlay.classList.remove("is-active");
      guardTarget = null;
      if (recorder) {
        discardRecording = false;
        stopRecording();
      }
    });
  }

  document.addEventListener(
    "click",
    (event) => {
      if (!recorder || guardBypass) return;
      const target = event.target;
      if (!target) return;
      if (guardOverlay && guardOverlay.contains(target)) return;
      if (target.closest("#export-stop, #export-cancel, #guard-continue, #guard-stop")) return;
      if (!target.closest("button, a, input, select, textarea, summary")) return;
      event.preventDefault();
      event.stopImmediatePropagation();
      guardTarget = target;
      if (guardOverlay) {
        guardOverlay.classList.add("is-active");
      }
    },
    true
  );

  document.addEventListener(
    "input",
    () => {
      if (recorder) return;
      clearStatus();
    },
    true
  );
  document.addEventListener(
    "change",
    () => {
      if (recorder) return;
      clearStatus();
    },
    true
  );
  document.addEventListener(
    "click",
    (event) => {
      if (recorder) return;
      const target = event.target;
      if (!target || target.closest("#export-guard")) return;
      if (!target.closest("button, a, input, select, textarea, summary")) return;
      clearStatus();
    },
    true
  );

  updateCaptureState();
}

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

window.addEventListener("DOMContentLoaded", () => {
  setupControlsToggle();
  setupHeaderMenu();
  setupStickyHeader();
  setupThemeToggle();
  setupExportControls();
  loadWasm().catch((error) => {
    console.error("Failed to load WebAssembly", error);
  });
});
