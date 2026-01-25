<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { Link } from '@inertiajs/vue3'
import Layout from '../components/Layout.vue'

// Recording state
const stream = ref<MediaStream | null>(null)
const mediaRecorder = ref<MediaRecorder | null>(null)
const recordedChunks = ref<Blob[]>([])
const isRecording = ref(false)
const isPaused = ref(false)
const hasRecording = ref(false)
const recordingTime = ref(0)
const recordingInterval = ref<ReturnType<typeof setInterval> | null>(null)

// Camera state
const cameraReady = ref(false)
const cameraError = ref<string | null>(null)
const showOverlay = ref(true)

// Preview state
const previewMode = ref(false)
const previewUrl = ref<string | null>(null)

// Form state
const title = ref('')
const description = ref('')
const tags = ref('')

// Upload state
const isUploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref<string | null>(null)
const showSuccess = ref(false)

// Refs
const videoRef = ref<HTMLVideoElement | null>(null)
// Preview video uses the same ref as live video

// Computed
const formattedTime = computed(() => {
  const mins = Math.floor(recordingTime.value / 60)
  const secs = recordingTime.value % 60
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
})

const isFormValid = computed(() => {
  return hasRecording.value && title.value.trim() && description.value.trim()
})

// Initialize camera
async function initCamera(): Promise<void> {
  cameraError.value = null

  try {
    // Check if getUserMedia is supported
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      throw new Error('Your browser does not support video recording')
    }

    // Check available devices
    const devices = await navigator.mediaDevices.enumerateDevices()
    const videoDevices = devices.filter(device => device.kind === 'videoinput')

    if (videoDevices.length === 0) {
      throw new Error('No camera found. Please connect a camera and try again.')
    }

    // Get media stream
    stream.value = await navigator.mediaDevices.getUserMedia({
      video: {
        width: { ideal: 1280 },
        height: { ideal: 720 },
        facingMode: 'user',
      },
      audio: true,
    })

    // Set video source
    if (videoRef.value) {
      videoRef.value.srcObject = stream.value
      await videoRef.value.play()
    }

    // Set up MediaRecorder
    let mimeType = 'video/webm;codecs=vp9,opus'
    if (!MediaRecorder.isTypeSupported(mimeType)) {
      mimeType = 'video/webm'
      if (!MediaRecorder.isTypeSupported(mimeType)) {
        mimeType = ''
      }
    }

    mediaRecorder.value = new MediaRecorder(stream.value, mimeType ? { mimeType } : {})

    mediaRecorder.value.ondataavailable = (event: BlobEvent) => {
      if (event.data.size > 0) {
        recordedChunks.value.push(event.data)
      }
    }

    mediaRecorder.value.onstop = () => {
      const blob = new Blob(recordedChunks.value, { type: 'video/webm' })
      previewUrl.value = URL.createObjectURL(blob)
      hasRecording.value = true
      previewMode.value = true
    }

    cameraReady.value = true
    showOverlay.value = false
  } catch (err) {
    console.error('Camera initialization error:', err)

    if (err instanceof Error) {
      if (err.name === 'NotAllowedError') {
        cameraError.value = 'Camera access denied. Please grant camera and microphone permissions and refresh the page.'
      } else if (err.name === 'NotFoundError') {
        cameraError.value = 'No camera found. Please check your device connections.'
      } else if (err.name === 'NotSupportedError') {
        cameraError.value = 'Your browser does not support video recording.'
      } else {
        cameraError.value = err.message
      }
    } else {
      cameraError.value = 'Failed to access camera'
    }
  }
}

// Start recording
function startRecording(): void {
  if (!mediaRecorder.value || isRecording.value) return

  recordedChunks.value = []
  recordingTime.value = 0
  previewMode.value = false
  hasRecording.value = false

  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = null
  }

  // Restore live preview
  if (videoRef.value && stream.value) {
    videoRef.value.srcObject = stream.value
    videoRef.value.play()
  }

  mediaRecorder.value.start()
  isRecording.value = true
  isPaused.value = false

  // Start timer
  recordingInterval.value = setInterval(() => {
    recordingTime.value++
  }, 1000)
}

// Stop recording
function stopRecording(): void {
  if (!mediaRecorder.value || !isRecording.value) return

  mediaRecorder.value.stop()
  isRecording.value = false
  isPaused.value = false

  if (recordingInterval.value) {
    clearInterval(recordingInterval.value)
    recordingInterval.value = null
  }
}

// Pause/resume recording
function togglePause(): void {
  if (!mediaRecorder.value || !isRecording.value) return

  if (isPaused.value) {
    mediaRecorder.value.resume()
    isPaused.value = false
    recordingInterval.value = setInterval(() => {
      recordingTime.value++
    }, 1000)
  } else {
    mediaRecorder.value.pause()
    isPaused.value = true
    if (recordingInterval.value) {
      clearInterval(recordingInterval.value)
      recordingInterval.value = null
    }
  }
}

// Discard recording and start over
function discardRecording(): void {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = null
  }

  recordedChunks.value = []
  hasRecording.value = false
  previewMode.value = false
  recordingTime.value = 0
  title.value = ''
  description.value = ''
  tags.value = ''

  // Restore live preview
  if (videoRef.value && stream.value) {
    videoRef.value.srcObject = stream.value
    videoRef.value.play()
  }
}

// Upload recording
async function uploadRecording(): Promise<void> {
  if (!hasRecording.value || !isFormValid.value) return

  isUploading.value = true
  uploadProgress.value = 0
  uploadError.value = null

  const blob = new Blob(recordedChunks.value, { type: 'video/webm' })
  const formData = new FormData()
  formData.append('video', blob, 'recording.webm')
  formData.append('title', title.value)
  formData.append('description', description.value)
  formData.append('tags', tags.value)

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
      uploadError.value = 'Upload failed. Please check your connection.'
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

// Reset for another recording
function resetForAnother(): void {
  showSuccess.value = false
  discardRecording()
}

// Cleanup
function cleanup(): void {
  if (recordingInterval.value) {
    clearInterval(recordingInterval.value)
  }

  if (stream.value) {
    stream.value.getTracks().forEach(track => track.stop())
  }

  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
  }
}

onMounted(() => {
  initCamera()
})

onUnmounted(() => {
  cleanup()
})
</script>

<template>
  <Layout>
    <div class="mx-auto max-w-4xl space-y-6">
      <!-- Header -->
      <div class="text-center">
        <h1 class="text-3xl font-bold text-foreground">Recording Studio</h1>
        <p class="mt-2 text-foreground/70">Record and share Ojibwe language and cultural content</p>
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
          <h2 class="mt-4 text-2xl font-bold text-foreground">Recording Uploaded!</h2>
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
              Record Another
            </button>
          </div>
        </div>
      </div>

      <!-- Camera/Preview Container -->
      <div class="overflow-hidden rounded-2xl bg-earth-black shadow-2xl">
        <!-- Video Area -->
        <div class="relative aspect-video">
          <!-- Loading Overlay -->
          <div
            v-if="showOverlay && !cameraError"
            class="absolute inset-0 z-10 flex flex-col items-center justify-center bg-earth-black"
          >
            <div class="h-12 w-12 animate-spin rounded-full border-4 border-accent border-t-transparent" />
            <p class="mt-4 text-snow-white">Initializing camera...</p>
          </div>

          <!-- Error State -->
          <div
            v-if="cameraError"
            class="absolute inset-0 z-10 flex flex-col items-center justify-center bg-earth-black p-8 text-center"
          >
            <svg class="h-16 w-16 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <p class="mt-4 text-lg text-snow-white">{{ cameraError }}</p>
            <button
              @click="initCamera"
              class="mt-4 rounded-lg bg-accent px-6 py-2 font-medium text-white hover:bg-accent/90"
            >
              Try Again
            </button>
          </div>

          <!-- Live Video Feed -->
          <video
            v-show="!previewMode"
            ref="videoRef"
            autoplay
            muted
            playsinline
            class="h-full w-full object-cover"
          />

          <!-- Recording Preview -->
          <video
            v-show="previewMode && previewUrl"
            ref="previewVideoRef"
            :src="previewUrl || undefined"
            controls
            class="h-full w-full object-cover"
          />

          <!-- Recording Indicator -->
          <div
            v-if="isRecording"
            class="absolute left-4 top-4 flex items-center gap-2 rounded-full bg-red-600 px-4 py-2 text-white"
          >
            <span class="h-3 w-3 animate-pulse rounded-full bg-white" />
            <span class="font-mono font-medium">{{ formattedTime }}</span>
            <span v-if="isPaused" class="text-sm">(Paused)</span>
          </div>
        </div>

        <!-- Controls -->
        <div class="flex items-center justify-center gap-4 bg-earth-black/90 p-6">
          <!-- Not recording: Show record button -->
          <template v-if="!isRecording && !previewMode">
            <button
              @click="startRecording"
              :disabled="!cameraReady"
              class="flex h-16 w-16 items-center justify-center rounded-full bg-red-500 text-white transition-all hover:bg-red-600 hover:scale-110 disabled:opacity-50 disabled:hover:scale-100"
              title="Start Recording"
            >
              <svg class="h-8 w-8" fill="currentColor" viewBox="0 0 24 24">
                <circle cx="12" cy="12" r="8" />
              </svg>
            </button>
          </template>

          <!-- Recording: Show pause and stop -->
          <template v-if="isRecording">
            <button
              @click="togglePause"
              class="flex h-12 w-12 items-center justify-center rounded-full bg-white/20 text-white transition-colors hover:bg-white/30"
              :title="isPaused ? 'Resume' : 'Pause'"
            >
              <svg v-if="isPaused" class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                <path d="M8 5v14l11-7z" />
              </svg>
              <svg v-else class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
              </svg>
            </button>

            <button
              @click="stopRecording"
              class="flex h-16 w-16 items-center justify-center rounded-full bg-red-500 text-white transition-all hover:bg-red-600"
              title="Stop Recording"
            >
              <svg class="h-8 w-8" fill="currentColor" viewBox="0 0 24 24">
                <rect x="6" y="6" width="12" height="12" rx="2" />
              </svg>
            </button>
          </template>

          <!-- Preview mode: Show re-record and continue -->
          <template v-if="previewMode && !isRecording">
            <button
              @click="discardRecording"
              class="flex items-center gap-2 rounded-lg bg-white/20 px-6 py-3 font-medium text-white transition-colors hover:bg-white/30"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Re-record
            </button>

            <span class="text-snow-white/70">Preview your recording above</span>
          </template>
        </div>
      </div>

      <!-- Upload Form (shown after recording) -->
      <div
        v-if="hasRecording"
        class="space-y-6 rounded-2xl border border-foreground/10 bg-snow-white p-6 shadow-lg"
      >
        <h2 class="text-xl font-bold text-foreground">Add Details</h2>

        <!-- Error message -->
        <div v-if="uploadError" class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-700">
          {{ uploadError }}
        </div>

        <form @submit.prevent="uploadRecording" class="space-y-4">
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
              rows="3"
              class="w-full resize-none rounded-lg border border-foreground/20 bg-background px-4 py-2 focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
              placeholder="Describe your recording"
            />
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
              placeholder="language, culture, story (comma separated)"
            />
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
              Upload Recording
            </template>
          </button>
        </form>
      </div>

      <!-- Tips -->
      <div class="rounded-2xl border border-foreground/10 bg-snow-white p-6 shadow-lg">
        <h3 class="font-bold text-foreground">Recording Tips</h3>
        <ul class="mt-3 space-y-2 text-sm text-foreground/70">
          <li class="flex items-start gap-2">
            <svg class="mt-0.5 h-4 w-4 flex-shrink-0 text-accent" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            Find a quiet space with good lighting
          </li>
          <li class="flex items-start gap-2">
            <svg class="mt-0.5 h-4 w-4 flex-shrink-0 text-accent" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            Speak clearly and at a natural pace
          </li>
          <li class="flex items-start gap-2">
            <svg class="mt-0.5 h-4 w-4 flex-shrink-0 text-accent" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            You can pause and resume during recording
          </li>
          <li class="flex items-start gap-2">
            <svg class="mt-0.5 h-4 w-4 flex-shrink-0 text-accent" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            Preview your recording before uploading
          </li>
        </ul>
      </div>

      <!-- Alternative: Upload instead -->
      <div class="text-center">
        <p class="text-foreground/70">
          Have an existing video file?
          <Link href="/upload" class="font-medium text-accent hover:underline">
            Upload a File
          </Link>
        </p>
      </div>
    </div>
  </Layout>
</template>
