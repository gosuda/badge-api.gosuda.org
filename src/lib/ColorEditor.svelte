<script>
  import {
    colorModes,
    hexToRgb,
    rgbFromColorMode,
    rgbToHex,
    valuesForColorMode
  } from './color.js';

  let {
    label,
    value,
    help = 'Every recipe is converted to RGB before the badge link is made.',
    oncolorchange = () => {}
  } = $props();

  let mode = $state('hex');
  let channels = $state({});
  let activeField = $state(null);
  let touched = $state(false);
  let invalid = $state(false);

  const editorId = $derived(label.toLowerCase().replace(/[^a-z0-9]+/g, '-'));
  const selectedMode = $derived(colorModes.find((candidate) => candidate.id === mode) ?? colorModes[0]);

  $effect(() => {
    if (activeField === null) channels = valuesForColorMode(value, mode);
  });

  function selectMode(event) {
    mode = event.currentTarget.value;
    channels = valuesForColorMode(value, mode);
    activeField = null;
    invalid = false;
  }

  function chooseNative(event) {
    const next = hexToRgb(event.currentTarget.value);
    if (!next) return;
    oncolorchange(next);
    channels = valuesForColorMode(next, mode);
    invalid = false;
  }

  function editChannel(channel, event) {
    channels = { ...channels, [channel.key]: event.currentTarget.value };
    const next = rgbFromColorMode(channels, mode);
    if (next) {
      oncolorchange(next);
      invalid = false;
    } else if (touched) {
      invalid = true;
    }
  }

  function blurChannel() {
    activeField = null;
    touched = true;
    const next = rgbFromColorMode(channels, mode);
    invalid = !next;
    if (!next) return;
    oncolorchange(next);
    channels = valuesForColorMode(next, mode);
  }
</script>

<div class:error={invalid} class="color-editor">
  <div class="color-editor__head">
    <span class="field-label">{label}</span>
    <select class="space-select" aria-label={`${label} color model`} value={mode} onchange={selectMode}>
      {#each colorModes as colorMode}
        <option value={colorMode.id}>{colorMode.label}</option>
      {/each}
    </select>
  </div>
  <div class="color-editor__row">
    <input
      class="native-color"
      type="color"
      value={rgbToHex(value)}
      aria-label={`Open the ${label} swatch`}
      oninput={chooseNative}
    />
    <div
      class:single={selectedMode.channels.length === 1}
      class="color-channel-grid"
      role="group"
      aria-label={`${label} ${selectedMode.label} channels`}
    >
      {#each selectedMode.channels as channel}
        <label class="color-channel" for={`${editorId}-${mode}-${channel.key}`}>
          <span>{channel.label}</span>
          <span class="color-channel__control">
            {#if channel.prefix}<span aria-hidden="true">{channel.prefix}</span>{/if}
            <input
              id={`${editorId}-${mode}-${channel.key}`}
              class="color-channel__input"
              type={channel.type}
              inputmode={channel.inputmode}
              min={channel.min}
              max={channel.max}
              step={channel.step}
              maxlength={channel.maxlength}
              value={channels[channel.key] ?? ''}
              aria-label={`${label} ${channel.label}`}
              aria-invalid={invalid}
              aria-describedby={`${editorId}-help`}
              onfocus={() => (activeField = channel.key)}
              oninput={(event) => editChannel(channel, event)}
              onblur={blurChannel}
              spellcheck="false"
              autocomplete="off"
            />
            {#if channel.suffix}<span aria-hidden="true">{channel.suffix}</span>{/if}
          </span>
        </label>
      {/each}
    </div>
  </div>
  <p id={`${editorId}-help`} class="field-help" class:error-text={invalid}>
    {invalid ? 'That channel is outside its usable range.' : help}
  </p>
</div>
