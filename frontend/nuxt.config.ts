// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: false },

  devServer: {
    port: 5454,
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || process.env.API_BASE || '',
    },
  },

  app: {
    head: {
      title: 'Enterprise Sizing & Capacity Planning Hub — Bank Mandiri',
      meta: [
        { name: 'description', content: 'Enterprise Sizing, Power-Law scaling, and EAWG Capacity Planning Hub' },
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&family=JetBrains+Mono:wght@400;500&display=swap',
        },
      ],
    },
  },

  css: [
    'katex/dist/katex.min.css',
    '~/assets/css/main.css',
  ],
})
