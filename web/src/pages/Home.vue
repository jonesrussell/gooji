<script setup lang="ts">
import { Link } from '@inertiajs/vue3'
import Container from '@/components/ui/Container.vue'
</script>

<template>
  <div class="min-h-screen">
    <!-- Hero -->
    <div class="bg-gradient-to-br from-primary-600 to-primary-800 text-white py-20">
      <Container>
        <h1 class="text-4xl md:text-5xl font-bold">Gooji</h1>
        <p class="mt-4 text-xl text-primary-100 max-w-2xl">
          Preserving Ojibwe language and culture through community video sharing.
        </p>
        <div class="mt-8 flex gap-4">
          <Link href="/gallery" class="btn btn-primary bg-white text-primary-700 hover:bg-primary-50">
            Browse Videos
          </Link>
          <Link href="/upload" class="btn border border-white/30 hover:bg-white/10">
            Share Your Story
          </Link>
        </div>
      </Container>
    </div>

    <!-- Features -->
    <Container class="py-16">
      <h2 class="text-2xl font-bold text-gray-900">Community Video Sharing</h2>
      <p class="mt-2 text-gray-600 max-w-2xl">
        Record or upload videos to share Ojibwe language, stories, and cultural knowledge with the community.
      </p>
    </Container>
  </div>
</template>
