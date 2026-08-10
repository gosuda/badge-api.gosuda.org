<script>
  import SiteFooter from '$lib/SiteFooter.svelte';
  import SiteNav from '$lib/SiteNav.svelte';

  let bannerDismissed = $state(false);
  let billing = $state('monthly');
  let annual = $derived(billing === 'annual');

  const plans = [
    {
      name: 'Solo',
      description: 'For one person with one excellent idea and at least three favorite colors.',
      perSeat: false,
      tone: 'pear',
      cta: 'Choose Solo',
      extras: ['Extra independent vibe', 'Extra “I made this” vibe']
    },
    {
      name: 'Team',
      description: 'For people who enjoy sharing links, opinions, and the occasional polite nudge.',
      perSeat: true,
      tone: 'cyan',
      featured: true,
      cta: 'Gather the team',
      extras: ['Extra together vibe', 'Extra quick-sync vibe', 'Extra good-lunch vibe']
    },
    {
      name: 'Studio',
      description: 'For a room full of taste, tabs, and someone asking whether the green feels greener.',
      perSeat: false,
      tone: 'lavender',
      cta: 'Go Studio',
      extras: ['Extra polished vibe', 'Extra good-chair vibe', 'Extra one-more-version vibe', 'Extra tasteful-nod vibe']
    },
    {
      name: 'Enterprise',
      description: 'For very large feelings, very small badges, and a purchasing process with a calendar invite.',
      perSeat: true,
      tone: 'ink',
      enterprise: true,
      cta: 'Choose Enterprise',
      extras: ['Extra enterprise vibe', 'Extra procurement vibe', 'Extra lanyard vibe', 'Extra quarterly vibe', 'Extra approved-font vibe']
    }
  ];

  const sharedBenefits = ['Unlimited badges', 'Every look and color tool', 'Links that are ready to travel'];
</script>

<svelte:head>
  <title>Pricing — Tiny Badge</title>
  <meta name="description" content="Four Tiny Badge plans. Every one costs exactly zero dollars." />
  <link rel="icon" href="/favicon.svg" />
</svelte:head>

<SiteNav active="pricing" bind:dismissed={bannerDismissed} />

<main id="top" class:banner-dismissed={bannerDismissed} class="pricing-page">
  <section class="pricing-hero section-shell">
    <div>
      <h1>Four plans.<br />One suspiciously identical price.</h1>
      <p>
        Choose by team size, personal mythology, or whichever card looks nicest next to your current browser tab.
      </p>
    </div>
    <div class="billing-box">
      <div class="billing-switch" role="group" aria-label="Billing period">
        <button class:active={billing === 'monthly'} type="button" aria-pressed={billing === 'monthly'} onclick={() => (billing = 'monthly')}>Monthly</button>
        <button class:active={billing === 'annual'} type="button" aria-pressed={billing === 'annual'} onclick={() => (billing = 'annual')}>Annual <span class="billing-discount">40% off</span></button>
      </div>
      <p>
        {annual
          ? 'Annual billing takes 40% off. Your new total is still $0. Beautiful.'
          : 'Monthly billing is $0 now, $0 next month, and $0 after that.'}
      </p>
    </div>
  </section>

  <section class="plan-grid section-shell" aria-label="Pricing plans">
    {#each plans as plan}
      <article class:featured={plan.featured} class:enterprise={plan.enterprise} class={`plan-card plan-card--${plan.tone}`}>
        {#if plan.featured}<span class="plan-badge">Most popular at lunch</span>{/if}
        {#if plan.enterprise}<span class="plan-badge">Most enterprise</span>{/if}
        <header>
          <h2>{plan.name}</h2>
          <p>{plan.description}</p>
        </header>
        <div class="plan-price">
          <strong>$0</strong>
          <span>{plan.perSeat ? 'per seat / month' : 'per month'}</span>
          {#if annual}<small>40% annual discount included</small>{/if}
        </div>
        <div class="plan-benefits">
          <h3>Everything useful</h3>
          <ul>
            {#each sharedBenefits as benefit}<li>{benefit}</li>{/each}
          </ul>
          <h3>The plan difference</h3>
          <ul>
            {#each plan.extras as extra}<li>{extra}</li>{/each}
          </ul>
        </div>
        <a class:btn--inverse={plan.enterprise} class="btn plan-cta" href="/#designer">{plan.cta}</a>
      </article>
    {/each}
  </section>

  <section class="pricing-note section-shell">
    <div>
      <h2>The math is simple because there is hardly any math.</h2>
      <p>
        Monthly is $0. Annual is 40% less. Team and Enterprise are priced per seat, so every additional seat costs another $0.
      </p>
    </div>
    <a class="btn btn--push" href="/#designer">Make the first badge</a>
  </section>
</main>

<SiteFooter />
