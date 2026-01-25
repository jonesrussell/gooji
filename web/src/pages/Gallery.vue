<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Link } from '@inertiajs/vue3'
import Layout from '../components/Layout.vue'

interface Video {
  id: string
  title: string
  description: string
  tags: string[]
  createdAt?: string
  duration?: number
  thumbnail?: string
}

// Props from Inertia
defineProps<{
  videos?: Video[]
}>()

// State
const videoList = ref<Video[]>([])
const searchQuery = ref('')
const selectedTag = ref('')
const sortBy = ref('newest')
const isLoading = ref(false)
const currentPage = ref(1)
const hasMore = ref(true)

// Modal state
const showModal = ref(false)
const selectedVideo = ref<Video | null>(null)

// Available tags (could be passed from server)
const availableTags = ['language', 'culture', 'story', 'history', 'music', 'crafts', 'ceremony']

// Format duration from seconds to MM:SS
function formatDuration(seconds: number): string {
  if (!seconds) return ''
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

// Format date
function formatDate(dateString: string): string {
  if (!dateString) return ''
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = Math.abs(now.getTime() - date.getTime())
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays === 1) return 'Today'
  if (diffDays === 2) return 'Yesterday'
  if (diffDays <= 7) return `${diffDays - 1} days ago`
  if (diffDays <= 30) return `${Math.ceil(diffDays / 7)} weeks ago`
  if (diffDays <= 365) return `${Math.ceil(diffDays / 30)} months ago`
  return date.toLocaleDateString()
}

// Load videos from API
async function loadVideos(page = 1, append = false): Promise<void> {
  if (isLoading.value) return

  isLoading.value = true

  try {
    const params = new URLSearchParams({
      page: page.toString(),
      search: searchQuery.value,
      tag: selectedTag.value,
      sort: sortBy.value,
    })

    const response = await fetch(`/api/videos?${params}`)
    if (!response.ok) throw new Error('Failed to load videos')

    const videos = await response.json() as Video[]
    hasMore.value = videos.length === 10

    if (append) {
      videoList.value = [...videoList.value, ...videos]
    } else {
      videoList.value = videos
    }

    currentPage.value = page
  } catch (err) {
    console.error('Error loading videos:', err)
  } finally {
    isLoading.value = false
  }
}

// Load more videos
function loadMore(): void {
  if (hasMore.value && !isLoading.value) {
    loadVideos(currentPage.value + 1, true)
  }
}

// Reset and reload on filter change
function resetAndReload(): void {
  currentPage.value = 1
  hasMore.value = true
  loadVideos(1)
}

// Watch for filter changes with debounce
let searchTimeout: ReturnType<typeof setTimeout> | null = null
watch(searchQuery, () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(resetAndReload, 300)
})

watch([selectedTag, sortBy], resetAndReload)

// Open video modal
function openModal(video: Video): void {
  selectedVideo.value = video
  showModal.value = true
}

// Close video modal
function closeModal(): void {
  showModal.value = false
  selectedVideo.value = null
}

// Delete video
async function deleteVideo(video: Video): Promise<void> {
  if (!confirm(`Are you sure you want to delete "${video.title}"?\n\nThis action cannot be undone.`)) {
    return
  }

  try {
    const response = await fetch(`/api/videos/${video.id}`, { method: 'DELETE' })
    if (!response.ok) throw new Error('Failed to delete video')

    videoList.value = videoList.value.filter(v => v.id !== video.id)
  } catch (err) {
    console.error('Error deleting video:', err)
    alert('Failed to delete video. Please try again.')
  }
}

// Fallback thumbnail
const fallbackThumbnail = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIwIiBoZWlnaHQ9IjE4MCIgdmlld0JveD0iMCAwIDMyMCAxODAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIzMjAiIGhlaWdodD0iMTgwIiBmaWxsPSIjRjNGNEY2Ii8+CjwvcmVjdD4KPC9zdmc+'

// Handle thumbnail error
function handleThumbnailError(event: Event): void {
  const img = event.target as HTMLImageElement
  img.src = fallbackThumbnail
}

onMounted(() => {
  loadVideos(1)
})
</script>

<template>
  <Layout>
    <div class="space-y-8">
      <!-- Header -->
      <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-3xl font-bold text-foreground">Video Gallery</h1>
          <p class="mt-1 text-foreground/70">Browse Ojibwe language and cultural videos</p>
        </div>
        <Link
          href="/record"
          class="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 font-medium text-white hover:bg-accent/90 transition-colors"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="10" stroke-width="2" />
            <circle cx="12" cy="12" r="4" fill="currentColor" />
          </svg>
          Record Video
        </Link>
      </div>

      <!-- Filters -->
      <div class="flex flex-col gap-4 sm:flex-row sm:items-center">
        <!-- Search -->
        <div class="relative flex-1">
          <svg
            class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-foreground/50"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search videos..."
            class="w-full rounded-lg border border-foreground/20 bg-background py-2 pl-10 pr-4 focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
          />
        </div>

        <!-- Tag filter -->
        <select
          v-model="selectedTag"
          class="rounded-lg border border-foreground/20 bg-background px-4 py-2 focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
        >
          <option value="">All Categories</option>
          <option v-for="tag in availableTags" :key="tag" :value="tag">
            {{ tag.charAt(0).toUpperCase() + tag.slice(1) }}
          </option>
        </select>

        <!-- Sort -->
        <select
          v-model="sortBy"
          class="rounded-lg border border-foreground/20 bg-background px-4 py-2 focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
        >
          <option value="newest">Newest First</option>
          <option value="oldest">Oldest First</option>
          <option value="title">Title A-Z</option>
        </select>
      </div>

      <!-- Video Grid -->
      <div v-if="videoList.length > 0" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="video in videoList"
          :key="video.id"
          class="group overflow-hidden rounded-2xl border border-foreground/10 bg-snow-white shadow-lg transition-all hover:-translate-y-1 hover:shadow-xl"
        >
          <!-- Thumbnail -->
          <div class="relative aspect-video cursor-pointer overflow-hidden" @click="openModal(video)">
            <img
              :src="`/api/thumbnails?id=${video.id}`"
              :alt="video.title"
              class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
              @error="handleThumbnailError"
            />
            <!-- Play overlay -->
            <div class="absolute inset-0 flex items-center justify-center bg-black/20 transition-colors group-hover:bg-black/30">
              <div class="flex h-16 w-16 items-center justify-center rounded-full bg-white/90 transition-all group-hover:scale-110 group-hover:bg-white">
                <svg class="ml-1 h-8 w-8 text-accent" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M8 5v14l11-7z" />
                </svg>
              </div>
            </div>
            <!-- Duration -->
            <div
              v-if="video.duration"
              class="absolute bottom-2 right-2 rounded-lg bg-black/70 px-2 py-1 text-xs font-medium text-white"
            >
              {{ formatDuration(video.duration) }}
            </div>
          </div>

          <!-- Info -->
          <div class="p-6">
            <h3 class="mb-2 line-clamp-2 text-lg font-bold text-foreground transition-colors group-hover:text-accent">
              {{ video.title }}
            </h3>
            <p class="mb-4 line-clamp-2 text-sm leading-relaxed text-foreground/70">
              {{ video.description }}
            </p>

            <!-- Tags -->
            <div class="mb-4 flex flex-wrap gap-2">
              <span
                v-for="tag in video.tags"
                :key="tag"
                class="rounded-full bg-accent/10 px-3 py-1 text-xs font-medium text-accent transition-colors group-hover:bg-accent/20"
              >
                {{ tag }}
              </span>
            </div>

            <!-- Footer -->
            <div class="flex items-center justify-between text-xs text-foreground/50">
              <span class="flex items-center gap-1">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ video.createdAt ? formatDate(video.createdAt) : 'Recently added' }}
              </span>
              <button
                @click.stop="deleteVideo(video)"
                class="flex items-center rounded p-1 text-red-500 transition-colors hover:bg-red-50 hover:text-red-700"
                title="Delete video"
              >
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else-if="!isLoading" class="flex flex-col items-center justify-center py-16 text-center">
        <svg class="h-16 w-16 text-foreground/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
        <h3 class="mt-4 text-lg font-medium text-foreground">No videos found</h3>
        <p class="mt-2 text-foreground/70">Try adjusting your search or filters</p>
      </div>

      <!-- Loading indicator -->
      <div v-if="isLoading" class="flex justify-center py-8">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-accent border-t-transparent" />
      </div>

      <!-- Load more -->
      <div v-if="hasMore && videoList.length > 0 && !isLoading" class="flex justify-center">
        <button
          @click="loadMore"
          class="rounded-lg border border-foreground/20 px-6 py-2 font-medium text-foreground transition-colors hover:border-accent hover:text-accent"
        >
          Load More Videos
        </button>
      </div>
    </div>

    <!-- Video Modal -->
    <Teleport to="body">
      <div
        v-if="showModal && selectedVideo"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4"
        @click.self="closeModal"
      >
        <div class="relative w-full max-w-4xl overflow-hidden rounded-2xl bg-snow-white">
          <!-- Close button -->
          <button
            @click="closeModal"
            class="absolute right-4 top-4 z-10 rounded-full bg-black/50 p-2 text-white transition-colors hover:bg-black/70"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>

          <!-- Video player -->
          <video
            :src="`/api/videos/${selectedVideo.id}`"
            controls
            autoplay
            class="aspect-video w-full bg-black"
          />

          <!-- Video info -->
          <div class="p-6">
            <h2 class="text-2xl font-bold text-foreground">{{ selectedVideo.title }}</h2>
            <p class="mt-2 text-foreground/70">{{ selectedVideo.description }}</p>
            <div class="mt-4 flex flex-wrap gap-2">
              <span
                v-for="tag in selectedVideo.tags"
                :key="tag"
                class="rounded-full bg-accent/10 px-3 py-1 text-sm font-medium text-accent"
              >
                {{ tag }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </Layout>
</template>
