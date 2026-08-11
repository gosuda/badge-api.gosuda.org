<script>
  import { onMount } from 'svelte';
  import ColorEditor from '$lib/ColorEditor.svelte';
  import { rgbToHex } from '$lib/color.js';
  import Seo from '$lib/Seo.svelte';
  import SiteFooter from '$lib/SiteFooter.svelte';
  import SiteNav from '$lib/SiteNav.svelte';
  import StylePicker from '$lib/StylePicker.svelte';

  const minSize = 50;
  const maxSize = 300;
  const minLetterSpacing = -1;
  const maxLetterSpacing = 3;
  const styleHeights = {
    flat: 20,
    'flat-square': 20,
    plastic: 18,
    round: 24,
    outline: 22,
    neon: 24,
    glass: 24,
    flatbar: 28,
    'old-school': 15,
    'click-here': 31,
    'best-viewed': 31
  };
  const defaultColors = {
    labelColor: { r: 41, g: 39, b: 36 },
    color: { r: 214, g: 239, b: 83 },
    labelTextColor: { r: 255, g: 255, b: 255 },
    textColor: { r: 41, g: 39, b: 36 }
  };

  const seoDescription = 'Create customizable SVG badges with eleven styles, including classic 80×15 and 88×31 web buttons with auto-scaling vector text, exact 50–300% sizing, adjustable letter spacing, RGB, HSL, and OKLCH controls, Unicode-aware width, immutable caching, and shareable URLs.';
  const homeStructuredData = {
    '@context': 'https://schema.org',
    '@graph': [
      {
        '@type': 'Organization',
        '@id': 'https://badge-api.gosuda.org/#organization',
        name: 'Tiny Badge',
        url: 'https://badge-api.gosuda.org/',
        logo: 'https://badge-api.gosuda.org/favicon.svg'
      },
      {
        '@type': 'WebSite',
        '@id': 'https://badge-api.gosuda.org/#website',
        url: 'https://badge-api.gosuda.org/',
        name: 'Tiny Badge',
        description: seoDescription,
        inLanguage: 'en',
        publisher: { '@id': 'https://badge-api.gosuda.org/#organization' }
      },
      {
        '@type': 'WebApplication',
        '@id': 'https://badge-api.gosuda.org/#application',
        name: 'Tiny Badge',
        url: 'https://badge-api.gosuda.org/',
        description: seoDescription,
        applicationCategory: 'DeveloperApplication',
        operatingSystem: 'Any',
        browserRequirements: 'A modern web browser with JavaScript for the visual maker; generated SVG URLs work as ordinary images.',
        offers: {
          '@type': 'Offer',
          price: '0',
          priceCurrency: 'USD'
        },
        featureList: [
          'Eleven SVG badge styles, including fixed 80 by 15 and 88 by 31 pixel buttons',
          'Exact size scaling from 50 to 300 percent',
          'Adjustable letter spacing from -1 to 3 SVG pixels',
          'Hex, RGB, HSL, and OKLCH color controls',
          'Unicode-aware text measurement',
          'Immutable badge URLs with ETags'
        ],
        publisher: { '@id': 'https://badge-api.gosuda.org/#organization' }
      },
      {
        '@type': 'WebPage',
        '@id': 'https://badge-api.gosuda.org/#webpage',
        url: 'https://badge-api.gosuda.org/',
        name: 'SVG Badge Generator & API — Tiny Badge',
        description: seoDescription,
        isPartOf: { '@id': 'https://badge-api.gosuda.org/#website' },
        about: { '@id': 'https://badge-api.gosuda.org/#application' },
        primaryImageOfPage: {
          '@type': 'ImageObject',
          url: 'https://badge-api.gosuda.org/og-image.png',
          width: 1200,
          height: 630
        },
        inLanguage: 'en'
      }
    ]
  };

  let label = $state('tiny');
  let message = $state('but mighty');
  let style = $state('flat');
  let size = $state(100);
  let sizeText = $state('100');
  let sizeInputFocused = $state(false);
  let letterSpacing = $state(0);
  let letterSpacingText = $state('0');
  let letterSpacingInputFocused = $state(false);
  let colors = $state(copyColors(defaultColors));
  let requestState = $state({
    label: 'tiny',
    message: 'but mighty',
    style: 'flat',
    size: 100,
    letterSpacing: 0,
    colors: copyColors(defaultColors)
  });
  let previewZoom = $state(2);
  let previewPending = $state(false);
  let origin = $state('');
  let copyState = $state('idle');
  let copyFormat = $state('url');
  let bannerDismissed = $state(false);
  let burst = $state(false);

  const previewHeight = $derived((styleHeights[requestState.style] * requestState.size) / 100);
  const exactSizeValid = $derived(Number.isInteger(Number(sizeText)) && Number(sizeText) >= minSize && Number(sizeText) <= maxSize);
  const exactLetterSpacingValid = $derived(
    letterSpacingText.trim() !== '' &&
      Number.isFinite(Number(letterSpacingText)) &&
      Number(letterSpacingText) >= minLetterSpacing &&
      Number(letterSpacingText) <= maxLetterSpacing
  );

  $effect(() => {
    if (!sizeInputFocused) sizeText = String(size);
  });
  $effect(() => {
    if (!letterSpacingInputFocused) letterSpacingText = formatLetterSpacing(letterSpacing);
  });


  $effect(() => {
    const next = { label, message, style, size, letterSpacing, colors: copyColors(colors) };
    previewPending = true;
    const timeout = setTimeout(() => {
      requestState = next;
      previewPending = false;
    }, 220);
    return () => clearTimeout(timeout);
  });

  const badgePath = $derived(buildBadgePath({ label, message, style, size, letterSpacing, colors }));
  const previewPath = $derived(buildBadgePath(requestState));
  const badgeURL = $derived(origin ? `${origin}${badgePath}` : badgePath);
  const markdownEmbed = $derived(`![${markdownAlt(label, message)}](${badgeURL})`);

  onMount(() => {
    origin = window.location.origin;
  });

  function copyColors(source) {
    return {
      labelColor: { ...source.labelColor },
      color: { ...source.color },
      labelTextColor: { ...source.labelTextColor },
      textColor: { ...source.textColor }
    };
  }

  function buildBadgePath(options) {
    const query = new URLSearchParams({
      label: options.label,
      message: options.message,
      style: options.style,
      size: String(options.size),
      letterSpacing: String(options.letterSpacing),
      labelColor: rgbToHex(options.colors.labelColor).slice(1),
      color: rgbToHex(options.colors.color).slice(1),
      labelTextColor: rgbToHex(options.colors.labelTextColor).slice(1),
      textColor: rgbToHex(options.colors.textColor).slice(1)
    });
    return `/badge.svg?${query}`;
  }


  function setColor(key, value) {
    colors[key] = value;
  }

  function setSizeFromRange(event) {
    size = Number(event.currentTarget.value);
  }

  function editExactSize(event) {
    sizeText = event.currentTarget.value;
    const next = Number(sizeText);
    if (Number.isInteger(next) && next >= minSize && next <= maxSize) size = next;
  }

  function blurExactSize() {
    sizeInputFocused = false;
    sizeText = String(size);
  }
  function setLetterSpacingFromRange(event) {
    letterSpacing = Number(event.currentTarget.value);
  }

  function editExactLetterSpacing(event) {
    letterSpacingText = event.currentTarget.value;
    const next = Number(letterSpacingText);
    if (letterSpacingText.trim() !== '' && Number.isFinite(next) && next >= minLetterSpacing && next <= maxLetterSpacing) {
      letterSpacing = next;
    }
  }

  function blurExactLetterSpacing() {
    letterSpacingInputFocused = false;
    letterSpacingText = formatLetterSpacing(letterSpacing);
  }

  function formatLetterSpacing(value) {
    return String(Number(value.toFixed(2)));
  }


  function formatPixels(value) {
    return Number.isInteger(value) ? String(value) : value.toFixed(1);
  }

  function markdownAlt(labelText, messageText) {
    const text = [labelText, messageText].filter(Boolean).join(': ') || 'Tiny Badge';
    return text.replace(/([\\\[\]])/g, '\\$1');
  }

  function copyLabel(format, idleLabel) {
    if (copyFormat !== format || copyState === 'idle') return idleLabel;
    if (copyState === 'loading') return 'Copying…';
    if (copyState === 'success') return format === 'markdown' ? 'Markdown copied' : 'URL copied';
    return format === 'markdown' ? 'Try Markdown again' : 'Try URL again';
  }

  async function copyBadge(format) {
    if (copyState === 'loading') return;
    copyFormat = format;
    copyState = 'loading';
    try {
      await navigator.clipboard.writeText(format === 'markdown' ? markdownEmbed : badgeURL);
      copyState = 'success';
      burst = true;
      window.setTimeout(() => (burst = false), 460);
      window.setTimeout(() => (copyState = 'idle'), 2500);
    } catch {
      copyState = 'error';
      window.setTimeout(() => (copyState = 'idle'), 3000);
    }
  }
</script>

<Seo
  title="SVG Badge Generator & API — Tiny Badge"
  description={seoDescription}
  path="/"
  structuredData={homeStructuredData}
/>

<SiteNav active="maker" bind:dismissed={bannerDismissed} />

<main id="top" class:banner-dismissed={bannerDismissed}>
  <section class="hero section-shell">
    <div class="hero__copy">
      <h1>A tiny SVG badge<br />with a lot to say.</h1>
      <p>
        Create a crisp SVG badge from a label, message, style, exact size, letter spacing, and color recipe—then copy one immutable URL for Markdown, project pages, profiles, or docs.
      </p>
      <div class="hero__actions">
        <a class="btn btn--modern" href="#designer">
          <span>Make a badge</span>
          <span class="btn--modern__icon" aria-hidden="true">↗</span>
        </a>
        <a class="text-link" href="#ideas">See where it can go</a>
      </div>
    </div>
    <div class="hero__proof" aria-label="A few cheerful facts">
      <div class="proof-item proof-item--pear">
        <strong>11</strong>
        <span>looks to try on</span>
      </div>
      <div class="proof-item proof-item--cyan">
        <strong>4</strong>
        <span>ways to mix a color</span>
      </div>
      <div class="proof-item proof-item--lavender">
        <strong>1</strong>
        <span>little link to keep</span>
      </div>
    </div>
  </section>

  <section id="designer" class="designer-band">
    <div class="section-shell designer-heading">
      <div>
        <h2>Make your SVG badge</h2>
        <p>Tweak the words, choose one of eleven styles, set an exact size and letter spacing, and tune every color channel.</p>
      </div>
      <code>No wrong turns</code>
    </div>

    <div class="section-shell designer-shell">
      <form class="control-panel" onsubmit={(event) => event.preventDefault()}>
        <section class="control-section">
          <div class="section-title-row">
            <h3>Words</h3>
          </div>
          <div class="field-grid">
            <div class="field">
              <label class="field-label" for="label">First bit</label>
              <input id="label" class="text-input" bind:value={label} maxlength="64" placeholder="tiny" />
              <p class="field-help">A short hello on the left. Sent as <code>label</code>.</p>
            </div>
            <div class="field">
              <label class="field-label" for="message">Main bit</label>
              <input id="message" class="text-input" bind:value={message} maxlength="128" required placeholder="but mighty" />
              <p class="field-help">The good part goes on the right. Sent as <code>message</code>.</p>
            </div>
          </div>
        </section>

        <section id="styles" class="control-section">
          <div class="section-title-row">
            <h3>Outfit</h3>
          </div>
          <StylePicker bind:value={style} />
        </section>

        <section class="control-section">
          <div class="section-title-row">
            <h3>Size &amp; spacing</h3>
          </div>
          <p id="size-help" class="control-intro">Set the overall scale, then give the letters more room—or pull them a little closer.</p>
          <div class="size-control">
            <span id="badge-size-label" class="field-label">Badge scale</span>
            <input
              id="badge-size"
              class="size-slider"
              type="range"
              min={minSize}
              max={maxSize}
              step="5"
              value={size}
              aria-labelledby="badge-size-label"
              aria-describedby="size-help"
              oninput={setSizeFromRange}
            />
            <span class:error={!exactSizeValid} class="size-exact-control">
              <input
                id="badge-size-exact"
                type="number"
                min={minSize}
                max={maxSize}
                step="1"
                value={sizeText}
                aria-label="Exact badge scale percentage"
                aria-invalid={!exactSizeValid}
                aria-describedby="size-help"
                onfocus={() => (sizeInputFocused = true)}
                oninput={editExactSize}
                onblur={blurExactSize}
              />
              <span aria-hidden="true">%</span>
            </span>
            <span class="size-range"><span>{minSize}%</span><span>{maxSize}%</span></span>
          </div>
          <div class="size-control letter-spacing-control">
            <span id="letter-spacing-label" class="field-label">Letter spacing</span>
            <input
              id="letter-spacing"
              class="size-slider"
              type="range"
              min={minLetterSpacing}
              max={maxLetterSpacing}
              step="0.05"
              value={letterSpacing}
              aria-labelledby="letter-spacing-label"
              aria-describedby="letter-spacing-help"
              oninput={setLetterSpacingFromRange}
            />
            <span class:error={!exactLetterSpacingValid} class="size-exact-control">
              <input
                id="letter-spacing-exact"
                type="number"
                min={minLetterSpacing}
                max={maxLetterSpacing}
                step="0.05"
                value={letterSpacingText}
                aria-label="Exact letter spacing in pixels"
                aria-invalid={!exactLetterSpacingValid}
                aria-describedby="letter-spacing-help"
                onfocus={() => (letterSpacingInputFocused = true)}
                oninput={editExactLetterSpacing}
                onblur={blurExactLetterSpacing}
              />
              <span aria-hidden="true">px</span>
            </span>
            <span id="letter-spacing-help" class="size-range"><span>{minLetterSpacing} px</span><span>{maxLetterSpacing} px</span></span>
          </div>
        </section>

        <section class="control-section">
          <div class="section-title-row">
            <h3>Colors</h3>
          </div>
          <p class="control-intro">Pick a recipe and tune each channel separately. The maker quietly converts every result to RGB.</p>
          <div class="color-grid">
            <ColorEditor label="First bit backdrop" value={colors.labelColor} oncolorchange={(value) => setColor('labelColor', value)} />
            <ColorEditor label="Main bit backdrop" value={colors.color} oncolorchange={(value) => setColor('color', value)} />
            <ColorEditor label="First bit words" value={colors.labelTextColor} oncolorchange={(value) => setColor('labelTextColor', value)} />
            <ColorEditor label="Main bit words" value={colors.textColor} oncolorchange={(value) => setColor('textColor', value)} />
          </div>
        </section>
      </form>

      <aside class="preview-column">
        <div class="preview-card">
          <div class="preview-card__head">
            <div>
              <span>Right this minute</span>
              <strong>Your badge</strong>
            </div>
            <span class:pending={previewPending} class="status-dot" role="status" aria-live="polite">{previewPending ? 'Updating…' : 'Live preview'}</span>
          </div>
          <div class="preview-stage">
            <div class="preview-zoom-switch" role="group" aria-label="Preview magnification">
              <button class:active={previewZoom === 1} type="button" aria-pressed={previewZoom === 1} onclick={() => (previewZoom = 1)}>1×</button>
              <button class:active={previewZoom === 2} type="button" aria-pressed={previewZoom === 2} onclick={() => (previewZoom = 2)}>2×</button>
            </div>
            <img class="badge-preview" src={previewPath} alt={`${requestState.label}: ${requestState.message}`} style={`--preview-zoom: ${previewZoom}`} />
          </div>
          <div class="preview-meta">
            <span>{requestState.size}% scale</span>
            <span>{formatPixels(previewHeight)} px tall</span>
            <span>{formatLetterSpacing(requestState.letterSpacing)} px tracking</span>
          </div>
        </div>

        <div class="url-card">
          <div class="url-card__head">
            <div>
              <span>Take it with you</span>
              <strong>Your little link</strong>
            </div>
            <a class="download-link" href={badgePath} download="badge.svg">Save a copy</a>
          </div>
          <a class="url-output" href={badgePath} target="_blank" rel="noreferrer">
            <span class="url-output__preview">
              <img src={previewPath} alt={`${requestState.label}: ${requestState.message}`} />
            </span>
            <span class="url-output__copy">
              <strong>Open the live SVG</strong>
              <code>{badgeURL}</code>
            </span>
            <span class="url-output__arrow" aria-hidden="true">↗</span>
          </a>
          <div class="markdown-output">
            <span>Markdown</span>
            <code>{markdownEmbed}</code>
          </div>
          <div class="copy-actions">
            <button
              class:success={copyState === 'success' && copyFormat === 'url'}
              class:error={copyState === 'error' && copyFormat === 'url'}
              class:loading={copyState === 'loading' && copyFormat === 'url'}
              class="btn btn--copy"
              type="button"
              disabled={copyState === 'loading'}
              onclick={() => copyBadge('url')}
            >
              {copyLabel('url', 'Copy URL')}
              {#if burst && copyFormat === 'url'}<span class="star-burst" aria-hidden="true"></span>{/if}
            </button>
            <button
              class:success={copyState === 'success' && copyFormat === 'markdown'}
              class:error={copyState === 'error' && copyFormat === 'markdown'}
              class:loading={copyState === 'loading' && copyFormat === 'markdown'}
              class="btn btn--copy btn--copy--secondary"
              type="button"
              disabled={copyState === 'loading'}
              onclick={() => copyBadge('markdown')}
            >
              {copyLabel('markdown', 'Copy Markdown')}
              {#if burst && copyFormat === 'markdown'}<span class="star-burst" aria-hidden="true"></span>{/if}
            </button>
          </div>
          <p class="visually-hidden" role="status" aria-live="polite">
            {copyState === 'success' ? `${copyFormat === 'markdown' ? 'Markdown' : 'URL'} copied to clipboard` : copyState === 'error' ? 'Copy failed. Try again.' : ''}
          </p>
          <p class="cache-note">Use the URL anywhere, or paste the Markdown straight into a README.</p>
        </div>
      </aside>
    </div>
  </section>

  <section id="ideas" class="api-section section-shell">
    <div class="api-section__copy">
      <h2>Small SVG badge. Plenty of places to go.</h2>
      <p>
        Tiny Badge is a free SVG badge generator and HTTP API for project pages, profiles, Markdown, documentation, newsletters, and anywhere an image URL can travel.
      </p>
    </div>
    <div class="api-reference">
      <div class="api-row">
        <span>Project pages</span>
        <code>A tiny sign for what is happening now.</code>
      </div>
      <div class="api-row">
        <span>Profiles</span>
        <code>A bit of color without making a speech.</code>
      </div>
      <div class="api-row">
        <span>Notes and newsletters</span>
        <code>Small enough to tuck between two sentences.</code>
      </div>
      <div class="api-row">
        <span>Anywhere with a link</span>
        <code>If there’s room for a link, your badge can come along.</code>
      </div>
      <div class="api-example">
        <span>Our only advice</span>
        <code>Keep it short. Pick colors that can see each other. Have fun.</code>
      </div>
    </div>
  </section>
</main>

<SiteFooter />
