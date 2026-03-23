<template>
  <div class="app-shell">
    <!-- Desktop sidebar -->
    <nav class="sidebar">
      <div class="sidebar-brand">
        <span class="brand-icon">🪂</span>
        <span class="brand-name">SkyLog</span>
      </div>
      <div class="nav-links">
        <router-link to="/log" class="nav-link" active-class="nav-link--active">
          <span class="nav-icon">📋</span>
          <span>Logbook</span>
        </router-link>
        <router-link to="/stats" class="nav-link" active-class="nav-link--active">
          <span class="nav-icon">📊</span>
          <span>Stats</span>
        </router-link>
        <router-link to="/docs" class="nav-link" active-class="nav-link--active">
          <span class="nav-icon">📁</span>
          <span>Documents</span>
        </router-link>
      </div>
    </nav>

    <!-- Main content -->
    <main class="main-content">
      <slot />
    </main>

    <!-- Mobile bottom nav -->
    <nav class="bottom-nav">
      <router-link to="/log" class="bottom-nav-item" active-class="bottom-nav-item--active">
        <span class="nav-icon">📋</span>
        <span class="bottom-nav-label">Log</span>
      </router-link>
      <router-link to="/stats" class="bottom-nav-item" active-class="bottom-nav-item--active">
        <span class="nav-icon">📊</span>
        <span class="bottom-nav-label">Stats</span>
      </router-link>
      <router-link to="/docs" class="bottom-nav-item" active-class="bottom-nav-item--active">
        <span class="nav-icon">📁</span>
        <span class="bottom-nav-label">Docs</span>
      </router-link>
    </nav>
  </div>
</template>

<style scoped>
.app-shell {
  display: flex;
  height: 100dvh;
  overflow: hidden;
}

/* === Sidebar (desktop) === */
.sidebar {
  display: none;
  flex-direction: column;
  width: 220px;
  min-width: 220px;
  background: var(--sky-900);
  border-right: 1px solid var(--sky-800);
  padding: 1.5rem 1rem;
  gap: 0.5rem;
}

@media (min-width: 768px) {
  .sidebar { display: flex; }
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem 0.75rem 1.5rem;
}

.brand-icon { font-size: 1.5rem; }
.brand-name {
  font-size: 1.25rem;
  font-weight: 700;
  background: linear-gradient(135deg, var(--sky-accent), var(--sunset-400));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.nav-links { display: flex; flex-direction: column; gap: 0.25rem; }

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.625rem 0.75rem;
  border-radius: 0.5rem;
  color: #94a3b8;
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.15s;
}
.nav-link:hover { background: var(--sky-800); color: #e2e8f0; }
.nav-link--active { background: linear-gradient(135deg, rgba(56,189,248,.15), rgba(249,115,22,.1)); color: var(--sky-accent); }

.nav-icon { font-size: 1.1rem; }

/* === Main content === */
.main-content {
  flex: 1;
  overflow-y: auto;
  background: var(--sky-950);
  padding-bottom: 4.5rem; /* room for bottom nav on mobile */
}

@media (min-width: 768px) {
  .main-content { padding-bottom: 0; }
}

/* === Bottom nav (mobile) === */
.bottom-nav {
  display: flex;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: var(--sky-900);
  border-top: 1px solid var(--sky-800);
  z-index: 50;
}

@media (min-width: 768px) {
  .bottom-nav { display: none; }
}

.bottom-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  padding: 0.5rem 0;
  gap: 0.2rem;
  color: #64748b;
  text-decoration: none;
  transition: color 0.15s;
}
.bottom-nav-item--active { color: var(--sky-accent); }

.bottom-nav-label { font-size: 0.7rem; font-weight: 500; }
</style>
