import astroSingleFile from 'astro-single-file'
import { defineConfig } from 'astro/config'
import tailwind from '@astrojs/tailwind'

export default defineConfig({
  integrations: [astroSingleFile(), tailwind()]
})
