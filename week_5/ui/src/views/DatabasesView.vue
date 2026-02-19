<script setup>
import { computed, onMounted, ref } from 'vue'
import { apiClient } from '@/services/api'

const databases = ref([])
const loading = ref(false)
const creating = ref(false)
const deletingName = ref('')
const updatingName = ref('')
const error = ref('')
const selectedDb = ref(null)

const newDbName = ref('')
const newDbInstances = ref(1)
const newDbStorage = ref('1Gi')
const editDbName = ref('')
const editDbInstances = ref(1)
const editDbStorage = ref('1Gi')
const originalEditDbInstances = ref(1)
const originalEditDbStorage = ref('1Gi')

const databaseCount = computed(() => `${databases.value.length} total`)

const rowMeta = (db) => {
  const status = db?.status ?? 'unknown'
  const instances = db?.spec?.instances ?? 'n/a'
  const storage = db?.spec?.storage ?? 'n/a'
  return `${status} · ${instances} replicas · ${storage}`
}

const fetchDatabases = async () => {
  loading.value = true
  error.value = ''

  try {
    const response = await apiClient.get('/databases')
    databases.value = response.data
  } catch (err) {
    console.error(err)
    error.value = 'Failed to load databases.'
  } finally {
    loading.value = false
  }
}

const createDatabase = async () => {
  creating.value = true
  error.value = ''

  try {
    await apiClient.post('/databases', {
      name: newDbName.value,
      instances: newDbInstances.value,
      storage: newDbStorage.value,
    })
    newDbName.value = ''
    newDbInstances.value = 1
    newDbStorage.value = '1Gi'
    await fetchDatabases()
  } catch (err) {
    console.error(err)
    error.value = 'Failed to create database.'
  } finally {
    creating.value = false
  }
}

const deleteDatabase = async (name) => {
  deletingName.value = name
  error.value = ''

  try {
    await apiClient.delete(`/databases/${name}`)
    if (selectedDb.value?.name === name) {
      selectedDb.value = null
    }
    await fetchDatabases()
  } catch (err) {
    console.error(err)
    error.value = 'Failed to delete database.'
  } finally {
    deletingName.value = ''
  }
}

const isEditing = (name) => editDbName.value === name

const startEdit = (db) => {
  editDbName.value = db.name
  editDbInstances.value = db?.spec?.instances ?? 1
  editDbStorage.value = db?.spec?.storage ?? '1Gi'
  originalEditDbInstances.value = editDbInstances.value
  originalEditDbStorage.value = editDbStorage.value
  error.value = ''
}

const cancelEdit = () => {
  editDbName.value = ''
  editDbInstances.value = 1
  editDbStorage.value = '1Gi'
  originalEditDbInstances.value = 1
  originalEditDbStorage.value = '1Gi'
}

const updateDatabase = async (name) => {
  updatingName.value = name
  error.value = ''

  const payload = {}

  if (editDbInstances.value !== originalEditDbInstances.value) {
    payload.instances = editDbInstances.value
  }

  if (editDbStorage.value.trim() !== originalEditDbStorage.value.trim()) {
    payload.storage = editDbStorage.value.trim()
  }

  if (Object.keys(payload).length === 0) {
    error.value = 'No changes to save.'
    updatingName.value = ''
    return
  }

  try {
    await apiClient.patch(`/databases/${name}`, payload)
    await fetchDatabases()
    if (selectedDb.value?.name === name) {
      await viewConnection(name)
    }
    cancelEdit()
  } catch (err) {
    console.error(err)
    error.value = 'Failed to update database.'
  } finally {
    updatingName.value = ''
  }
}

const viewConnection = async (name) => {
  error.value = ''
  try {
    const response = await apiClient.get(`/databases/${name}`)
    selectedDb.value = response.data
  } catch (err) {
    console.error(err)
    error.value = 'Failed to load connection info.'
  }
}

const closeConnection = () => {
  selectedDb.value = null
}

onMounted(fetchDatabases)
</script>

<template>
  <section class="databases-page">
    <header class="hero">
      <p class="eyebrow">Control Plane</p>
      <h1>My Databases</h1>
      <p class="subhead">
        Provision new clusters, inspect endpoints, and retire instances without touching YAML.
      </p>
    </header>

    <p v-if="error" class="status status--error" role="alert">{{ error }}</p>

    <div class="content-grid">
      <section class="panel">
        <div class="panel-head">
          <h2>Create New Database</h2>
        </div>

        <form class="create-form" @submit.prevent="createDatabase">
          <label class="field">
            <span>Database Name</span>
            <input v-model="newDbName" type="text" placeholder="Database Name" required />
          </label>

          <label class="field">
            <span>Instances</span>
            <input v-model.number="newDbInstances" type="number" min="1" placeholder="Instances" required />
          </label>

          <label class="field">
            <span>Storage</span>
            <input v-model="newDbStorage" type="text" placeholder="Storage (e.g., 1Gi)" required />
          </label>

          <button type="submit" class="primary-button" :disabled="creating">
            {{ creating ? 'Creating...' : 'Create' }}
          </button>
        </form>
      </section>

      <section class="panel">
        <div class="panel-head">
          <h2>Databases</h2>
          <span class="panel-meta">{{ databaseCount }}</span>
        </div>

        <p v-if="loading" class="status">Loading...</p>

        <ul v-else-if="databases.length > 0" class="db-list">
          <li v-for="db in databases" :key="db.name" class="db-item">
            <div class="db-main">
              <div>
                <h3>{{ db.name }}</h3>
                <p class="db-meta">{{ rowMeta(db) }}</p>
              </div>

              <div class="actions">
                <button
                  class="secondary-button"
                  :disabled="updatingName === db.name || deletingName === db.name"
                  @click="startEdit(db)"
                >
                  Edit
                </button>
                <button class="secondary-button" @click="viewConnection(db.name)">View Connection Info</button>
                <button
                  class="danger-button"
                  :disabled="deletingName === db.name || updatingName === db.name"
                  @click="deleteDatabase(db.name)"
                >
                  {{ deletingName === db.name ? 'Deleting...' : 'Delete' }}
                </button>
              </div>
            </div>

            <form v-if="isEditing(db.name)" class="edit-form" @submit.prevent="updateDatabase(db.name)">
              <label class="field">
                <span>Instances</span>
                <input v-model.number="editDbInstances" type="number" min="1" required />
              </label>
              <label class="field">
                <span>Storage</span>
                <input v-model="editDbStorage" type="text" placeholder="Storage (e.g., 2Gi)" required />
              </label>
              <div class="edit-actions">
                <button type="submit" class="primary-button" :disabled="updatingName === db.name">
                  {{ updatingName === db.name ? 'Saving...' : 'Save Changes' }}
                </button>
                <button type="button" class="secondary-button" :disabled="updatingName === db.name" @click="cancelEdit">
                  Cancel
                </button>
              </div>
            </form>
          </li>
        </ul>

        <div v-else class="status status--empty">
          <h3>No databases yet</h3>
          <p>Create your first PostgreSQL cluster using the form.</p>
        </div>
      </section>
    </div>

    <section v-if="selectedDb" class="panel connection-info">
      <div class="panel-head">
        <h2>Connection Info for {{ selectedDb.name }}</h2>
        <button class="secondary-button" @click="closeConnection">Close</button>
      </div>

      <div class="connection-grid">
        <div class="connection-row">
          <span class="key">Host:</span>
          <code>{{ selectedDb.connection?.host ?? '-' }}</code>
        </div>
        <div class="connection-row">
          <span class="key">Port:</span>
          <code>{{ selectedDb.connection?.port ?? '-' }}</code>
        </div>
        <div class="connection-row">
          <span class="key">Username:</span>
          <code>{{ selectedDb.connection?.username ?? '-' }}</code>
        </div>
        <div class="connection-row">
          <span class="key">Password:</span>
          <code>{{ selectedDb.connection?.password ?? '-' }}</code>
        </div>
        <div class="connection-row">
          <span class="key">Database:</span>
          <code>{{ selectedDb.connection?.database ?? '-' }}</code>
        </div>
      </div>
    </section>
  </section>
</template>

<style scoped>
.databases-page {
  display: grid;
  gap: var(--space-6);
}

.hero {
  max-width: 860px;
}

.eyebrow {
  font-size: 13px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--color-text-soft);
}

h1 {
  margin-top: var(--space-2);
  font-size: var(--text-h1);
  line-height: 1.08;
  letter-spacing: -0.03em;
  color: var(--color-text-strong);
}

.subhead {
  margin-top: var(--space-4);
  color: var(--color-text-soft);
  max-width: 62ch;
}

.content-grid {
  display: grid;
  gap: var(--space-5);
}

.panel {
  background: linear-gradient(180deg, #fff, #fefefe);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  box-shadow: var(--shadow-soft);
  padding: clamp(20px, 2vw, 28px);
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
}

h2 {
  font-size: var(--text-h2);
  line-height: 1.2;
  color: var(--color-text-strong);
  letter-spacing: -0.02em;
}

.panel-meta {
  font-size: 13px;
  color: var(--color-text-soft);
  background: var(--color-surface-muted);
  border: 1px solid var(--color-border);
  padding: 6px 10px;
  border-radius: var(--radius-pill);
}

.create-form {
  display: grid;
  gap: var(--space-4);
  margin-top: var(--space-5);
}

.field {
  display: grid;
  gap: var(--space-2);
}

.field span {
  font-size: var(--text-label);
  color: var(--color-text);
  font-weight: 560;
}

input {
  width: 100%;
  height: 44px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: var(--color-surface);
  color: var(--color-text-strong);
  padding: 0 var(--space-4);
  transition: border-color 180ms ease, box-shadow 180ms ease;
}

input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 4px rgb(15 131 253 / 14%);
}

.primary-button,
.secondary-button,
.danger-button {
  border-radius: var(--radius-pill);
  border: 1px solid transparent;
  min-height: 40px;
  padding: 0 14px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 180ms ease, border-color 180ms ease, transform 180ms ease;
}

.primary-button {
  color: #fff;
  background: var(--color-primary);
}

.primary-button:hover:enabled {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
}

.secondary-button {
  color: var(--color-text);
  background: var(--color-surface);
  border-color: var(--color-border);
}

.secondary-button:hover:enabled {
  border-color: var(--color-border-strong);
  background: #f8fafc;
}

.danger-button {
  color: var(--color-danger);
  background: #fff;
  border-color: #fecaca;
}

.danger-button:hover:enabled {
  background: var(--color-danger-soft);
}

.primary-button:disabled,
.secondary-button:disabled,
.danger-button:disabled {
  cursor: not-allowed;
  opacity: 0.65;
}

.db-list {
  margin-top: var(--space-5);
  display: grid;
  gap: var(--space-3);
  padding: 0;
  list-style: none;
}

.db-item {
  border: 1px solid var(--color-border);
  border-radius: 14px;
  background: var(--color-surface-muted);
  padding: var(--space-4);
  display: grid;
  gap: var(--space-4);
}

.db-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
}

h3 {
  font-size: 18px;
  color: var(--color-text-strong);
  letter-spacing: -0.01em;
}

.db-meta {
  margin-top: 2px;
  color: var(--color-text-soft);
  font-size: 14px;
}

.actions {
  display: inline-flex;
  gap: var(--space-2);
  flex-wrap: wrap;
  justify-content: flex-end;
}

.edit-form {
  display: grid;
  gap: var(--space-3);
  border-top: 1px dashed var(--color-border-strong);
  padding-top: var(--space-4);
}

.edit-actions {
  display: inline-flex;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.status {
  margin-top: var(--space-5);
  color: var(--color-text-soft);
}

.status--error {
  margin-top: 0;
  color: var(--color-danger);
  background: var(--color-danger-soft);
  border: 1px solid #fecaca;
  border-radius: 12px;
  padding: var(--space-3) var(--space-4);
}

.status--empty {
  margin-top: var(--space-5);
  border: 1px dashed var(--color-border-strong);
  border-radius: 12px;
  padding: var(--space-5);
  background: var(--color-surface-muted);
}

.status--empty p {
  margin-top: var(--space-2);
}

.connection-info {
  display: grid;
  gap: var(--space-5);
}

.connection-grid {
  display: grid;
  gap: var(--space-3);
}

.connection-row {
  display: grid;
  gap: var(--space-1);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  padding: var(--space-3) var(--space-4);
  background: var(--color-surface-muted);
}

.key {
  font-size: var(--text-label);
  font-weight: 600;
  color: var(--color-text-soft);
}

code {
  color: var(--color-text-strong);
  font-size: 14px;
  overflow-wrap: anywhere;
}

@media (min-width: 980px) {
  .content-grid {
    grid-template-columns: minmax(320px, 380px) minmax(0, 1fr);
    align-items: start;
  }
}

@media (max-width: 720px) {
  .db-main {
    flex-direction: column;
    align-items: flex-start;
  }

  .actions {
    width: 100%;
  }
}
</style>
