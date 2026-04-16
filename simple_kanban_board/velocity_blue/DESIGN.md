# The Editorial Engine: Design System Documentation

## 1. Overview & Creative North Star
**The Creative North Star: The Architectural Workspace**
This design system moves away from the "sticky note on a wall" metaphor of traditional Kanban boards. Instead, it adopts an **Architectural Workspace** philosophy. We treat the mobile screen as a curated editorial layout where productivity is driven by clarity, tonal depth, and intentional white space rather than rigid borders and heavy shadows. 

By leveraging asymmetrical header placements, overlapping surfaces, and a sophisticated "No-Line" aesthetic, we create a tool that feels like a premium workspace rather than a generic utility. The goal is to reduce cognitive load by using "visual breathing room" and a high-contrast typographic scale to guide the eye toward the most critical tasks.

## 2. Colors & Surface Logic
The palette is rooted in professional blues and greys, but its execution is what defines the premium feel.

*   **The "No-Line" Rule:** Under no circumstances should 1px solid borders be used to separate columns or sections. Boundaries are created through tonal shifts. For example, the main application background uses `surface`, while a Kanban column "track" uses `surface-container-low`.
*   **Surface Hierarchy & Nesting:** Use the surface tiers to create physical depth.
    *   **Base Layer:** `surface` (The canvas).
    *   **Middle Layer:** `surface-container-low` (The Kanban Column container).
    *   **Top Layer:** `surface-container-lowest` (The individual Task Card).
*   **The "Glass & Gradient" Rule:** Floating elements, such as the bottom navigation or a "Quick Add" button, should utilize a Glassmorphic effect. Use `surface_container_lowest` at 85% opacity with a `backdrop-blur(20px)` to allow task colors to bleed through subtly.
*   **Signature Textures:** For primary actions (e.g., "Complete Task"), do not use flat fills. Use a subtle linear gradient from `primary` to `primary_container` to provide a sense of "liquid" depth.

## 3. Typography
We pair the structural precision of **Inter** with the editorial character of **Manrope**.

*   **Display & Headlines (Manrope):** Use `headline-sm` for column titles. These should be set with slightly tighter letter spacing to feel like a magazine header.
*   **The Title Scale (Inter):** Task titles must use `title-sm` or `title-md`. They are the "hero" of the card.
*   **Body & Labels (Inter):** Descriptions use `body-md`. For metadata (dates, tags), use `label-sm` in `on-surface-variant`.
*   **Hierarchy via Scale:** To emphasize importance without adding weight, jump two levels in the scale (e.g., pairing a `title-lg` header with `label-md` metadata) to create a high-end, asymmetrical rhythm.

## 4. Elevation & Depth
In this design system, elevation is a property of light and material, not just "drop shadows."

*   **Tonal Layering:** Avoid traditional shadows for cards. A `surface-container-lowest` card sitting on a `surface-container-low` column provides enough contrast for the eye.
*   **Ambient Shadows:** If a card is being "dragged," apply an extra-diffused shadow: `box-shadow: 0 20px 40px rgba(25, 28, 30, 0.06)`. Note the use of the `on-surface` color for the shadow tint rather than pure black.
*   **The "Ghost Border" Fallback:** If a task requires a "Selected" state, do not use a thick stroke. Use a "Ghost Border" by applying a 1px stroke of `outline-variant` at 20% opacity.
*   **Tactile Feedback:** When a user interacts with a card, the surface should transition from `surface-container-lowest` to `surface-bright`, mimicking the way light hits a raised object.

## 5. Components

### Task Cards
*   **Container:** `surface-container-lowest`, rounded at `xl` (0.75rem).
*   **Layout:** No dividers. Use 16px of internal padding (`body-md` spacing). 
*   **Priority Indicators:** Instead of large banners, use a "Glass Pill." A small, semi-transparent chip using `tertiary_container` for High Priority, with text in `on-tertiary_fixed_variant`.

### Kanban Columns
*   **Background:** `surface-container-low`, rounded at `xl`. 
*   **Header:** Asymmetrical alignment. The title (`headline-sm`) should be left-aligned with a 4px vertical accent bar of `primary` to the left.
*   **Spacing:** Use `1.5rem` (24px) between cards to maintain the "Editorial" airy feel.

### Action Buttons
*   **Primary:** Gradient of `primary` to `primary_container`, `full` roundedness, white text (`on_primary`).
*   **Secondary:** `surface-container-highest` background with `on-surface` text. No border.

### Status Chips
*   **Structure:** Use `label-md` for text. 
*   **Coloring:** Use "Muted Alerts." For a 'Done' status, use `primary_fixed` background with `on_primary_fixed_variant` text. It signals success without the "vibrating" neon green of standard UI.

### Input Fields
*   **Style:** Minimalist. No bounding box. Use a `surface-variant` bottom-only highlight (2px) that expands to `primary` on focus.
*   **Placeholder:** `on-surface-variant` at 50% opacity in `body-md`.

## 6. Do's and Don'ts

*   **DO:** Use white space as a structural element. If two items feel cluttered, increase the gap before reaching for a divider line.
*   **DO:** Use `surface-tint` sparingly to highlight active states or progress bars.
*   **DON'T:** Use pure black (#000000) for text. Always use `on-surface` or `on-surface-variant` to maintain the soft, high-end look.
*   **DON'T:** Use default "Medium" or "Large" shadows. Stick to the tonal layering principles defined in Section 4.
*   **DO:** Ensure that all interactive elements have a minimum touch target of 44x44dp, even if the visual element (like a priority dot) is smaller.