<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/utils'

interface Props {
  variant?: 'default' | 'outline' | 'ghost' | 'destructive'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'default',
  disabled: false,
})

const classes = computed(() =>
  cn(
    'inline-flex items-center justify-center rounded-md font-medium transition-colors',
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
    'disabled:pointer-events-none disabled:opacity-50',
    {
      'bg-primary text-snow-white hover:bg-primary/90': props.variant === 'default',
      'border border-foreground/20 bg-transparent hover:bg-foreground/5': props.variant === 'outline',
      'hover:bg-foreground/5': props.variant === 'ghost',
      'bg-red-600 text-white hover:bg-red-700': props.variant === 'destructive',
    },
    {
      'h-10 px-4 py-2': props.size === 'default',
      'h-8 px-3 text-sm': props.size === 'sm',
      'h-12 px-6 text-lg': props.size === 'lg',
      'h-10 w-10': props.size === 'icon',
    }
  )
)
</script>

<template>
  <button :class="classes" :disabled="disabled">
    <slot />
  </button>
</template>
