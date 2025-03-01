document.addEventListener('DOMContentLoaded', () => {
    // Mobile menu toggle
    const menuToggle = document.querySelector('.menu-toggle');
    const nav = document.querySelector('nav');
    
    if (menuToggle) {
        menuToggle.addEventListener('click', () => {
            menuToggle.classList.toggle('active');
            nav.classList.toggle('active');
            document.body.classList.toggle('menu-open');
        });
    }
    
    // Add classes for fade animations
    const fadeElements = document.querySelectorAll('.hero-content, .hero-code, .section-header, .feature-card, .workflow-item, .doc-card');
    for (const el of fadeElements) {
        el.classList.add('fadeUp');
    }
    
    // Code tabs functionality
    const tabButtons = document.querySelectorAll('.tab-btn');
    for (const button of tabButtons) {
        button.addEventListener('click', () => {
            // Get the tab ID from data attribute
            const tabId = button.getAttribute('data-tab');
            
            // Remove active class from all buttons and panes
            const allButtons = document.querySelectorAll('.tab-btn');
            for (const btn of allButtons) {
                btn.classList.remove('active');
            }
            
            const allPanes = document.querySelectorAll('.tab-pane');
            for (const pane of allPanes) {
                pane.classList.remove('active');
            }
            
            // Add active class to clicked button and corresponding tab
            button.classList.add('active');
            document.getElementById(tabId).classList.add('active');
        });
    }
    
    // Smooth scrolling for anchor links
    const anchorLinks = document.querySelectorAll('a[href^="#"]');
    
    for (const link of anchorLinks) {
        link.addEventListener('click', (e) => {
            const targetId = link.getAttribute('href');
            
            // Skip if href is just "#"
            if (targetId === '#') return;
            
            const targetElement = document.getElementById(targetId.substring(1));
            
            if (targetElement) {
                e.preventDefault();
                
                // Close mobile menu if it's open
                if (menuToggle?.classList.contains('active')) {
                    menuToggle.classList.remove('active');
                    nav.classList.remove('active');
                    document.body.classList.remove('menu-open');
                }
                
                // Get header height for offset
                const headerHeight = document.querySelector('header').offsetHeight;
                
                // Calculate position with offset
                const targetPosition = targetElement.getBoundingClientRect().top + window.scrollY - headerHeight;
                
                // Scroll to target
                window.scrollTo({
                    top: targetPosition,
                    behavior: 'smooth'
                });
            }
        });
    }
    
    // Animation on scroll
    const animateOnScroll = () => {
        const animatedElements = document.querySelectorAll('.feature-card, .workflow-item, .doc-card, .installation');
        
        for (const el of animatedElements) {
            const elementTop = el.getBoundingClientRect().top;
            const windowHeight = window.innerHeight;
            
            // If element is in viewport (with offset)
            if (elementTop < windowHeight * 0.85) {
                el.style.opacity = '1';
                el.style.transform = 'translateY(0)';
            }
        }
    };
    
    // Add initial styles
    const animatedElements = document.querySelectorAll('.feature-card, .workflow-item, .doc-card, .installation');
    for (const el of animatedElements) {
        el.style.opacity = '0';
        el.style.transform = 'translateY(20px)';
        el.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
    }
    
    // Run on scroll
    window.addEventListener('scroll', animateOnScroll);
    
    // Run on initial load
    animateOnScroll();
    
    // Handle staggered animations for cards
    const staggerCards = () => {
        const cards = document.querySelectorAll('.feature-card');
        let delay = 0;
        
        for (const card of cards) {
            card.style.transitionDelay = `${delay}s`;
            delay += 0.1;
        }
    };
    
    staggerCards();
    
    // Parallax effect for blob backgrounds
    const parallaxBlobs = () => {
        const blob1 = document.querySelector('.blob-1');
        const blob2 = document.querySelector('.blob-2');
        
        if (blob1 && blob2) {
            window.addEventListener('mousemove', (e) => {
                const moveX = (e.clientX - window.innerWidth / 2) * 0.01;
                const moveY = (e.clientY - window.innerHeight / 2) * 0.01;
                
                blob1.style.transform = `translate(${moveX}px, ${moveY}px)`;
                blob2.style.transform = `translate(${-moveX}px, ${-moveY}px)`;
            });
        }
    };
    
    parallaxBlobs();
    
    // Add animation to gradient background
    const animateGradient = () => {
        const svg = document.querySelector('.gradient-bg svg');
        if (svg) {
            svg.innerHTML += `
                <filter id="turbulence" x="0" y="0" width="100%" height="100%">
                    <feTurbulence type="fractalNoise" baseFrequency="0.01 0.02" numOctaves="2" seed="0" stitchTiles="stitch" result="turbulence"/>
                    <feDisplacementMap in="SourceGraphic" in2="turbulence" scale="30" xChannelSelector="R" yChannelSelector="G"/>
                </filter>
                <g filter="url(#turbulence)">
                    <rect width="100%" height="100%" fill="url(#a)" opacity="0.8"/>
                </g>
            `;
        }
    };
    
    // Only animate gradient on powerful devices
    if (!window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
        animateGradient();
    }
}); 