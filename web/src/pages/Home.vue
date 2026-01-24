<script setup lang="ts">
defineProps<{
  message: string
}>()
</script>

<template>
  <div>
    <h1>Gooji</h1>
    <p>{{ message }}</p>
  </div>
</template>
