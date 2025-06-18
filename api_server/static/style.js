document.addEventListener("DOMContentLoaded", () => {
  const trendsDiv = document.getElementById("trends");
  const roleSelect = document.getElementById("role");
  const crawlBtn = document.getElementById("crawl-btn");
  const statusDiv = document.getElementById("status");
  const countdownDiv = document.getElementById("countdown");

  let skillsChart, experienceChart, locationsChart, companiesChart;

  const initCharts = () => {
    // Destroy existing charts if they exist
    if (skillsChart) skillsChart.destroy();
    if (experienceChart) experienceChart.destroy();
    if (locationsChart) locationsChart.destroy();
    if (companiesChart) companiesChart.destroy();

    // Initialize empty charts
    const ctxExperience = document.getElementById("experienceChart");

    // Skills Chart (now horizontal bar)
    skillsChart = new ApexCharts(document.querySelector("#skillsChart"), {
      series: [],
      chart: {
        type: "bar",
        height: 350,
        toolbar: { show: false },
      },
      plotOptions: {
        bar: {
          borderRadius: 4,
          horizontal: true,
          dataLabels: { position: "bottom" },
        },
      },
      dataLabels: {
        enabled: true,
        formatter: function (val) {
          return val;
        },
        style: { fontSize: "12px" },
      },
      xaxis: { categories: [] },
      colors: ["#4361ee"],
      tooltip: { enabled: false },
    });
    skillsChart.render();

    // Experience Chart (Bar)
    experienceChart = new Chart(ctxExperience, {
      type: "bar",
      data: {
        labels: [],
        datasets: [{ data: [], backgroundColor: "#4361ee" }],
      },
      options: {
        responsive: true,
        scales: {
          y: { beginAtZero: true },
        },
        plugins: {
          legend: { display: false },
        },
      },
    });

    // Locations Chart (Horizontal Bar)
    locationsChart = new ApexCharts(document.querySelector("#locationsChart"), {
      series: [],
      chart: { type: "bar", height: 350, toolbar: { show: false } },
      plotOptions: { bar: { borderRadius: 4, horizontal: true } },
      dataLabels: { enabled: false },
      xaxis: { categories: [] },
      colors: ["#4cc9f0"],
    });
    locationsChart.render();

    // Companies Chart (Horizontal Bar)
    companiesChart = new ApexCharts(document.querySelector("#companiesChart"), {
      series: [],
      chart: { type: "bar", height: 350, toolbar: { show: false } },
      plotOptions: { bar: { borderRadius: 4, horizontal: true } },
      dataLabels: { enabled: false },
      xaxis: { categories: [] },
      colors: ["#f8961e"],
    });
    companiesChart.render();
  };

  const updateCharts = (data) => {
    // Get top 10 skills only
    const topSkills = data.top_skills.slice(0, 10);

    // Skills Chart (Horizontal Bar)
    skillsChart.updateOptions({
      series: [
        {
          data: topSkills.map((s) => s.count),
        },
      ],
      xaxis: {
        categories: topSkills.map((s) => s.value),
      },
    });

    // Experience Chart (Bar)
    experienceChart.data.labels = data.experience_distribution.map(
      (e) => e.value
    );
    experienceChart.data.datasets[0].data = data.experience_distribution.map(
      (e) => e.count
    );
    experienceChart.update();

    // Locations Chart (Horizontal Bar)
    locationsChart.updateOptions({
      series: [
        {
          data: data.top_locations.map((l) => l.count),
        },
      ],
      xaxis: {
        categories: data.top_locations.map((l) => l.value),
      },
    });

    // Companies Chart (Horizontal Bar)
    companiesChart.updateOptions({
      series: [
        {
          data: data.top_companies.map((c) => c.count),
        },
      ],
      xaxis: {
        categories: data.top_companies.map((c) => c.value),
      },
    });
  };

  const fetchTrends = async (role) => {
    try {
      statusDiv.textContent = "Loading trends...";
      trendsDiv.style.display = "none";

      const res = await fetch(`/api/trends?role=${encodeURIComponent(role)}`);
      const data = await res.json();

      updateCharts(data);

      statusDiv.textContent = "Analysis complete!";
      trendsDiv.style.display = "block";
    } catch (err) {
      console.error(err);
      statusDiv.textContent = "Failed to load trends.";
      statusDiv.style.color = "#ef233c";
    }
  };

  // Initialize empty charts on page load
  initCharts();

  // On page load, show general trends
  fetchTrends("");

  // Crawl button handler
  crawlBtn.addEventListener("click", async () => {
    const role = roleSelect.value;
    if (!role) {
      statusDiv.textContent = "Please select a role first.";
      statusDiv.style.color = "#ef233c";
      return;
    }

    // Start crawling
    statusDiv.textContent = `Analyzing "${role}" market trends...`;
    statusDiv.style.color = "var(--primary)";
    countdownDiv.textContent = "This typically takes about 60 seconds...";
    trendsDiv.style.display = "none";

    await fetch(`/api/crawl?role=${encodeURIComponent(role)}`);

    // Countdown logic
    let seconds = 80;
    const interval = setInterval(() => {
      seconds--;
      countdownDiv.textContent = `⏳ ${seconds} seconds remaining...`;

      if (seconds <= 0) {
        clearInterval(interval);
        statusDiv.textContent = "Processing data...";
        fetchTrends(role);
        countdownDiv.textContent = "";
      }
    }, 1000);
  });

  // Role select change handler
  roleSelect.addEventListener("change", () => {
    const role = roleSelect.value;
    if (role) {
      crawlBtn.querySelector(".button-text").textContent = `Analyze ${role}`;
    } else {
      crawlBtn.querySelector(".button-text").textContent = "Analyze Role";
    }
  });
});
