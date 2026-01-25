<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Link } from '@inertiajs/vue3'
import Layout from '../components/Layout.vue'

interface Video {
  id: string
  title: string
  description: string
  tags: string[]
  duration?: number
}

// Props from Inertia
defineProps<{
  message?: string
  featuredVideos?: Video[]
}>()

// Recent videos
const recentVideos = ref<Video[]>([])
const isLoading = ref(true)

// Format duration
function formatDuration(seconds: number): string {
  if (!seconds) return ''
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

// Load recent videos
async function loadRecentVideos(): Promise<void> {
  try {
    const response = await fetch('/api/videos?limit=3')
    if (response.ok) {
      recentVideos.value = await response.json()
    }
  } catch (err) {
    console.error('Failed to load recent videos:', err)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadRecentVideos()
})
</script>

<template>
  <Layout>
    <div class="space-y-16">
      <!-- Hero Section -->
      <section class="text-center">
        <h1 class="text-4xl font-bold tracking-tight text-foreground sm:text-5xl md:text-6xl">
          <span class="block">Gooji</span>
          <span class="mt-2 block text-2xl font-medium text-foreground/70 sm:text-3xl">
            Preserving Ojibwe Language & Culture
          </span>
        </h1>
        <p class="mx-auto mt-6 max-w-2xl text-lg text-foreground/70">
          A community platform for recording, sharing, and preserving Ojibwe language lessons, cultural stories, and traditional knowledge through video.
        </p>
        <div class="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
          <Link
            href="/record"
            class="inline-flex items-center gap-2 rounded-lg bg-accent px-8 py-4 text-lg font-medium text-white shadow-lg transition-all hover:bg-accent/90 hover:shadow-xl"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <circle cx="12" cy="12" r="10" stroke-width="2" />
              <circle cx="12" cy="12" r="4" fill="currentColor" />
            </svg>
            Start Recording
          </Link>
          <Link
            href="/gallery"
            class="inline-flex items-center gap-2 rounded-lg border-2 border-foreground/20 px-8 py-4 text-lg font-medium text-foreground transition-colors hover:border-accent hover:text-accent"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            Browse Videos
          </Link>
        </div>
      </section>

      <!-- Features Section -->
      <section class="grid gap-8 sm:grid-cols-3">
        <!-- Record -->
        <div class="rounded-2xl border border-foreground/10 bg-snow-white p-8 shadow-lg transition-all hover:shadow-xl">
          <div class="flex h-14 w-14 items-center justify-center rounded-xl bg-accent/10">
            <svg class="h-8 w-8 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <circle cx="12" cy="12" r="10" stroke-width="2" />
              <circle cx="12" cy="12" r="4" fill="currentColor" />
            </svg>
          </div>
          <h3 class="mt-6 text-xl font-bold text-foreground">Record</h3>
          <p class="mt-2 text-foreground/70">
            Capture language lessons, stories, and cultural knowledge directly from your browser with our easy-to-use recorder.
          </p>
          <Link href="/record" class="mt-4 inline-flex items-center text-accent hover:underline">
            Start recording
            <svg class="ml-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </Link>
        </div>

        <!-- Upload -->
        <div class="rounded-2xl border border-foreground/10 bg-snow-white p-8 shadow-lg transition-all hover:shadow-xl">
          <div class="flex h-14 w-14 items-center justify-center rounded-xl bg-accent/10">
            <svg class="h-8 w-8 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
          </div>
          <h3 class="mt-6 text-xl font-bold text-foreground">Upload</h3>
          <p class="mt-2 text-foreground/70">
            Share existing video content by uploading files. Support for MP4, WebM, MOV, and other common formats.
          </p>
          <Link href="/upload" class="mt-4 inline-flex items-center text-accent hover:underline">
            Upload video
            <svg class="ml-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </Link>
        </div>

        <!-- Browse -->
        <div class="rounded-2xl border border-foreground/10 bg-snow-white p-8 shadow-lg transition-all hover:shadow-xl">
          <div class="flex h-14 w-14 items-center justify-center rounded-xl bg-accent/10">
            <svg class="h-8 w-8 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </div>
          <h3 class="mt-6 text-xl font-bold text-foreground">Browse</h3>
          <p class="mt-2 text-foreground/70">
            Explore the video library to learn Ojibwe language, discover cultural stories, and connect with traditions.
          </p>
          <Link href="/gallery" class="mt-4 inline-flex items-center text-accent hover:underline">
            View gallery
            <svg class="ml-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </Link>
        </div>
      </section>

      <!-- Recent Videos Section -->
      <section>
        <div class="flex items-center justify-between">
          <h2 class="text-2xl font-bold text-foreground">Recent Videos</h2>
          <Link href="/gallery" class="text-accent hover:underline">
            View all
          </Link>
        </div>

        <!-- Loading -->
        <div v-if="isLoading" class="mt-6 flex justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-4 border-accent border-t-transparent" />
        </div>

        <!-- Videos Grid -->
        <div v-else-if="recentVideos.length > 0" class="mt-6 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <Link
            v-for="video in recentVideos"
            :key="video.id"
            :href="`/video/${video.id}`"
            class="group overflow-hidden rounded-2xl border border-foreground/10 bg-snow-white shadow-lg transition-all hover:-translate-y-1 hover:shadow-xl"
          >
            <!-- Thumbnail -->
            <div class="relative aspect-video overflow-hidden">
              <img
                :src="`/api/thumbnails?id=${video.id}`"
                :alt="video.title"
                class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
              />
              <div class="absolute inset-0 flex items-center justify-center bg-black/20 opacity-0 transition-opacity group-hover:opacity-100">
                <div class="flex h-14 w-14 items-center justify-center rounded-full bg-white/90">
                  <svg class="ml-1 h-6 w-6 text-accent" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M8 5v14l11-7z" />
                  </svg>
                </div>
              </div>
              <div
                v-if="video.duration"
                class="absolute bottom-2 right-2 rounded bg-black/70 px-2 py-1 text-xs text-white"
              >
                {{ formatDuration(video.duration) }}
              </div>
            </div>

            <!-- Info -->
            <div class="p-4">
              <h3 class="font-bold text-foreground group-hover:text-accent">{{ video.title }}</h3>
              <p class="mt-1 line-clamp-2 text-sm text-foreground/70">{{ video.description }}</p>
            </div>
          </Link>
        </div>

        <!-- Empty State -->
        <div v-else class="mt-6 rounded-2xl border border-dashed border-foreground/20 p-12 text-center">
          <svg class="mx-auto h-12 w-12 text-foreground/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <h3 class="mt-4 font-medium text-foreground">No videos yet</h3>
          <p class="mt-2 text-foreground/70">Be the first to contribute!</p>
          <Link
            href="/record"
            class="mt-4 inline-flex items-center gap-2 rounded-lg bg-accent px-6 py-2 font-medium text-white hover:bg-accent/90"
          >
            Record a Video
          </Link>
        </div>
      </section>

      <!-- About Section -->
      <section class="rounded-2xl bg-earth-black p-8 text-snow-white sm:p-12">
        <div class="mx-auto max-w-3xl text-center">
          <h2 class="text-2xl font-bold sm:text-3xl">About Gooji</h2>
          <p class="mt-4 text-snow-white/80">
            "Gooji" means "video" in Ojibwe. This platform was created to support the preservation and revitalization of Ojibwe language and culture through accessible video content.
          </p>
          <p class="mt-4 text-snow-white/80">
            Whether you're a fluent speaker sharing knowledge, a learner practicing pronunciation, or an elder passing down stories, Gooji provides a simple way to contribute to this living archive.
          </p>
          <div class="mt-8 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
            <Link
              href="/record"
              class="inline-flex items-center gap-2 rounded-lg bg-accent px-6 py-3 font-medium text-white hover:bg-accent/90"
            >
              Contribute Now
            </Link>
            <Link
              href="/gallery"
              class="inline-flex items-center gap-2 rounded-lg border border-snow-white/30 px-6 py-3 font-medium text-snow-white hover:bg-snow-white/10"
            >
              Explore Content
            </Link>
          </div>
        </div>
      </section>
    </div>
  </Layout>
</template>
