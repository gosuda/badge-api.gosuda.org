<script>
  import { colorModes, formatColor, parseColor } from './color.js';

  const colorModeLabels = {
    hex: 'Quick',
    rgb: 'Light mix',
    hsl: 'Hue mix',
    oklch: 'Even light'
  };

  let {
    label,
    value,
    help = 'Pick the one that feels right.',
    oncolorchange = () => {}
  } = $props();

  let mode = $state('hex');
  let text = $state('');
  let touched = $state(false);
  let focused = $state(false);
  let invalid = $state(false);

  $effect(() => {
    if (!focused) text = formatColor(value, mode);
  });

  function selectMode(event) {
    mode = event.currentTarget.value;
    text = formatColor(value, mode);
    invalid = false;
  }

  function chooseNative(event) {
    const next = event.currentTarget.value;
    oncolorchange(next);
    text = formatColor(next, mode);
    invalid = false;
  }

  function editValue(event) {
    text = event.currentTarget.value;
    const next = parseColor(text, mode);
    if (next) {
      oncolorchange(next);
      invalid = false;
    } else if (touched) {
      invalid = true;
    }
  }

  function blurValue() {
    focused = false;
    touched = true;
    const next = parseColor(text, mode);
    invalid = !next;
    if (!next) return;
    oncolorchange(next);
    text = formatColor(next, mode);
  }
</script>

<div class:error={invalid} class="color-editor">
  <div class="color-editor__head">
    <label class="field-label" for={`${label}-value`}>{label}</label>
    <select class="space-select" aria-label={`${label} recipe`} value={mode} onchange={selectMode}>
      {#each colorModes as colorMode}
        <option value={colorMode}>{colorModeLabels[colorMode]}</option>
      {/each}
    </select>
  </div>
  <div class="color-editor__row">
    <input
      class="native-color"
      type="color"
      value={value}
      aria-label={`Open the ${label} swatch`}
      oninput={chooseNative}
    />
    <input
      id={`${label}-value`}
      class="text-input color-editor__value"
      type="text"
      value={text}
      aria-invalid={invalid}
      aria-describedby={`${label}-help`}
      onfocus={() => (focused = true)}
      oninput={editValue}
      onblur={blurValue}
      spellcheck="false"
      autocomplete="off"
    />
  </div>
  <p id={`${label}-help`} class="field-help" class:error-text={invalid}>
    {invalid ? 'We couldn’t read that color. Try the swatch or check the numbers.' : help}
  </p>
</div>
