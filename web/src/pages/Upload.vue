<script setup lang="ts">
import { ref, computed } from 'vue'
import { Link } from '@inertiajs/vue3'
import Layout from '../components/Layout.vue'

// Form state
const selectedFile = ref<File | null>(null)
const title = ref('')
const description = ref('')
const category = ref('')
const tags = ref('')
const language = ref('oji')
const isPublic = ref(true)

// Upload state
const isUploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref<string | null>(null)
const showSuccess = ref(false)

// Drag state
const isDragging = ref(false)

// Categories
const categories = [
  { value: '', label: 'Select a category' },
  { value: 'language', label: 'Language & Learning' },
  { value: 'culture', label: 'Culture & Traditions' },
  { value: 'story', label: 'Stories & Legends' },
  { value: 'history', label: 'History' },
  { value: 'music', label: 'Music & Song' },
  { value: 'crafts', label: 'Traditional Crafts' },
  { value: 'ceremony', label: 'Ceremony & Spirituality' },
]

// Languages
const languages = [
  { value: 'oji', label: 'Ojibwe' },
  { value: 'en', label: 'English' },
  { value: 'both', label: 'Both' },
]

// Tag suggestions based on category
const tagSuggestions: Record<string, string> = {
  language: 'ojibwe, language, learning, pronunciation',
  culture: 'tradition, culture, ceremony, customs',
  story: 'story, legend, tale, narrative',
  history: 'history, historical, past, ancestors',
  music: 'music, song, singing, melody',
  crafts: 'craft, traditional, making, art',
  ceremony: 'ceremony, ritual, sacred, spiritual',
}

// Computed
const isFormValid = computed(() => {
  return selectedFile.value && title.value.trim() && description.value.trim() && category.value
})

const fileSizeFormatted = computed(() => {
  if (!selectedFile.value) return ''
  return formatFileSize(selectedFile.value.size)
})

const currentTagSuggestion = computed(() => {
  return tagSuggestions[category.value] || ''
})

// Format file size
function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

// Validate file
function isValidVideoFile(file: File): boolean {
  const allowedTypes = ['video/mp4', 'video/webm', 'video/quicktime', 'video/x-msvideo']
  return allowedTypes.includes(file.type)
}

function isValidFileSize(file: File): boolean {
  const maxSize = 500 * 1024 * 1024 // 500MB
  return file.size <= maxSize
}

// Handle file selection
function handleFileSelect(file: File): void {
  uploadError.value = null

  if (!isValidVideoFile(file)) {
    uploadError.value = 'Please select a valid video file (MP4, WebM, MOV, AVI).'
    return
  }

  if (!isValidFileSize(file)) {
    uploadError.value = 'File size must be less than 500MB.'
    return
  }

  selectedFile.value = file
}

// Drag and drop handlers
function handleDragEnter(e: DragEvent): void {
  e.preventDefault()
  isDragging.value = true
}

function handleDragLeave(e: DragEvent): void {
  e.preventDefault()
  isDragging.value = false
}

function handleDragOver(e: DragEvent): void {
  e.preventDefault()
}

function handleDrop(e: DragEvent): void {
  e.preventDefault()
  isDragging.value = false

  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    handleFileSelect(files[0])
  }
}

// File input change
function handleFileInput(e: Event): void {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    handleFileSelect(file)
  }
}

// Reset file
function resetFile(): void {
  selectedFile.value = null
  uploadError.value = null
}

// Submit upload
async function handleSubmit(): Promise<void> {
  if (!selectedFile.value || !isFormValid.value) return

  isUploading.value = true
  uploadProgress.value = 0
  uploadError.value = null

  const formData = new FormData()
  formData.append('video', selectedFile.value)
  formData.append('title', title.value)
  formData.append('description', description.value)
  formData.append('category', category.value)
  formData.append('tags', tags.value)
  formData.append('language', language.value)
  formData.append('public', isPublic.value.toString())

  try {
    const xhr = new XMLHttpRequest()

    xhr.upload.addEventListener('progress', (e: ProgressEvent) => {
      if (e.lengthComputable) {
        uploadProgress.value = Math.round((e.loaded / e.total) * 100)
      }
    })

    xhr.addEventListener('load', () => {
      if (xhr.status === 200) {
        showSuccess.value = true
        isUploading.value = false
      } else {
        uploadError.value = 'Upload failed. Please try again.'
        isUploading.value = false
      }
    })

    xhr.addEventListener('error', () => {
      uploadError.value = 'Upload failed. Please check your connection and try again.'
      isUploading.value = false
    })

    xhr.open('POST', '/api/videos')
    xhr.send(formData)
  } catch (err) {
    console.error('Upload error:', err)
    uploadError.value = 'Upload failed. Please try again.'
    isUploading.value = false
  }
}

// Reset form for another upload
function resetForAnother(): void {
  showSuccess.value = false
  selectedFile.value = null
  title.value = ''
  description.value = ''
  category.value = ''
  tags.value = ''
  language.value = 'oji'
  isPublic.value = true
  uploadProgress.value = 0
}
</script>

<template>
  <Layout>
    <div class="mx-auto max-w-2xl space-y-8">
      <!-- Header -->
      <div class="text-center">
        <h1 class="text-3xl font-bold text-foreground">Upload Video</h1>
        <p class="mt-2 text-foreground/70">Share Ojibwe language and cultural content with the community</p>
      </div>

      <!-- Success Modal -->
      <div
        v-if="showSuccess"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
      >
        <div class="w-full max-w-md rounded-2xl bg-snow-white p-8 text-center shadow-2xl">
          <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
            <svg class="h-8 w-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h2 class="mt-4 text-2xl font-bold text-foreground">Upload Successful!</h2>
          <p class="mt-2 text-foreground/70">Your video has been uploaded and is being processed.</p>
          <div class="mt-6 flex flex-col gap-3">
            <Link
              href="/gallery"
              class="rounded-lg bg-accent px-6 py-3 font-medium text-white transition-colors hover:bg-accent/90"
            >
              View Gallery
            </Link>
            <button
              @click="resetForAnother"
              class="rounded-lg border border-foreground/20 px-6 py-3 font-medium text-foreground transition-colors hover:border-accent hover:text-accent"
            >
              Upload Another
            </button>
          </div>
        </div>
      </div>

      <!-- Upload Form -->
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Drop Zone -->
        <div
          class="relative rounded-2xl border-2 border-dashed p-8 text-center transition-colors"
          :class="[
            isDragging ? 'border-accent bg-accent/5' : 'border-foreground/20',
            selectedFile ? 'border-green-400 bg-green-50' : 'hover:border-accent hover:bg-accent/5'
          ]"
          @dragenter="handleDragEnter"
          @dragleave="handleDragLeave"
          @dragover="handleDragOver"
          @drop="handleDrop"
        >
          <!-- No file selected -->
          <div v-if="!selectedFile" class="space-y-4">
            <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-accent/10">
              <svg class="h-8 w-8 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>
            </div>
            <div>
              <p class="text-lg font-medium text-foreground">Drag and drop your video here</p>
              <p class="text-foreground/70">or</p>
            </div>
            <label class="inline-block cursor-pointer rounded-lg bg-accent px-6 py-2 font-medium text-white transition-colors hover:bg-accent/90">
              Browse Files
              <input
                type="file"
                accept="video/mp4,video/webm,video/quicktime,video/x-msvideo"
                class="hidden"
                @change="handleFileInput"
              />
            </label>
            <p class="text-sm text-foreground/50">MP4, WebM, MOV, AVI up to 500MB</p>
          </div>

          <!-- File selected -->
          <div v-else class="space-y-4">
            <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
              <svg class="h-8 w-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div>
              <p class="text-lg font-medium text-foreground">{{ selectedFile.name }}</p>
              <p class="text-foreground/70">{{ fileSizeFormatted }}</p>
            </div>
            <button
              type="button"
              @click="resetFile"
              class="text-sm text-accent hover:underline"
            >
              Choose a different file
            </button>
          </div>
        </div>

        <!-- Error message -->
        <div v-if="uploadError" class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-700">
          {{ uploadError }}
        </div>

        <!-- Form fields -->
        <div class="space-y-4 rounded-2xl border border-foreground/10 bg-snow-white p-6 shadow-lg">
          <!-- Title -->
          <div>
            <label for="title" class="mb-1 block text-sm font-medium text-foreground">
              Title <span class="text-red-500">*</span>
            </label>
            <input
              id="title"
              v-model="title"
              type="text"
              required
              class="w-full rounded-lg border border-foreground/20 bg-background px-4 py-2 focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
              placeholder="Enter a descriptive title"
            />
          </div>

          <!-- Description -->
          <div>
            <label for="description" class="mb-1 block text-sm font-medium text-foreground">
              Description <span class="text-red-500">*</span>
            </label>
            <textarea
              id="description"
              v-model="description"
              required
              rows="4"
              class="w-full resize-none rounded-lg border border-foreground/20 bg-background px-4 py-2 focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
              placeholder="Describe your video content"
            />
          </div>

          <!-- Category -->
          <div>
            <label for="category" class="mb-1 block text-sm font-medium text-foreground">
              Category <span class="text-red-500">*</span>
            </label>
            <select
              id="category"
              v-model="category"
              required
              class="w-full rounded-lg border border-foreground/20 bg-background px-4 py-2 focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
            >
              <option v-for="cat in categories" :key="cat.value" :value="cat.value">
                {{ cat.label }}
              </option>
            </select>
          </div>

          <!-- Tags -->
          <div>
            <label for="tags" class="mb-1 block text-sm font-medium text-foreground">
              Tags
            </label>
            <input
              id="tags"
              v-model="tags"
              type="text"
              class="w-full rounded-lg border border-foreground/20 bg-background px-4 py-2 focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
              :placeholder="currentTagSuggestion ? `Suggested: ${currentTagSuggestion}` : 'Enter tags separated by commas'"
            />
            <p class="mt-1 text-xs text-foreground/50">Separate tags with commas</p>
          </div>

          <!-- Language -->
          <div>
            <label for="language" class="mb-1 block text-sm font-medium text-foreground">
              Primary Language
            </label>
            <select
              id="language"
              v-model="language"
              class="w-full rounded-lg border border-foreground/20 bg-background px-4 py-2 focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
            >
              <option v-for="lang in languages" :key="lang.value" :value="lang.value">
                {{ lang.label }}
              </option>
            </select>
          </div>

          <!-- Public checkbox -->
          <div class="flex items-center gap-3">
            <input
              id="public"
              v-model="isPublic"
              type="checkbox"
              class="h-5 w-5 rounded border-foreground/20 text-accent focus:ring-accent"
            />
            <label for="public" class="text-sm text-foreground">
              Make this video publicly visible
            </label>
          </div>
        </div>

        <!-- Upload Progress -->
        <div v-if="isUploading" class="space-y-2">
          <div class="flex items-center justify-between text-sm text-foreground">
            <span>Uploading...</span>
            <span>{{ uploadProgress }}%</span>
          </div>
          <div class="h-3 overflow-hidden rounded-full bg-foreground/10">
            <div
              class="h-full rounded-full bg-accent transition-all"
              :style="{ width: `${uploadProgress}%` }"
            />
          </div>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="!isFormValid || isUploading"
          class="flex w-full items-center justify-center gap-3 rounded-lg bg-accent px-6 py-3 font-medium text-white transition-colors hover:bg-accent/90 disabled:cursor-not-allowed disabled:opacity-50"
        >
          <template v-if="isUploading">
            <div class="h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent" />
            Uploading...
          </template>
          <template v-else>
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
            Upload Video
          </template>
        </button>
      </form>

      <!-- Alternative: Record instead -->
      <div class="text-center">
        <p class="text-foreground/70">
          Want to record a new video instead?
          <Link href="/record" class="font-medium text-accent hover:underline">
            Go to Recording Studio
          </Link>
        </p>
      </div>
    </div>
  </Layout>
</template>
