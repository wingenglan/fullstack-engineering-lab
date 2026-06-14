import {
  defineConfig,
  presetUno,
  presetAttributify,
  presetIcons,
} from 'unocss'

export default defineConfig({
  presets: [
    presetUno(),
    presetAttributify(),
    presetIcons({
      collections: {
        lucide: () => import('@iconify-json/lucide').then((i) => i.default),
      },
    }),
  ],
  theme: {
    colors: {
      brand: {
        DEFAULT: '#6366F1',
        light: '#818CF8',
        dark: '#4F46E5',
      },
      accent: '#8B5CF6',
      success: '#10B981',
      warning: '#F59E0B',
      danger: '#EF4444',
      bg: {
        DEFAULT: '#09090B',
        card: '#18181B',
        elevated: '#27272A',
      },
      text: {
        primary: '#F4F4F5',
        secondary: '#A1A1AA',
        muted: '#71717A',
      },
      border: 'rgba(255, 255, 255, 0.08)',
    },
  },
  shortcuts: {
    'glow-card': 'relative bg-bg-card rounded-xl border border-border transition-all duration-300 hover:shadow-[0_0_20px_rgba(99,102,241,0.1),0_4px_16px_rgba(0,0,0,0.3)] hover:border-brand/25',
    'btn-primary': 'px-6 py-2.5 bg-brand text-white rounded-lg font-medium transition-all duration-200 hover:bg-brand-dark active:scale-95',
    'btn-ghost': 'px-6 py-2.5 text-text-secondary border border-border rounded-lg font-medium transition-all duration-200 hover:text-text-primary hover:border-brand/50',
    'text-gradient': 'bg-gradient-to-r from-brand to-accent bg-clip-text text-transparent',
    'section-title': 'text-3xl font-bold text-text-primary mb-2',
    'section-desc': 'text-text-secondary text-lg',
  },
})
