<script>
  import Seo from '$lib/Seo.svelte';
  import CopyTextButton from '$lib/CopyTextButton.svelte';
  import SiteFooter from '$lib/SiteFooter.svelte';
  import SiteNav from '$lib/SiteNav.svelte';

  let bannerDismissed = $state(false);

  const productionOrigin = 'https://badge-api.gosuda.org';
  const queryExample = `${productionOrigin}/badge.svg?label=build&message=passing&style=flatbar&size=125&letterSpacing=0.5`;
  const pathExample = `${productionOrigin}/badge/release/stable.svg`;
  const markdownExample = `![Build status](${productionOrigin}/badge.svg?label=build&message=passing)`;
  const htmlExample = `<img src="${productionOrigin}/badge.svg?label=build&message=passing" alt="Build: passing">`;
  const curlExample = `curl -o badge.svg '${productionOrigin}/badge.svg?label=build&message=passing&style=flatbar'`;
  const responseHeaders = `Content-Type: image/svg+xml; charset=utf-8
Cache-Control: public, max-age=315360000, immutable
CDN-Cache-Control: public, max-age=315360000, immutable
Surrogate-Control: public, max-age=315360000, immutable
Access-Control-Allow-Origin: *
ETag: "<content hash>"`;

  const parameters = [
    { name: 'message', required: 'Required', fallback: '—', detail: 'Main badge text. Maximum 128 characters.' },
    { name: 'label', required: 'Optional', fallback: 'Empty', detail: 'Left-side badge text. Maximum 64 characters.' },
    { name: 'style', required: 'Optional', fallback: 'flat', detail: 'Selects one of the eleven rendering styles.' },
    { name: 'size', required: 'Optional', fallback: '100', detail: 'Scales the whole badge from 50 to 300 percent. Integer values only.' },
    { name: 'letterSpacing', required: 'Optional', fallback: '0', detail: 'Adds -1 to 3 SVG pixels between grapheme clusters. Decimal values are accepted.' },
    { name: 'labelColor', required: 'Optional', fallback: '555555', detail: 'Background color for the label segment.' },
    { name: 'color', required: 'Optional', fallback: '44cc11', detail: 'Background color for the message segment.' },
    { name: 'labelTextColor', required: 'Optional', fallback: 'ffffff', detail: 'Text color for the label segment.' },
    { name: 'textColor', required: 'Optional', fallback: 'ffffff', detail: 'Text color for the message segment.' }
  ];

  const styles = [
    ['flat', '20 px, gently rounded, and the default.'],
    ['flat-square', 'The flat style with square corners.'],
    ['plastic', 'A compact badge with a highlight layer.'],
    ['round', 'A taller capsule with fully rounded ends.'],
    ['outline', 'A framed surface with a brighter inner edge.'],
    ['neon', 'A focused glow around each color segment.'],
    ['glass', 'A layered sheen with a quiet border.'],
    ['flatbar', 'A 28 px uppercase bar with relaxed 0.45 px built-in tracking.'],
    ['old-school', 'Old School 80×15: a fixed split-panel button with customizable colors and resolution-independent text that automatically shrinks to fit each panel.'],
    ['click-here', 'Click Here 88×31: the supplied raised gray artwork for its reference phrase, with high-resolution auto-fit vector text for custom copy.'],
    ['best-viewed', 'Best Viewed 88×31: the supplied BEST rail and Chrome artwork, with separately auto-scaled vector text for both custom lines.']
  ];

  function styleSampleURL(style) {
    const presets = {
      'old-school': { label: 'pixel', message: 'button', size: '200', labelColor: 'ff5a18', color: 'a8a979' },
      'click-here': { label: 'click', message: 'here', size: '200', labelColor: '555555', color: '44cc11' },
      'best-viewed': { label: 'viewed with', message: 'chrome', size: '200', labelColor: '555555', color: '44cc11' }
    };
    const preset = presets[style[0]] ?? { label: 'style', message: style[0], size: '100', labelColor: '555555', color: '44cc11' };
    const query = new URLSearchParams({
      label: preset.label,
      message: preset.message,
      style: style[0],
      size: preset.size,
      labelColor: preset.labelColor,
      color: preset.color,
      labelTextColor: 'ffffff',
      textColor: 'ffffff'
    });
    return `/badge.svg?${query}`;
  }

  const namedColors = ['brightgreen', 'green', 'yellowgreen', 'yellow', 'orange', 'red', 'blue', 'grey', 'gray', 'lightgrey', 'lightgray', 'success', 'important', 'critical', 'informational', 'inactive'];

  const sections = [
    ['overview', 'Overview'],
    ['endpoints', 'Endpoints'],
    ['parameters', 'Parameters'],
    ['styles', 'Styles'],
    ['colors', 'Colors'],
    ['examples', 'Examples'],
    ['responses', 'Responses'],
    ['errors', 'Errors']
  ];

  const seoDescription = 'Read the Tiny Badge SVG API reference for endpoints, parameters, eleven styles including classic 80×15 and 88×31 buttons with auto-scaling vector text, colors, 50–300% sizing, adjustable letter spacing, Unicode width, ETags, immutable caching, examples, and errors.';
  const referenceStructuredData = {
    '@context': 'https://schema.org',
    '@graph': [
      {
        '@type': 'TechArticle',
        '@id': 'https://badge-api.gosuda.org/reference/#article',
        headline: 'Tiny Badge SVG Badge API Reference',
        description: seoDescription,
        url: 'https://badge-api.gosuda.org/reference/',
        mainEntityOfPage: 'https://badge-api.gosuda.org/reference/',
        dateModified: '2026-08-11',
        inLanguage: 'en',
        author: { '@id': 'https://badge-api.gosuda.org/#organization' },
        publisher: { '@id': 'https://badge-api.gosuda.org/#organization' },
        image: 'https://badge-api.gosuda.org/og-image.png',
        about: [
          { '@type': 'Thing', name: 'SVG badge generation' },
          { '@type': 'Thing', name: 'HTTP API' },
          { '@type': 'Thing', name: 'Immutable caching and ETags' }
        ]
      },
      {
        '@type': 'BreadcrumbList',
        itemListElement: [
          {
            '@type': 'ListItem',
            position: 1,
            name: 'Tiny Badge',
            item: 'https://badge-api.gosuda.org/'
          },
          {
            '@type': 'ListItem',
            position: 2,
            name: 'API Reference',
            item: 'https://badge-api.gosuda.org/reference/'
          }
        ]
      }
    ]
  };
</script>

<Seo
  title="SVG Badge API Reference — Tiny Badge"
  description={seoDescription}
  path="/reference/"
  type="article"
  modifiedTime="2026-08-11T00:00:00Z"
  structuredData={referenceStructuredData}
/>

<SiteNav active="reference" bind:dismissed={bannerDismissed} />

<main id="top" class:banner-dismissed={bannerDismissed} class="reference-page">
  <section class="reference-hero section-shell">
    <h1>Everything the SVG badge API understands.</h1>
    <p>
      Query and path endpoints, eleven rendering styles including three fixed-size classics with auto-scaling vector text, exact 50–300% sizing, adjustable letter spacing, four color controls, Unicode-aware width, immutable cache headers, ETags, and practical examples.
    </p>
  </section>

  <div class="reference-layout section-shell">
    <aside class="reference-toc">
      <strong>On this page</strong>
      <nav aria-label="API reference sections">
        {#each sections as section}<a href={`#${section[0]}`}>{section[1]}</a>{/each}
      </nav>
      <a class="btn reference-maker-link" href="/#designer">Open the maker</a>
    </aside>

    <div class="reference-content">
      <section id="overview" class="reference-section">
        <h2>Overview</h2>
        <p>
          Tiny Badge returns a deterministic SVG image from the label, message, style, size, letter spacing, and colors encoded in the request URL. Unicode graphemes are measured with <code>go-runewidth</code> and UAX #29 segmentation, successful responses use long-lived immutable caching and content-derived ETags, and the same URL always describes the same badge.
        </p>
        <div class="reference-callout">
          <strong>Production origin</strong>
          <code>{productionOrigin}</code>
          <p>Every example on this page is ready to use against the live Tiny Badge service.</p>
        </div>
      </section>

      <section id="endpoints" class="reference-section">
        <h2>Endpoints</h2>
        <h3>Query endpoint</h3>
        <p>Use query parameters when text contains spaces, punctuation, or a literal <code>.svg</code> suffix.</p>
        <div class="snippet"><pre><code>GET /badge.svg
HEAD /badge.svg

{queryExample}</code></pre><CopyTextButton text={queryExample} label="Copy URL" /></div>
        <h3>Path endpoint</h3>
        <p>Use the shorter path form for simple label and message values.</p>
        <div class="snippet"><pre><code>GET /badge/:label/:message
HEAD /badge/:label/:message

{pathExample}</code></pre><CopyTextButton text={pathExample} label="Copy URL" /></div>
        <p>Use <code>_</code> as the path label for a badge without a left segment.</p>
        <div class="inline-example"><code>/badge/_/available.svg</code></div>
        <p>The path form removes a final <code>.svg</code> suffix. The query form preserves it.</p>
      </section>

      <section id="parameters" class="reference-section">
        <h2>Parameters</h2>
        <div class="parameter-list">
          {#each parameters as parameter}
            <article>
              <div><code>{parameter.name}</code><span>{parameter.required}</span></div>
              <p>{parameter.detail}</p>
              <small>Default: {parameter.fallback}</small>
            </article>
          {/each}
        </div>
        <p>Spaces and punctuation must be URL encoded. For example, <code>ready to ship</code> becomes <code>ready%20to%20ship</code>.</p>
        <p><code>letterSpacing</code> is measured in SVG user units: one unit is one output pixel at 100% scale. Nonzero spacing replaces the fixed word artwork in matching <code>click-here</code> and <code>best-viewed</code> reference phrases with adjustable vector text.</p>
      </section>

      <section id="styles" class="reference-section">
        <h2>Styles</h2>
        <div class="reference-style-grid">
          {#each styles as style}
            <article>
              <img src={styleSampleURL(style)} alt={`${style[0]} badge example`} />
              <h3><code>{style[0]}</code></h3>
              <p>{style[1]}</p>
            </article>
          {/each}
        </div>
      </section>

      <section id="colors" class="reference-section">
        <h2>Colors</h2>
        <p>The API uses 3-digit or 6-digit hexadecimal values without the leading <code>#</code>.</p>
        <div class="inline-example"><code>labelColor=292724&amp;color=d6ef53&amp;labelTextColor=ffffff&amp;textColor=292724</code></div>
        <p>The maker keeps colors internally as RGB. Hex, HSL, and OKLCH channel inputs are converted to RGB before the URL is assembled.</p>
        <p>If a leading <code>#</code> is included, encode it as <code>%23</code> so it does not become a URL fragment.</p>
        <h3>Named colors</h3>
        <div class="color-name-list">{#each namedColors as color}<code>{color}</code>{/each}</div>
      </section>

      <section id="examples" class="reference-section">
        <h2>Examples</h2>
        <h3>Markdown</h3>
        <div class="snippet"><pre><code>{markdownExample}</code></pre><CopyTextButton text={markdownExample} /></div>
        <h3>HTML</h3>
        <div class="snippet"><pre><code>{htmlExample}</code></pre><CopyTextButton text={htmlExample} /></div>
        <h3>cURL</h3>
        <div class="snippet"><pre><code>{curlExample}</code></pre><CopyTextButton text={curlExample} /></div>
      </section>

      <section id="responses" class="reference-section">
        <h2>Responses and caching</h2>
        <p>Successful requests return an SVG image with long-lived immutable caching and a content-derived ETag.</p>
        <div class="snippet"><pre><code>{responseHeaders}</code></pre><CopyTextButton text={responseHeaders} label="Copy headers" /></div>
        <p>Send the ETag in <code>If-None-Match</code>. An unchanged badge returns <code>304 Not Modified</code> without a response body.</p>
      </section>

      <section id="errors" class="reference-section">
        <h2>Errors and limits</h2>
        <p>Invalid input returns <code>400 Bad Request</code> with <code>Cache-Control: no-store</code>.</p>
        <div class="error-grid">
          <article><code>message</code><p>Required and limited to 128 characters.</p></article>
          <article><code>label</code><p>Optional and limited to 64 characters.</p></article>
          <article><code>style</code><p>Must match one of the documented style names.</p></article>
          <article><code>size</code><p>Must be an integer percentage from 50 through 300.</p></article>
          <article><code>letterSpacing</code><p>Must be a decimal number from -1 through 3.</p></article>
          <article><code>colors</code><p>Must be a supported name or a 3-digit or 6-digit hex value.</p></article>
        </div>
        <h3>Health check</h3>
        <div class="inline-example"><code>GET /healthz → 200 OK → ok</code></div>
      </section>
    </div>
  </div>
</main>

<SiteFooter />
