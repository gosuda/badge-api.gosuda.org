<script>
  import { onMount } from 'svelte';
  import ColorEditor from '$lib/ColorEditor.svelte';
  import { rgbToHex } from '$lib/color.js';
  import SiteFooter from '$lib/SiteFooter.svelte';
  import SiteNav from '$lib/SiteNav.svelte';
  import StylePicker from '$lib/StylePicker.svelte';

  const minSize = 50;
  const maxSize = 300;
  const styleHeights = {
    flat: 20,
    'flat-square': 20,
    plastic: 18,
    round: 24,
    outline: 22,
    neon: 24,
    glass: 24,
    flatbar: 28
  };
  const defaultColors = {
    labelColor: { r: 41, g: 39, b: 36 },
    color: { r: 214, g: 239, b: 83 },
    labelTextColor: { r: 255, g: 255, b: 255 },
    textColor: { r: 41, g: 39, b: 36 }
  };

  let label = $state('tiny');
  let message = $state('but mighty');
  let style = $state('flat');
  let size = $state(100);
  let sizeText = $state('100');
  let sizeInputFocused = $state(false);
  let colors = $state(copyColors(defaultColors));
  let requestState = $state({
    label: 'tiny',
    message: 'but mighty',
    style: 'flat',
    size: 100,
    colors: copyColors(defaultColors)
  });
  let previewZoom = $state(2);
  let previewPending = $state(false);
  let origin = $state('');
  let copyState = $state('idle');
  let bannerDismissed = $state(false);
  let burst = $state(false);

  const previewHeight = $derived((styleHeights[requestState.style] * requestState.size) / 100);
  const exactSizeValid = $derived(Number.isInteger(Number(sizeText)) && Number(sizeText) >= minSize && Number(sizeText) <= maxSize);

  $effect(() => {
    if (!sizeInputFocused) sizeText = String(size);
  });

  $effect(() => {
    const next = { label, message, style, size, colors: copyColors(colors) };
    previewPending = true;
    const timeout = setTimeout(() => {
      requestState = next;
      previewPending = false;
    }, 220);
    return () => clearTimeout(timeout);
  });

  const badgePath = $derived(buildBadgePath({ label, message, style, size, colors }));
  const previewPath = $derived(buildBadgePath(requestState));
  const badgeURL = $derived(origin ? `${origin}${badgePath}` : badgePath);

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

  function formatPixels(value) {
    return Number.isInteger(value) ? String(value) : value.toFixed(1);
  }

  async function copyURL() {
    if (copyState === 'loading') return;
    copyState = 'loading';
    try {
      await navigator.clipboard.writeText(badgeURL);
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

<svelte:head>
  <title>Tiny Badge — Small signs, good manners</title>
  <meta
    name="description"
    content="Make a small, colorful badge and send it off into the world."
  />
  <link rel="icon" href="/favicon.svg" />
</svelte:head>

<SiteNav active="maker" bind:dismissed={bannerDismissed} />

<main id="top" class:banner-dismissed={bannerDismissed}>
  <section class="hero section-shell">
    <div class="hero__copy">
      <h1>A tiny badge<br />with a lot to say.</h1>
      <p>
        Pick the words, fuss over the colors, and leave with a little link that is ready to wander.
      </p>
      <div class="hero__actions">
        <a class="btn btn--push" href="#designer">Make a badge</a>
        <a class="text-link" href="#ideas">See where it can go</a>
      </div>
    </div>
    <div class="hero__proof" aria-label="A few cheerful facts">
      <div class="proof-item proof-item--pear">
        <strong>8</strong>
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
        <h2>Make it yours</h2>
        <p>Tweak a word. Nudge a color. Find the size that fits. Change your mind twice.</p>
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
              <p class="field-help">A short hello on the left.</p>
            </div>
            <div class="field">
              <label class="field-label" for="message">Main bit</label>
              <input id="message" class="text-input" bind:value={message} maxlength="128" required placeholder="but mighty" />
              <p class="field-help">The good part goes on the right.</p>
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
            <h3>Size</h3>
          </div>
          <p id="size-help" class="control-intro">Slide in broad strokes, then type the exact percentage when you want to be picky.</p>
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
            <span class:pending={previewPending} class="status-dot">{previewPending ? 'Updating…' : 'Live preview'}</span>
          </div>
          <div class="preview-stage">
            <div class="preview-zoom-switch" role="group" aria-label="Preview magnification">
              <button class:active={previewZoom === 1} type="button" aria-pressed={previewZoom === 1} onclick={() => (previewZoom = 1)}>1×</button>
              <button class:active={previewZoom === 2} type="button" aria-pressed={previewZoom === 2} onclick={() => (previewZoom = 2)}>2×</button>
            </div>
            {#key previewPath}
              <img class="badge-preview" src={previewPath} alt={`${requestState.label}: ${requestState.message}`} style={`--preview-zoom: ${previewZoom}`} />
            {/key}
          </div>
          <div class="preview-meta">
            <span>{requestState.size}% scale</span>
            <span>{formatPixels(previewHeight)} px tall</span>
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
          <div class="url-output">
            <p>Your badge is packed and ready to go.</p>
          </div>
          <button
            class:success={copyState === 'success'}
            class:error={copyState === 'error'}
            class:loading={copyState === 'loading'}
            class="btn btn--copy"
            type="button"
            disabled={copyState === 'loading'}
            onclick={copyURL}
          >
            {#if copyState === 'loading'}
              Picking it up…
            {:else if copyState === 'success'}
              All yours
            {:else if copyState === 'error'}
              Try that once more
            {:else}
              Copy my link
            {/if}
            {#if burst}<span class="star-burst" aria-hidden="true"></span>{/if}
          </button>
          <p class="cache-note">Paste it wherever your badge deserves a tiny bit of attention.</p>
        </div>
      </aside>
    </div>
  </section>

  <section id="ideas" class="api-section section-shell">
    <div class="api-section__copy">
      <h2>Small badge. Plenty of places to go.</h2>
      <p>
        Tuck it into the corners that could use a little clarity, color, or cheerful showing off.
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
