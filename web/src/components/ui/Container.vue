<script setup lang="ts">
defineProps<{
  size?: 'sm' | 'md' | 'lg' | 'xl'
}>()
</script>

<template>
  <div
    class="mx-auto px-4 sm:px-6 lg:px-8"
    :class="{
      'max-w-2xl': size === 'sm',
      'max-w-4xl': size === 'md',
      'max-w-6xl': size === 'lg' || !size,
      'max-w-7xl': size === 'xl',
    }"
  >
    <slot />
  </div>
</template>
