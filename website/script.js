document.addEventListener('DOMContentLoaded', () => {
    const root = document.documentElement;
    const themeToggle = document.getElementById('theme-toggle');

    const syncThemePressed = () => {
        if (!themeToggle) return;
        themeToggle.setAttribute('aria-pressed', root.classList.contains('dark-mode') ? 'true' : 'false');
    };
    syncThemePressed();

    if (themeToggle) {
        themeToggle.addEventListener('click', () => {
            const isDark = root.classList.toggle('dark-mode');
            try {
                localStorage.setItem('theme', isDark ? 'dark' : 'light');
            } catch (e) {}
            syncThemePressed();
            if (window.Prism) {
                Prism.highlightAll();
            }
        });
    }

    const menuToggle = document.querySelector('.mobile-menu-toggle');
    const mainNav = document.querySelector('.main-nav');

    if (menuToggle && mainNav) {
        menuToggle.addEventListener('click', () => {
            const open = menuToggle.classList.toggle('active');
            mainNav.classList.toggle('active', open);
            document.body.classList.toggle('menu-open', open);
            menuToggle.setAttribute('aria-expanded', open ? 'true' : 'false');
        });
    }

    const navLinks = document.querySelectorAll('.main-nav a');
    for (const link of navLinks) {
        link.addEventListener('click', () => {
            if (menuToggle && mainNav) {
                menuToggle.classList.remove('active');
                mainNav.classList.remove('active');
                document.body.classList.remove('menu-open');
                menuToggle.setAttribute('aria-expanded', 'false');
            }
        });
    }

    const tabButtons = document.querySelectorAll('.tab-button');
    for (const button of tabButtons) {
        button.addEventListener('click', () => {
            const tabId = button.getAttribute('data-tab');
            const tabs = button.closest('.tabs');
            if (!tabs || !tabId) return;
            for (const btn of tabs.querySelectorAll('.tab-button')) {
                btn.classList.remove('active');
            }
            for (const pane of tabs.querySelectorAll('.tab-pane')) {
                pane.classList.remove('active');
            }
            button.classList.add('active');
            const pane = document.getElementById(tabId);
            if (pane) {
                pane.classList.add('active');
                if (window.Prism) Prism.highlightAllUnder(pane);
            }
        });
    }

    const anchorLinks = document.querySelectorAll('a[href^="#"]:not(.annotation-nav a)');
    for (const link of anchorLinks) {
        link.addEventListener('click', (e) => {
            const targetId = link.getAttribute('href');
            if (!targetId || targetId === '#') return;
            const targetElement = document.querySelector(targetId);
            if (!targetElement) return;
            e.preventDefault();
            const header = document.querySelector('header');
            const headerHeight = header ? header.offsetHeight : 0;
            const top = targetElement.getBoundingClientRect().top + window.scrollY - headerHeight;
            const reduce = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
            window.scrollTo({ top, behavior: reduce ? 'auto' : 'smooth' });
        });
    }

    const annotationNavLinks = document.querySelectorAll('.annotation-nav a');
    const annotationSections = document.querySelectorAll('.annotation-section');
    for (const link of annotationNavLinks) {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            for (const item of annotationNavLinks) item.classList.remove('active');
            for (const section of annotationSections) section.classList.remove('active');
            link.classList.add('active');
            const target = document.querySelector(link.getAttribute('href'));
            if (target) target.classList.add('active');
        });
    }

    const copyButtons = document.querySelectorAll('[data-copy]');
    for (const button of copyButtons) {
        button.addEventListener('click', async () => {
            const figure = button.closest('figure');
            const code = figure ? figure.querySelector('code') : null;
            if (!code) return;
            const text = code.textContent || '';
            try {
                await navigator.clipboard.writeText(text);
            } catch (e) {
                return;
            }
            const previous = button.textContent;
            button.textContent = button.lang === 'zh-CN' || document.documentElement.lang === 'zh-CN' ? '已复制' : 'Copied';
            window.setTimeout(() => {
                button.textContent = previous;
            }, 1400);
        });
    }
});
