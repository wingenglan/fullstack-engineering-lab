import { defineConfig } from 'vitepress'

export default defineConfig({
  base: '/docs/',
  srcDir: 'docs',
  title: '全栈工程实验室',
  description: '可运行、可体验、可学习的全栈工程实践案例库',
  ignoreDeadLinks: true,
  themeConfig: {
    nav: [
      { text: '指南', items: [
        { text: '快速开始', link: '/guide/quick-start' },
        { text: '项目架构', link: '/guide/architecture' },
        { text: '技术栈', link: '/guide/tech-stack' },
      ]},
      { text: '案例', items: [
        { text: 'JWT 认证授权', link: '/cases/jwt-auth' },
      ]},
      { text: '路线图', link: '/roadmap' },
      { text: '贡献指南', link: '/contributing' },
    ],
    sidebar: [
      {
        text: '指南',
        items: [
          { text: '快速开始', link: '/guide/quick-start' },
          { text: '项目架构', link: '/guide/architecture' },
          { text: '技术栈', link: '/guide/tech-stack' },
        ],
      },
      {
        text: '案例',
        items: [
          { text: 'JWT 认证授权', link: '/cases/jwt-auth' },
        ],
      },
      {
        text: '更多',
        items: [
          { text: '路线图', link: '/roadmap' },
          { text: '贡献指南', link: '/contributing' },
        ],
      },
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com' },
    ],
    search: {
      provider: 'local',
    },
    footer: {
      message: '基于 MIT 许可证发布',
      copyright: 'Copyright 2026 全栈工程实验室',
    },
    outline: {
      label: '页面导航',
    },
    docFooter: {
      prev: '上一页',
      next: '下一页',
    },
  },
})
