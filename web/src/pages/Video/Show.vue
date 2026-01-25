<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { Link } from '@inertiajs/vue3'
import Layout from '../../components/Layout.vue'

interface Video {
  id: string
  title: string
  description: string
  tags: string[]
  createdAt?: string
  duration?: number
}

// Props from Inertia (server passes video data)
const props = defineProps<{
  video?: Video
  id?: string
}>()

// State
const videoData = ref<Video | null>(props.video || null)
const isLoading = ref(!props.video)
const error = ref<string | null>(null)
const videoRef = ref<HTMLVideoElement | null>(null)
const isPlaying = ref(false)
const currentTime = ref(0)
const duration = ref(0)
const volume = ref(1)
const isMuted = ref(false)
const isFullscreen = ref(false)

// Computed video ID (from props or URL)
const videoId = computed(() => props.video?.id || props.id || '')

// Format time for display
function formatTime(seconds: number): string {
  if (!seconds || isNaN(seconds)) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

// Format date
function formatDate(dateString: string): string {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

// Progress percentage
const progressPercent = computed(() => {
  if (!duration.value) return 0
  return (currentTime.value / duration.value) * 100
})

// Load video data if not provided by props
async function loadVideo(): Promise<void> {
  if (!videoId.value) {
    error.value = 'No video ID provided'
    return
  }

  if (props.video) {
    videoData.value = props.video
    return
  }

  isLoading.value = true
  error.value = null

  try {
    const response = await fetch(`/api/videos/${videoId.value}/meta`)
    if (!response.ok) {
      if (response.status === 404) {
        error.value = 'Video not found'
      } else {
        error.value = 'Failed to load video'
      }
      return
    }

    videoData.value = await response.json()
  } catch (err) {
    console.error('Error loading video:', err)
    error.value = 'Failed to load video'
  } finally {
    isLoading.value = false
  }
}

// Video controls
function togglePlay(): void {
  if (!videoRef.value) return
  if (isPlaying.value) {
    videoRef.value.pause()
  } else {
    videoRef.value.play()
  }
}

function handlePlay(): void {
  isPlaying.value = true
}

function handlePause(): void {
  isPlaying.value = false
}

function handleTimeUpdate(): void {
  if (!videoRef.value) return
  currentTime.value = videoRef.value.currentTime
}

function handleLoadedMetadata(): void {
  if (!videoRef.value) return
  duration.value = videoRef.value.duration
}

function handleVolumeChange(): void {
  if (!videoRef.value) return
  volume.value = videoRef.value.volume
  isMuted.value = videoRef.value.muted
}

function seek(event: MouseEvent): void {
  if (!videoRef.value) return
  const target = event.currentTarget as HTMLElement
  const rect = target.getBoundingClientRect()
  const percent = (event.clientX - rect.left) / rect.width
  videoRef.value.currentTime = percent * duration.value
}

function toggleMute(): void {
  if (!videoRef.value) return
  videoRef.value.muted = !videoRef.value.muted
}

function setVolume(event: Event): void {
  if (!videoRef.value) return
  const target = event.target as HTMLInputElement
  videoRef.value.volume = parseFloat(target.value)
}

function toggleFullscreen(): void {
  if (!videoRef.value) return

  if (!document.fullscreenElement) {
    videoRef.value.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
}

function skip(seconds: number): void {
  if (!videoRef.value) return
  videoRef.value.currentTime = Math.max(0, Math.min(duration.value, videoRef.value.currentTime + seconds))
}

// Keyboard shortcuts
function handleKeydown(event: KeyboardEvent): void {
  if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) return

  switch (event.key) {
    case ' ':
    case 'k':
      event.preventDefault()
      togglePlay()
      break
    case 'ArrowLeft':
      event.preventDefault()
      skip(-10)
      break
    case 'ArrowRight':
      event.preventDefault()
      skip(10)
      break
    case 'm':
      event.preventDefault()
      toggleMute()
      break
    case 'f':
      event.preventDefault()
      toggleFullscreen()
      break
  }
}

onMounted(() => {
  loadVideo()
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <Layout>
    <!-- Loading -->
    <div v-if="isLoading" class="flex items-center justify-center py-16">
      <div class="h-12 w-12 animate-spin rounded-full border-4 border-accent border-t-transparent" />
    </div>

    <!-- Error -->
    <div v-else-if="error" class="flex flex-col items-center justify-center py-16 text-center">
      <svg class="h-16 w-16 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <h2 class="mt-4 text-xl font-bold text-foreground">{{ error }}</h2>
      <Link href="/gallery" class="mt-4 text-accent hover:underline">
        Back to Gallery
      </Link>
    </div>

    <!-- Video Content -->
    <div v-else-if="videoData" class="space-y-6">
      <!-- Back link -->
      <Link href="/gallery" class="inline-flex items-center gap-2 text-foreground/70 hover:text-foreground">
        <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        Back to Gallery
      </Link>

      <!-- Video Player Container -->
      <div class="overflow-hidden rounded-2xl bg-black shadow-2xl">
        <!-- Video Element -->
        <div class="relative aspect-video">
          <video
            ref="videoRef"
            :src="`/api/videos/${videoId}`"
            class="h-full w-full"
            @play="handlePlay"
            @pause="handlePause"
            @timeupdate="handleTimeUpdate"
            @loadedmetadata="handleLoadedMetadata"
            @volumechange="handleVolumeChange"
            @click="togglePlay"
          />

          <!-- Play/Pause Overlay (shown when paused) -->
          <div
            v-if="!isPlaying"
            class="absolute inset-0 flex cursor-pointer items-center justify-center bg-black/30"
            @click="togglePlay"
          >
            <div class="flex h-20 w-20 items-center justify-center rounded-full bg-white/90 transition-transform hover:scale-110">
              <svg class="ml-1 h-10 w-10 text-accent" fill="currentColor" viewBox="0 0 24 24">
                <path d="M8 5v14l11-7z" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Custom Controls -->
        <div class="bg-earth-black p-4">
          <!-- Progress bar -->
          <div
            class="mb-4 h-2 cursor-pointer rounded-full bg-white/20"
            @click="seek"
          >
            <div
              class="h-full rounded-full bg-accent transition-all"
              :style="{ width: `${progressPercent}%` }"
            />
          </div>

          <div class="flex items-center justify-between">
            <!-- Left controls -->
            <div class="flex items-center gap-4">
              <!-- Play/Pause -->
              <button
                @click="togglePlay"
                class="text-white transition-colors hover:text-accent"
              >
                <svg v-if="isPlaying" class="h-8 w-8" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
                </svg>
                <svg v-else class="h-8 w-8" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M8 5v14l11-7z" />
                </svg>
              </button>

              <!-- Skip buttons -->
              <button @click="skip(-10)" class="text-white/70 transition-colors hover:text-white" title="Skip back 10s">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12.066 11.2a1 1 0 000 1.6l5.334 4A1 1 0 0019 16V8a1 1 0 00-1.6-.8l-5.333 4zM4.066 11.2a1 1 0 000 1.6l5.334 4A1 1 0 0011 16V8a1 1 0 00-1.6-.8l-5.334 4z" />
                </svg>
              </button>
              <button @click="skip(10)" class="text-white/70 transition-colors hover:text-white" title="Skip forward 10s">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.933 12.8a1 1 0 000-1.6L6.6 7.2A1 1 0 005 8v8a1 1 0 001.6.8l5.333-4zM19.933 12.8a1 1 0 000-1.6l-5.333-4A1 1 0 0013 8v8a1 1 0 001.6.8l5.333-4z" />
                </svg>
              </button>

              <!-- Volume -->
              <div class="flex items-center gap-2">
                <button @click="toggleMute" class="text-white/70 transition-colors hover:text-white">
                  <svg v-if="isMuted || volume === 0" class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z" clip-rule="evenodd" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2" />
                  </svg>
                  <svg v-else class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z" />
                  </svg>
                </button>
                <input
                  type="range"
                  min="0"
                  max="1"
                  step="0.1"
                  :value="volume"
                  @input="setVolume"
                  class="h-1 w-20 cursor-pointer appearance-none rounded-full bg-white/20 accent-accent"
                />
              </div>

              <!-- Time -->
              <span class="text-sm text-white/70">
                {{ formatTime(currentTime) }} / {{ formatTime(duration) }}
              </span>
            </div>

            <!-- Right controls -->
            <div class="flex items-center gap-4">
              <!-- Fullscreen -->
              <button
                @click="toggleFullscreen"
                class="text-white/70 transition-colors hover:text-white"
                title="Toggle fullscreen (f)"
              >
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Video Info -->
      <div class="rounded-2xl border border-foreground/10 bg-snow-white p-6 shadow-lg">
        <h1 class="text-2xl font-bold text-foreground">{{ videoData.title }}</h1>

        <p v-if="videoData.createdAt" class="mt-2 text-sm text-foreground/50">
          Published on {{ formatDate(videoData.createdAt) }}
        </p>

        <p class="mt-4 text-foreground/70">{{ videoData.description }}</p>

        <!-- Tags -->
        <div v-if="videoData.tags?.length" class="mt-4 flex flex-wrap gap-2">
          <span
            v-for="tag in videoData.tags"
            :key="tag"
            class="rounded-full bg-accent/10 px-3 py-1 text-sm font-medium text-accent"
          >
            {{ tag }}
          </span>
        </div>
      </div>

      <!-- Keyboard shortcuts hint -->
      <div class="text-center text-sm text-foreground/50">
        Keyboard shortcuts: Space/K to play/pause, Arrow keys to skip, M to mute, F for fullscreen
      </div>
    </div>
  </Layout>
</template>
