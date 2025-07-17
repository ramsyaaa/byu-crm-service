// MAU Dashboard JavaScript
window.MAUDashboard = (function () {
  let charts = {
    dailyActiveUsers: null,
    topUsers: null,
  };
  let currentFilters = {
    startDate: "",
    endDate: "",
    userEmail: "",
  };

  // DOM Elements
  let elements = {};

  function init() {
    // No authentication required
    initializeElements();
    initializeDatePickers();
    initializeEventListeners();
    setDefaultDateRange();
    loadUsersList();
    fetchMAUData();
  }

  function initializeElements() {
    elements = {
      // Filters
      startDate: document.getElementById("mauStartDate"),
      endDate: document.getElementById("mauEndDate"),
      userFilter: document.getElementById("mauUserFilter"),
      applyFilters: document.getElementById("applyMAUFilters"),
      clearFilters: document.getElementById("clearMAUFilters"),

      // Quick filter buttons
      thisMonthBtn: document.getElementById("thisMonthBtn"),
      lastMonthBtn: document.getElementById("lastMonthBtn"),
      last3MonthsBtn: document.getElementById("last3MonthsBtn"),

      // Action buttons
      refreshBtn: document.getElementById("mauRefreshBtn"),
      exportBtn: document.getElementById("mauExportBtn"),

      // Metrics
      activeUsersCount: document.getElementById("activeUsersCount"),
      totalRequestsCount: document.getElementById("totalRequestsCount"),
      avgRequestsPerUser: document.getElementById("avgRequestsPerUser"),
      peakActivityDay: document.getElementById("peakActivityDay"),

      // Table
      topUsersTableBody: document.getElementById("topUsersTableBody"),
    };
  }

  function initializeDatePickers() {
    flatpickr(elements.startDate, {
      dateFormat: "Y-m-d",
      maxDate: "today",
      onChange: function () {
        updateFilters();
      },
    });

    flatpickr(elements.endDate, {
      dateFormat: "Y-m-d",
      maxDate: "today",
      onChange: function () {
        updateFilters();
      },
    });
  }

  function initializeEventListeners() {
    // Filter buttons
    elements.applyFilters.addEventListener("click", fetchMAUData);
    elements.clearFilters.addEventListener("click", clearFilters);

    // Quick filter buttons
    elements.thisMonthBtn.addEventListener("click", () =>
      setDateRange("thisMonth")
    );
    elements.lastMonthBtn.addEventListener("click", () =>
      setDateRange("lastMonth")
    );
    elements.last3MonthsBtn.addEventListener("click", () =>
      setDateRange("last3Months")
    );

    // Action buttons
    elements.refreshBtn.addEventListener("click", fetchMAUData);
    elements.exportBtn.addEventListener("click", exportMAUData);

    // User filter change
    elements.userFilter.addEventListener("change", updateFilters);
  }

  function setDefaultDateRange() {
    const now = new Date();
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
    const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0);

    elements.startDate.value = startOfMonth.toISOString().split("T")[0];
    elements.endDate.value = endOfMonth.toISOString().split("T")[0];
    updateFilters();
  }

  function setDateRange(period) {
    const now = new Date();
    let startDate, endDate;

    switch (period) {
      case "thisMonth":
        startDate = new Date(now.getFullYear(), now.getMonth(), 1);
        endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0);
        break;
      case "lastMonth":
        startDate = new Date(now.getFullYear(), now.getMonth() - 1, 1);
        endDate = new Date(now.getFullYear(), now.getMonth(), 0);
        break;
      case "last3Months":
        startDate = new Date(now.getFullYear(), now.getMonth() - 2, 1);
        endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0);
        break;
    }

    elements.startDate.value = startDate.toISOString().split("T")[0];
    elements.endDate.value = endDate.toISOString().split("T")[0];
    updateFilters();
    fetchMAUData();
  }

  function updateFilters() {
    currentFilters.startDate = elements.startDate.value;
    currentFilters.endDate = elements.endDate.value;
    currentFilters.userEmail = elements.userFilter.value;
  }

  function clearFilters() {
    elements.userFilter.value = "";
    setDefaultDateRange();
    fetchMAUData();
  }

  async function loadUsersList() {
    try {
      const response = await fetch("/api-logs/users");

      const data = await response.json();

      if (data.status === "success") {
        populateUsersDropdown(data.data);
      }
    } catch (error) {
      console.error("Error loading users list:", error);
    }
  }

  function populateUsersDropdown(users) {
    // Clear existing options except "All Users"
    elements.userFilter.innerHTML = '<option value="">All Users</option>';

    users.forEach((user) => {
      const option = document.createElement("option");
      option.value = user.user_email;
      option.textContent = user.user_email;
      elements.userFilter.appendChild(option);
    });
  }

  async function fetchMAUData() {
    try {
      showLoading(true);

      const params = new URLSearchParams();
      if (currentFilters.startDate)
        params.append("start_date", currentFilters.startDate);
      if (currentFilters.endDate)
        params.append("end_date", currentFilters.endDate);
      if (currentFilters.userEmail)
        params.append("user_email", currentFilters.userEmail);

      const response = await fetch(`/api-logs/mau?${params}`);

      const data = await response.json();

      if (data.status === "success") {
        updateMetrics(data.data);
        updateCharts(data.data);
        updateTopUsersTable(data.data.top_users);
      } else {
        showError("Failed to fetch MAU data: " + data.message);
      }
    } catch (error) {
      console.error("Error fetching MAU data:", error);
      showError("Failed to fetch MAU data");
    } finally {
      showLoading(false);
    }
  }

  function updateMetrics(data) {
    elements.activeUsersCount.textContent =
      data.active_users_count.toLocaleString();
    elements.totalRequestsCount.textContent =
      data.total_requests.toLocaleString();

    const avgPerUser =
      data.active_users_count > 0
        ? Math.round(data.total_requests / data.active_users_count)
        : 0;
    elements.avgRequestsPerUser.textContent = avgPerUser.toLocaleString();

    // Find peak activity day
    if (data.daily_active_users && data.daily_active_users.length > 0) {
      const peakDay = data.daily_active_users.reduce((max, day) =>
        day.count > max.count ? day : max
      );
      const date = new Date(peakDay.date);
      elements.peakActivityDay.textContent = date.toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
      });
    } else {
      elements.peakActivityDay.textContent = "-";
    }
  }

  function updateCharts(data) {
    updateDailyActiveUsersChart(data.daily_active_users);
    updateTopUsersChart(data.top_users);
  }

  function updateDailyActiveUsersChart(dailyData) {
    const ctx = document
      .getElementById("dailyActiveUsersChart")
      .getContext("2d");

    if (charts.dailyActiveUsers) {
      charts.dailyActiveUsers.destroy();
    }

    const labels = dailyData.map((item) => {
      const date = new Date(item.date);
      return date.toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
      });
    });
    const values = dailyData.map((item) => item.count);

    charts.dailyActiveUsers = new Chart(ctx, {
      type: "line",
      data: {
        labels: labels,
        datasets: [
          {
            label: "Daily Active Users",
            data: values,
            borderColor: "rgb(0, 178, 229)",
            backgroundColor: "rgba(0, 178, 229, 0.1)",
            borderWidth: 3,
            fill: true,
            tension: 0.4,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: false,
          },
        },
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              precision: 0,
            },
          },
        },
      },
    });
  }

  function updateTopUsersChart(topUsers) {
    const ctx = document.getElementById("topUsersChart").getContext("2d");

    if (charts.topUsers) {
      charts.topUsers.destroy();
    }

    const top5Users = topUsers.slice(0, 5);
    const labels = top5Users.map((user) => user.user_email.split("@")[0]); // Show only username part
    const values = top5Users.map((user) => user.request_count);

    charts.topUsers = new Chart(ctx, {
      type: "bar",
      data: {
        labels: labels,
        datasets: [
          {
            label: "Requests",
            data: values,
            backgroundColor: [
              "rgba(0, 178, 229, 0.8)",
              "rgba(16, 192, 243, 0.8)",
              "rgba(34, 197, 94, 0.8)",
              "rgba(245, 158, 11, 0.8)",
              "rgba(239, 68, 68, 0.8)",
            ],
            borderColor: [
              "rgb(0, 178, 229)",
              "rgb(16, 192, 243)",
              "rgb(34, 197, 94)",
              "rgb(245, 158, 11)",
              "rgb(239, 68, 68)",
            ],
            borderWidth: 2,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: false,
          },
        },
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              precision: 0,
            },
          },
        },
      },
    });
  }

  function updateTopUsersTable(topUsers) {
    elements.topUsersTableBody.innerHTML = "";

    topUsers.forEach((user, index) => {
      const row = document.createElement("tr");
      row.classList.add("hover:bg-gray-50");

      const lastActive = new Date(user.last_active).toLocaleDateString(
        "en-US",
        {
          year: "numeric",
          month: "short",
          day: "numeric",
          hour: "2-digit",
          minute: "2-digit",
        }
      );

      row.innerHTML = `
        <td class="px-6 py-4 text-sm font-medium text-gray-900">
          <span class="inline-flex items-center justify-center w-8 h-8 rounded-full ${getRankBadgeColor(
            index
          )} text-white font-bold">
            ${index + 1}
          </span>
        </td>
        <td class="px-6 py-4 text-sm text-gray-900">${user.user_email}</td>
        <td class="px-6 py-4 text-sm text-gray-600">
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
            ${user.request_count.toLocaleString()}
          </span>
        </td>
        <td class="px-6 py-4 text-sm text-gray-600">${lastActive}</td>
      `;

      elements.topUsersTableBody.appendChild(row);
    });
  }

  function getRankBadgeColor(index) {
    const colors = [
      "bg-yellow-500", // Gold
      "bg-gray-400", // Silver
      "bg-yellow-600", // Bronze
      "bg-blue-500", // Blue
      "bg-purple-500", // Purple
    ];
    return colors[index] || "bg-gray-500";
  }

  async function exportMAUData() {
    try {
      showLoading(true);

      const params = new URLSearchParams();
      if (currentFilters.startDate)
        params.append("start_date", currentFilters.startDate);
      if (currentFilters.endDate)
        params.append("end_date", currentFilters.endDate);
      if (currentFilters.userEmail)
        params.append("user_email", currentFilters.userEmail);

      const response = await fetch(`/api-logs/mau?${params}`);

      const data = await response.json();

      if (data.status === "success") {
        exportToExcel(data.data);
      } else {
        showError("Failed to export MAU data");
      }
    } catch (error) {
      console.error("Error exporting MAU data:", error);
      showError("Failed to export MAU data");
    } finally {
      showLoading(false);
    }
  }

  function exportToExcel(data) {
    const wb = XLSX.utils.book_new();

    // Summary sheet
    const summaryData = [
      ["Metric", "Value"],
      ["Active Users", data.active_users_count],
      ["Total Requests", data.total_requests],
      [
        "Average Requests per User",
        Math.round(data.total_requests / data.active_users_count),
      ],
      ["Period Start", data.period.start_date],
      ["Period End", data.period.end_date],
    ];

    const summaryWs = XLSX.utils.aoa_to_sheet(summaryData);
    XLSX.utils.book_append_sheet(wb, summaryWs, "Summary");

    // Top Users sheet
    const topUsersData = [
      ["Rank", "User Email", "Request Count", "Last Active"],
    ];

    data.top_users.forEach((user, index) => {
      topUsersData.push([
        index + 1,
        user.user_email,
        user.request_count,
        user.last_active,
      ]);
    });

    const topUsersWs = XLSX.utils.aoa_to_sheet(topUsersData);
    XLSX.utils.book_append_sheet(wb, topUsersWs, "Top Users");

    // Daily Activity sheet
    const dailyData = [["Date", "Active Users"]];

    data.daily_active_users.forEach((day) => {
      dailyData.push([day.date, day.count]);
    });

    const dailyWs = XLSX.utils.aoa_to_sheet(dailyData);
    XLSX.utils.book_append_sheet(wb, dailyWs, "Daily Activity");

    // Save file
    const fileName = `MAU_Report_${data.period.start_date}_to_${data.period.end_date}.xlsx`;
    XLSX.writeFile(wb, fileName);
  }

  function showLoading(show) {
    if (window.adminDashboard) {
      window.adminDashboard.showLoading(show);
    }
  }

  function showError(message) {
    if (window.adminDashboard) {
      window.adminDashboard.showError(message);
    } else {
      console.error(message);
    }
  }

  // Public API
  return {
    init: init,
  };
})();
