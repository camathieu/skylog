<template>
  <div class="jump-form-page">
    <div class="form-header">
      <router-link to="/log" class="back-btn">← Back</router-link>
      <h1 class="page-title">{{ isEdit ? 'Edit Jump' : 'New Jump' }}</h1>
    </div>

    <div v-if="loading" class="state-center"><div class="spinner"></div></div>
    <div v-else-if="fetchError" class="state-center error-box">{{ fetchError }}</div>

    <form v-else class="form-card" @submit.prevent="submit">
      <!-- Row 1: Jump # and Date -->
      <div class="form-row">
        <div class="form-group" v-if="isEdit">
          <label for="jumpNumber">Jump #</label>
          <input id="jumpNumber" :value="form.jumpNumber" type="number" disabled class="input-disabled" />
        </div>
        <div class="form-group" v-else>
          <label>Jump #</label>
          <input type="text" value="Auto-assigned" disabled class="input-disabled" />
        </div>
        <div class="form-group">
          <label for="date">Date *</label>
          <input id="date" v-model="form.date" type="date" required />
        </div>
      </div>

      <!-- Row 2: Location and Dropzone -->
      <div class="form-row">
        <div class="form-group">
          <label for="location">Location</label>
          <input id="location" v-model="form.location" type="text" placeholder="City, Country" />
        </div>
        <div class="form-group">
          <label for="dropzone">Dropzone</label>
          <input id="dropzone" v-model="form.dropzone" type="text" placeholder="e.g. Skydive Empuriabrava" />
        </div>
      </div>

      <!-- Row 3: Aircraft and Jump type -->
      <div class="form-row">
        <div class="form-group">
          <label for="aircraft">Aircraft</label>
          <input id="aircraft" v-model="form.aircraft" type="text" placeholder="e.g. Cessna 208" />
        </div>
        <div class="form-group">
          <label for="jumpType">Jump Type</label>
          <select id="jumpType" v-model="form.jumpType">
            <option value="">Select type…</option>
            <option v-for="t in jumpTypes" :key="t" :value="t">{{ t }}</option>
          </select>
        </div>
      </div>

      <!-- Row 4: Altitudes -->
      <div class="form-row">
        <div class="form-group">
          <label for="exitAlt">Exit Altitude (ft)</label>
          <input id="exitAlt" v-model.number="form.exitAltitude" type="number" min="0" placeholder="e.g. 13000" />
        </div>
        <div class="form-group">
          <label for="deployAlt">Deploy Altitude (ft)</label>
          <input id="deployAlt" v-model.number="form.deploymentAltitude" type="number" min="0" placeholder="e.g. 3500" />
        </div>
      </div>

      <!-- Row 5: Freefall and Canopy -->
      <div class="form-row">
        <div class="form-group">
          <label for="freefallTime">Freefall Time (s)</label>
          <input id="freefallTime" v-model.number="form.freefallTime" type="number" min="0" placeholder="e.g. 55" />
        </div>
        <div class="form-group">
          <label for="canopySize">Canopy / Size</label>
          <input id="canopySize" v-model="form.canopySize" type="text" placeholder="e.g. Pilot 188" />
        </div>
      </div>

      <!-- Notes -->
      <div class="form-group form-group--full">
        <label for="notes">Notes</label>
        <textarea id="notes" v-model="form.notes" rows="3" placeholder="Anything worth remembering…"></textarea>
      </div>

      <!-- Submit -->
      <div class="form-actions">
        <router-link to="/log" class="btn btn-secondary">Cancel</router-link>
        <button type="submit" class="btn btn-primary" :disabled="saving">
          {{ saving ? 'Saving…' : (isEdit ? 'Save Changes' : 'Log Jump') }}
        </button>
      </div>

      <p v-if="saveError" class="error-msg">{{ saveError }}</p>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { jumpsApi } from '../api/jumps.js'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const loading = ref(false)
const saving = ref(false)
const fetchError = ref(null)
const saveError = ref(null)

const jumpTypes = [
  'Fun', 'Formation (FS/VFS)', 'AFF', 'Tandem', 'Wingsuit',
  'CRW', 'Freefly', 'Angle', 'Tracking', 'Hop & Pop', 'Demo', 'Other',
]

const form = ref({
  jumpNumber: null,
  date: new Date().toISOString().slice(0, 10),
  location: '',
  dropzone: '',
  aircraft: '',
  exitAltitude: null,
  deploymentAltitude: null,
  freefallTime: null,
  canopySize: '',
  jumpType: '',
  notes: '',
})

onMounted(async () => {
  if (!isEdit.value) return
  loading.value = true
  try {
    const jump = await jumpsApi.get(route.params.id)
    form.value = {
      ...jump,
      date: jump.date ? jump.date.slice(0, 10) : '',
    }
  } catch (e) {
    fetchError.value = e.message
  } finally {
    loading.value = false
  }
})

async function submit() {
  saving.value = true
  saveError.value = null
  try {
    const payload = {
      ...form.value,
      date: form.value.date ? new Date(form.value.date).toISOString() : null,
    }
    if (isEdit.value) {
      await jumpsApi.update(route.params.id, payload)
    } else {
      await jumpsApi.create(payload)
    }
    router.push('/log')
  } catch (e) {
    saveError.value = e.message
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.jump-form-page {
  padding: 1.5rem 1rem;
  max-width: 640px;
  margin: 0 auto;
}
@media (min-width: 768px) {
  .jump-form-page { padding: 2rem; }
}

.form-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
}
.back-btn {
  color: #64748b;
  text-decoration: none;
  font-size: 0.9rem;
  transition: color 0.15s;
  white-space: nowrap;
}
.back-btn:hover { color: var(--sky-accent); }
.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(135deg, var(--sky-accent), var(--sunset-400));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.form-card {
  background: var(--sky-900);
  border: 1px solid var(--sky-800);
  border-radius: 0.75rem;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-row {
  display: grid;
  gap: 1rem;
  grid-template-columns: 1fr;
}
@media (min-width: 480px) {
  .form-row { grid-template-columns: 1fr 1fr; }
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}
.form-group--full { grid-column: 1 / -1; }

label {
  font-size: 0.8rem;
  font-weight: 600;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

input, select, textarea {
  background: var(--sky-800);
  border: 1px solid var(--sky-700);
  border-radius: 0.5rem;
  color: #e2e8f0;
  font-size: 0.9rem;
  padding: 0.5rem 0.75rem;
  font-family: inherit;
  transition: border-color 0.15s;
  width: 100%;
}
input:focus:not(:disabled), select:focus:not(:disabled), textarea:focus:not(:disabled) {
  outline: none;
  border-color: var(--sky-accent);
}
.input-disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
input::placeholder, textarea::placeholder { color: #475569; }
textarea { resize: vertical; }

/* States */
.state-center {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4rem;
}
.spinner {
  width: 40px; height: 40px;
  border: 3px solid var(--sky-800);
  border-top-color: var(--sky-accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
.error-box { color: #f87171; }

/* Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding-top: 0.5rem;
  border-top: 1px solid var(--sky-800);
}
.error-msg { color: #f87171; font-size: 0.85rem; text-align: center; margin: 0; }

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 1.25rem;
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
.btn-primary:hover:not(:disabled) { opacity: 0.9; }
.btn-secondary {
  background: var(--sky-800);
  color: #94a3b8;
  border: 1px solid var(--sky-700);
}
.btn-secondary:hover { background: var(--sky-700); color: #e2e8f0; }
</style>
