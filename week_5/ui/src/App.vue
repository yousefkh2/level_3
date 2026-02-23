<script setup>
import { onMounted, ref } from 'vue'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'

const token = ref(localStorage.getItem('token') || '')
const loginUsername = ref('')
const loginPassword = ref('')
const loginBusy = ref(false)
const loginError = ref('')

const databases = ref([])
const loadingDatabases = ref(false)
const busyAction = ref('')
const errorMessage = ref('')

const newName = ref('')
const newInstances = ref(1)
const newStorage = ref('1Gi')

const editingName = ref('')
const editInstances = ref(1)
const editStorage = ref('1Gi')

const selectedDb = ref(null)
const copyMessage = ref('')

function setToken(value) {
  token.value = value

  if (value) {
    localStorage.setItem('token', value)
    return
  }

  localStorage.removeItem('token')
}

async function request(path, options = {}) {
  const config = {
    method: options.method || 'GET',
    headers: {},
  }

  if (options.body !== undefined) {
    config.headers['Content-Type'] = 'application/json'
    config.body = JSON.stringify(options.body)
  }

  if (token.value) {
    config.headers.Authorization = `Bearer ${token.value}`
  }

  const response = await fetch(`${API_BASE_URL}${path}`, config)

  if (response.status === 204) {
    return null
  }

  let data = null

  try {
    data = await response.json()
  } catch {
    data = null
  }

  if (!response.ok) {
    const message = data?.error || data?.message || `Request failed (${response.status})`
    throw new Error(message)
  }

  return data
}

async function login() {
  loginBusy.value = true
  loginError.value = ''

  try {
    const data = await request('/auth/login', {
      method: 'POST',
      body: {
        username: loginUsername.value,
        password: loginPassword.value,
      },
    })

    setToken(data.token)
    loginUsername.value = ''
    loginPassword.value = ''
    await loadDatabases()
  } catch (error) {
    loginError.value = error.message
  } finally {
    loginBusy.value = false
  }
}

function logout() {
  setToken('')
  databases.value = []
  selectedDb.value = null
  editingName.value = ''
  errorMessage.value = ''
}

async function loadDatabases() {
  loadingDatabases.value = true
  errorMessage.value = ''

  try {
    const data = await request('/databases')
    databases.value = Array.isArray(data) ? data : []
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loadingDatabases.value = false
  }
}

async function createDatabase() {
  busyAction.value = 'create'
  errorMessage.value = ''

  try {
    await request('/databases', {
      method: 'POST',
      body: {
        name: newName.value,
        instances: Number(newInstances.value),
        storage: newStorage.value,
      },
    })

    newName.value = ''
    newInstances.value = 1
    newStorage.value = '1Gi'
    await loadDatabases()
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    busyAction.value = ''
  }
}

function startEdit(db) {
  editingName.value = db.name
  editInstances.value = db?.spec?.instances ?? 1
  editStorage.value = db?.spec?.storage ?? '1Gi'
}

function cancelEdit() {
  editingName.value = ''
  editInstances.value = 1
  editStorage.value = '1Gi'
}

async function saveEdit(name) {
  busyAction.value = `edit:${name}`
  errorMessage.value = ''

  try {
    await request(`/databases/${name}`, {
      method: 'PATCH',
      body: {
        instances: Number(editInstances.value),
        storage: editStorage.value,
      },
    })

    cancelEdit()
    await loadDatabases()

    if (selectedDb.value?.name === name) {
      await showConnection(name)
    }
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    busyAction.value = ''
  }
}

async function deleteDatabase(name) {
  busyAction.value = `delete:${name}`
  errorMessage.value = ''

  try {
    await request(`/databases/${name}`, { method: 'DELETE' })

    if (selectedDb.value?.name === name) {
      selectedDb.value = null
    }

    await loadDatabases()
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    busyAction.value = ''
  }
}

async function showConnection(name) {
  busyAction.value = `connection:${name}`
  errorMessage.value = ''

  try {
    selectedDb.value = await request(`/databases/${name}`)
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    busyAction.value = ''
  }
}

function connectionValue(key) {
  return selectedDb.value?.connection?.[key] || ''
}

function rwHost() {
  return connectionValue('host')
}

function roHost() {
  const host = rwHost()
  if (!host) {
    return ''
  }
  if (host.includes('-rw')) {
    return host.replace('-rw', '-ro')
  }
  return `${selectedDb.value?.name || ''}-ro`
}

function postgresURL(host) {
  const username = encodeURIComponent(connectionValue('username'))
  const password = encodeURIComponent(connectionValue('password'))
  const port = connectionValue('port')
  const database = encodeURIComponent(connectionValue('database'))

  if (!host || !username || !port || !database) {
    return ''
  }
  return `postgresql://${username}:${password}@${host}:${port}/${database}`
}

function psqlCommand(host) {
  const username = connectionValue('username')
  const password = connectionValue('password')
  const port = connectionValue('port')
  const database = connectionValue('database')

  if (!host || !username || !port || !database) {
    return ''
  }
  return `psql "host=${host} port=${port} user=${username} password=${password} dbname=${database} sslmode=disable"`
}

async function copyText(value, label) {
  if (!value) {
    return
  }
  try {
    await navigator.clipboard.writeText(value)
    copyMessage.value = `${label} copied`
    setTimeout(() => {
      copyMessage.value = ''
    }, 1500)
  } catch {
    copyMessage.value = `Failed to copy ${label.toLowerCase()}`
  }
}

onMounted(async () => {
  if (token.value) {
    await loadDatabases()
  }
})
</script>

<template>
  <main class="page">
    <h1>PaaS Control Plane</h1>

    <section v-if="!token" class="panel">
      <h2>Login</h2>
      <form class="form" @submit.prevent="login">
        <label>
          Username
          <input v-model="loginUsername" placeholder="Username" required />
        </label>

        <label>
          Password
          <input v-model="loginPassword" type="password" placeholder="Password" required />
        </label>

        <button type="submit" :disabled="loginBusy">
          {{ loginBusy ? 'Signing in...' : 'Login' }}
        </button>
      </form>

      <p v-if="loginError" class="error">{{ loginError }}</p>
    </section>

    <section v-else class="panel">
      <div class="header-row">
        <h2>My Databases</h2>
        <div class="buttons-inline">
          <button type="button" @click="loadDatabases" :disabled="loadingDatabases">
            {{ loadingDatabases ? 'Refreshing...' : 'Refresh' }}
          </button>
          <button type="button" @click="logout">Logout</button>
        </div>
      </div>

      <p v-if="errorMessage" class="error">{{ errorMessage }}</p>

      <form class="form create-form" @submit.prevent="createDatabase">
        <h3>Create</h3>
        <label>
          Name
          <input v-model="newName" placeholder="Database Name" required />
        </label>

        <label>
          Instances
          <input v-model.number="newInstances" type="number" min="1" placeholder="Instances" required />
        </label>

        <label>
          Storage
          <input v-model="newStorage" placeholder="Storage (e.g., 1Gi)" required />
        </label>

        <button type="submit" :disabled="busyAction === 'create'">
          {{ busyAction === 'create' ? 'Creating...' : 'Create' }}
        </button>
      </form>

      <h3>List</h3>
      <p v-if="loadingDatabases">Loading...</p>
      <p v-else-if="databases.length === 0">No databases yet.</p>

      <ul v-else class="db-list">
        <li v-for="db in databases" :key="db.name" class="db-item">
          <div>
            <strong>{{ db.name }}</strong>
            <div class="meta">
              status: {{ db.status || 'unknown' }} | instances: {{ db.spec?.instances || 'n/a' }} | storage:
              {{ db.spec?.storage || 'n/a' }}
            </div>
          </div>

          <div class="buttons-inline">
            <button type="button" @click="startEdit(db)">Edit</button>
            <button type="button" @click="showConnection(db.name)" :disabled="busyAction === `connection:${db.name}`">
              View Connection Info
            </button>
            <button type="button" @click="deleteDatabase(db.name)" :disabled="busyAction === `delete:${db.name}`">
              {{ busyAction === `delete:${db.name}` ? 'Deleting...' : 'Delete' }}
            </button>
          </div>

          <form v-if="editingName === db.name" class="form edit-form" @submit.prevent="saveEdit(db.name)">
            <label>
              Instances
              <input v-model.number="editInstances" type="number" min="1" required />
            </label>

            <label>
              Storage
              <input v-model="editStorage" placeholder="Storage (e.g., 2Gi)" required />
            </label>

            <div class="buttons-inline">
              <button type="submit" :disabled="busyAction === `edit:${db.name}`">
                {{ busyAction === `edit:${db.name}` ? 'Saving...' : 'Save Changes' }}
              </button>
              <button type="button" @click="cancelEdit">Cancel</button>
            </div>
          </form>
        </li>
      </ul>

      <section v-if="selectedDb" class="connection-info">
        <div class="header-row">
          <h3>Connection Info for {{ selectedDb.name }}</h3>
          <button type="button" @click="selectedDb = null">Close</button>
        </div>
        <p v-if="copyMessage" class="copy-message">{{ copyMessage }}</p>
        <p>Port: {{ selectedDb.connection?.port || '-' }}</p>
        <p>Username: {{ selectedDb.connection?.username || '-' }}</p>
        <p>Password: {{ selectedDb.connection?.password || '-' }}</p>
        <p>Database: {{ selectedDb.connection?.database || '-' }}</p>

        <div class="connection-block">
          <h4>Service Endpoints</h4>
          <div class="connection-line">
            <span><strong>Read/Write:</strong> {{ rwHost() || '-' }}</span>
            <button type="button" @click="copyText(rwHost(), 'Read/write host')">Copy</button>
          </div>
          <div class="connection-line">
            <span><strong>Read-Only:</strong> {{ roHost() || '-' }}</span>
            <button type="button" @click="copyText(roHost(), 'Read-only host')">Copy</button>
          </div>
        </div>

        <div class="connection-block">
          <h4>Connection Strings</h4>
          <div class="connection-line">
            <code>{{ postgresURL(rwHost()) || '-' }}</code>
            <button type="button" @click="copyText(postgresURL(rwHost()), 'PostgreSQL URL')">Copy</button>
          </div>
          <div class="connection-line">
            <code>{{ psqlCommand(rwHost()) || '-' }}</code>
            <button type="button" @click="copyText(psqlCommand(rwHost()), 'psql command')">Copy</button>
          </div>
        </div>

        <p class="note">
          SSL note: these endpoints are internal cluster services. In-cluster clients can use
          <code>sslmode=disable</code>.
        </p>
        <p class="warning">
          Internal-only access: these hosts are only reachable from inside the Kubernetes cluster/VPN.
        </p>
      </section>
    </section>
  </main>
</template>

<style scoped>
.page {
  max-width: 900px;
  margin: 0 auto;
  padding: 24px;
}

h1,
h2,
h3,
p {
  margin: 0;
}

h1 {
  margin-bottom: 16px;
}

.panel {
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  padding: 16px;
  display: grid;
  gap: 16px;
}

.form {
  display: grid;
  gap: 12px;
}

.form label {
  display: grid;
  gap: 6px;
  font-size: 14px;
}

input {
  border: 1px solid #9ca3af;
  border-radius: 6px;
  padding: 8px;
}

button {
  border: 1px solid #9ca3af;
  background: #f9fafb;
  border-radius: 6px;
  padding: 8px 12px;
  cursor: pointer;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.buttons-inline {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.error {
  color: #b91c1c;
}

.create-form {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 12px;
}

.db-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 12px;
}

.db-item {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 12px;
  display: grid;
  gap: 10px;
}

.meta {
  margin-top: 4px;
  color: #4b5563;
  font-size: 14px;
}

.edit-form {
  border-top: 1px solid #e5e7eb;
  padding-top: 12px;
}

.connection-info {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 12px;
  display: grid;
  gap: 8px;
}

.connection-block {
  border-top: 1px solid #e5e7eb;
  padding-top: 8px;
  display: grid;
  gap: 8px;
}

.connection-line {
  display: flex;
  gap: 8px;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
}

code {
  background: #f3f4f6;
  border-radius: 6px;
  padding: 6px 8px;
  word-break: break-all;
}

.copy-message {
  color: #065f46;
  font-size: 14px;
}

.note {
  font-size: 14px;
  color: #374151;
}

.warning {
  font-size: 14px;
  color: #92400e;
}
</style>
