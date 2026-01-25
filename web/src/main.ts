import { createApp, h, DefineComponent } from 'vue'
import { createInertiaApp } from '@inertiajs/vue3'
import './assets/index.css'

type PageModule = { default: DefineComponent }

createInertiaApp({
  resolve: (name) => {
    const pages = import.meta.glob<PageModule>('./pages/**/*.vue', { eager: true })
    return pages[`./pages/${name}.vue`]
  },
  setup({ el, App, props, plugin }) {
    createApp({ render: () => h(App, props) })
      .use(plugin)
      .mount(el)
  },
})
