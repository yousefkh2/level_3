<template>
  <div class="databases-container">
    <h1>My Databases</h1>
    
    <!-- Create new database form -->
    <div class="create-form">
      <h2>Create New Database</h2>
      <form @submit.prevent="createDatabase">
        <input v-model="newDbName" type="text" placeholder="Database Name" required />
        <input v-model.number="newDbInstances" type="number" placeholder="Instances" required />
        <input v-model="newDbStorage" type="text" placeholder="Storage (e.g., 1Gi)" required />
        <button type="submit">Create</button>
      </form>
    </div>

    <!-- Loading state -->
    <p v-if="loading">Loading...</p>
    
    <!-- Error message -->
    <p v-if="error" class="error">{{ error }}</p>
    
    <!-- List of databases -->
    <div v-if="!loading && databases.length > 0" class="databases-list">
      <h2>Databases</h2>
      <ul>
        <li v-for="db in databases" :key="db.id">
          <span>{{ db.name }}</span>
          <button @click="viewConnection(db.name)">View Connection Info</button>
          <button @click="deleteDatabase(db.name)">Delete</button>
        </li>
      </ul>
    </div>
    
    <!-- No databases message -->
    <p v-if="!loading && databases.length === 0">No databases yet</p>

    <!-- Connection info section -->
    <div v-if="selectedDb" class="connection-info">
      <h2>Connection Info for {{ selectedDb.name }}</h2>
      <p><strong>Host:</strong> {{ selectedDb.connection.host }}</p>
      <p><strong>Port:</strong> {{ selectedDb.connection.port }}</p>
      <p><strong>Username:</strong> {{ selectedDb.connection.username }}</p>
      <p><strong>Password:</strong> {{ selectedDb.connection.password }}</p>
      <p><strong>Database:</strong> {{ selectedDb.connection.database }}</p>
      <button @click="selectedDb = null">Close</button>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { apiClient } from '@/services/api'

export default {
  setup() {
    // state for displaying databases
    const databases = ref([])
    const loading = ref(false)
    const error = ref(null)

    // state for create form
    const newDbName = ref('')
    const newDbInstances = ref(1)
    const newDbStorage = ref('1Gi')
    const selectedDb = ref(null) // for showing connection info

    // Fetch databases from API
    const fetchDatabases = async () => {
      loading.value = true
      error.value = null
      try {
        const response = await apiClient.get('/databases')
        databases.value = response.data
      } catch (err) {
        error.value = 'Failed to load databases'
        console.error(err)
      } finally {
        loading.value = false
      }
    }

    const createDatabase = async () => {
        try {
            await apiClient.post('/databases', {
                name: newDbName.value,
                instances: newDbInstances.value,
                storage: newDbStorage.value
            })
            // reset form
            newDbName.value = ''
            newDbInstances.value = 1
            newDbStorage.value = '1Gi'
            // refresh list
            await fetchDatabases()
        } catch (err) {
            error.value = 'Failed to create database'
            console.error(err)
        }
    }

    const deleteDatabase = async (name) => {
        try {
            await apiClient.delete(`/databases/${name}`)
            //refresh list
            await fetchDatabases()
            // clear selected db if it was deleted
            if (selectedDb.value?.name === name) {
              selectedDb.value = null
            }
        } catch(err) {
            error.value = 'Failed to delete database'
            console.error(err)
        }
    }

    // Fetch connection info for a specific database
    const viewConnection = async (name) => {
      try {
        const response = await apiClient.get(`/databases/${name}`)
        selectedDb.value = response.data
      } catch (err) {
        error.value = 'Failed to load connection info'
        console.error(err)
      }
    }

    // Call fetchDatabases when component loads
    onMounted(() => {
      fetchDatabases()
    })

    return { databases, loading, error, newDbName, newDbInstances, newDbStorage, createDatabase, deleteDatabase, selectedDb, viewConnection }
  }
}
</script>