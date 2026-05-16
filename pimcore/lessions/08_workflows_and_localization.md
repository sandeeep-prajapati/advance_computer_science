# Chapter 8: Workflows & Localization

---

## Part 1: Workflows

### What are Workflows?

Workflows in Pimcore define a **series of states** that content moves through before being published. They enforce quality control and governance.

- Built on **Symfony's workflow component**
- Can be applied to Data Objects, Assets, or Documents
- Define who can do what at each stage

**Example use cases:**
- Product content: Draft → Review → Approved → Published
- Marketing images: Submitted → QA Review → Approved for use
- Multi-language content: Translated → Native Speaker Review → Published

---

### Workflow Components

| Component | Description |
|---|---|
| **Places (States)** | Stages of the process (Draft, Review, Approved) |
| **Transitions** | Rules to move from one state to another |
| **Roles & Permissions** | Who can execute each transition |
| **UI Actions** | Buttons that appear in the Pimcore Admin panel |

---

### Example: Product Approval Workflow

**Goal:** Product moves through: `Draft → Review → Approved → Published`

**File:** `var/config/workflow.yaml`

```yaml
workflows:
  product_approval:
    label: "Product Approval Workflow"
    type: 'state_machine'
    supports:
      - 'Pimcore\Model\DataObject\Product'
    places:
      - draft
      - review
      - approved
      - published
    initial_place: draft
    transitions:
      submit_for_review:
        from: draft
        to: review
      approve:
        from: review
        to: approved
      publish:
        from: approved
        to: published
      send_back:
        from: review
        to: draft
    metadata:
      submit_for_review:
        label: "Submit for Review"
      approve:
        label: "Approve"
      publish:
        label: "Publish"
      send_back:
        label: "Send Back to Draft"
```

---

### Assign Workflow to Users

In **Settings → Users/Roles**, assign permissions:

| Role | Allowed Transitions |
|---|---|
| Content Editor | Draft → Review (`submit_for_review`) |
| Manager | Review → Approved (`approve`) |
| Admin | Approved → Published (`publish`) |

---

### Workflow in Action

1. A new Product is created → starts in **Draft**
2. Editor clicks **Submit for Review** → moves to **Review**
3. Manager reviews content → clicks **Approve** → moves to **Approved**
4. Admin clicks **Publish** → Product goes **live**

In the Product editor screen, users see workflow action buttons based on their permissions.

---

### Benefits of Workflows

- Ensures **quality control** before content goes live
- Provides **audit trails** (who approved what, when)
- Makes **content governance** easier for distributed teams

---

## Part 2: Localization & Multilingual

### How Pimcore Handles Localization

| Component | How It Works |
|---|---|
| Data Objects | Fields can be made **localized** (per language) |
| Documents | Language-specific page trees (`/en`, `/de`, `/fr`) |
| Assets | Metadata translations (alt text, titles) |
| System | Configurable language list |

---

### Step 1: Configure System Languages

1. Admin → **Settings → System → Localization & Internationalization**
2. Add languages: `en`, `de`, `fr`
3. Save

---

### Step 2: Add Localized Fields to Product Class

1. Go to Settings → Data Objects → Classes → `Product`
2. Add a **Localized Fields** container
3. Inside it, add fields like:
   - `name` (Input)
   - `description` (WYSIWYG)

---

### Step 3: Fill Data Per Language

Open a Product (e.g., Laptop X1). You'll see language tabs (EN, DE, FR):

| Language | Name | Description |
|---|---|---|
| EN | "Laptop X1" | "Lightweight business laptop" |
| DE | "Laptop X1" | "Leichter Business-Laptop" |
| FR | "Laptop X1" | "Ordinateur portable léger pour entreprise" |

The same product object holds all translations — no duplication needed.

---

### Step 4: Multilingual Document Trees

```
/documents
  /en
    /about-us
    /products/laptop-x1
  /de
    /ueber-uns
    /produkte/laptop-x1
  /fr
    /a-propos
    /produits/laptop-x1
```

- Each language folder = one localized website tree
- Pimcore auto-detects language via URL prefix (`/en/`, `/de/`)

---

### Step 5: Language Switcher in Twig

```twig
<ul>
    <li>
        <a href="{{ pimcore_url(document, {'_locale': 'en'}) }}">English</a>
    </li>
    <li>
        <a href="{{ pimcore_url(document, {'_locale': 'de'}) }}">Deutsch</a>
    </li>
    <li>
        <a href="{{ pimcore_url(document, {'_locale': 'fr'}) }}">Français</a>
    </li>
</ul>
```

---

### Step 6: Asset Metadata Translations

1. Select an image in Assets
2. Go to **Metadata** tab
3. Enter translated values per language:
   - Alt Text (EN): "Laptop X1 Product Image"
   - Alt Text (DE): "Laptop X1 Produktbild"

---

### Localization URL Examples

| URL | Language | Content |
|---|---|---|
| `/en/products/laptop-x1` | English | Product in English |
| `/de/produkte/laptop-x1` | German | Same product, German translation |
| `/fr/produits/laptop-x1` | French | Same product, French translation |

All powered by **one Pimcore instance** with structured multilingual data — no content duplication.
