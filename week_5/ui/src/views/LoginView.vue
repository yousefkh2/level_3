<template>
    <div class="login-container">
        <h1>Login</h1>
        <form @submit.prevent="login"> <!--prevent full page reload on form submission-->
            <input v-model="username" type="text" placeholder="Username" />
            <input v-model="password" type="password" placeholder="Password" />
            <button type="submit">Login</button>
        </form>
    </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiClient, setToken } from '@/services/api'

export default {
  setup() {
    const username = ref('')
    const password = ref('')
    const router = useRouter()
    const login = async () => {
      try {
        const response = await apiClient.post('/auth/login', {
          username: username.value,
          password: password.value
        })
        setToken(response.data.token)
        router.push('/databases')
      } catch (error) {
        console.error('Login failed:', error)
      }
    }

    return { username, password, login }
  }
}
</script>

<style scoped> 
.login-container {
  max-width: 400px;
  margin: 50px auto;
}
</style>