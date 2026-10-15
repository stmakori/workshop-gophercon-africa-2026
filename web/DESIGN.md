# Design: Safari Rally Kenya · Wildlife stage

A WRC broadcast and timing screen crossed with the Kenyan flag. The sightings are "wildlife spotted
along the rally route": Kenya only, the animals of the Naivasha, Hell's Gate and Soysambu stages.
Dark and fast, but easy to read on a phone.

## Colours (`web/assets/css/app.css`, `:root`)

| Token | Hex | Use |
|-------|-----|-----|
| `--carbon` | `#0B0B0C` | Page background, under a fine carbon-fibre weave |
| `--carbon-2` / `--carbon-3` | `#141416` / `#1D1D21` | Panels and entry rows |
| `--line` | `#2C2C32` | Borders |
| `--flag-black` | `#000000` | Timing board, footer |
| `--flag-red` | `#BB0000` | Fills only: buttons, section numbers, live cell, chips. Always white text on it |
| `--flag-green` | `#006600` | Fills only: entry edge, year bars, button hover |
| `--flag-white` | `#FFFFFF` | Fimbriation in the stripes |
| `--earth` / `--dust` | `#A4471E` / `#E89A5F` | Red-earth gravel. `--dust` is the link and accent text colour |
| `--text` / `--muted` | `#F4F1EC` / `#A9A39B` | Body text and labels |
| `--red-text` / `--green-text` | `#FF5A4E` / `#4CC46A` | When red or green must be text on dark (the flag hexes are too dark to read) |

## Lettering

- Headings, numbers, nav, chips: **Barlow Condensed**, italic 700–800, uppercase. It's slanted
  like car liveries and stage boards.
- Body: **Barlow** 400–600.
- Numbers use `font-variant-numeric: tabular-nums`, so they line up like a timing board.

## Shapes and motifs

- `--slant: -10deg`. Applied with `skewX` to section numbers, the "Kenya" title block, the year bars and buttons.
- Buttons are skewed boxes with sharp corners (not `clip-path`, which cuts the outline button's sides). `--radius` is 4px wherever there's a radius.
- Livery strip under the header: diagonal black, red and green flag stripes with white edges.
- Chequered flag: a small square before each section heading, and a band above the footer.
- Sections are numbered like special stages: SS1 Route, SS2 Results, SS3 Standings, SS4 The rally.

## Components

- **Hero**: event name and dates, "Naivasha · Hell's Gate · Soysambu", and a four-cell timing board
  (days, stages, km, on the route). The "on the route" cell is red, as the "live" cell. The photo is the
  2024 Safari (Radziejka, CC BY-SA 4.0), served at 960 or 1280 px through `srcset`. It has a slanted
  frame with red stripes behind it, and is credited under it.
- **Entry** (`components/sighting_card.templ`): a results row, ranked by distance to the nearest stage.
  It shows the position (P01… on the route, a grey "OFF" with a grey edge off it), a stage badge, the first
  photo (4:3, click to enlarge), the animal, its scientific name, the date, place and crew (observers),
  the count, and the stage conditions (weather) on a dust-edged strip. Up to three more photos sit
  below as thumbnails. The left edge is green, and red on hover or when targeted from the map. On a
  phone it stacks.
- **Map**: Leaflet with dark tiles. Chequered stage flags have a green ring, the Power Stage a red ring and
  the service park a green dot. A dashed dust-coloured line joins them in running order. Sightings on the
  route are red dots and those off it grey. It frames the stages, not the whole country. A legend sits
  under it.
- **Stage badge** (`components/stage_badge.templ`): a green chip with a chequered square, such as
  "2.0 km from Soysambu". Off the route it's grey: "Off route · 63 km from Kedong".
- **Standings**: "Most seen" and "Busiest stages" bars with positions, and "Results by year" as a table with skewed
  green bars. The best year is red.
- **From the rally**: freely licensed Safari Rally photos from Wikimedia Commons, each credited.

- **Stage conditions** (`components/conditions.templ`): the weather as a WRC conditions panel, with
  the sky, temperature, rain and wind in a black strip with a dust-coloured edge. The stage report
  shows the full version, which adds feels-like, gusts, humidity and cloud.
- **3D wheel** (`assets/js/scenes/rally-wheel.js`): a gravel rally wheel with a white six-spoke rim,
  a red hub and a green caliper. It rolls out of the banner photo's corner. It's low poly, pauses off
  screen, uses low-power WebGL, and draws one still frame under reduced motion.

## Pages

| URL | Page | Look |
|-----|------|------|
| `/` | Home | Banner, route map, stage results, standings summary, rally photos |
| `/sightings/{id}` | Stage report | Big photo with a red underline, a timing board (seen, count), facts, full stage conditions, a small dark map, more photos |
| `/stats` | Championship standings | A six-cell timing board of totals, then panels: most seen (wide), results by year, busiest stages, top crews |
| `/gallery` | Photo wall | Square tiles, each captioned with the animal (linking to its report), the photographer and the licence. The caption edges cycle green, red and white |
| anything else | 404 | "Off the road, wrong stage", a button home, a chequered flag |

The nav holds Home, Stats and Gallery. The current page is red with a green underline (`aria-current="page"`).

## Data

- `routeFilter()` in `web/handlers/home.go` lists Kenya (`country.Kenya`) and the route animals: zebra,
  giraffe, buffalo, impala, warthog, wildebeest, hippo, greater flamingo, spotted hyena and rhinoceros.
  No limit. Home, stats and gallery all use it.
- The stages are in `web/views/components/rally.go`, in running order. Each position is the centre of
  its area, not its roads, and the footer and README say so. A sighting is on the route within
  `RouteKm` (25 km) of a stage.

## Rules

- Rally photos must be free to reuse (CC0 / CC BY / CC BY-SA), with credit and a licence link.
  No press or agency photos.
- Never red text in `#BB0000` or green text in `#006600` on the dark background: they fail contrast.
- Motion respects `prefers-reduced-motion`.
- Pages declare what they need in `layout.BaseProps` (`Map`, `Scene`, `Filter`), so Leaflet, three.js
  and the filter script load only where they're used.
- No `aria-label` on `p`, `dl`, `dd` or `div` without a role. Use visually hidden text instead.
- The footer says it's a fan site, not affiliated with the WRC, the FIA or the organisers.
