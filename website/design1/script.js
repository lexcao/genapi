document.addEventListener('DOMContentLoaded', () => {
    // Code tabs functionality
    const tabButtons = document.querySelectorAll('.tab-btn');
    
    for (const button of tabButtons) {
        button.addEventListener('click', function() {
            // Remove active class from all buttons and panes
            const allButtons = document.querySelectorAll('.tab-btn');
            for (const btn of allButtons) {
                btn.classList.remove('active');
            }
            
            const allPanes = document.querySelectorAll('.tab-pane');
            for (const pane of allPanes) {
                pane.classList.remove('active');
            }
            
            // Add active class to clicked button
            this.classList.add('active');
            
            // Show the corresponding tab pane
            const tabId = this.getAttribute('data-tab');
            document.getElementById(tabId).classList.add('active');
        });
    }

    // Smooth scrolling for anchor links
    const anchors = document.querySelectorAll('a[href^="#"]');
    for (const anchor of anchors) {
        anchor.addEventListener('click', function(e) {
            e.preventDefault();
            
            const targetId = this.getAttribute('href');
            if (targetId === '#') return;
            
            const targetElement = document.querySelector(targetId);
            if (targetElement) {
                window.scrollTo({
                    top: targetElement.offsetTop - 80, // Offset for fixed header
                    behavior: 'smooth'
                });
            }
        });
    }

    // Add animation on scroll
    const animateOnScroll = () => {
        const sections = document.querySelectorAll('section');
        
        for (const section of sections) {
            const sectionTop = section.getBoundingClientRect().top;
            const windowHeight = window.innerHeight;
            
            if (sectionTop < windowHeight * 0.75) {
                section.classList.add('animate');
            }
        }
    };
    
    // Run on load
    animateOnScroll();
    
    // Run on scroll
    window.addEventListener('scroll', animateOnScroll);

    // Mobile menu toggle
    const mobileBreakpoint = 768;
    
    const checkScreenSize = () => {
        if (window.innerWidth <= mobileBreakpoint) {
            const nav = document.querySelector('nav');
            const navList = nav.querySelector('ul');
            
            if (!nav.querySelector('.mobile-toggle')) {
                // Create mobile toggle button
                const toggleBtn = document.createElement('button');
                toggleBtn.classList.add('mobile-toggle');
                toggleBtn.innerHTML = '☰';
                nav.insertBefore(toggleBtn, navList);
                
                // Add toggle functionality
                toggleBtn.addEventListener('click', () => {
                    navList.classList.toggle('active');
                    toggleBtn.classList.toggle('active');
                    toggleBtn.innerHTML = toggleBtn.classList.contains('active') ? '✕' : '☰';
                });
            }
        }
    };
    
    // Run on load and resize
    checkScreenSize();
    window.addEventListener('resize', checkScreenSize);
}); 