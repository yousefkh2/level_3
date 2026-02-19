<script setup>
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { ref, onBeforeUnmount, onMounted } from 'vue'
import { getToken, clearToken } from '@/services/api'

const router = useRouter()
const isLoggedIn = ref(false)
let removeAfterEachHook = null

const checkAuth = () => {
  isLoggedIn.value = !!getToken()
}

const logout = () => {
  clearToken()
  isLoggedIn.value = false
  router.push('/login')
}

onMounted(() => {
  checkAuth()
  removeAfterEachHook = router.afterEach(() => {
    checkAuth()
  })
})

onBeforeUnmount(() => {
  if (typeof removeAfterEachHook === 'function') {
    removeAfterEachHook()
  }
})
</script>

<template>
  <div class="app-shell">
    <header class="top-bar">
      <div class="container top-bar__inner">
        <RouterLink class="brand" to="/databases">
          <span class="brand__dot" />
          <span class="brand__name">PaaS Control Plane</span>
        </RouterLink>

        <nav class="top-nav" aria-label="Primary">
          <RouterLink v-if="!isLoggedIn" class="nav-link" to="/login">Login</RouterLink>
          <RouterLink v-if="isLoggedIn" class="nav-link" to="/databases">Databases</RouterLink>
          <button v-if="isLoggedIn" type="button" class="ghost-button" @click="logout">Logout</button>
        </nav>
      </div>
    </header>

    <main class="main-content">
      <div class="container">
        <RouterView v-slot="{ Component }">
          <Transition name="page" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </div>
    </main>
  </div>
</template>

<style scoped>
.app-shell {
  min-height: 100vh;
}

.top-bar {
  position: sticky;
  top: 0;
  z-index: 20;
  backdrop-filter: blur(8px);
  background: rgb(243 245 248 / 80%);
  border-bottom: 1px solid var(--color-border);
}

.top-bar__inner {
  min-height: 72px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-4);
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: var(--space-3);
  color: var(--color-text-strong);
}

.brand__dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: linear-gradient(145deg, #33a0ff, var(--color-primary));
  box-shadow: 0 0 0 6px rgb(15 131 253 / 14%);
}

.brand__name {
  font-size: 15px;
  font-weight: 620;
  letter-spacing: -0.02em;
}

.top-nav {
  display: inline-flex;
  gap: var(--space-2);
  align-items: center;
}

.nav-link,
.ghost-button {
  border: 1px solid transparent;
  border-radius: var(--radius-pill);
  font-size: var(--text-label);
  font-weight: 560;
  color: var(--color-text);
  padding: 8px 14px;
  background: transparent;
  transition: all 180ms ease;
}

.nav-link:hover,
.ghost-button:hover {
  border-color: var(--color-border);
  background: var(--color-surface);
}

.nav-link.router-link-exact-active {
  color: var(--color-text-strong);
  background: var(--color-surface);
  border-color: var(--color-border-strong);
}

.ghost-button {
  cursor: pointer;
}

.main-content {
  padding-block: var(--space-7) var(--space-8);
}

@media (max-width: 640px) {
  .top-bar__inner {
    min-height: 64px;
  }

  .brand__name {
    font-size: 14px;
  }
}
</style>
