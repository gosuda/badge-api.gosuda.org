<script>
  let { value = $bindable('flat') } = $props();

  const styles = [
    { id: 'flat', name: 'Flat', note: 'Easygoing', sample: 'easy' },
    { id: 'flat-square', name: 'Flat square', note: 'No-nonsense', sample: 'steady' },
    { id: 'plastic', name: 'Plastic', note: 'A little retro', sample: 'retro' },
    { id: 'round', name: 'Round', note: 'Soft at the edges', sample: 'soft' },
    { id: 'outline', name: 'Outline', note: 'Neatly dressed', sample: 'tidy' },
    { id: 'neon', name: 'Neon', note: 'Out after dark', sample: 'glow' },
    { id: 'glass', name: 'Glass', note: 'Feeling fancy', sample: 'fancy' },
    { id: 'flatbar', name: 'Flatbar', note: 'Big badge energy', sample: 'loud' }
  ];

  function sampleURL(style) {
    const query = new URLSearchParams({
      label: 'mood',
      message: style.sample,
      style,
      labelColor: '292724',
      color: style === 'flatbar' ? '7c5cff' : 'd6ef53',
      labelTextColor: 'ffffff',
      textColor: style === 'flatbar' ? 'ffffff' : '292724'
    });
    return `/badge.svg?${query}`;
  }
</script>

<fieldset class="style-picker">
  <legend class="field-label">Pick an outfit</legend>
  <div class="style-picker__grid">
    {#each styles as style}
      <label class:selected={value === style.id} class="style-option">
        <input type="radio" name="badge-style" value={style.id} bind:group={value} />
        <span class="style-option__preview"><img src={sampleURL(style)} alt="" /></span>
        <span class="style-option__copy">
          <strong>{style.name}</strong>
          <small>{style.note}</small>
        </span>
      </label>
    {/each}
  </div>
</fieldset>
