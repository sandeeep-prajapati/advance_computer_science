# Chapter 4: Assets & Digital Asset Management (DAM)

## What are Assets?

**Assets** are all the media files in Pimcore — images, videos, PDFs, documents, audio, and any binary file. They are stored in Pimcore's **Digital Asset Management (DAM)** system, organized in a hierarchical folder tree just like a file system.

Every asset is a Pimcore element with:
- A unique ID
- Metadata fields (alt text, copyright, license, custom fields)
- Version history (Pimcore auto-saves older versions when you re-upload)
- Permission controls (role-based access)

---

## Step 1: Upload Assets

1. In Admin Panel → **Assets** (left sidebar)
2. Right-click → **Add Folder** (e.g., `Product Images`)
3. Inside the folder, **drag & drop** files from your computer

> You can upload: `laptop.jpg`, `manual.pdf`, `promo.mp4`, etc.

---

## Step 2: Metadata & Tagging

Each asset can have custom metadata:

1. Select an uploaded asset (e.g., `laptop.jpg`)
2. Go to the **Properties / Metadata** tab
3. Add fields:
   - `Alt Text` = "Laptop X1 Product Image"
   - `Copyright` = "© Company 2025"
   - `Tags` = Laptop, Electronics, Product

> Metadata makes assets searchable and reusable across projects.

---

## Step 3: Organize Assets

Best practices for large media libraries:

| Practice | Why |
|---|---|
| Use logical folder hierarchies | `/products/electronics/phones` avoids dumping everything in `/` |
| Define required metadata | Enforce license, usage rights fields |
| Use tags for quick search | Find all "PDF manuals" in one click |
| Use versions | Pimcore keeps older versions when you re-upload |
| Enable CDN | Serve assets through CDN for better performance |
| Use Image Thumbnails | Never serve original huge images directly |
| Use S3 / MinIO for huge files | Videos, CAD models belong in object storage |

---

## Step 4: Image Thumbnails

Pimcore has a powerful **Thumbnail** system that lets you define image transformation presets.

### Define a Thumbnail in Admin UI

1. Go to **Settings → Thumbnails**
2. Click **Add** → name it `product_preview`
3. Add transformations:
   - Resize to 800×600
   - Crop center
   - Output format: **WebP** (best for web)
4. Save

### Use in Twig Templates

```twig
<img src="{{ asset.thumbnail('product_preview') }}" alt="{{ product.name }}">
```

### Use in PHP

```php
$asset = \Pimcore\Model\Asset\Image::getById(123);

// Use a named preset
$thumbPath = $asset->getThumbnail('product_preview')->getPath();
echo $thumbPath;

// Or define on-the-fly (not stored in settings)
$thumbConfig = new \Pimcore\Model\Asset\Image\Thumbnail\Config();
$thumbConfig->addItem('resize', ['width' => 400, 'height' => 300, 'forceResize' => true]);
$thumbConfig->setFormat('webp');
$thumbConfig->setQuality(85);

$thumbPath = $asset->getThumbnail($thumbConfig)->getPath();
```

---

## Step 5: Connect Assets to Data Objects

### Method A: Asset Field in Product Class

1. Go to Settings → Data Objects → Classes → `Product`
2. Add fields:
   - `Image` → **Image** type (for one main image)
   - `Gallery` → **Image Gallery** (for multiple images)
   - `Documents` → **Many-to-Many Object Relation** with Asset (for PDFs, videos)
3. Save & Publish the class

Now open a Product object and drag & drop assets directly into the image fields.

### Method B: Check Usage from the Asset Side

1. Open an Asset (e.g., `laptop.jpg`)
2. Go to the **Relations** tab
3. See which Products, Pages, or Objects are using this asset

---

## Serving Large Media Libraries

| Strategy | Detail |
|---|---|
| Always use thumbnails | Never serve full-resolution originals to the browser |
| Different sizes for different contexts | Small `100×100` for grid preview, `1200×800` for web display |
| Use WebP/AVIF | Significantly smaller file sizes than JPEG/PNG |
| Cache headers + CDN | CloudFront, Akamai, Fastly for global distribution |
| S3 for large files | Videos, CAD models → Pimcore S3 Asset Store plugin |
| Elasticsearch indexing | Enable for searching/browsing large asset libraries |

---

## Example: A Complete Product Media Setup

Product: **Laptop X1**
```
SKU: LAP123
Price: 1200
Image: laptop.jpg (Asset)
Manual: manual.pdf (Asset)
Promo Video: promo.mp4 (Asset)
```

In the frontend (CMS page or API), you can fetch and display:
- Laptop Name + Price (from Product Object)
- Product Image (from Asset — resized to thumbnail)
- Downloadable Manual PDF (from Asset)
- Promo Video embed (from Asset)

---

## Summary

- **Assets** = the media library (files, images, videos, PDFs)
- **Metadata** = makes assets searchable and manageable
- **Thumbnails** = define reusable image transformation presets
- **Relations** = link assets to products, pages, and any object
- **Best Practice:** organized folder structure + CDN + WebP thumbnails = fast, clean media delivery
