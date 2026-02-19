<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiClient, setToken } from '@/services/api'

const router = useRouter()
const username = ref('')
const password = ref('')
const submitting = ref(false)
const error = ref('')

const login = async () => {
  submitting.value = true
  error.value = ''

  try {
    const response = await apiClient.post('/auth/login', {
      username: username.value,
      password: password.value,
    })
    setToken(response.data.token)
    router.push('/databases')
  } catch (err) {
    console.error('Login failed:', err)
    error.value = 'Login failed. Check your username and password.'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <section class="login-page">
    <div class="login-intro">
      <p class="eyebrow">Database Platform</p>
      <h1>Operate PostgreSQL with less friction.</h1>
      <p class="subhead">
        Sign in to provision, inspect, and remove clusters from one clean control plane.
      </p>
    </div>

    <div class="login-card">
      <h2>Login</h2>
      <p class="card-copy">Use the API credentials configured for the control plane.</p>

      <form class="login-form" @submit.prevent="login">
        <label class="field">
          <span>Username</span>
          <input v-model="username" type="text" placeholder="Username" autocomplete="username" required />
        </label>

        <label class="field">
          <span>Password</span>
          <input
            v-model="password"
            type="password"
            placeholder="Password"
            autocomplete="current-password"
            required
          />
        </label>

        <button type="submit" :disabled="submitting">
          {{ submitting ? 'Signing in...' : 'Login' }}
        </button>
      </form>

      <p v-if="error" class="error-text" role="alert">{{ error }}</p>
    </div>
  </section>
</template>

<style scoped>
.login-page {
  display: grid;
  gap: var(--space-6);
  padding-top: var(--space-4);
  align-items: start;
}

.login-intro {
  max-width: 720px;
}

.eyebrow {
  font-size: 13px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--color-text-soft);
  margin-bottom: var(--space-3);
}

h1 {
  font-size: var(--text-hero);
  line-height: 1.05;
  letter-spacing: -0.03em;
  color: var(--color-text-strong);
  max-width: 14ch;
}

.subhead {
  margin-top: var(--space-5);
  max-width: 56ch;
  color: var(--color-text-soft);
}

.login-card {
  background: linear-gradient(180deg, #fff, #fefefe);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  box-shadow: var(--shadow-soft);
  max-width: 520px;
  padding: clamp(24px, 3vw, 36px);
}

h2 {
  font-size: var(--text-h2);
  letter-spacing: -0.02em;
  color: var(--color-text-strong);
}

.card-copy {
  margin-top: var(--space-2);
  color: var(--color-text-soft);
}

.login-form {
  margin-top: var(--space-6);
  display: grid;
  gap: var(--space-4);
}

.field {
  display: grid;
  gap: var(--space-2);
  color: var(--color-text);
}

.field span {
  font-size: var(--text-label);
  font-weight: 560;
}

input {
  width: 100%;
  height: 44px;
  border-radius: 12px;
  border: 1px solid var(--color-border);
  padding: 0 var(--space-4);
  background: var(--color-surface);
  color: var(--color-text-strong);
  transition: border-color 180ms ease, box-shadow 180ms ease;
}

input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 4px rgb(15 131 253 / 14%);
}

button {
  margin-top: var(--space-2);
  height: 46px;
  border: none;
  border-radius: 999px;
  font-weight: 620;
  color: #fff;
  background: var(--color-primary);
  cursor: pointer;
  transition: transform 160ms ease, background-color 160ms ease;
}

button:hover:enabled {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
}

button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.error-text {
  margin-top: var(--space-4);
  color: var(--color-danger);
  background: var(--color-danger-soft);
  border: 1px solid #fecaca;
  padding: var(--space-3) var(--space-4);
  border-radius: 12px;
}

@media (min-width: 940px) {
  .login-page {
    grid-template-columns: 1.4fr 1fr;
    gap: var(--space-8);
    min-height: calc(100vh - 170px);
    align-items: center;
  }
}
</style>
