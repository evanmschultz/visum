const go = new Go();

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
  loadWasm().catch((error) => {
    console.error("Failed to load WebAssembly", error);
  });
});
