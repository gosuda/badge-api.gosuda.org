import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    adapter: adapter({
      pages: 'dist/frontend',
      assets: 'dist/frontend',
      precompress: true,
      strict: true
    }),
    prerender: {
      handleHttpError: ({ path, message }) => {
        if (path.startsWith('/badge.svg')) return;
        throw new Error(message);
      }
    }
  }
};

export default config;
