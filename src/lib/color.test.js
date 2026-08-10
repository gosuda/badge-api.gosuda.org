import assert from 'node:assert/strict';
import test from 'node:test';

import {
  colorModes,
  hexToRgb,
  rgbFromColorMode,
  rgbToHex,
  valuesForColorMode
} from './color.js';

test('all color models resolve to canonical integer RGB', () => {
  const expected = { r: 214, g: 239, b: 83 };
  assert.deepEqual(hexToRgb('#d6ef53'), expected);
  assert.deepEqual(rgbFromColorMode({ r: 214, g: 239, b: 83 }, 'rgb'), expected);
  assert.deepEqual(rgbFromColorMode(valuesForColorMode(expected, 'hsl'), 'hsl'), expected);
  assert.deepEqual(rgbFromColorMode(valuesForColorMode(expected, 'oklch'), 'oklch'), expected);
  assert.equal(rgbToHex(expected), '#d6ef53');
});

test('channel-based modes expose separate editable values', () => {
  assert.deepEqual(colorModes.find((mode) => mode.id === 'rgb').channels.map((channel) => channel.key), ['r', 'g', 'b']);
  assert.deepEqual(colorModes.find((mode) => mode.id === 'hsl').channels.map((channel) => channel.key), ['h', 's', 'l']);
  assert.deepEqual(colorModes.find((mode) => mode.id === 'oklch').channels.map((channel) => channel.key), ['l', 'c', 'h']);
});

test('invalid channel values do not produce an RGB color', () => {
  assert.equal(rgbFromColorMode({ r: 256, g: 0, b: 0 }, 'rgb'), null);
  assert.equal(rgbFromColorMode({ h: 20, s: 101, l: 50 }, 'hsl'), null);
  assert.equal(rgbFromColorMode({ l: 80, c: 0.8, h: 90 }, 'oklch'), null);
});
