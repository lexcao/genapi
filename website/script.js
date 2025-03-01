document.addEventListener('DOMContentLoaded', () => {
    // Theme toggling
    const themeToggle = document.getElementById('theme-toggle');
    const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');
    
    // Function to add/remove PrismJS dark theme
    const updatePrismTheme = (isDark) => {
        // For custom styling, we rely on the .dark-mode class and our CSS overrides
        // This provides a smooth transition between themes
        if (isDark) {
            document.documentElement.setAttribute('data-theme', 'dark');
        } else {
            document.documentElement.setAttribute('data-theme', 'light');
        }
        
        // Re-highlight all code blocks to apply new theme
        if (window.Prism) {
            Prism.highlightAll();
        }
    };
    
    // Check for saved theme preference or use the system preference
    const savedTheme = localStorage.getItem('theme');
    
    if (savedTheme === 'dark' || (!savedTheme && prefersDarkScheme.matches)) {
        document.body.classList.add('dark-mode');
        updatePrismTheme(true);
    }
    
    // Toggle theme when button is clicked
    themeToggle.addEventListener('click', () => {
        const isDarkMode = document.body.classList.toggle('dark-mode');
        updatePrismTheme(isDarkMode);
        
        // Save preference to localStorage
        localStorage.setItem('theme', isDarkMode ? 'dark' : 'light');
    });
    
    // Mobile menu toggle
    const menuToggle = document.querySelector('.mobile-menu-toggle');
    const mainNav = document.querySelector('.main-nav');
    
    if (menuToggle && mainNav) {
        menuToggle.addEventListener('click', () => {
            menuToggle.classList.toggle('active');
            mainNav.classList.toggle('active');
            document.body.classList.toggle('menu-open');
        });
    }
    
    // Close mobile menu when clicking on a link
    const navLinks = document.querySelectorAll('.main-nav a');
    
    for (const link of navLinks) {
        link.addEventListener('click', () => {
            if (menuToggle && mainNav) {
                menuToggle.classList.remove('active');
                mainNav.classList.remove('active');
                document.body.classList.remove('menu-open');
            }
        });
    }
    
    // Tab functionality
    const tabButtons = document.querySelectorAll('.tab-button');
    
    for (const button of tabButtons) {
        button.addEventListener('click', () => {
            // Get the tab ID from the data attribute
            const tabId = button.getAttribute('data-tab');
            
            // Remove active class from all buttons and panes
            const allButtons = document.querySelectorAll('.tab-button');
            for (const btn of allButtons) {
                btn.classList.remove('active');
            }
            
            const allPanes = document.querySelectorAll('.tab-pane');
            for (const pane of allPanes) {
                pane.classList.remove('active');
            }
            
            // Add active class to current button and pane
            button.classList.add('active');
            document.getElementById(tabId).classList.add('active');
        });
    }
    
    // Smooth scrolling for anchor links
    const anchorLinks = document.querySelectorAll('a[href^="#"]:not(.annotation-nav a)');
    
    for (const link of anchorLinks) {
        link.addEventListener('click', (e) => {
            const targetId = link.getAttribute('href');
            
            // Skip if the href is just "#"
            if (targetId === '#') return;
            
            const targetElement = document.querySelector(targetId);
            
            if (targetElement) {
                e.preventDefault();
                
                const headerHeight = document.querySelector('header').offsetHeight;
                const targetPosition = targetElement.getBoundingClientRect().top + window.scrollY - headerHeight;
                
                window.scrollTo({
                    top: targetPosition,
                    behavior: 'smooth'
                });
            }
        });
    }
    
    // Annotation navigation functionality
    const annotationNavLinks = document.querySelectorAll('.annotation-nav a');
    const annotationSections = document.querySelectorAll('.annotation-section');
    
    for (const link of annotationNavLinks) {
        link.addEventListener('click', function (e) {
            e.preventDefault();

            // Remove active class from all links and sections
            for (const l of annotationNavLinks) {
                l.classList.remove('active');
            }
            
            for (const s of annotationSections) {
                s.classList.remove('active');
            }

            // Add active class to clicked link
            this.classList.add('active');

            // Show corresponding section
            const targetId = this.getAttribute('href');
            document.querySelector(targetId).classList.add('active');
        });
    }
    
    // Scroll-triggered animations
    const animateElements = () => {
        const elements = document.querySelectorAll('.card, .step, .annotation-card, .getting-started-step');
        
        for (const el of elements) {
            const rect = el.getBoundingClientRect();
            const windowHeight = window.innerHeight;
            
            // If element is in viewport
            if (rect.top <= windowHeight * 0.8 && rect.bottom >= 0) {
                el.classList.add('animate-in');
            }
        }
    };
    
    // Run on initial load
    animateElements();
    
    // Run on scroll
    window.addEventListener('scroll', animateElements);
    
    // Add some style for the animation
    const style = document.createElement('style');
    style.textContent = `
        .card, .step, .annotation-card, .getting-started-step {
            opacity: 0;
            transform: translateY(20px);
            transition: opacity 0.6s ease, transform 0.6s ease;
        }
        
        .animate-in {
            opacity: 1;
            transform: translateY(0);
        }
        
        @media (prefers-reduced-motion: reduce) {
            .card, .step, .annotation-card, .getting-started-step {
                opacity: 1;
                transform: none;
                transition: none;
            }
        }
    `;
    document.head.appendChild(style);
}); 