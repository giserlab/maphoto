<template>
  <div class="auth-layout">
    <div class="auth-container">
      <router-view />
    </div>
  </div>
</template>

<style scoped>
.auth-layout {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(180deg, #76abe1 0%, #3882cd 100%);
}

.auth-container {
  width: 100%;
  max-width: 420px;
  padding: 20px;
}
</style>
