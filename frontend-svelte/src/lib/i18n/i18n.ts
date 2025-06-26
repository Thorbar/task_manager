// i18n.ts

import { register, init, getLocaleFromNavigator, locale  } from 'svelte-i18n';
import { browser } from '$app/environment';

register('es', () => import('./es.json'));
register('ca', () => import('./ca.json'));
register('en', () => import('./en.json'));
register('ja', () => import('./ja.json'));

init({
  fallbackLocale: 'en',
  initialLocale: browser ? localStorage.getItem('lang') || getLocaleFromNavigator() : 'en'
});

if (browser) {
  locale.subscribe((lang) => {
    if (lang) localStorage.setItem('lang', lang);
  });
}

locale.set(getLocaleFromNavigator());

export { locale };
