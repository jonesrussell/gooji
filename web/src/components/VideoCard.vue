<script setup lang="ts">
import type { Video } from '@/types'

defineProps<{
  video: Video
}>()

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}
</script>

<template>
  <article class="group relative bg-white rounded-lg shadow-sm overflow-hidden hover:shadow-md transition-shadow">
    <div class="aspect-video bg-gray-200 relative">
      <img
        v-if="video.thumbnailPath"
        :src="video.thumbnailPath"
        :alt="video.title"
        class="w-full h-full object-cover"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
        <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
      </div>
      <div class="absolute bottom-2 right-2 bg-black/75 text-white text-xs px-1.5 py-0.5 rounded">
        {{ formatDuration(video.durationSeconds) }}
      </div>
    </div>
    <div class="p-4">
      <h3 class="font-medium text-gray-900 group-hover:text-primary-600 transition-colors line-clamp-1">
        {{ video.title }}
      </h3>
      <p v-if="video.description" class="mt-1 text-sm text-gray-500 line-clamp-2">
        {{ video.description }}
      </p>
    </div>
  </article>
</template>
