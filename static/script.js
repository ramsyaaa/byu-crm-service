document.addEventListener("DOMContentLoaded", function () {
  // DOM Elements - removed legacy file support

  // Filter elements
  const searchFilter = document.getElementById("searchFilter");
  const methodFilter = document.getElementById("methodFilter");
  const statusFilter = document.getElementById("statusFilter");
  const userEmailFilter = document.getElementById("userEmailFilter");
  const startDate = document.getElementById("startDate");
  const endDate = document.getElementById("endDate");
  const minResponseTime = document.getElementById("minResponseTime");
  const maxResponseTime = document.getElementById("maxResponseTime");

  // Quick filter buttons
  const errorsOnlyBtn = document.getElementById("errorsOnlyBtn");
  const slowRequestsBtn = document.getElementById("slowRequestsBtn");
  const todayOnlyBtn = document.getElementById("todayOnlyBtn");
  const clearFilters = document.getElementById("clearFilters");

  // Table and pagination
  const logList = document.getElementById("logList");
  const limitSelect = document.getElementById("limitSelect");
  const resultsCount = document.getElementById("resultsCount");
  const loadingState = document.getElementById("loadingState");
  const emptyState = document.getElementById("emptyState");

  // Pagination elements
  const firstPage = document.getElementById("firstPage");
  const prevPage = document.getElementById("prevPage");
  const nextPage = document.getElementById("nextPage");
  const lastPage = document.getElementById("lastPage");
  const pageNumbers = document.getElementById("pageNumbers");
  const showingStart = document.getElementById("showingStart");
  const showingEnd = document.getElementById("showingEnd");
  const totalItems = document.getElementById("totalItems");

  // Action buttons
  const refreshBtn = document.getElementById("refreshBtn");
  const exportButton = document.getElementById("exportButton");

  // Modal elements
  const logModal = document.getElementById("logModal");
  const closeModal = document.getElementById("closeModal");
  const modalContent = document.getElementById("modalContent");

  // Analytics elements
  const analyticsContainer = document.getElementById("analyticsContainer");
  const chartsContainer = document.getElementById("chartsContainer");
  const totalRequests = document.getElementById("totalRequests");
  const successRate = document.getElementById("successRate");
  const avgResponse = document.getElementById("avgResponse");
  const errorRate = document.getElementById("errorRate");

  // State variables
  let currentPage = 1;
  let currentLimit = 50;
  let totalPages = 0;
  let totalCount = 0;
  let currentSort = { field: "accessed_at", direction: "desc" };
  let logs = [];
  let charts = {
    requestsChart: null,
    statusChart: null,
  };
  let debounceTimer = null;
  let lastUpdated = new Date();

  // Initialize date pickers with default to current day
  function initializeDatePickers() {
    const today = new Date().toISOString().split("T")[0];

    flatpickr(startDate, {
      dateFormat: "Y-m-d",
      maxDate: "today",
      defaultDate: today,
      onChange: function () {
        debouncedFetchLogs();
      },
    });

    flatpickr(endDate, {
      dateFormat: "Y-m-d",
      maxDate: "today",
      defaultDate: today,
      onChange: function () {
        debouncedFetchLogs();
      },
    });

    // Set default values
    startDate.value = today;
    endDate.value = today;
  }

  // Initialize Select2 dropdowns
  function initializeSelect2() {
    $(methodFilter).select2({
      placeholder: "All Methods",
      allowClear: true,
    });

    $(statusFilter).select2({
      placeholder: "All Status Codes",
      allowClear: true,
    });
  }

  // Debounced fetch function to avoid too many API calls
  function debouncedFetchLogs() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(fetchLogs, 300);
  }

  // Main function to fetch logs based on current settings
  function fetchLogs() {
    showLoading(true);
    fetchDatabaseLogs();
  }

  // Fetch database logs with filters
  function fetchDatabaseLogs() {
    const params = new URLSearchParams({
      page: currentPage,
      limit: currentLimit,
      api_only: "true", // Filter only API endpoints
    });

    // Add filters
    if (searchFilter.value) params.append("search", searchFilter.value);
    if (methodFilter.value) params.append("method", methodFilter.value);
    if (statusFilter.value) params.append("status_code", statusFilter.value);
    if (userEmailFilter.value)
      params.append("user_email", userEmailFilter.value);
    if (startDate.value) params.append("start_date", startDate.value);
    if (endDate.value) params.append("end_date", endDate.value);
    if (minResponseTime.value)
      params.append("min_response_time", minResponseTime.value);
    if (maxResponseTime.value)
      params.append("max_response_time", maxResponseTime.value);

    fetch(`/admin/api-logs?${params}`)
      .then((response) => response.json())
      .then((data) => {
        if (data.status === "success") {
          logs = data.data.logs;
          totalCount = data.data.pagination.total_items;
          totalPages = data.data.pagination.total_pages;
          loadDatabaseLogs();
          updatePagination();
          fetchAnalytics();
          updateLastUpdated();
        } else {
          showError("Failed to fetch logs: " + data.message);
        }
      })
      .catch((error) => {
        console.error("Error fetching database logs:", error);
        showError("Failed to fetch logs");
      })
      .finally(() => {
        showLoading(false);
      });
  }

  // Update last updated timestamp
  function updateLastUpdated() {
    lastUpdated = new Date();
    const lastUpdatedElement = document.getElementById("lastUpdated");
    if (lastUpdatedElement) {
      lastUpdatedElement.textContent = lastUpdated.toLocaleTimeString();
    }
  }

  // Load database logs into table
  function loadDatabaseLogs() {
    logList.innerHTML = "";

    if (logs.length === 0) {
      emptyState.classList.remove("hidden");
      return;
    }

    emptyState.classList.add("hidden");

    logs.forEach((log, index) => {
      const row = document.createElement("tr");
      row.classList.add("hover:bg-gray-50", "cursor-pointer");
      row.setAttribute("data-log-index", index);
      row.onclick = (e) => {
        // Don't trigger if clicking on the action button
        if (!e.target.closest(".view-log-btn")) {
          showLogDetails(log);
        }
      };

      const statusClass = getStatusClass(log.status_code);
      const responseTime = log.response_time_ms
        ? `${log.response_time_ms}ms`
        : "-";
      const userEmail = log.auth_user_email || "-";
      const timestamp = new Date(log.accessed_at).toLocaleString();

      row.innerHTML = `
        <td class="px-6 py-4 text-sm text-gray-900">${timestamp}</td>
        <td class="px-6 py-4">
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${statusClass}">
            ${log.status_code}
          </span>
        </td>
        <td class="px-6 py-4">
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
            ${log.method}
          </span>
        </td>
        <td class="px-6 py-4 text-sm text-gray-900 max-w-xs truncate" title="${
          log.endpoint
        }">
          ${log.endpoint}
        </td>
        <td class="px-6 py-4 text-sm text-gray-600">${responseTime}</td>
        <td class="px-6 py-4 text-sm text-gray-600 max-w-xs truncate" title="${userEmail}">
          ${userEmail}
        </td>
        <td class="px-6 py-4 text-sm text-gray-600">${log.ip_address}</td>
        <td class="px-6 py-4 text-sm text-gray-600 max-w-xs">
          ${
            log.error_message
              ? `<span class="text-red-600 font-medium truncate block" title="${log.error_message}">${log.error_message}</span>`
              : '<span class="text-gray-400">-</span>'
          }
        </td>
        <td class="px-6 py-4 text-sm">
          <button class="view-log-btn text-blue-600 hover:text-blue-900 p-2 rounded-lg hover:bg-blue-50 transition-colors" data-log-id="${
            log.id
          }">
            <i class="fas fa-eye"></i>
          </button>
        </td>
      `;
      logList.appendChild(row);
    });

    updateResultsCount();
  }

  // Utility functions
  function getStatusClass(statusCode) {
    const code = parseInt(statusCode);
    if (code >= 200 && code < 300) {
      return "bg-green-100 text-green-800";
    } else if (code >= 400 && code < 500) {
      return "bg-yellow-100 text-yellow-800";
    } else if (code >= 500) {
      return "bg-red-100 text-red-800";
    }
    return "bg-gray-100 text-gray-800";
  }

  function showLoading(show) {
    if (show) {
      loadingState.classList.remove("hidden");
      emptyState.classList.add("hidden");
    } else {
      loadingState.classList.add("hidden");
    }
  }

  function showError(message) {
    console.error(message);
    // You could add a toast notification here
  }

  function formatJsonPayload(payload) {
    try {
      if (typeof payload === "string") {
        const parsed = JSON.parse(payload);
        return JSON.stringify(parsed, null, 2);
      }
      return JSON.stringify(payload, null, 2);
    } catch (e) {
      // If JSON parsing fails, return the raw payload
      return payload || "Invalid JSON";
    }
  }

  function updateResultsCount() {
    const start = (currentPage - 1) * currentLimit + 1;
    const end = Math.min(currentPage * currentLimit, totalCount);

    showingStart.textContent = totalCount > 0 ? start : 0;
    showingEnd.textContent = end;
    totalItems.textContent = totalCount;
    resultsCount.textContent = `Showing ${totalCount} results`;
  }

  // Enhanced pagination
  function updatePagination() {
    updateResultsCount();

    // Update button states
    firstPage.disabled = currentPage === 1;
    prevPage.disabled = currentPage === 1;
    nextPage.disabled = currentPage === totalPages;
    lastPage.disabled = currentPage === totalPages;

    // Generate page numbers
    generatePageNumbers();
  }

  function generatePageNumbers() {
    pageNumbers.innerHTML = "";

    const maxVisible = 5;
    let startPage = Math.max(1, currentPage - Math.floor(maxVisible / 2));
    let endPage = Math.min(totalPages, startPage + maxVisible - 1);

    if (endPage - startPage + 1 < maxVisible) {
      startPage = Math.max(1, endPage - maxVisible + 1);
    }

    for (let i = startPage; i <= endPage; i++) {
      const button = document.createElement("button");
      button.textContent = i;
      button.className = `px-3 py-2 text-sm rounded-lg ${
        i === currentPage
          ? "bg-blue-500 text-white"
          : "bg-gray-100 text-gray-600 hover:bg-gray-200"
      }`;
      button.onclick = () => goToPage(i);
      pageNumbers.appendChild(button);
    }
  }

  function goToPage(page) {
    currentPage = page;
    fetchLogs();
  }

  // Analytics functions
  function fetchAnalytics() {
    // Build params for current day analytics
    const params = new URLSearchParams({
      api_only: "true",
    });

    // Add current date filters for analytics
    if (startDate.value) params.append("start_date", startDate.value);
    if (endDate.value) params.append("end_date", endDate.value);

    fetch(`/admin/api-logs/stats?${params}`)
      .then((response) => response.json())
      .then((data) => {
        if (data.status === "success") {
          updateAnalytics(data.data);
        }
      })
      .catch((error) => console.error("Error fetching analytics:", error));

    // Fetch chart data
    fetchChartData();
  }

  // Fetch chart data
  function fetchChartData() {
    const params = new URLSearchParams({
      api_only: "true",
    });

    // Add current date filters
    if (startDate.value) params.append("start_date", startDate.value);
    if (endDate.value) params.append("end_date", endDate.value);

    // Fetch requests over time data
    fetch(`/admin/api-logs/chart-data/requests-over-time?${params}`)
      .then((response) => response.json())
      .then((data) => {
        if (data.status === "success") {
          updateRequestsChart(data.data);
        }
      })
      .catch((error) =>
        console.error("Error fetching requests chart data:", error)
      );

    // Fetch status code distribution data
    fetch(`/admin/api-logs/chart-data/status-distribution?${params}`)
      .then((response) => response.json())
      .then((data) => {
        if (data.status === "success") {
          updateStatusChart(data.data);
        }
      })
      .catch((error) =>
        console.error("Error fetching status chart data:", error)
      );
  }

  function updateAnalytics(stats) {
    totalRequests.textContent = stats.total_logs?.toLocaleString() || "0";

    const total = stats.total_logs || 0;
    const success = stats.success_count || 0;
    const errors = stats.server_error_count || 0;

    const successRateValue =
      total > 0 ? ((success / total) * 100).toFixed(1) : "0";
    const errorRateValue =
      total > 0 ? ((errors / total) * 100).toFixed(1) : "0";

    successRate.textContent = `${successRateValue}%`;
    errorRate.textContent = `${errorRateValue}%`;
    avgResponse.textContent = stats.avg_response_time_ms
      ? `${Math.round(stats.avg_response_time_ms)}ms`
      : "-";
  }

  // Initialize charts
  function initializeCharts() {
    // Requests over time chart
    const requestsCtx = document
      .getElementById("requestsChart")
      .getContext("2d");
    charts.requestsChart = new Chart(requestsCtx, {
      type: "line",
      data: {
        labels: [],
        datasets: [
          {
            label: "API Requests",
            data: [],
            borderColor: "rgb(59, 130, 246)",
            backgroundColor: "rgba(59, 130, 246, 0.1)",
            borderWidth: 3,
            fill: true,
            tension: 0.4,
            pointBackgroundColor: "rgb(59, 130, 246)",
            pointBorderColor: "rgb(59, 130, 246)",
            pointRadius: 4,
            pointHoverRadius: 6,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
          intersect: false,
          mode: "index",
        },
        plugins: {
          legend: {
            display: false,
          },
          tooltip: {
            backgroundColor: "rgba(0, 0, 0, 0.8)",
            titleColor: "white",
            bodyColor: "white",
            borderColor: "rgb(59, 130, 246)",
            borderWidth: 1,
            cornerRadius: 8,
            displayColors: false,
            callbacks: {
              title: function (context) {
                return `Time: ${context[0].label}`;
              },
              label: function (context) {
                return `Requests: ${context.parsed.y}`;
              },
            },
          },
        },
        scales: {
          y: {
            beginAtZero: true,
            grid: {
              color: "rgba(0, 0, 0, 0.1)",
              drawBorder: false,
            },
            ticks: {
              color: "rgba(0, 0, 0, 0.6)",
              font: {
                size: 12,
              },
            },
          },
          x: {
            grid: {
              color: "rgba(0, 0, 0, 0.1)",
              drawBorder: false,
            },
            ticks: {
              color: "rgba(0, 0, 0, 0.6)",
              font: {
                size: 12,
              },
            },
          },
        },
        elements: {
          line: {
            borderJoinStyle: "round",
          },
        },
      },
    });

    // Status code distribution chart
    const statusCtx = document.getElementById("statusChart").getContext("2d");
    charts.statusChart = new Chart(statusCtx, {
      type: "doughnut",
      data: {
        labels: ["Success (2xx)", "Client Error (4xx)", "Server Error (5xx)"],
        datasets: [
          {
            data: [0, 0, 0],
            backgroundColor: [
              "rgba(34, 197, 94, 0.8)",
              "rgba(251, 191, 36, 0.8)",
              "rgba(239, 68, 68, 0.8)",
            ],
            borderColor: [
              "rgb(34, 197, 94)",
              "rgb(251, 191, 36)",
              "rgb(239, 68, 68)",
            ],
            borderWidth: 2,
            hoverBackgroundColor: [
              "rgba(34, 197, 94, 1)",
              "rgba(251, 191, 36, 1)",
              "rgba(239, 68, 68, 1)",
            ],
            hoverBorderWidth: 3,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        cutout: "60%",
        plugins: {
          legend: {
            position: "bottom",
            labels: {
              padding: 15,
              usePointStyle: true,
              font: {
                size: 12,
              },
              color: "rgba(0, 0, 0, 0.7)",
            },
          },
          tooltip: {
            backgroundColor: "rgba(0, 0, 0, 0.8)",
            titleColor: "white",
            bodyColor: "white",
            borderColor: "rgba(255, 255, 255, 0.2)",
            borderWidth: 1,
            cornerRadius: 8,
            callbacks: {
              label: function (context) {
                const total = context.dataset.data.reduce((a, b) => a + b, 0);
                const percentage =
                  total > 0 ? ((context.parsed / total) * 100).toFixed(1) : 0;
                return `${context.label}: ${context.parsed} (${percentage}%)`;
              },
            },
          },
        },
        animation: {
          animateRotate: true,
          animateScale: true,
          duration: 1000,
        },
      },
    });
  }

  // Update requests over time chart
  function updateRequestsChart(data) {
    if (!charts.requestsChart) return;

    charts.requestsChart.data.labels = data.labels || [];
    charts.requestsChart.data.datasets[0].data = data.values || [];
    charts.requestsChart.update();
  }

  // Update status code distribution chart
  function updateStatusChart(data) {
    if (!charts.statusChart) return;

    charts.statusChart.data.datasets[0].data = [
      data.success_count || 0,
      data.client_error_count || 0,
      data.server_error_count || 0,
    ];
    charts.statusChart.update();
  }

  // Modal functions
  function showLogDetails(log) {
    const content = `
      <div class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700">Timestamp</label>
            <p class="mt-1 text-sm text-gray-900">${new Date(
              log.accessed_at || log.timestamp
            ).toLocaleString()}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Status Code</label>
            <p class="mt-1"><span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusClass(
              log.status_code
            )}">${log.status_code}</span></p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Method</label>
            <p class="mt-1 text-sm text-gray-900">${log.method}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Response Time</label>
            <p class="mt-1 text-sm text-gray-900">${
              log.response_time_ms
                ? `${log.response_time_ms}ms`
                : log.latency || "-"
            }</p>
          </div>
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700">Endpoint</label>
            <p class="mt-1 text-sm text-gray-900 break-all">${
              log.endpoint || log.path
            }</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">IP Address</label>
            <p class="mt-1 text-sm text-gray-900">${
              log.ip_address || log.ip
            }</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">User Email</label>
            <p class="mt-1 text-sm text-gray-900">${
              log.auth_user_email || "-"
            }</p>
          </div>
        </div>

        ${
          log.user_agent
            ? `
          <div>
            <label class="block text-sm font-medium text-gray-700">User Agent</label>
            <p class="mt-1 text-sm text-gray-900 break-all">${log.user_agent}</p>
          </div>
        `
            : ""
        }

        ${
          log.request_payload
            ? `
          <div>
            <label class="block text-sm font-medium text-gray-700">Request Payload</label>
            <pre class="mt-1 text-sm text-gray-900 bg-gray-50 p-3 rounded-lg overflow-auto max-h-40">${formatJsonPayload(
              log.request_payload
            )}</pre>
          </div>
        `
            : ""
        }

        ${
          log.response_payload
            ? `
          <div>
            <label class="block text-sm font-medium text-gray-700">Response Payload</label>
            <pre class="mt-1 text-sm text-gray-900 bg-gray-50 p-3 rounded-lg overflow-auto max-h-40">${formatJsonPayload(
              log.response_payload
            )}</pre>
          </div>
        `
            : ""
        }

        ${
          log.error_message
            ? `
          <div>
            <label class="block text-sm font-medium text-gray-700">Error Message</label>
            <p class="mt-1 text-sm text-red-600 bg-red-50 p-3 rounded-lg">${log.error_message}</p>
          </div>
        `
            : ""
        }
      </div>
    `;

    modalContent.innerHTML = content;
    logModal.classList.remove("hidden");
  }

  // Event listeners for modal buttons
  document.addEventListener("click", function (e) {
    if (e.target.closest(".view-log-btn")) {
      e.preventDefault();
      e.stopPropagation();
      const button = e.target.closest(".view-log-btn");
      const logId = button.getAttribute("data-log-id");

      // Find the log in our current logs array
      const log = logs.find((l) => l.id == logId);
      if (log) {
        showLogDetails(log);
      }
    }
  });

  // Apply filters button - main way to apply filters
  const applyFilters = document.getElementById("applyFilters");
  const filterIndicator = document.getElementById("filterIndicator");

  applyFilters.addEventListener("click", function () {
    currentPage = 1;
    fetchLogs();
    // Hide the filter indicator after applying
    filterIndicator.classList.add("hidden");
  });

  // Function to show filter indicator when filters change
  function showFilterIndicator() {
    filterIndicator.classList.remove("hidden");
  }

  // Add event listeners to show indicator when filters change
  methodFilter.addEventListener("change", showFilterIndicator);
  statusFilter.addEventListener("change", showFilterIndicator);
  userEmailFilter.addEventListener("input", showFilterIndicator);
  minResponseTime.addEventListener("input", showFilterIndicator);
  maxResponseTime.addEventListener("input", showFilterIndicator);

  // Only search filter has auto-refresh for better UX
  searchFilter.addEventListener("input", debouncedFetchLogs);

  // Date picker changes trigger auto-refresh since they're commonly used
  // Other filters require clicking "Apply Filters" button

  // Quick filter buttons - set values but don't auto-apply
  errorsOnlyBtn.addEventListener("click", function () {
    statusFilter.value = "500";
    $(statusFilter).trigger("change");
    showFilterIndicator();
  });

  slowRequestsBtn.addEventListener("click", function () {
    minResponseTime.value = "1000";
    showFilterIndicator();
  });

  todayOnlyBtn.addEventListener("click", function () {
    const today = new Date().toISOString().split("T")[0];
    startDate.value = today;
    endDate.value = today;
    // Auto-apply for date filters since they're commonly used
    currentPage = 1;
    fetchLogs();
  });

  clearFilters.addEventListener("click", function () {
    searchFilter.value = "";
    methodFilter.value = "";
    statusFilter.value = "";
    userEmailFilter.value = "";
    startDate.value = "";
    endDate.value = "";
    minResponseTime.value = "";
    maxResponseTime.value = "";
    $(methodFilter).trigger("change");
    $(statusFilter).trigger("change");
    currentPage = 1;
    fetchLogs();
    // Hide the filter indicator after clearing
    filterIndicator.classList.add("hidden");
  });

  // Pagination event listeners
  firstPage.addEventListener("click", () => goToPage(1));
  prevPage.addEventListener("click", () =>
    goToPage(Math.max(1, currentPage - 1))
  );
  nextPage.addEventListener("click", () =>
    goToPage(Math.min(totalPages, currentPage + 1))
  );
  lastPage.addEventListener("click", () => goToPage(totalPages));

  // Limit change
  limitSelect.addEventListener("change", function () {
    currentLimit = parseInt(this.value);
    currentPage = 1;
    fetchLogs();
  });

  // Action buttons
  refreshBtn.addEventListener("click", fetchLogs);

  exportButton.addEventListener("click", function () {
    const wb = XLSX.utils.book_new();
    const ws = XLSX.utils.table_to_sheet(document.getElementById("logTable"));
    XLSX.utils.book_append_sheet(wb, ws, "Logs");

    const filename = `api-logs-${new Date().toISOString().split("T")[0]}.xlsx`;
    XLSX.writeFile(wb, filename);
  });

  // Modal event listeners
  closeModal.addEventListener("click", function () {
    logModal.classList.add("hidden");
  });

  logModal.addEventListener("click", function (e) {
    if (e.target === logModal) {
      logModal.classList.add("hidden");
    }
  });

  // Initialize everything
  initializeDatePickers();
  initializeSelect2();
  initializeCharts();

  // Handle window resize for charts
  window.addEventListener("resize", function () {
    if (charts.requestsChart) {
      charts.requestsChart.resize();
    }
    if (charts.statusChart) {
      charts.statusChart.resize();
    }
  });

  // Start with database logs
  fetchLogs();
});
