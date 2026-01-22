<script setup lang="ts">
import type { Video } from '@/types'
import Container from '@/components/ui/Container.vue'
import VideoCard from '@/components/VideoCard.vue'

defineProps<{
  videos: Video[]
}>()
</script>

<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <Container>
      <header class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Video Gallery</h1>
        <p class="mt-2 text-gray-600">
          Browse videos shared by the community.
        </p>
      </header>

      <!-- Empty state -->
      <div
        v-if="videos.length === 0"
        class="text-center py-16 bg-white rounded-lg border-2 border-dashed border-gray-300"
      >
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
        <h3 class="mt-4 text-lg font-medium text-gray-900">No videos yet</h3>
        <p class="mt-2 text-gray-500">Be the first to share a video with the community.</p>
      </div>

      <!-- Video grid -->
      <div
        v-else
        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
      >
        <VideoCard
          v-for="video in videos"
          :key="video.id"
          :video="video"
        />
      </div>
    </Container>
  </div>
</template>
