const clamp = (value, min = 0, max = 1) => Math.min(max, Math.max(min, value));
const round = (value, digits = 0) => {
  const factor = 10 ** digits;
  return Math.round(value * factor) / factor;
};

export const colorModes = [
  {
    id: 'hex',
    label: 'Hex',
    channels: [
      { key: 'hex', label: 'Hex', type: 'text', inputmode: 'text', prefix: '#', maxlength: 6 }
    ]
  },
  {
    id: 'rgb',
    label: 'RGB',
    channels: [
      { key: 'r', label: 'R', type: 'number', inputmode: 'numeric', min: 0, max: 255, step: 1 },
      { key: 'g', label: 'G', type: 'number', inputmode: 'numeric', min: 0, max: 255, step: 1 },
      { key: 'b', label: 'B', type: 'number', inputmode: 'numeric', min: 0, max: 255, step: 1 }
    ]
  },
  {
    id: 'hsl',
    label: 'HSL',
    channels: [
      { key: 'h', label: 'H', type: 'number', inputmode: 'decimal', min: 0, max: 360, step: 0.1, suffix: '°' },
      { key: 's', label: 'S', type: 'number', inputmode: 'decimal', min: 0, max: 100, step: 0.1, suffix: '%' },
      { key: 'l', label: 'L', type: 'number', inputmode: 'decimal', min: 0, max: 100, step: 0.1, suffix: '%' }
    ]
  },
  {
    id: 'oklch',
    label: 'OKLCH',
    channels: [
      { key: 'l', label: 'L', type: 'number', inputmode: 'decimal', min: 0, max: 100, step: 0.1, suffix: '%' },
      { key: 'c', label: 'C', type: 'number', inputmode: 'decimal', min: 0, max: 0.5, step: 0.001 },
      { key: 'h', label: 'H', type: 'number', inputmode: 'decimal', min: 0, max: 360, step: 0.1, suffix: '°' }
    ]
  }
];

export function normalizeHex(value) {
  const input = String(value).trim().replace(/^#/, '');
  if (/^[0-9a-f]{3}$/i.test(input)) {
    return `#${input.split('').map((part) => part + part).join('').toLowerCase()}`;
  }
  if (/^[0-9a-f]{6}$/i.test(input)) {
    return `#${input.toLowerCase()}`;
  }
  return null;
}

export function normalizeRgb({ r, g, b }) {
  const channels = [r, g, b].map(Number);
  if (channels.some((value) => !Number.isFinite(value))) return null;
  return {
    r: Math.round(clamp(channels[0], 0, 255)),
    g: Math.round(clamp(channels[1], 0, 255)),
    b: Math.round(clamp(channels[2], 0, 255))
  };
}

export function hexToRgb(hex) {
  const normalized = normalizeHex(hex);
  if (!normalized) return null;
  return {
    r: Number.parseInt(normalized.slice(1, 3), 16),
    g: Number.parseInt(normalized.slice(3, 5), 16),
    b: Number.parseInt(normalized.slice(5, 7), 16)
  };
}

export function rgbToHex(rgb) {
  const normalized = normalizeRgb(rgb) ?? { r: 0, g: 0, b: 0 };
  const channel = (value) => value.toString(16).padStart(2, '0');
  return `#${channel(normalized.r)}${channel(normalized.g)}${channel(normalized.b)}`;
}

export function rgbToHsl(rgb) {
  const normalized = normalizeRgb(rgb) ?? { r: 0, g: 0, b: 0 };
  const red = normalized.r / 255;
  const green = normalized.g / 255;
  const blue = normalized.b / 255;
  const max = Math.max(red, green, blue);
  const min = Math.min(red, green, blue);
  const lightness = (max + min) / 2;
  const delta = max - min;
  if (delta === 0) return { h: 0, s: 0, l: lightness * 100 };
  const saturation = delta / (1 - Math.abs(2 * lightness - 1));
  let hue;
  if (max === red) hue = 60 * (((green - blue) / delta) % 6);
  else if (max === green) hue = 60 * ((blue - red) / delta + 2);
  else hue = 60 * ((red - green) / delta + 4);
  if (hue < 0) hue += 360;
  return { h: hue, s: saturation * 100, l: lightness * 100 };
}

export function hslToRgb({ h, s, l }) {
  const hue = ((Number(h) % 360) + 360) % 360;
  const saturation = clamp(Number(s) / 100);
  const lightness = clamp(Number(l) / 100);
  const chroma = (1 - Math.abs(2 * lightness - 1)) * saturation;
  const x = chroma * (1 - Math.abs(((hue / 60) % 2) - 1));
  const match = lightness - chroma / 2;
  let red = 0;
  let green = 0;
  let blue = 0;
  if (hue < 60) [red, green] = [chroma, x];
  else if (hue < 120) [red, green] = [x, chroma];
  else if (hue < 180) [green, blue] = [chroma, x];
  else if (hue < 240) [green, blue] = [x, chroma];
  else if (hue < 300) [red, blue] = [x, chroma];
  else [red, blue] = [chroma, x];
  return normalizeRgb({ r: (red + match) * 255, g: (green + match) * 255, b: (blue + match) * 255 });
}

const srgbToLinear = (value) => {
  const channel = value / 255;
  return channel <= 0.04045 ? channel / 12.92 : ((channel + 0.055) / 1.055) ** 2.4;
};

const linearToSrgb = (value) => {
  const channel = value <= 0.0031308 ? 12.92 * value : 1.055 * value ** (1 / 2.4) - 0.055;
  return clamp(channel) * 255;
};

export function rgbToOklch(rgb) {
  const normalized = normalizeRgb(rgb) ?? { r: 0, g: 0, b: 0 };
  const red = srgbToLinear(normalized.r);
  const green = srgbToLinear(normalized.g);
  const blue = srgbToLinear(normalized.b);
  const l = Math.cbrt(0.4122214708 * red + 0.5363325363 * green + 0.0514459929 * blue);
  const m = Math.cbrt(0.2119034982 * red + 0.6806995451 * green + 0.1073969566 * blue);
  const s = Math.cbrt(0.0883024619 * red + 0.2817188376 * green + 0.6299787005 * blue);
  const lightness = 0.2104542553 * l + 0.793617785 * m - 0.0040720468 * s;
  const a = 1.9779984951 * l - 2.428592205 * m + 0.4505937099 * s;
  const axisB = 0.0259040371 * l + 0.7827717662 * m - 0.808675766 * s;
  const chroma = Math.sqrt(a * a + axisB * axisB);
  let hue = Math.atan2(axisB, a) * 180 / Math.PI;
  if (hue < 0) hue += 360;
  return { l: lightness * 100, c: chroma, h: hue };
}

export function oklchToRgb({ l, c, h }) {
  const lightness = clamp(Number(l) / 100);
  const radians = Number(h) * Math.PI / 180;
  const a = Math.max(0, Number(c)) * Math.cos(radians);
  const axisB = Math.max(0, Number(c)) * Math.sin(radians);
  const lRoot = lightness + 0.3963377774 * a + 0.2158037573 * axisB;
  const mRoot = lightness - 0.1055613458 * a - 0.0638541728 * axisB;
  const sRoot = lightness - 0.0894841775 * a - 1.291485548 * axisB;
  const lValue = lRoot ** 3;
  const mValue = mRoot ** 3;
  const sValue = sRoot ** 3;
  return normalizeRgb({
    r: linearToSrgb(4.0767416621 * lValue - 3.3077115913 * mValue + 0.2309699292 * sValue),
    g: linearToSrgb(-1.2684380046 * lValue + 2.6097574011 * mValue - 0.3413193965 * sValue),
    b: linearToSrgb(-0.0041960863 * lValue - 0.7034186147 * mValue + 1.707614701 * sValue)
  });
}

export function valuesForColorMode(rgb, mode) {
  const normalized = normalizeRgb(rgb) ?? { r: 0, g: 0, b: 0 };
  if (mode === 'hex') return { hex: rgbToHex(normalized).slice(1) };
  if (mode === 'rgb') return normalized;
  if (mode === 'hsl') {
    const hsl = rgbToHsl(normalized);
    return { h: round(hsl.h, 1), s: round(hsl.s, 1), l: round(hsl.l, 1) };
  }
  const oklch = rgbToOklch(normalized);
  return { l: round(oklch.l, 1), c: round(oklch.c, 3), h: round(oklch.h, 1) };
}

export function rgbFromColorMode(values, mode) {
  if (mode === 'hex') return hexToRgb(values.hex ?? '');
  const numbers = Object.fromEntries(Object.entries(values).map(([key, value]) => [key, Number(value)]));
  if (Object.values(numbers).some((value) => !Number.isFinite(value))) return null;
  if (mode === 'rgb') {
    if ([numbers.r, numbers.g, numbers.b].some((value) => value < 0 || value > 255)) return null;
    return normalizeRgb(numbers);
  }
  if (mode === 'hsl') {
    if (numbers.h < 0 || numbers.h > 360 || numbers.s < 0 || numbers.s > 100 || numbers.l < 0 || numbers.l > 100) return null;
    return hslToRgb(numbers);
  }
  if (mode === 'oklch') {
    if (numbers.l < 0 || numbers.l > 100 || numbers.c < 0 || numbers.c > 0.5 || numbers.h < 0 || numbers.h > 360) return null;
    return oklchToRgb(numbers);
  }
  return null;
}
