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
    { id: 'flatbar', name: 'Flatbar', note: 'Big badge energy', sample: 'loud' },
    { id: 'old-school', name: 'Old School 80×15', note: 'Split-panel pixel type', sample: 'button' },
    { id: 'click-here', name: 'Click Here 88×31', note: 'Raised, loud, clickable', sample: 'here' },
    { id: 'best-viewed', name: 'Best Viewed 88×31', note: 'Chrome-era nostalgia', sample: 'chrome' }
  ];

  function sampleURL(style) {
    const presets = {
      flatbar: { label: 'mood', message: 'loud', size: '100', labelColor: '292724', color: '7c5cff', textColor: 'ffffff' },
      'old-school': { label: 'pixel', message: 'button', size: '200', labelColor: 'ff5a18', color: 'a8a979', textColor: 'ffffff' },
      'click-here': { label: 'click', message: 'here', size: '200', labelColor: '292724', color: 'd6ef53', textColor: '292724' },
      'best-viewed': { label: 'viewed with', message: 'chrome', size: '200', labelColor: '292724', color: 'd6ef53', textColor: '292724' }
    };
    const preset = presets[style.id] ?? {
      label: 'mood',
      message: style.sample,
      size: '100',
      labelColor: '292724',
      color: 'd6ef53',
      textColor: '292724'
    };
    const query = new URLSearchParams({
      label: preset.label,
      message: preset.message,
      style: style.id,
      size: preset.size,
      labelColor: preset.labelColor,
      color: preset.color,
      labelTextColor: 'ffffff',
      textColor: preset.textColor
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
        <span class="style-option__preview"><img src={sampleURL(style)} alt="" loading="lazy" decoding="async" /></span>
        <span class="style-option__copy">
          <strong>{style.name}</strong>
          <small>{style.note}</small>
        </span>
      </label>
    {/each}
  </div>
</fieldset>
