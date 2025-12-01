let chartInstances = {};

function showSection(id) {
  document.querySelectorAll("section").forEach(sec => sec.classList.remove("active"));
  const el = document.getElementById(id);
  if (!el) return;
  el.classList.add("active");

  if (id === "dashboard") loadDashboard();
  if (id === "recommend") loadRecommend();
  if (id === "production") loadProduction();
  if (id === "analysis") loadAnalysis();
}

function makeChart(id, type, labels, values, title, colors) {
  const canvas = document.getElementById(id);
  if (!canvas) {
    console.warn("makeChart: canvas not found:", id);
    return;
  }

  if (chartInstances[id]) {
    chartInstances[id].destroy();
  }

  const bg = Array.isArray(colors) ? colors : (colors ? [colors] : [
    "rgba(54, 162, 235, 0.6)"
  ]);

  chartInstances[id] = new Chart(canvas.getContext('2d'), {
    type,
    data: {
      labels,
      datasets: [{
        label: title,
        data: values,
        backgroundColor: bg,
        borderColor: bg.map(c => {
          return typeof c === "string" ? c.replace(/0\.6\)$/, "1)") : c;
        }),
        borderWidth: 1
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: true, position: "bottom" },
        title: { display: !!title, text: title }
      },
      scales: { y: { beginAtZero: true } }
    }
  });
}

async function fetchJson(url) {
  try {
    const res = await fetch(url);
    if (!res.ok) throw new Error(`HTTP ${res.status} ${res.statusText}`);
    return await res.json();
  } catch (err) {
    console.error("fetchJson error:", url, err);
    throw err;
  }
}

async function loadDashboard() {
  try {
    const defMonth = document.getElementById("month")?.value || "Januari";
    const defRegion = document.getElementById("region")?.value || "Jember Utara";

    console.log("Loading dashboard with:", defMonth, defRegion); // Debug

    const [rec, prod, anal] = await Promise.all([
      fetchJson(`/recommend?month=${encodeURIComponent(defMonth)}&region=${encodeURIComponent(defRegion)}`),
      fetchJson("/production"),
      fetchJson("/analysis")
    ]);

    console.log("loadDashboard: recommend", rec);
    console.log("loadDashboard: production", prod);
    console.log("loadDashboard: analysis", anal);

    const topPlant = (rec.rekomendasi && rec.rekomendasi.length) ? rec.rekomendasi[0] : { tanaman: "-", skor: 0 };

    let topName = "-";
    let topValue = 0;
    for (const wilayah in prod) {
      for (const tanaman in prod[wilayah]) {
        const v = prod[wilayah][tanaman];
        if (typeof v === "number" && v > topValue) {
          topValue = v;
          topName = tanaman;
        }
      }
    }

    const summaryEl = document.getElementById("summary");
    if (summaryEl) {
      summaryEl.innerHTML = `
        <h3>Ringkasan Otomatis</h3>
        <p><b>Musim Saat Ini:</b> ${rec.musim || "-"}</p>
        <p><b>Tanaman Direkomendasikan:</b> ${topPlant.tanaman} (Skor ${topPlant.skor})</p>
        <p><b>Produksi Tertinggi:</b> ${topName} (${topValue} kuintal/ha)</p>
      `;
    }

    if (anal && typeof anal === "object") {
      const aUt = Array.isArray(anal.wilayah_utara) ? anal.wilayah_utara.length : 0;
      const aTe = Array.isArray(anal.wilayah_tengah) ? anal.wilayah_tengah.length : 0;
      const aSe = Array.isArray(anal.wilayah_selatan) ? anal.wilayah_selatan.length : 0;

      // Hancurkan chart yang ada sebelum membuat yang baru
      if (chartInstances["regionChart"]) {
        chartInstances["regionChart"].destroy();
      }

      makeChart("regionChart", "pie", ["Utara", "Tengah", "Selatan"], [aUt, aTe, aSe], "Distribusi Tanaman per Wilayah", [
        "rgba(160, 112, 0, 0.6)",
        "rgba(0, 151, 151, 0.6)",
        "rgba(52, 0, 156, 0.6)"
      ]);
    }

  } catch (err) {
    console.error("loadDashboard failed:", err);
    const summaryEl = document.getElementById("summary");
    if (summaryEl) summaryEl.innerHTML = `<p style="color:red;">Gagal memuat ringkasan: ${err.message}</p>`;
  }
}

async function loadRecommend() {
  try {
    const month = document.getElementById("month")?.value || "";
    const region = document.getElementById("region")?.value || "";
    const season = document.getElementById("season")?.value || "";

    console.log("Mengirim request dengan:", { month, region, season }); // Debug

    const q = `/recommend?month=${encodeURIComponent(month)}&region=${encodeURIComponent(region)}&season=${encodeURIComponent(season)}`;
    const data = await fetchJson(q);
    console.log("Response dari server:", data);

    const div = document.getElementById("result");
    if (!div) return;
    div.innerHTML = "";

    if (!data.rekomendasi || !data.rekomendasi.length) {
      div.innerHTML = "<p style='color:red;'>Tidak ada rekomendasi tanaman.</p>";
      if (chartInstances["recommendChart"]) chartInstances["recommendChart"].destroy();
      return;
    }

    const labels = [];
    const values = [];

    data.rekomendasi.forEach(r => {
      div.innerHTML += `
        <div class="card">
          <h4>${r.tanaman}</h4>
          <p>${r.deskripsi}</p>
          <p><b>Skor:</b> ${r.skor}</p>
        </div>
      `;
      labels.push(r.tanaman);
      values.push(r.skor);
    });

    makeChart("recommendChart", "bar", labels, values, "Skor Kecocokan Tanaman", [
      "rgba(41, 0, 122, 0.6)",
      "rgba(0, 123, 123, 0.6)",
      "rgba(123, 61, 0, 0.6)",
      "rgba(123, 0, 29, 0.6)",
    ]);

  } catch (err) {
    console.error("loadRecommend failed:", err);
    const div = document.getElementById("result");
    if (div) div.innerHTML = `<p style="color:red;">Gagal memuat rekomendasi</p>`;
  }
}

async function loadProduction() {
  try {
    const data = await fetchJson("/production");
    console.log("DEBUG RAW DATA:", data);

    const div = document.getElementById("productionResult");
    if (!div) return;
    div.innerHTML = "";

    const labels = [];
    const values = [];

    Object.keys(data).forEach((wilayah) => {
      const tanamanObj = data[wilayah] || {};

      let maxValue = 0;
      let maxPlant = "-";

      Object.keys(tanamanObj).forEach((tanaman) => {
        const v = tanamanObj[tanaman];
        if (typeof v === "number" && v > maxValue) {
          maxValue = v;
          maxPlant = tanaman;
        }
      });

      div.innerHTML += `
        <div class="card">
          <b>${wilayah}</b>: ${maxPlant} (${maxValue} kuintal/ha)
        </div>
      `;

      labels.push(wilayah);
      values.push(maxValue);
    });

    makeChart("productionChart", "bar", labels, values, "Produksi Tertinggi per Wilayah", [
      "rgba(136, 0, 29, 0.6)",
      "rgba(0, 147, 245, 0.6)",
      "rgba(136, 97, 0, 0.6)",
      "rgba(0, 117, 117, 0.6)",
      "rgba(42, 0, 128, 0.6)"
    ]);

  } catch (err) {
    console.error("FATAL ERROR loadProduction:", err);
    const div = document.getElementById("productionResult");
    if (div) div.innerHTML = `<p style="color:red;">Gagal memuat data produksi</p>`;
  }
}

async function loadCareFor(query) {
  if (!query) return;
  try {
    const res = await fetchJson(`/care?plant=${encodeURIComponent(query)}`);
    const div = document.getElementById("careResult");
    if (!div) return;

    if (res.error) {
      div.innerHTML = `<p style="color:red;">${res.error}</p>`;
    } else {
      const preview = res.panduan.length > 100 ? res.panduan.substring(0, 100) + "..." : res.panduan;
      div.innerHTML = `
        <div class="card">
          <h4>Panduan untuk ${res.tanaman}</h4>
          <p id="previewText">${preview}</p>
          <p id="fullText" style="display:none;">${res.panduan}</p>
          <button id="toggleDetail" class="detail-btn">Lihat Detail</button>
        </div>`;
      const btn = document.getElementById("toggleDetail");
      if (btn) {
        btn.addEventListener("click", () => {
          const previewEl = document.getElementById("previewText");
          const fullEl = document.getElementById("fullText");
          const showFull = fullEl.style.display === "none";
          fullEl.style.display = showFull ? "block" : "none";
          previewEl.style.display = showFull ? "none" : "block";
          btn.textContent = showFull ? "Sembunyikan" : "Lihat Detail";
        });
      }
    }
  } catch (err) {
    console.error("loadCareFor failed:", err);
  }
}

async function loadFertilizerFor(query) {
  if (!query) return;
  try {
    const res = await fetchJson(`/fertilizer?plant=${encodeURIComponent(query)}`);
    const div = document.getElementById("fertilizerResult");
    if (!div) return;
    if (res.error) {
      div.innerHTML = `<p style="color:red;">${res.error}</p>`;
    } else {
      div.innerHTML = `<div class="card"><h4>${res.tanaman}</h4><p>${res.pupuk}</p><p><i>${res.keterangan}</i></p></div>`;
    }
  } catch (err) {
    console.error("loadFertilizerFor failed:", err);
  }
}

async function loadAnalysis() {
  try {
    const res = await fetchJson("/analysis");
    console.log("loadAnalysis:", res);

    const el = document.getElementById("analysisResult");
    if (!el) return;

    const ut = Array.isArray(res.wilayah_utara) ? res.wilayah_utara.join(", ") : "-";
    const te = Array.isArray(res.wilayah_tengah) ? res.wilayah_tengah.join(", ") : "-";
    const se = Array.isArray(res.wilayah_selatan) ? res.wilayah_selatan.join(", ") : "-";

    el.innerHTML = `
      <div class="card">
        <p><b>Wilayah Utara:</b> ${ut}</p>
        <p><b>Wilayah Tengah:</b> ${te}</p>
        <p><b>Wilayah Selatan:</b> ${se}</p>
        <hr>
        <p><b>Produksi Tertinggi:</b> ${res.produksi_tertinggi?.tanaman || "-"} (${res.produksi_tertinggi?.nilai || "-"})</p>
        <p><b>Musim (terurut):</b> ${Array.isArray(res.musim_berurut) ? res.musim_berurut.join(", ") : "-"}</p>
        <p><b>Total Tanaman:</b> ${res.total_tanaman || "-"}</p>
      </div>
    `;
  } catch (err) {
    console.error("loadAnalysis failed:", err);
    const el = document.getElementById("analysisResult");
    if (el) el.innerHTML = `<p style="color:red;">Gagal memuat analisis</p>`;
  }
}

function setupEventListeners() {
  ["month", "region", "season"].forEach(id => {
    const el = document.getElementById(id);
    if (el) {
      el.addEventListener("change", () => {
        console.log("Changed:", id, "value:", el.value); // Debug
        loadRecommend(); // Update rekomendasi
        loadDashboard(); // Update dashboard
      });
    }
  });

  const fert = document.getElementById("fertilizerSearch");
  if (fert) {
    let t;
    fert.addEventListener("input", (e) => {
      clearTimeout(t);
      const q = e.target.value.trim();
      t = setTimeout(() => loadFertilizerFor(q), 300);
    });
  }

  const care = document.getElementById("careSearch");
  if (care) {
    let t;
    care.addEventListener("input", (e) => {
      clearTimeout(t);
      const q = e.target.value.trim();
      t = setTimeout(() => loadCareFor(q), 300);
    });
  }

  const darkToggle = document.getElementById("darkToggle");
  if (darkToggle) {
    darkToggle.addEventListener("click", () => {
      document.body.classList.toggle("dark-mode");
    });
  }
}

window.addEventListener("DOMContentLoaded", () => {
  console.log("script.js: DOM ready");
  setupEventListeners();

  loadDashboard();
});
