<script>
  const origin = 'https://badge-api.gosuda.org';
  const image = `${origin}/og-image.png`;

  let {
    title,
    description,
    path = '/',
    structuredData,
    type = 'website',
    modifiedTime = ''
  } = $props();

  const canonical = $derived(`${origin}${path}`);
  const jsonLd = $derived(JSON.stringify(structuredData).replaceAll('<', '\\u003c'));
</script>

<svelte:head>
  <title>{title}</title>
  <meta name="description" content={description} />
  <meta name="robots" content="index, follow, max-image-preview:large, max-snippet:-1, max-video-preview:-1" />
  <meta name="author" content="Tiny Badge" />
  <meta name="application-name" content="Tiny Badge" />
  <link rel="canonical" href={canonical} />
  <link rel="alternate" hreflang="en" href={canonical} />
  <link rel="alternate" hreflang="x-default" href={canonical} />
  <link rel="alternate" type="text/plain" href={`${origin}/llms.txt`} title="Tiny Badge LLM-readable summary" />
  <link rel="icon" href="/favicon.svg" />

  <meta property="og:title" content={title} />
  <meta property="og:description" content={description} />
  <meta property="og:type" content={type} />
  <meta property="og:url" content={canonical} />
  <meta property="og:site_name" content="Tiny Badge" />
  <meta property="og:locale" content="en_US" />
  <meta property="og:image" content={image} />
  <meta property="og:image:secure_url" content={image} />
  <meta property="og:image:type" content="image/png" />
  <meta property="og:image:width" content="1200" />
  <meta property="og:image:height" content="630" />
  <meta property="og:image:alt" content="Tiny Badge SVG badge generator and API" />
  {#if modifiedTime}<meta property="article:modified_time" content={modifiedTime} />{/if}

  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content={title} />
  <meta name="twitter:description" content={description} />
  <meta name="twitter:image" content={image} />
  <meta name="twitter:image:alt" content="Tiny Badge SVG badge generator and API" />

  {@html `<script type="application/ld+json">${jsonLd}</script>`}
</svelte:head>
