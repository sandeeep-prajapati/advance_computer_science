# Chapter 8: UI & Asset Integration

## Two UI Styles

PHP frameworks usually support two frontend styles.

```text
Server-rendered UI
├── Laravel Blade
├── Symfony/Pimcore Twig
└── React widgets mounted inside pages

Independent frontend
├── React SPA
├── Next.js app
├── TypeScript client
└── PHP backend as API
```

---

## npm and Vite

Composer handles PHP. npm handles frontend assets.

```text
composer install
→ PHP packages and autoloading

npm install
→ JS/CSS dependencies

npm run build
→ compiled assets in public/build
```

Example `package.json`:

```json
{
  "private": true,
  "scripts": {
    "dev": "vite",
    "build": "vite build"
  },
  "devDependencies": {
    "vite": "^5.0.0",
    "typescript": "^5.0.0",
    "react": "^18.0.0",
    "react-dom": "^18.0.0"
  }
}
```

---

## React Widget in PHP Page

Twig/Blade renders a mount point:

```html
<div id="product-editor" data-product-id="1001"></div>
<script type="module" src="/build/product-editor.js"></script>
```

React mounts there:

```tsx
import { createRoot } from 'react-dom/client';

const el = document.getElementById('product-editor');

if (el) {
  createRoot(el).render(
    <ProductEditor productId={el.dataset.productId!} />
  );
}
```

This is useful for Pimcore admin extensions and Laravel/Symfony admin pages.

---

## Asset Publishing

Packages/bundles may contain assets.

```text
BlogBundle/
└── Resources/assets/
    ├── blog-admin.tsx
    └── blog-admin.css
```

Command:

```bash
php console assets:publish
```

Output:

```text
public/bundles/blog/
├── blog-admin.js
└── blog-admin.css
```

---

## Backend-to-Frontend Contract

Do not let frontend depend on database tables directly.

Good:

```text
React -> HTTP API -> Application Service -> Repository
```

Bad:

```text
React assumes exact database columns and hidden backend rules
```

