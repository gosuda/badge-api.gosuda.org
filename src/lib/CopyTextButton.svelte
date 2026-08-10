<script>
  let { text, label = 'Copy' } = $props();
  let state = $state('idle');

  async function copyText() {
    if (state === 'loading') return;
    state = 'loading';
    try {
      await navigator.clipboard.writeText(text);
      state = 'success';
      window.setTimeout(() => (state = 'idle'), 2200);
    } catch {
      state = 'error';
      window.setTimeout(() => (state = 'idle'), 2800);
    }
  }
</script>

<button
  class:success={state === 'success'}
  class:error={state === 'error'}
  class="snippet-copy"
  type="button"
  disabled={state === 'loading'}
  onclick={copyText}
>
  {#if state === 'loading'}
    Copying…
  {:else if state === 'success'}
    Copied
  {:else if state === 'error'}
    Try again
  {:else}
    {label}
  {/if}
</button>
