// Base API URL configuration
const API_BASE = window.location.port === '4000' ? '' : 'http://localhost:4000';

// Global session cache (metadata only, NEVER the raw JWT)
let currentSession = null;

// Route switching
function showView(viewName) {
  document.querySelectorAll('.view-section').forEach(el => el.classList.add('hidden'));
  document.querySelectorAll('.nav-btn').forEach(el => el.classList.remove('active'));

  const targetView = document.getElementById(viewName + '-view');
  if (targetView) targetView.classList.remove('hidden');

  const navBtn = document.getElementById('nav-' + viewName);
  if (navBtn) navBtn.classList.add('active');

  if (viewName === 'dashboard') {
    loadDashboardSession();
  }
}

// Alert helper
function showAlert(containerId, message, type = 'info') {
  const container = document.getElementById(containerId);
  if (!container) return;
  container.innerHTML = `<div class="alert alert-${type}">${message}</div>`;
  container.classList.remove('hidden');
}

function clearAlert(containerId) {
  const container = document.getElementById(containerId);
  if (container) {
    container.innerHTML = '';
    container.classList.add('hidden');
  }
}

// Login Handler - Uses HttpOnly cookies for secure production session management
async function handleLogin(e) {
  if (e) e.preventDefault();
  clearAlert('login-alert');

  const email = document.getElementById('login-email').value.trim();
  const password = document.getElementById('login-password').value;

  if (!email || !password) {
    showAlert('login-alert', 'Please enter both email and password.', 'danger');
    return;
  }

  const submitBtn = document.getElementById('btn-login-submit');
  submitBtn.innerText = 'Authenticating...';
  submitBtn.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'same-origin', // Transmits secure HttpOnly cookie
      body: JSON.stringify({ email, password })
    });

    const data = await res.json();

    if (res.ok && data.code === 200) {
      // NOTE: We deliberately do NOT store raw JWT in localStorage or sessionStorage.
      // The browser automatically and securely stores the HttpOnly session cookie.
      showAlert('login-alert', '✅ Login successful! Establishing secure session...', 'success');
      setTimeout(() => {
        showView('dashboard');
      }, 600);
    } else {
      const errMsg = data.details || data.message || 'Login failed. Invalid credentials.';
      showAlert('login-alert', `❌ ${errMsg}`, 'danger');
    }
  } catch (err) {
    showAlert('login-alert', `Server connection error: ${err.message}.`, 'danger');
  } finally {
    submitBtn.innerText = 'Sign In';
    submitBtn.disabled = false;
  }
}

// Registration Handler
async function handleRegister(e) {
  if (e) e.preventDefault();
  clearAlert('register-alert');

  const username = document.getElementById('reg-username').value.trim();
  const email = document.getElementById('reg-email').value.trim();
  const password = document.getElementById('reg-password').value;
  const full_name = document.getElementById('reg-fullname').value.trim();
  const phone_number = document.getElementById('reg-phone').value.trim();

  if (!username || !email || !password) {
    showAlert('register-alert', 'Username, email and password are required.', 'danger');
    return;
  }

  const submitBtn = document.getElementById('btn-register-submit');
  submitBtn.innerText = 'Registering...';
  submitBtn.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/users`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, email, password, full_name, phone_number })
    });

    const data = await res.json();

    if (data.code === 201 || (data.message && data.message.includes('successfully'))) {
      showAlert('register-alert', `✅ Registration successful! User ID: ${data.details}. Redirecting to sign in...`, 'success');
      setTimeout(() => {
        document.getElementById('login-email').value = email;
        document.getElementById('login-password').value = password;
        showView('login');
      }, 1500);
    } else {
      const errMsg = data.details || data.message || 'Registration failed';
      showAlert('register-alert', `❌ ${errMsg}`, 'danger');
    }
  } catch (err) {
    showAlert('register-alert', `Server error: ${err.message}. Ensure backend is running.`, 'danger');
  } finally {
    submitBtn.innerText = 'Create Account';
    submitBtn.disabled = false;
  }
}

// Demo Token Handler - Sets secure cookie
async function handleDemoToken() {
  clearAlert('login-alert');
  try {
    const res = await fetch(`${API_BASE}/demo-token`, {
      credentials: 'same-origin'
    });
    const data = await res.json();
    if (res.ok && data.code === 200) {
      showAlert('login-alert', '✅ Demo session initialized! Loading Protected Dashboard...', 'success');
      setTimeout(() => showView('dashboard'), 600);
    } else {
      showAlert('login-alert', 'Failed to generate demo session: ' + data.details, 'danger');
    }
  } catch (err) {
    showAlert('login-alert', 'Error connecting to demo endpoint: ' + err.message, 'danger');
  }
}

// Dynamic Session Loader - Fetches safe metadata (NEVER returns raw JWT)
async function loadDashboardSession() {
  try {
    const res = await fetch(`${API_BASE}/api/v1/auth/session`, {
      method: 'GET',
      credentials: 'same-origin' // Sends HttpOnly cookie
    });

    if (res.status === 401 || !res.ok) {
      // Unauthenticated or expired session
      updateNavState(false);
      renderUnauthenticatedDashboard();
      showAlert('login-alert', 'Your session is inactive or expired. Please sign in.', 'danger');
      showView('login');
      return;
    }

    const data = await res.json();
    if (data.authenticated) {
      currentSession = data;
      updateNavState(true);
      renderAuthenticatedDashboard(data);
      testProtectedEndpoint(); // Verify protected route with session cookie
    } else {
      updateNavState(false);
      showView('login');
    }
  } catch (err) {
    console.error('Session retrieval error:', err);
    updateNavState(false);
    showAlert('login-alert', `Unable to verify session: ${err.message}`, 'danger');
    showView('login');
  }
}

// Render Authentication Details Component
function renderAuthenticatedDashboard(session) {
  // 1. User Identity Details
  if (session.user) {
    document.getElementById('dash-welcome-user').innerText = session.user.username || 'User';
    document.getElementById('dash-user-id').innerText = session.user.id || 'N/A';
    document.getElementById('dash-username').innerText = session.user.username || 'N/A';
    document.getElementById('dash-role').innerText = session.user.role || 'user';
    document.getElementById('dash-status').innerText = session.user.status || 'active';
  }

  // 2. Authentication Details Component
  const statusBadge = document.getElementById('auth-status-badge');
  statusBadge.innerText = session.session_status === 'active' ? 'Active' : 'Inactive';
  statusBadge.className = `claim-value status-badge ${session.session_status === 'active' ? 'active' : 'inactive'}`;

  document.getElementById('auth-token-type').innerText = session.token_type || 'Bearer JWT';

  // Dynamic Masked Token Display
  document.getElementById('auth-masked-token').innerText = session.masked_token || '••••••••••••••••';

  const tokenStatusBadge = document.getElementById('auth-token-status');
  tokenStatusBadge.innerText = session.token_status || 'Valid';
  tokenStatusBadge.className = `claim-value status-badge ${session.token_status === 'Valid' ? 'valid' : 'expired'}`;

  // Expiration time display (formatted local & UTC)
  if (session.expires_at) {
    const expDate = new Date(session.expires_at);
    document.getElementById('auth-expiration').innerText = expDate.toLocaleString() + ` (${expDate.toUTCString()})`;
  } else {
    document.getElementById('auth-expiration').innerText = 'N/A';
  }
}

function renderUnauthenticatedDashboard() {
  document.getElementById('dash-welcome-user').innerText = 'Guest';
  document.getElementById('auth-status-badge').innerText = 'Inactive';
  document.getElementById('auth-status-badge').className = 'claim-value status-badge inactive';
  document.getElementById('auth-masked-token').innerText = '••••••••••••••••';
  document.getElementById('auth-token-status').innerText = 'Missing';
  document.getElementById('auth-token-status').className = 'claim-value status-badge expired';
  document.getElementById('auth-expiration').innerText = 'No active session';
}

// Test Protected Endpoint
async function testProtectedEndpoint() {
  const resultBox = document.getElementById('protected-endpoint-result');
  resultBox.innerText = 'Calling GET /protected with session cookie...';

  try {
    const res = await fetch(`${API_BASE}/protected`, {
      method: 'GET',
      credentials: 'same-origin' // Sends HttpOnly cookie
    });

    const data = await res.json();
    resultBox.innerText = JSON.stringify(data, null, 2);
  } catch (err) {
    resultBox.innerText = `Error calling /protected: ${err.message}`;
  }
}

// Sign Out Handler
async function handleLogout() {
  try {
    await fetch(`${API_BASE}/api/v1/auth/logout`, {
      method: 'POST',
      credentials: 'same-origin'
    });
  } catch (e) {
    console.error('Logout error:', e);
  }

  currentSession = null;
  updateNavState(false);
  renderUnauthenticatedDashboard();
  showAlert('login-alert', 'You have been securely signed out.', 'info');
  showView('login');
}

function updateNavState(isAuthenticated) {
  const navLogout = document.getElementById('nav-logout');
  const navDashboard = document.getElementById('nav-dashboard');
  if (isAuthenticated) {
    if (navLogout) navLogout.classList.remove('hidden');
    if (navDashboard) navDashboard.classList.remove('hidden');
  } else {
    if (navLogout) navLogout.classList.add('hidden');
  }
}

// Initialization
document.addEventListener('DOMContentLoaded', async () => {
  // Check active session dynamically on initial page load / refresh
  try {
    const res = await fetch(`${API_BASE}/api/v1/auth/session`, {
      credentials: 'same-origin'
    });
    if (res.ok) {
      const data = await res.json();
      if (data.authenticated) {
        currentSession = data;
        updateNavState(true);
        if (window.location.hash !== '#register') {
          showView('dashboard');
          return;
        }
      }
    }
  } catch (e) {
    // Session inactive
  }

  updateNavState(false);
  const hash = window.location.hash;
  if (hash === '#register') {
    showView('register');
  } else {
    showView('login');
  }
});
