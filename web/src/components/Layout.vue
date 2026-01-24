<script setup lang="ts">
import { Link } from '@inertiajs/vue3'
</script>

<template>
  <div class="min-h-screen bg-background flex flex-col">
    <nav class="border-b border-foreground/10 bg-snow-white">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex h-16 items-center justify-between">
          <Link href="/" class="text-2xl font-bold text-foreground">
            Gooji
          </Link>
          <div class="flex gap-6">
            <Link
              href="/gallery"
              class="text-foreground/70 hover:text-foreground border-b-2 border-transparent hover:border-accent transition-colors"
            >
              Gallery
            </Link>
            <Link
              href="/record"
              class="text-foreground/70 hover:text-foreground border-b-2 border-transparent hover:border-accent transition-colors"
            >
              Record
            </Link>
            <Link
              href="/upload"
              class="text-foreground/70 hover:text-foreground border-b-2 border-transparent hover:border-accent transition-colors"
            >
              Upload
            </Link>
          </div>
        </div>
      </div>
    </nav>

    <main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8 flex-1">
      <slot />
    </main>

    <footer class="border-t border-foreground/10 bg-earth-black text-snow-white">
      <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <p class="text-center text-sm">
          Gooji - Preserving Ojibwe language and culture through video
        </p>
      </div>
    </footer>
  </div>
</template>
