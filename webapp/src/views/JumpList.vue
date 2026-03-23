<template>
  <div class="jump-list">
    <!-- Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Logbook</h1>
        <p v-if="!loading && jumps.length" class="page-subtitle">
          {{ total }} jump{{ total !== 1 ? 's' : '' }} logged
        </p>
      </div>
      <router-link to="/log/new" class="btn btn-primary">+ New Jump</router-link>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="state-center">
      <div class="spinner"></div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="state-center error-box">
      <p>{{ error }}</p>
      <button class="btn btn-secondary" @click="loadJumps">Retry</button>
    </div>

    <!-- Empty state -->
    <div v-else-if="jumps.length === 0" class="state-center empty-state">
      <div class="empty-icon">🪂</div>
      <h2>No jumps yet</h2>
      <p>Log your first skydive to get started.</p>
      <router-link to="/log/new" class="btn btn-primary">Log a Jump</router-link>
    </div>

    <!-- Jump cards -->
    <div v-else class="jump-grid">
      <div
        v-for="jump in jumps"
        :key="jump.id"
        class="jump-card"
      >
        <div class="jump-card-header">
          <span class="jump-number">#{{ jump.jumpNumber }}</span>
          <span class="jump-type-badge">{{ jump.jumpType || 'Fun' }}</span>
        </div>
        <div class="jump-card-body">
          <div class="jump-date">{{ formatDate(jump.date) }}</div>
          <div class="jump-dz">{{ jump.dropzone || jump.location || '—' }}</div>
        </div>
        <div class="jump-card-stats">
          <div class="jump-stat">
            <span class="stat-label">Exit Alt</span>
            <span class="stat-value">{{ jump.exitAltitude ? jump.exitAltitude.toLocaleString() + ' ft' : '—' }}</span>
          </div>
          <div class="jump-stat">
            <span class="stat-label">Freefall</span>
            <span class="stat-value">{{ formatFF(jump.freefallTime) }}</span>
          </div>
          <div class="jump-stat">
            <span class="stat-label">Canopy</span>
            <span class="stat-value">{{ jump.canopySize || '—' }}</span>
          </div>
        </div>
        <div class="jump-card-actions">
          <router-link :to="`/log/${jump.id}/edit`" class="action-btn">Edit</router-link>
          <button class="action-btn action-btn--danger" @click="deleteJump(jump)">Delete</button>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="total > perPage" class="pagination">
      <button class="btn btn-secondary" :disabled="page === 1" @click="page--; loadJumps()">← Prev</button>
      <span class="page-info">{{ page }} / {{ Math.ceil(total / perPage) }}</span>
      <button class="btn btn-secondary" :disabled="page >= Math.ceil(total / perPage)" @click="page++; loadJumps()">Next →</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { jumpsApi } from '../api/jumps.js'

const jumps = ref([])
const total = ref(0)
const page = ref(1)
const perPage = 20
const loading = ref(false)
const error = ref(null)

async function loadJumps() {
  loading.value = true
  error.value = null
  try {
    const res = await jumpsApi.list({ page: page.value, per_page: perPage, order: 'desc' })
    jumps.value = res.jumps || []
    total.value = res.total || 0
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

async function deleteJump(jump) {
  if (!confirm(`Delete jump #${jump.jumpNumber}?`)) return
  try {
    await jumpsApi.delete(jump.id)
    jumps.value = jumps.value.filter(j => j.id !== jump.id)
    total.value--
  } catch (e) {
    alert('Failed to delete jump: ' + e.message)
  }
}

function formatDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

function formatFF(seconds) {
  if (!seconds) return '—'
  return seconds >= 60
    ? `${Math.floor(seconds / 60)}m ${seconds % 60}s`
    : `${seconds}s`
}

onMounted(loadJumps)
</script>

<style scoped>
.jump-list {
  padding: 1.5rem 1rem;
  max-width: 960px;
  margin: 0 auto;
}

@media (min-width: 768px) {
  .jump-list { padding: 2rem 2rem; }
}

/* Header */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
  gap: 1rem;
  flex-wrap: wrap;
}
.page-title {
  font-size: 1.75rem;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(135deg, var(--sky-accent), var(--sunset-400));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.page-subtitle { margin: 0.25rem 0 0; color: #64748b; font-size: 0.9rem; }

/* States */
.state-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 4rem 1rem;
  text-align: center;
}
.empty-state .empty-icon { font-size: 4rem; }
.empty-state h2 { margin: 0; font-size: 1.5rem; color: #e2e8f0; }
.empty-state p { margin: 0; color: #64748b; }
.error-box { color: #f87171; }

/* Spinner */
.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--sky-800);
  border-top-color: var(--sky-accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* Jump grid */
.jump-grid {
  display: grid;
  gap: 1rem;
  grid-template-columns: 1fr;
}
@media (min-width: 600px) {
  .jump-grid { grid-template-columns: repeat(2, 1fr); }
}
@media (min-width: 960px) {
  .jump-grid { grid-template-columns: repeat(3, 1fr); }
}

/* Jump card */
.jump-card {
  background: var(--sky-900);
  border: 1px solid var(--sky-800);
  border-radius: 0.75rem;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  transition: border-color 0.15s, transform 0.15s;
}
.jump-card:hover {
  border-color: var(--sky-accent);
  transform: translateY(-2px);
}

.jump-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.jump-number { font-size: 1.1rem; font-weight: 700; color: var(--sky-accent); }
.jump-type-badge {
  font-size: 0.7rem;
  font-weight: 600;
  padding: 0.2rem 0.5rem;
  border-radius: 999px;
  background: rgba(56,189,248,.15);
  color: var(--sky-accent);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.jump-card-body { }
.jump-date { font-size: 0.85rem; color: #94a3b8; }
.jump-dz { font-weight: 600; color: #e2e8f0; }

.jump-card-stats {
  display: flex;
  gap: 0.75rem;
  border-top: 1px solid var(--sky-800);
  padding-top: 0.75rem;
}
.jump-stat { display: flex; flex-direction: column; gap: 0.15rem; flex: 1; }
.stat-label { font-size: 0.65rem; color: #64748b; text-transform: uppercase; letter-spacing: 0.05em; }
.stat-value { font-size: 0.85rem; font-weight: 500; color: #cbd5e1; }

.jump-card-actions {
  display: flex;
  gap: 0.5rem;
  border-top: 1px solid var(--sky-800);
  padding-top: 0.75rem;
}
.action-btn {
  flex: 1;
  padding: 0.375rem;
  border-radius: 0.375rem;
  font-size: 0.8rem;
  font-weight: 500;
  text-align: center;
  text-decoration: none;
  border: 1px solid var(--sky-700);
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  transition: all 0.15s;
}
.action-btn:hover { background: var(--sky-800); color: #e2e8f0; }
.action-btn--danger:hover { border-color: #ef4444; color: #f87171; background: rgba(239,68,68,.1); }

/* Pagination */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  margin-top: 2rem;
}
.page-info { color: #64748b; font-size: 0.9rem; }

/* Shared buttons */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.9rem;
  font-weight: 600;
  text-decoration: none;
  cursor: pointer;
  border: none;
  transition: all 0.15s;
}
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-primary {
  background: linear-gradient(135deg, var(--sky-accent-dim), var(--sunset-500));
  color: white;
}
.btn-primary:hover { opacity: 0.9; transform: translateY(-1px); }
.btn-secondary {
  background: var(--sky-800);
  color: #94a3b8;
  border: 1px solid var(--sky-700);
}
.btn-secondary:hover:not(:disabled) { background: var(--sky-700); color: #e2e8f0; }
</style>
