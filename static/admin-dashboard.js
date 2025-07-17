document.addEventListener("DOMContentLoaded", function () {
  // Authentication check
  const token = localStorage.getItem("admin_token");
  if (!token) {
    window.location.href = "/admin/login";
    return;
  }

  // DOM Elements
  const mauTab = document.getElementById("mauTab");
  const logTab = document.getElementById("logTab");
  const mauTabMobile = document.getElementById("mauTabMobile");
  const logTabMobile = document.getElementById("logTabMobile");
  const mauModule = document.getElementById("mauModule");
  const logModule = document.getElementById("logModule");
  const userEmail = document.getElementById("userEmail");
  const userMenuBtn = document.getElementById("userMenuBtn");
  const userDropdown = document.getElementById("userDropdown");
  const logoutBtn = document.getElementById("logoutBtn");
  const loadingOverlay = document.getElementById("loadingOverlay");

  // State
  let currentModule = "mau";
  let userInfo = null;

  // Initialize
  init();

  async function init() {
    try {
      showLoading(true);

      // Verify authentication and get user info
      await verifyAuth();

      // Load initial module (MAU Dashboard)
      await loadModule("mau");
    } catch (error) {
      console.error("Initialization error:", error);
      logout();
    } finally {
      showLoading(false);
    }
  }

  // Authentication verification
  async function verifyAuth() {
    try {
      const response = await fetch("/api/v1/users/profile", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const data = await response.json();

      // Handle the correct API response format with meta object
      if (data.meta && data.meta.status === "success") {
        userInfo = data.data;
        userEmail.textContent = userInfo.email;

        // Check if user has Super-Admin role
        if (userInfo.user_role !== "Super-Admin") {
          throw new Error("Access denied. Super-Admin privileges required.");
        }
      } else {
        throw new Error("Authentication failed");
      }
    } catch (error) {
      throw error;
    }
  }

  // Module switching
  function switchModule(module) {
    if (currentModule === module) return;

    currentModule = module;

    // Update tab states
    updateTabStates(module);

    // Show/hide modules with animation
    if (module === "mau") {
      mauModule.classList.remove("hidden");
      logModule.classList.add("hidden");

      setTimeout(() => {
        mauModule.classList.add("active");
        logModule.classList.remove("active");
      }, 50);
    } else {
      logModule.classList.remove("hidden");
      mauModule.classList.add("hidden");

      setTimeout(() => {
        logModule.classList.add("active");
        mauModule.classList.remove("active");
      }, 50);
    }

    // Load module content
    loadModule(module);
  }

  function updateTabStates(activeModule) {
    const tabs = [mauTab, mauTabMobile, logTab, logTabMobile];

    tabs.forEach((tab) => {
      tab.classList.remove("active");
    });

    if (activeModule === "mau") {
      mauTab.classList.add("active");
      mauTabMobile.classList.add("active");
    } else {
      logTab.classList.add("active");
      logTabMobile.classList.add("active");
    }
  }

  // Load module content
  async function loadModule(module) {
    try {
      if (module === "mau") {
        await loadMAUDashboard();
      } else {
        await loadLogViewer();
      }
    } catch (error) {
      console.error(`Error loading ${module} module:`, error);
      showError(`Failed to load ${module} module`);
    }
  }

  // Load MAU Dashboard
  async function loadMAUDashboard() {
    try {
      const response = await fetch("/static/mau-dashboard.html");
      const html = await response.text();
      document.getElementById("mauDashboardContent").innerHTML = html;

      // Initialize MAU dashboard functionality
      if (window.initMAUDashboard) {
        window.initMAUDashboard(token);
      }
    } catch (error) {
      document.getElementById("mauDashboardContent").innerHTML = `
        <div class="text-center py-12">
          <i class="fas fa-exclamation-triangle text-4xl text-red-400 mb-4"></i>
          <p class="text-red-600">Failed to load MAU Dashboard</p>
          <button onclick="loadModule('mau')" class="mt-4 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600">
            Retry
          </button>
        </div>
      `;
    }
  }

  // Load Log Viewer
  async function loadLogViewer() {
    try {
      const response = await fetch("/static/log-viewer-content.html");
      const html = await response.text();
      document.getElementById("logViewerContent").innerHTML = html;

      // Initialize log viewer functionality
      if (window.initLogViewer) {
        window.initLogViewer(token);
      }
    } catch (error) {
      document.getElementById("logViewerContent").innerHTML = `
        <div class="text-center py-12">
          <i class="fas fa-exclamation-triangle text-4xl text-red-400 mb-4"></i>
          <p class="text-red-600">Failed to load Log Viewer</p>
          <button onclick="loadModule('log')" class="mt-4 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600">
            Retry
          </button>
        </div>
      `;
    }
  }

  // Utility functions
  function showLoading(show) {
    if (show) {
      loadingOverlay.classList.remove("hidden");
    } else {
      loadingOverlay.classList.add("hidden");
    }
  }

  function showError(message) {
    // You could implement a toast notification system here
    console.error(message);
    alert(message); // Simple fallback
  }

  function logout() {
    localStorage.removeItem("admin_token");
    window.location.href = "/admin/login";
  }

  // Event listeners
  mauTab.addEventListener("click", () => switchModule("mau"));
  logTab.addEventListener("click", () => switchModule("log"));
  mauTabMobile.addEventListener("click", () => switchModule("mau"));
  logTabMobile.addEventListener("click", () => switchModule("log"));

  // User menu toggle
  userMenuBtn.addEventListener("click", function (e) {
    e.stopPropagation();
    userDropdown.classList.toggle("hidden");
  });

  // Close dropdown when clicking outside
  document.addEventListener("click", function () {
    userDropdown.classList.add("hidden");
  });

  // Logout
  logoutBtn.addEventListener("click", logout);

  // Expose functions globally for module scripts
  window.adminDashboard = {
    showLoading,
    showError,
    token,
    userInfo: () => userInfo,
    loadModule,
  };
});
