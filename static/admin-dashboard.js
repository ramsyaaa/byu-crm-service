// Admin Dashboard JavaScript
class AdminDashboard {
  constructor() {
    this.charts = {};
    this.currentFilters = {
      startDate: "",
      endDate: "",
      userEmail: "",
      apiOnly: true,
      businessHours: false,
      uniqueOnly: true,
    };

    this.init();
  }

  init() {
    this.setupFilters();
    this.setupEventListeners();
    this.setDefaultDates();
    this.loadDashboardData();
  }

  setupFilters() {
    // Set default dates to today only (for better performance)
    const now = new Date();

    document.getElementById("startDate").value = this.formatDate(now);
    document.getElementById("endDate").value = this.formatDate(now);

    this.currentFilters.startDate = this.formatDate(now);
    this.currentFilters.endDate = this.formatDate(now);
  }

  setupEventListeners() {
    document.getElementById("applyFilters").addEventListener("click", () => {
      this.updateFilters();
      this.loadDashboardData();
    });

    document.getElementById("resetFilters").addEventListener("click", () => {
      this.resetFilters();
      this.loadDashboardData();
    });
  }

  setDefaultDates() {
    const now = new Date();

    document.getElementById("startDate").value = this.formatDate(now);
    document.getElementById("endDate").value = this.formatDate(now);
  }

  formatDate(date) {
    return date.toISOString().split("T")[0];
  }

  updateFilters() {
    this.currentFilters.startDate = document.getElementById("startDate").value;
    this.currentFilters.endDate = document.getElementById("endDate").value;
    this.currentFilters.userEmail = document.getElementById("userFilter").value;
    this.currentFilters.apiOnly = document.getElementById("apiOnly").checked;
    this.currentFilters.businessHours =
      document.getElementById("businessHours").checked;
    this.currentFilters.uniqueOnly =
      document.getElementById("uniqueOnly").checked;
  }

  resetFilters() {
    this.setupFilters();
    document.getElementById("userFilter").value = "";
    document.getElementById("apiOnly").checked = true;
    document.getElementById("businessHours").checked = false;
    document.getElementById("uniqueOnly").checked = true;
    this.updateFilters();
  }

  showLoading() {
    document.getElementById("loadingOverlay").classList.remove("hidden");
  }

  hideLoading() {
    document.getElementById("loadingOverlay").classList.add("hidden");
  }

  async loadDashboardData() {
    this.showLoading();

    try {
      await Promise.all([
        this.loadMAUStats(),
        this.loadUsersList(),
        this.loadTopUsers(),
        this.loadDailyActiveUsers(),
      ]);
    } catch (error) {
      console.error("Error loading dashboard data:", error);
      this.showError("Failed to load dashboard data");
    } finally {
      this.hideLoading();
    }
  }

  async loadMAUStats() {
    const params = new URLSearchParams({
      start_date: this.currentFilters.startDate,
      end_date: this.currentFilters.endDate,
      api_only: this.currentFilters.apiOnly,
      business_hours: this.currentFilters.businessHours,
      unique_only: this.currentFilters.uniqueOnly,
    });

    const response = await fetch(`/admin/api/mau/stats?${params}`);
    const data = await response.json();

    if (data.status === "success" && data.data) {
      this.updateStatsCards(data.data);
      this.updateMonthlyComparisonChart(data.data);
    } else {
      console.error("Invalid MAU stats response:", data);
      this.showError("Failed to load MAU statistics");
    }
  }

  async loadUsersList() {
    const params = new URLSearchParams({
      start_date: this.currentFilters.startDate,
      end_date: this.currentFilters.endDate,
      api_only: this.currentFilters.apiOnly,
      business_hours: this.currentFilters.businessHours,
      unique_only: this.currentFilters.uniqueOnly,
    });

    const response = await fetch(`/admin/api/mau/users?${params}`);
    const data = await response.json();

    if (data.status === "success" && data.data && data.data.users) {
      this.updateUsersDropdown(data.data.users);
    } else {
      console.error("Invalid users list response:", data);
      this.updateUsersDropdown([]); // Fallback to empty array
    }
  }

  async loadTopUsers() {
    const params = new URLSearchParams({
      start_date: this.currentFilters.startDate,
      end_date: this.currentFilters.endDate,
      api_only: this.currentFilters.apiOnly,
      business_hours: this.currentFilters.businessHours,
      unique_only: this.currentFilters.uniqueOnly,
      limit: 10,
    });

    if (this.currentFilters.userEmail) {
      params.append("user_email", this.currentFilters.userEmail);
    }

    const response = await fetch(`/admin/api/mau/activity?${params}`);
    const data = await response.json();

    if (data.status === "success" && data.data && data.data.top_users) {
      this.updateTopUsersTable(data.data.top_users);
    } else {
      console.error("Invalid top users response:", data);
      this.updateTopUsersTable([]); // Fallback to empty array
    }
  }

  async loadDailyActiveUsers() {
    const params = new URLSearchParams({
      days: 30,
      api_only: this.currentFilters.apiOnly,
    });

    const response = await fetch(`/admin/api/mau/daily?${params}`);
    const data = await response.json();

    if (data.status === "success" && data.data && data.data.daily_activities) {
      this.updateDailyUsersChart(data.data.daily_activities);
    } else {
      console.error("Invalid daily active users response:", data);
      console.error("Expected: data.data.daily_activities, got:", data.data);
      this.updateDailyUsersChart([]); // Fallback to empty array
    }
  }

  updateStatsCards(data) {
    // Defensive programming: check if data exists and has expected properties
    if (!data || typeof data !== "object") {
      console.error("Invalid data structure for stats cards:", data);
      return;
    }

    // Safely update current MAU with fallback
    const currentMAU = data.current_mau || 0;
    document.getElementById("currentMAU").textContent = (
      typeof currentMAU === "number" ? currentMAU : parseInt(currentMAU) || 0
    ).toLocaleString();

    // Safely update previous MAU with fallback
    const previousMAU = data.previous_mau || 0;
    document.getElementById("previousMAU").textContent = (
      typeof previousMAU === "number" ? previousMAU : parseInt(previousMAU) || 0
    ).toLocaleString();

    // Safely update total API calls with fallback
    const totalAPICalls = data.total_api_calls || 0;
    document.getElementById("totalAPICalls").textContent = (
      typeof totalAPICalls === "number"
        ? totalAPICalls
        : parseInt(totalAPICalls) || 0
    ).toLocaleString();

    // Safely update growth percentage
    const growthElement = document.getElementById("growthPercentage");
    const growth = data.growth_percentage || 0;
    const growthValue =
      typeof growth === "number" ? growth : parseFloat(growth) || 0;

    growthElement.textContent = `${
      growthValue >= 0 ? "+" : ""
    }${growthValue.toFixed(1)}%`;
    growthElement.className = `text-2xl font-bold ${
      growthValue >= 0 ? "text-green-600" : "text-red-600"
    }`;

    // Show cache status if available
    if (data.cached !== undefined) {
      this.showCacheStatus(data.cached, data.business_hours, data.unique_only);
    }
  }

  showCacheStatus(cached, businessHours, uniqueOnly) {
    // Create or update cache status indicator
    let cacheIndicator = document.getElementById("cacheStatus");
    if (!cacheIndicator) {
      cacheIndicator = document.createElement("div");
      cacheIndicator.id = "cacheStatus";
      cacheIndicator.className = "text-xs mt-2 px-2 py-1 rounded";

      // Insert after the first stats card
      const firstCard = document.querySelector(
        ".bg-white.rounded-lg.shadow-sm.p-6"
      );
      if (firstCard) {
        firstCard.appendChild(cacheIndicator);
      }
    }

    if (cached) {
      cacheIndicator.className =
        "text-xs mt-2 px-2 py-1 rounded bg-green-100 text-green-800";
      cacheIndicator.innerHTML = `<i class="fas fa-bolt mr-1"></i>Data served from cache (12h TTL)`;
    } else {
      cacheIndicator.className =
        "text-xs mt-2 px-2 py-1 rounded bg-blue-100 text-blue-800";
      cacheIndicator.innerHTML = `<i class="fas fa-database mr-1"></i>Data fetched from database`;
    }

    // Add filter status
    const filters = [];
    if (businessHours) filters.push("Business Hours");
    if (uniqueOnly) filters.push("Unique Calls");

    if (filters.length > 0) {
      cacheIndicator.innerHTML += ` • Filters: ${filters.join(", ")}`;
    }
  }

  updateUsersDropdown(users) {
    const select = document.getElementById("userFilter");
    select.innerHTML = '<option value="">All Users</option>';

    // Defensive programming: check if users is an array
    if (!Array.isArray(users)) {
      console.error("Invalid users data structure - expected array:", users);
      return;
    }

    users.forEach((user) => {
      const option = document.createElement("option");

      // Handle both old format (string) and new format (object)
      if (typeof user === "string") {
        option.value = user;
        option.textContent = user;
      } else if (user && typeof user === "object") {
        // Enhanced user data with name and territory - show name as primary text
        option.value = user.email || user;
        option.textContent = user.name || user.email || user; // Show name as primary, fallback to email
        option.setAttribute("data-name", user.name || "");
        option.setAttribute("data-email", user.email || "");
        option.setAttribute("data-territory", user.territory_name || "");
      } else {
        // Fallback for unexpected data types
        option.value = String(user);
        option.textContent = String(user);
      }

      select.appendChild(option);
    });

    // Initialize or refresh Select2
    if ($(select).hasClass("select2-hidden-accessible")) {
      $(select).select2("destroy");
    }

    $(select).select2({
      placeholder: "Search users by name or email...",
      allowClear: true,
      width: "100%",
      theme: "default",
      dropdownAutoWidth: true,
      templateResult: function (user) {
        if (!user.id) return user.text;

        const $option = $(user.element);
        const name = $option.attr("data-name");
        const email = $option.attr("data-email");
        const territory = $option.attr("data-territory");

        if (name && email && territory) {
          return $(`
            <div class="user-option">
              <div class="font-medium text-gray-900">${name}</div>
              <div class="text-sm text-gray-500">${email}</div>
              <div class="text-xs text-blue-600">Territory: ${territory}</div>
            </div>
          `);
        }

        return $(`<div>${user.text}</div>`);
      },
      templateSelection: function (user) {
        if (!user.id) return user.text;

        const $option = $(user.element);
        const name = $option.attr("data-name");

        // Show just the name in the selection
        return name || user.text;
      },
      matcher: function (params, data) {
        // If there are no search terms, return all data
        if ($.trim(params.term) === "") {
          return data;
        }

        const term = params.term.toLowerCase();
        const text = data.text.toLowerCase();
        const $option = $(data.element);
        const name = ($option.attr("data-name") || "").toLowerCase();
        const email = ($option.attr("data-email") || "").toLowerCase();

        // Search in name, email, or full text
        if (
          text.indexOf(term) > -1 ||
          name.indexOf(term) > -1 ||
          email.indexOf(term) > -1
        ) {
          return data;
        }

        // Return null if the term should not be displayed
        return null;
      },
    });
  }

  updateTopUsersTable(users) {
    const tbody = document.getElementById("topUsersTable");
    tbody.innerHTML = "";

    // Defensive programming: check if users is an array
    if (!Array.isArray(users)) {
      console.error("Invalid users data structure - expected array:", users);
      // Show empty state message
      const row = document.createElement("tr");
      row.innerHTML = `
        <td colspan="4" class="px-6 py-4 text-center text-gray-500">
          No user data available
        </td>
      `;
      tbody.appendChild(row);
      return;
    }

    users.forEach((user, index) => {
      // Defensive check for user object
      if (!user || typeof user !== "object") {
        console.warn("Invalid user object at index", index, ":", user);
        return;
      }

      const row = document.createElement("tr");
      row.className = "hover:bg-gray-50";

      const rankDisplay =
        index < 3
          ? `<i class="fas fa-medal text-${
              index === 0 ? "yellow" : index === 1 ? "gray" : "orange"
            }-500 mr-2"></i>${index + 1}`
          : `${index + 1}`;

      // Handle both old format and new enhanced format with safe property access
      const userEmail = user.auth_user_email || user.email || "Unknown";
      const userName = user.name || userEmail;
      const territoryName = user.territory_name || "Not Assigned";

      const userDisplayHtml = user.name
        ? `
        <div class="user-display-container">
          <div class="font-semibold text-gray-900">${userName}</div>
          <div class="text-sm text-gray-500">(${userEmail})</div>
          <div class="text-xs text-blue-600">Territory: ${territoryName}</div>
        </div>
      `
        : `
        <div class="font-medium text-gray-900">${userEmail}</div>
      `;

      // Safely handle call_count and last_activity with fallbacks
      const callCount = user.call_count || 0;
      const callCountValue =
        typeof callCount === "number" ? callCount : parseInt(callCount) || 0;

      const lastActivity = user.last_activity || new Date();
      const lastActivityDate = new Date(lastActivity);
      const lastActivityString = isNaN(lastActivityDate.getTime())
        ? "Unknown"
        : lastActivityDate.toLocaleString();

      row.innerHTML = `
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    ${rankDisplay}
                </td>
                <td class="px-6 py-4 text-sm">
                    ${userDisplayHtml}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                        ${callCountValue.toLocaleString()}
                    </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    ${lastActivityString}
                </td>
            `;
      tbody.appendChild(row);
    });
  }

  updateDailyUsersChart(dailyData) {
    // Defensive programming: check if dailyData is an array
    if (!Array.isArray(dailyData)) {
      console.error(
        "Invalid daily data structure - expected array:",
        dailyData
      );
      return;
    }

    const ctx = document.getElementById("dailyUsersChart").getContext("2d");

    if (this.charts.dailyUsers) {
      this.charts.dailyUsers.destroy();
    }

    // Handle empty data gracefully
    if (dailyData.length === 0) {
      console.warn("No daily activity data available");
      // Create empty chart with placeholder data
      this.charts.dailyUsers = new Chart(ctx, {
        type: "line",
        data: {
          labels: ["No Data"],
          datasets: [
            {
              label: "Daily Active Users",
              data: [0],
              borderColor: "#00B2E5",
              backgroundColor: "rgba(0, 178, 229, 0.1)",
              tension: 0.4,
            },
          ],
        },
        options: {
          responsive: true,
          plugins: {
            title: {
              display: true,
              text: "No daily activity data available",
            },
          },
        },
      });
      return;
    }

    const labels = dailyData.map((item) => {
      // Safe date handling
      const date = item && item.date ? new Date(item.date) : new Date();
      return isNaN(date.getTime()) ? "Invalid Date" : date.toLocaleDateString();
    });

    const data = dailyData.map((item) => {
      // Safe active users handling
      const activeUsers =
        item && item.active_users !== undefined ? item.active_users : 0;
      return typeof activeUsers === "number"
        ? activeUsers
        : parseInt(activeUsers) || 0;
    });

    this.charts.dailyUsers = new Chart(ctx, {
      type: "line",
      data: {
        labels: labels,
        datasets: [
          {
            label: "Daily Active Users",
            data: data,
            borderColor: "#00B2E5",
            backgroundColor: "rgba(0, 178, 229, 0.1)",
            borderWidth: 2,
            fill: true,
            tension: 0.4,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              stepSize: 1,
            },
          },
        },
        plugins: {
          legend: {
            display: false,
          },
        },
      },
    });
  }

  updateMonthlyComparisonChart(data) {
    // Defensive programming: check if data exists and has expected properties
    if (!data || typeof data !== "object") {
      console.error(
        "Invalid data structure for monthly comparison chart:",
        data
      );
      return;
    }

    const ctx = document
      .getElementById("monthlyComparisonChart")
      .getContext("2d");

    if (this.charts.monthlyComparison) {
      this.charts.monthlyComparison.destroy();
    }

    // Safely extract data with fallbacks
    const previousMAU = data.previous_mau || 0;
    const currentMAU = data.current_mau || 0;
    const previousMAUValue =
      typeof previousMAU === "number"
        ? previousMAU
        : parseInt(previousMAU) || 0;
    const currentMAUValue =
      typeof currentMAU === "number" ? currentMAU : parseInt(currentMAU) || 0;

    this.charts.monthlyComparison = new Chart(ctx, {
      type: "bar",
      data: {
        labels: ["Previous Month", "Current Month"],
        datasets: [
          {
            label: "Monthly Active Users",
            data: [previousMAUValue, currentMAUValue],
            backgroundColor: ["#10C0F3", "#00B2E5"],
            borderColor: ["#10C0F3", "#00B2E5"],
            borderWidth: 1,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              stepSize: 1,
            },
          },
        },
        plugins: {
          legend: {
            display: false,
          },
        },
      },
    });
  }

  showError(message) {
    // Simple error display - could be enhanced with a proper notification system
    alert(message);
  }
}

// Initialize dashboard when DOM is loaded
document.addEventListener("DOMContentLoaded", () => {
  new AdminDashboard();
});
