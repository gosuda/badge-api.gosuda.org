<script>
  import CopyTextButton from '$lib/CopyTextButton.svelte';
  import SiteFooter from '$lib/SiteFooter.svelte';
  import SiteNav from '$lib/SiteNav.svelte';

  let bannerDismissed = $state(false);

  const queryExample = 'https://badge.example.com/badge.svg?label=build&message=passing&style=flatbar';
  const pathExample = 'https://badge.example.com/badge/release/stable.svg';
  const markdownExample = '![Build status](https://badge.example.com/badge.svg?label=build&message=passing)';
  const htmlExample = '<img src="https://badge.example.com/badge.svg?label=build&message=passing" alt="Build: passing">';
  const curlExample = "curl -o badge.svg 'https://badge.example.com/badge.svg?label=build&message=passing&style=flatbar'";
  const responseHeaders = `Content-Type: image/svg+xml; charset=utf-8
Cache-Control: public, max-age=315360000, immutable
CDN-Cache-Control: public, max-age=315360000, immutable
Surrogate-Control: public, max-age=315360000, immutable
Access-Control-Allow-Origin: *
ETag: "<content hash>"`;

  const parameters = [
    { name: 'message', required: 'Required', fallback: '—', detail: 'Main badge text. Maximum 128 characters.' },
    { name: 'label', required: 'Optional', fallback: 'Empty', detail: 'Left-side badge text. Maximum 64 characters.' },
    { name: 'style', required: 'Optional', fallback: 'flat', detail: 'Selects one of the eight rendering styles.' },
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
    ['flatbar', 'A 28 px uppercase bar inspired by larger badge styles.']
  ];

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
</script>

<svelte:head>
  <title>API Reference — Tiny Badge</title>
  <meta name="description" content="Detailed Tiny Badge API endpoints, parameters, styles, colors, examples, caching, and errors." />
  <link rel="icon" href="/favicon.svg" />
</svelte:head>

<SiteNav active="reference" bind:dismissed={bannerDismissed} />

<main id="top" class:banner-dismissed={bannerDismissed} class="reference-page">
  <section class="reference-hero section-shell">
    <h1>Everything the badge link understands.</h1>
    <p>
      Two endpoint shapes, eight looks, four color controls, and enough examples to keep guesswork out of the address bar.
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
          Tiny Badge returns an SVG image from the words, style, and colors encoded in the request URL. The same URL always describes the same badge, which makes it suitable for profiles, project pages, Markdown, and ordinary image tags.
        </p>
        <div class="reference-callout">
          <strong>Placeholder origin</strong>
          <code>https://badge.example.com</code>
          <p>Replace this origin with the hostname of your deployment.</p>
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
      </section>

      <section id="styles" class="reference-section">
        <h2>Styles</h2>
        <div class="reference-style-grid">
          {#each styles as style}
            <article>
              <img src={`/badge.svg?label=style&message=${style[0]}&style=${style[0]}`} alt={`${style[0]} badge example`} />
              <h3><code>{style[0]}</code></h3>
              <p>{style[1]}</p>
            </article>
          {/each}
        </div>
      </section>

      <section id="colors" class="reference-section">
        <h2>Colors</h2>
        <p>Use 3-digit or 6-digit hexadecimal values without the leading <code>#</code>.</p>
        <div class="inline-example"><code>labelColor=292724&amp;color=d6ef53&amp;labelTextColor=ffffff&amp;textColor=292724</code></div>
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
          <article><code>colors</code><p>Must be a supported name or a 3-digit or 6-digit hex value.</p></article>
        </div>
        <h3>Health check</h3>
        <div class="inline-example"><code>GET /healthz → 200 OK → ok</code></div>
      </section>
    </div>
  </div>
</main>

<SiteFooter />
