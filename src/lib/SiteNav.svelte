<script>
  import { onMount } from 'svelte';

  let { active = 'maker', dismissed = $bindable(false) } = $props();
  let compact = $state(false);

  const links = [
    { id: 'maker', href: '/', label: 'Maker' },
    { id: 'pricing', href: '/pricing/', label: 'Pricing' },
    { id: 'reference', href: '/reference/', label: 'API Reference' }
  ];

  onMount(() => {
    let previousY = window.scrollY;
    const handleScroll = () => {
      const currentY = window.scrollY;
      compact = !dismissed && currentY > 48 && currentY > previousY;
      if (currentY < 20 || currentY < previousY) compact = false;
      previousY = currentY;
    };
    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  });

  function dismissBanner() {
    dismissed = true;
    compact = false;
  }
</script>

<a class="skip-link" href="#top">Skip to main content</a>
<header class:compact class:dismissed class="site-nav">
  <div class="announcement">
    <p>Every plan is $0. Finance has stopped returning our calls. <a href="/pricing/">See pricing</a></p>
    <button class="announcement__dismiss" type="button" aria-label="Dismiss announcement" onclick={dismissBanner}>×</button>
  </div>
  <div class="nav-bar">
    <a class="brand" href="/" aria-label="Tiny Badge home">
      tiny<span class="brand-mark" aria-hidden="true"></span>badge
    </a>
    <nav class="nav-links" aria-label="Primary navigation">
      {#each links as link}
        <a href={link.href} aria-current={active === link.id ? 'page' : undefined}>{link.label}</a>
      {/each}
    </nav>
    <a class="nav-action" href="/#designer">Make one</a>
  </div>
</header>
