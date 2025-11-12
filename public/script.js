let charts = {};

function makeChart(id, type, labels, data, title, color = "rgba(75,192,192,0.6)") {
  const canvas = document.getElementById(id);
  if (!canvas) return; 
  const ctx = canvas.getContext("2d");
  if (charts[id]) charts[id].destroy();

  charts[id] = new Chart(ctx, {
    type,
    data: {
      labels,
      datasets: [{
        label: title,
        data,
        backgroundColor: color,
        borderColor: color.replace("0.6", "1"),
        borderWidth: 1
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: true, position: "bottom" },
        title: { display: true, text: title }
      },
      scales: { y: { beginAtZero: true } }
    }
  });
}

function showSection(id) {
  document.querySelectorAll("section").forEach(sec => sec.classList.remove("active"));
  document.getElementById(id).classList.add("active");
  window.scrollTo({ top: 0, behavior: "smooth" }); // scroll ke atas tiap pindah halaman
}

document.getElementById("darkToggle").addEventListener("click", () => {
  document.body.classList.toggle("dark-mode");
});

async function loadDashboard() {
  const resRec = await fetch("/recommend?month=Januari&region=Jember Utara");
  const dataRec = await resRec.json();
  const resProd = await fetch("/production");
  const dataProd = await resProd.json();
  const resAnal = await fetch("/analysis");
  const dataAnal = await resAnal.json();

  const topPlant = dataRec.rekomendasi[0] || { tanaman: "-", skor: 0 };
  const topProd = Object.entries(dataProd).sort((a, b) => b[1] - a[1])[0] || ["-", 0];

  document.getElementById("summary").innerHTML = `
    <h3>📋 Ringkasan Otomatis</h3>
    <p><b>Musim Saat Ini:</b> ${dataRec.musim}</p>
    <p><b>Tanaman Direkomendasikan:</b> ${topPlant.tanaman} (Skor ${topPlant.skor})</p>
    <p><b>Produksi Tertinggi:</b> ${topProd[0]} (${topProd[1]} kuintal/ha)</p>
  `;

  makeChart("regionChart", "pie",
    ["Utara", "Tengah", "Selatan"],
    [dataAnal.wilayah_utara.length, dataAnal.wilayah_tengah.length, dataAnal.wilayah_selatan.length],
    "Distribusi Tanaman per Wilayah",
    "rgba(255,205,86,0.6)"
  );
}

["month", "region", "season"].forEach(id =>
  document.getElementById(id).addEventListener("change", getRekomendasi)
);

async function getRekomendasi() {
  const month = document.getElementById("month").value;
  const region = document.getElementById("region").value;
  const season = document.getElementById("season").value;
  const res = await fetch(`/recommend?month=${month}&region=${region}&season=${season}`);
  const data = await res.json();

  const div = document.getElementById("result");
  div.innerHTML = "";

  if (!data.rekomendasi?.length) {
    div.innerHTML = "<p style='color:red;'>Tidak ada rekomendasi tanaman.</p>";
    return;
  }

  const labels = [], values = [];
  data.rekomendasi.forEach(r => {
    div.innerHTML += `
      <div class="card">
        <h4>${r.tanaman}</h4>
        <p>${r.deskripsi}</p>
        <p><b>Skor:</b> ${r.skor}</p>
      </div>`;
    labels.push(r.tanaman);
    values.push(r.skor);
  });

  makeChart("recommendChart", "bar", labels, values, "Skor Kecocokan Tanaman", "rgba(153,102,255,0.6)");
}

async function loadProduction() {
  const res = await fetch("/production");
  const data = await res.json();
  const div = document.getElementById("productionResult");
  div.innerHTML = "";
  const labels = Object.keys(data);
  const values = Object.values(data);
  labels.forEach((n, i) => {
    div.innerHTML += `<div class="card"><b>${n}</b>: ${values[i]} kuintal/ha</div>`;
  });
  makeChart("productionChart", "bar", labels, values, "Produksi Tanaman (kuintal/ha)", "rgba(54,162,235,0.6)");
}

async function loadWeather() {
  const res = await fetch("/weather?month=Januari");
  const data = await res.json();
  document.getElementById("weatherResult").innerHTML = `
    <div class="card"><p><b>Musim:</b> ${data.musim}</p><p>${data.info}</p></div>`;
}

document.getElementById("fertilizerSearch").addEventListener("input", async (e) => {
  const query = e.target.value.trim();
  const div = document.getElementById("fertilizerResult");
  div.innerHTML = "";
  if (!query) return;

  const res = await fetch(`/fertilizer?plant=${query}`);
  const data = await res.json();

  div.innerHTML = data.error
    ? `<p style="color:red;">${data.error}</p>`
    : `<div class="card"><h4>${data.tanaman}</h4><p>${data.pupuk}</p><p><i>${data.keterangan}</i></p></div>`;
});

document.getElementById("careSearch").addEventListener("input", async (e) => {
  const query = e.target.value.trim();
  const div = document.getElementById("careResult");
  div.innerHTML = "";
  if (!query) return;

  const res = await fetch(`/care?plant=${query}`);
  const data = await res.json();

  if (data.error) {
    div.innerHTML = `<p style="color:red;">${data.error}</p>`;
  } else {
    const preview = data.panduan.length > 100 ? data.panduan.substring(0, 100) + "..." : data.panduan;
    div.innerHTML = `
      <div class="card">
        <h4>Panduan untuk ${data.tanaman}</h4>
        <p id="previewText">${preview}</p>
        <p id="fullText" style="display:none;">${data.panduan}</p>
        <button id="toggleDetail" class="detail-btn">Lihat Detail</button>
      </div>`;
    document.getElementById("toggleDetail").addEventListener("click", () => {
      const previewEl = document.getElementById("previewText");
      const fullEl = document.getElementById("fullText");
      const btn = document.getElementById("toggleDetail");
      const showFull = fullEl.style.display === "none";
      fullEl.style.display = showFull ? "block" : "none";
      previewEl.style.display = showFull ? "none" : "block";
      btn.textContent = showFull ? "Sembunyikan" : "Lihat Detail";
    });
  }
});

async function loadAnalysis() {
  const res = await fetch("/analysis");
  const data = await res.json();
  document.getElementById("analysisResult").innerHTML = `
    <div class="card">
      <p><b>Wilayah Utara:</b> ${data.wilayah_utara.join(", ")}</p>
      <p><b>Wilayah Tengah:</b> ${data.wilayah_tengah.join(", ")}</p>
      <p><b>Wilayah Selatan:</b> ${data.wilayah_selatan.join(", ")}</p>
    </div>`;
}

window.onload = () => {
  loadDashboard();
  getRekomendasi();
  loadProduction();
  loadWeather();
  loadAnalysis();
};
