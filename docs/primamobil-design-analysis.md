# PrimaMobil.id Design Analysis

**Analysis Date:** 2025-11-16
**Source:** https://primamobil.id/
**Purpose:** Adopt design patterns into auto-lmk public website

---

## Executive Summary

PrimaMobil.id features a clean, modern design focused on user experience with minimal distractions. The design emphasizes:
- Simple navigation with prominent search
- Hero-driven homepage with clear call-to-action
- Feature cards highlighting value propositions
- Clean car grid layout with essential info
- Detailed car pages with large photos and WhatsApp CTA
- Mobile-responsive design throughout

---

## Color Scheme

### Primary Colors
- **Background**: `oklch(0.99 0 0)` - Very light background, almost white (#FCFCFC)
- **Text**: `oklch(0.15 0.01 250)` - Dark text for readability (#1F1F1F)
- **Primary Accent**: Orange (`#FF5722` or similar) - Used in logo, badges
- **Success/CTA**: Green (`#25D366`) - WhatsApp button color

### Secondary Colors
- **Card Background**: White with subtle shadow
- **Border**: Light gray for separators
- **Badge Background**: Light blue/gray for tags/pills
- **Hover States**: Subtle opacity or color shifts

---

## Typography

### Font Stack
```css
font-family: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji"
```

### Font Sizes & Weights
- **H1 (Hero)**: 30px, font-weight: 700 (bold)
- **H2 (Section headings)**: ~24px, font-weight: 700
- **H3 (Car titles)**: ~18-20px, font-weight: 600
- **Body**: 16px, font-weight: 400
- **Small text**: 14px for metadata

### Text Hierarchy
- Large hero titles with bold weight
- Section headings in bold
- Car titles semi-bold
- Body text regular weight
- Metadata (year, km, etc.) in smaller gray text

---

## Layout & Spacing

### Container Width
- Max width: ~1200px centered
- Responsive breakpoints for mobile/tablet/desktop

### Spacing Patterns
- Generous whitespace between sections
- Consistent padding: 16px, 24px, 32px, 48px
- Card spacing: 16-24px gaps in grid

### Grid Layouts
- **Homepage car grid**: 2 columns on desktop, 1 on mobile
- **Car listing**: 2 columns with sidebar filter
- **Responsive**: Stack on smaller screens

---

## Component Analysis

### 1. Navigation Bar
**Structure:**
- Logo (orange "A" icon + "AutoLeads" text) on left
- Centered search bar with placeholder "Search cars"
- "Cari Mobil" link on right
- Fixed/sticky on scroll
- Clean white background with subtle shadow

**Key Features:**
- Minimal navigation items
- Prominent search functionality
- Mobile-responsive (hamburger menu)

---

### 2. Hero Section (Homepage)

**Layout:**
- Full-width section with subtle gradient or solid color
- Centered content with max-width constraint

**Content:**
- Small eyebrow text: "Mobil Bekas Berkualitas"
- Large H1: "Temukan Mobil Impian Kamu"
- Subtitle: "Jelajahi koleksi mobil bekas pilihan kami..."
- Primary CTA button: "Lihat Semua Mobil" (orange/primary color)

**Visual Design:**
- Clean, minimal design
- Strong typography hierarchy
- Clear call-to-action

---

### 3. Feature Cards (3-Column)

**Layout:**
- 3 equal-width cards in a row
- Responsive: Stack on mobile

**Card Content:**
- Icon (checkmark, money, grid icons)
- H3 heading: "Kualitas Terjamin", "Harga Terbaik", "Pilihan Lengkap"
- Description text below

**Visual Design:**
- White background
- Subtle shadow or border
- Icon at top, centered or left-aligned
- Consistent padding

---

### 4. Car Card (Grid Item)

**Layout:**
- Vertical card with image on top
- Content section below
- Hover effect: subtle shadow or scale

**Content Structure:**
```
┌─────────────────────┐
│                     │
│     Car Image       │
│                     │
│  #MZ01 (badge)      │
└─────────────────────┘
  Mazda Mazda 2 R Hatchback
  2014 • Merah Metalik

  Rp 135.000.000

  82.000 km    Matic

  🔧 Mesin SKYACTIV 1.5L Irit BBM
  ⚙️ Transmisi Matic 4-Speed Smooth
  +2 lainnya
```

**Visual Design:**
- Image: aspect ratio ~4:3, object-fit cover
- Badge: top-right overlay on image (dark with white text)
- Title: bold, truncate if too long
- Metadata: smaller gray text
- Feature tags: blue/gray pills
- "+X lainnya" for additional features

---

### 5. Car Detail Page

**Layout:**
- Large photo gallery with main image + thumbnail strip
- Photo counter: "1 / 11"
- Navigation arrows for prev/next
- Specifications table below
- Sidebar or separate section for "Tertarik?" CTA

**Specifications Table:**
```
┌────────────────────────────┐
│ Kode      │ #MZ01          │
│ Tahun     │ 2014           │
│ Kilometer │ 82.000 km      │
│ Transmisi │ Matic          │
│ Warna     │ Merah Metalik  │
│ Bahan Bakar│ Bensin        │
└────────────────────────────┘
```

**Feature List:**
- Bullet points or pills for main features
- "Catatan Kondisi" section with condition notes

**CTA Box:**
- Heading: "Tertarik?"
- Description: "Hubungi kami untuk informasi lebih lanjut..."
- Large green WhatsApp button
- Trust indicators:
  - "Biasanya merespons dalam hitungan menit"
  - "Pelayanan profesional terjamin"

---

### 6. Filter Sidebar (Car Listing)

**Layout:**
- Left sidebar on desktop, collapsible on mobile
- Heading: "Filter"

**Filter Options:**
- **Urutkan**: Dropdown (Terbaru, Harga Tertinggi, etc.)
- **Merek**: Dropdown (Semua Merek, Toyota, etc.)
- **Transmisi**: Dropdown (Semua, Matic, Manual)
- **Rentang Tahun**: Two dropdowns (from - to)
- **Rentang Harga**: Two dropdowns (from - to)

**Visual Design:**
- Compact dropdowns with consistent styling
- Labels above each dropdown
- White background, subtle borders

---

### 7. Footer

**Content:**
- Copyright: "© 2025 AutoLeads. All rights reserved."
- Branding: "Powered by Lumiku AutoLeads • Trusted by thousands"
- Minimal design, centered text

**Visual Design:**
- Light background or white
- Small text
- Link to Lumiku website
- No complex footer navigation

---

## UI Patterns & Best Practices

### 1. Consistent Card Design
- All cards (feature cards, car cards) use similar shadow/border
- Consistent padding and spacing
- Hover effects for interactive elements

### 2. Badge System
- Car code badge (#MZ01) prominently displayed
- Positioned top-right on car images
- Dark background with white text for contrast

### 3. Price Display
- Large, bold price in Indonesian Rupiah format
- Blue color to stand out
- Positioned prominently on car cards

### 4. Feature Pills/Tags
- Blue/gray rounded pills for features
- "+X lainnya" to indicate more features
- Truncation to avoid clutter

### 5. WhatsApp Integration
- Prominent green WhatsApp button
- Pre-filled message with car details
- Trust indicators near CTA

### 6. Image Optimization
- WebP format for images
- Responsive images with proper sizing
- Lazy loading for performance

### 7. Responsive Design
- Mobile-first approach
- Stacked layout on small screens
- Touch-friendly button sizes

---

## Component Mapping: PrimaMobil → auto-lmk

### Files to Modify/Create

| PrimaMobil Component | auto-lmk Template | Status |
|---------------------|-------------------|--------|
| Navigation Bar | `templates/layouts/base.html` | Needs update |
| Hero Section | `templates/pages/home.html` | Needs creation |
| Feature Cards (3) | `templates/pages/home.html` | Needs creation |
| Car Grid | `templates/pages/home.html` | Partial (needs styling) |
| Car Card | `templates/components/car_card.html` | Needs creation |
| Filter Sidebar | `templates/pages/cars.html` | Needs creation |
| Car Detail Hero | `templates/pages/car_detail.html` | Needs major update |
| Photo Gallery | `templates/pages/car_detail.html` | Needs creation |
| Specs Table | `templates/pages/car_detail.html` | Partial (needs styling) |
| WhatsApp CTA Box | `templates/components/whatsapp_cta.html` | Needs creation |
| Footer | `templates/layouts/base.html` | Needs update |

---

## Implementation Priority

### Phase 1: Foundation (1-2 hours)
1. ✅ Update color scheme in Tailwind config or CSS variables
2. ✅ Update typography in base layout
3. ✅ Simplify navigation bar
4. ✅ Update footer to minimal design

### Phase 2: Homepage (2-3 hours)
1. ✅ Create hero section component
2. ✅ Create 3 feature cards section
3. ✅ Update car grid layout
4. ✅ Create car card component
5. ✅ Add bottom CTA section

### Phase 3: Car Listing (2 hours)
1. ✅ Create filter sidebar
2. ✅ Update car grid for listing page
3. ✅ Add sort/filter functionality

### Phase 4: Car Detail (2-3 hours)
1. ✅ Create photo gallery with lightbox
2. ✅ Update specifications table styling
3. ✅ Create WhatsApp CTA box
4. ✅ Add feature pills/tags
5. ✅ Add "Catatan Kondisi" section

### Phase 5: Testing & Polish (1 hour)
1. ✅ Test responsive design
2. ✅ Test all interactive elements
3. ✅ Optimize images
4. ✅ Performance check

**Total Estimated Time:** 8-11 hours

---

## Technical Implementation Notes

### Tailwind CSS Classes

Based on the design analysis, key Tailwind classes to use:

```css
/* Container */
.container { max-width: 1200px; margin: 0 auto; padding: 0 1rem; }

/* Typography */
.hero-title { @apply text-3xl font-bold text-gray-900; }
.section-title { @apply text-2xl font-bold text-gray-900; }
.car-title { @apply text-lg font-semibold text-gray-800; }

/* Cards */
.card { @apply bg-white rounded-lg shadow-sm border border-gray-100 p-6; }
.card:hover { @apply shadow-md; }

/* Buttons */
.btn-primary { @apply bg-orange-500 text-white px-6 py-3 rounded-lg font-medium hover:bg-orange-600; }
.btn-whatsapp { @apply bg-green-500 text-white px-6 py-3 rounded-lg font-medium hover:bg-green-600; }

/* Badge */
.badge { @apply absolute top-2 right-2 bg-gray-900 text-white px-3 py-1 rounded text-sm font-medium; }

/* Feature Pills */
.feature-pill { @apply inline-block bg-blue-50 text-blue-700 px-3 py-1 rounded-full text-sm mr-2 mb-2; }
```

### Alpine.js for Interactivity

```javascript
// Photo gallery
Alpine.data('photoGallery', () => ({
  currentPhoto: 0,
  totalPhotos: 11,
  next() { this.currentPhoto = (this.currentPhoto + 1) % this.totalPhotos },
  prev() { this.currentPhoto = (this.currentPhoto - 1 + this.totalPhotos) % this.totalPhotos },
  select(index) { this.currentPhoto = index }
}))

// Filter sidebar (mobile toggle)
Alpine.data('filterSidebar', () => ({
  open: false,
  toggle() { this.open = !this.open }
}))
```

### Go Template Variables Needed

Ensure these variables are passed from handlers:

```go
// Homepage
type HomePageData struct {
    FeaturedCars []Car
    ShowHero     bool
    // ... existing fields
}

// Car Detail
type CarDetailData struct {
    Car           Car
    Photos        []string
    Features      []string
    WhatsAppLink  string
    // ... existing fields
}

// Car Listing
type CarListingData struct {
    Cars         []Car
    Filters      FilterOptions
    TotalCount   int
    // ... existing fields
}
```

---

## Design Principles to Maintain

1. **Simplicity First**: Don't overcomplicate the UI
2. **Mobile-Responsive**: Test on all screen sizes
3. **Fast Loading**: Optimize images, minimize CSS/JS
4. **Clear CTAs**: Make WhatsApp contact obvious
5. **Trust Indicators**: Show quality, reliability messaging
6. **Consistent Spacing**: Use Tailwind spacing scale
7. **Readable Typography**: Maintain good contrast and font sizes
8. **Accessible**: Ensure semantic HTML and ARIA labels

---

## Conclusion

The PrimaMobil.id design is clean, modern, and focused on conversions. Key takeaways:

- **Minimal navigation** reduces friction
- **Hero section** immediately communicates value
- **Feature cards** build trust
- **Clean car cards** show essential info at a glance
- **Large photos** are crucial for car sales
- **WhatsApp CTA** makes contact easy
- **Responsive design** works on all devices

By adopting these patterns, auto-lmk will have a professional, user-friendly interface optimized for car sales.
