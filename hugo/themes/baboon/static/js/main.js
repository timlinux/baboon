// Baboon Documentation - Main JavaScript

// Theme toggle
function toggleTheme() {
  const html = document.documentElement;
  const currentTheme = html.getAttribute('data-theme');
  const newTheme = currentTheme === 'light' ? 'dark' : 'light';

  html.setAttribute('data-theme', newTheme);
  localStorage.setItem('baboon-theme', newTheme);
}

// Initialize theme from localStorage or system preference
function initTheme() {
  const savedTheme = localStorage.getItem('baboon-theme');
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  const theme = savedTheme || (prefersDark ? 'dark' : 'light');

  document.documentElement.setAttribute('data-theme', theme);
}

// Mobile menu toggle
function toggleMobileMenu() {
  const mobileNav = document.getElementById('mobile-nav');
  mobileNav.classList.toggle('active');
}

// Smooth scroll for anchor links
function initSmoothScroll() {
  document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
      const href = this.getAttribute('href');
      if (href === '#') return;

      const target = document.querySelector(href);
      if (target) {
        e.preventDefault();
        target.scrollIntoView({
          behavior: 'smooth',
          block: 'start'
        });

        // Update URL without scrolling
        history.pushState(null, null, href);
      }
    });
  });
}

// Add copy buttons to code blocks
function initCodeCopy() {
  document.querySelectorAll('pre code').forEach((codeBlock) => {
    const pre = codeBlock.parentElement;
    const wrapper = document.createElement('div');
    wrapper.style.position = 'relative';
    pre.parentNode.insertBefore(wrapper, pre);
    wrapper.appendChild(pre);

    const copyButton = document.createElement('button');
    copyButton.className = 'copy-button';
    copyButton.textContent = 'Copy';
    copyButton.style.cssText = `
      position: absolute;
      top: 0.5rem;
      right: 0.5rem;
      padding: 0.25rem 0.75rem;
      background: var(--bg-tertiary);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 0.375rem;
      color: var(--text-muted);
      font-size: 0.75rem;
      cursor: pointer;
      opacity: 0;
      transition: opacity 0.2s ease;
    `;

    wrapper.appendChild(copyButton);

    wrapper.addEventListener('mouseenter', () => {
      copyButton.style.opacity = '1';
    });

    wrapper.addEventListener('mouseleave', () => {
      copyButton.style.opacity = '0';
    });

    copyButton.addEventListener('click', async () => {
      try {
        await navigator.clipboard.writeText(codeBlock.textContent);
        copyButton.textContent = 'Copied!';
        copyButton.style.color = 'var(--accent-green)';
        setTimeout(() => {
          copyButton.textContent = 'Copy';
          copyButton.style.color = 'var(--text-muted)';
        }, 2000);
      } catch (err) {
        copyButton.textContent = 'Failed';
        setTimeout(() => {
          copyButton.textContent = 'Copy';
        }, 2000);
      }
    });
  });
}

// Highlight current section in TOC
function initTocHighlight() {
  const toc = document.querySelector('.toc');
  if (!toc) return;

  const headings = document.querySelectorAll('h2[id], h3[id]');
  const tocLinks = toc.querySelectorAll('a');

  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const id = entry.target.getAttribute('id');
        tocLinks.forEach(link => {
          link.classList.remove('active');
          if (link.getAttribute('href') === `#${id}`) {
            link.classList.add('active');
          }
        });
      }
    });
  }, {
    rootMargin: '-100px 0px -80% 0px'
  });

  headings.forEach(heading => observer.observe(heading));
}

// Initialize all functionality
document.addEventListener('DOMContentLoaded', () => {
  initTheme();
  initSmoothScroll();
  initCodeCopy();
  initTocHighlight();
});

// Listen for system theme changes
window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
  if (!localStorage.getItem('baboon-theme')) {
    document.documentElement.setAttribute('data-theme', e.matches ? 'dark' : 'light');
  }
});
